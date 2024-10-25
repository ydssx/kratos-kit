package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

type WatermillPubSub struct {
	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewWatermillPubSub(cli *redis.Client) *WatermillPubSub {
	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: cli,
	}, watermill.NewStdLogger(true, true))
	if err != nil {
		panic(err)
	}
	
	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: cli,
	}, watermill.NewStdLogger(true, true))
	if err != nil {
		panic(err)
	}
	
	return &WatermillPubSub{
		publisher:  publisher,
		subscriber: subscriber,
	}
}

func (ps *WatermillPubSub) PublishMessage(ctx context.Context, subject string, payload interface{}, opts ...Option) error {
	event, err := NewEvent(ctx, payload, opts...)
	if err != nil {
		return err
	}

	data, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)
	return ps.publisher.Publish(subject, msg)
}

func (ps *WatermillPubSub) Subscribe(ctx context.Context, subject string, handler EventHandler, opts ...Subscription) error {
	messages, err := ps.subscriber.Subscribe(ctx, subject)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			event := new(CloudEvent)
			err := event.UnmarshalJSON(msg.Payload)
			if err != nil {
				// 处理解码错误
				continue
			}

			handler(ctx, event)
			msg.Ack()
		}
	}()

	return nil
}

func (ps *WatermillPubSub) SubscribeAsync(ctx context.Context, subject string, handler EventHandler, opts ...Subscription) error {
	return ps.Subscribe(ctx, subject, handler, opts...)
}

func (ps *WatermillPubSub) Unsubscribe(ctx context.Context, subject string) error {
	// Watermill 不支持取消订阅单个主题
	return nil
}

func (ps *WatermillPubSub) UnsubscribeAll(ctx context.Context) error {
	return ps.subscriber.Close()
}

func (ps *WatermillPubSub) Close() error {
	if err := ps.publisher.Close(); err != nil {
		return err
	}
	return ps.subscriber.Close()
}
