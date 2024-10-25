package pubsub

import (
	"context"
	"sync"

	"github.com/ydssx/kratos-kit/pkg/errors"
	"github.com/ydssx/kratos-kit/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	cli  *redis.Client
	subs map[string]*redis.PubSub
}

// NewRedisPubSub 创建RedisPubSub对象
func NewRedisPubSub(cli *redis.Client) *RedisPubSub {
	return &RedisPubSub{cli: cli, subs: make(map[string]*redis.PubSub)}
}

// PublishMessage publishes a message to the given topic.
// It returns an error if the publish failed.
func (ps *RedisPubSub) PublishMessage(ctx context.Context, subject string, payload interface{}, opts ...Option) error {
	event, err := NewEvent(ctx, payload, opts...)
	if err != nil {
		return err
	}

	message, err := event.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "序列化消息失败")
	}

	err = ps.cli.Publish(context.Background(), subject, message).Err()
	if err != nil {
		return errors.Wrap(err, "发布消息失败")
	}
	return nil
}

// SubscribeToTopic subscribes to the given topic and calls the handler
// function whenever a new message is received on that topic.
func (ps *RedisPubSub) SubscribeToTopic(ctx context.Context, topic string, handler EventHandler, maxConcurrency int) error {
	sub := ps.cli.Subscribe(context.Background(), topic)
	ps.subs[topic] = sub

	ch := sub.Channel()
	semaphore := make(chan struct{}, maxConcurrency)

	go func() {
		for msg := range ch {
			if msg == nil {
				continue
			}

			select {
			case semaphore <- struct{}{}:
				go func(msg *redis.Message) {
					defer func() { <-semaphore }()

					data := new(CloudEvent)
					err := data.UnmarshalJSON([]byte(msg.Payload))
					if err != nil {
						logger.Errorf(ctx, "Failed to unmarshal message: %s", err.Error())
						return
					}

					handler(ctx, data)
				}(msg)
			default:
				// 如果无法立即获取信号量，记录日志并继续
				logger.Warnf(ctx, "Max concurrency reached for topic [%s], message processing delayed", topic)
			}
		}
		logger.Infof(ctx, "Stopped subscribing to messages on topic [%s]", topic)
	}()

	return nil
}

func (ps *RedisPubSub) SubscribeToQueue(ctx context.Context, queue string, handler EventHandler) error {
	// 创建一个同步池来重用 CloudEvent 实例
	pool := &sync.Pool{
		New: func() interface{} {
			return new(CloudEvent)
		},
	}

	// 创建一个错误通道用于处理解码错误
	errCh := make(chan error, 100)
	go func() {
		for err := range errCh {
			// 将错误推送到另一个队列或进行其他处理
			logger.Errorf(ctx, "Failed to unmarshal message: %s", err.Error())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// 上下文被取消,退出循环
			return ctx.Err()
		default:
			// 从对象池获取 CloudEvent 实例
			data := pool.Get().(*CloudEvent)
			msg := ps.cli.BLPop(ctx, 0, queue).Val()
			if msg != nil {
				msg := msg[0]
				err := data.UnmarshalJSON([]byte(msg))
				if err != nil {
					// 将解码错误推送到错误通道
					errCh <- err
					// 将 CloudEvent 实例放回对象池
					pool.Put(data)
					continue
				}

				handler(ctx, data)
				// 将 CloudEvent 实例放回对象池
				pool.Put(data)
			}
		}
	}
}

func (ps *RedisPubSub) Subscribe(ctx context.Context, subject string, handler EventHandler, opts ...Subscription) error {
	consumer := new(Consumer)
	for _, opt := range opts {
		opt(consumer)
	}

	if consumer.Concurrency <= 0 {
		consumer.Concurrency = 10
	}
	if consumer.Type == SubscribeTypeQueue {
		return ps.SubscribeToQueue(ctx, subject, handler) // 订阅队列
	}
	return ps.SubscribeToTopic(ctx, subject, handler, consumer.Concurrency) // 订阅主题
}

func (ps *RedisPubSub) SubscribeAsync(ctx context.Context, subject string, handler EventHandler, opts ...Subscription) error {
	return ps.Subscribe(ctx, subject, handler, opts...)
}

// Close 关闭RedisPubSub对象
func (ps *RedisPubSub) Close() error {
	var errs []error
	for t, v := range ps.subs {
		err := v.Close()
		errs = append(errs, errors.Wrap(err, "关闭主题["+t+"]的订阅失败"))
	}
	return errors.Join(errs...)
}

func (ps *RedisPubSub) Unsubscribe(ctx context.Context, subject string) error {
	sub, ok := ps.subs[subject]
	if !ok {
		return errors.New("主题[" + subject + "]不存在")
	}
	err := sub.Close()
	if err != nil {
		return errors.Wrap(err, "关闭主题["+subject+"]的订阅失败")
	}
	delete(ps.subs, subject)
	return nil
}

func (ps *RedisPubSub) UnsubscribeAll(ctx context.Context) error {
	return ps.Close()
}
