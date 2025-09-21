package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	Brokers      []string
	BatchSize    int
	BatchTimeout time.Duration
	Async        bool
}

// NewProducer creates a new Kafka producer with the given configuration
func NewProducer(cfg ProducerConfig) (*Producer, error) {
	if cfg.BatchTimeout == 0 {
		cfg.BatchTimeout = time.Millisecond * 100
	}
	if cfg.BatchSize == 0 {
		cfg.BatchSize = 100
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		RequiredAcks: kafka.RequireAll,
		Async:        cfg.Async,
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

// SendMessage sends a message using the Producer
func (p *Producer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Key:   key,
			Value: value,
		},
	)
}

// SendMessages sends multiple messages using the Producer
func (p *Producer) SendMessages(ctx context.Context, topic string, messages []kafka.Message) error {
	return p.writer.WriteMessages(ctx, messages...)
}

type Consumer struct {
	reader      *kafka.Reader
	handler     func(message []byte)
	errorHandler func(err error)
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	Brokers      []string
	GroupID      string
	Topics       []string
	MinBytes     int
	MaxBytes     int
	MaxWait      time.Duration
	StartOffset  int64
	ErrorHandler func(err error)
}

// NewConsumer creates a new consumer with the given configuration
func NewConsumer(cfg ConsumerConfig, handler func(message []byte)) (*Consumer, error) {
	if cfg.MinBytes == 0 {
		cfg.MinBytes = 10e3 // 10KB
	}
	if cfg.MaxBytes == 0 {
		cfg.MaxBytes = 10e6 // 10MB
	}
	if cfg.MaxWait == 0 {
		cfg.MaxWait = time.Second
	}
	if cfg.StartOffset == 0 {
		cfg.StartOffset = kafka.FirstOffset
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(err error) {}
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		GroupID:     cfg.GroupID,
		GroupTopics: cfg.Topics,
		MinBytes:    cfg.MinBytes,
		MaxBytes:    cfg.MaxBytes,
		MaxWait:     cfg.MaxWait,
		StartOffset: cfg.StartOffset,
	})

	return &Consumer{
		reader:       reader,
		handler:      handler,
		errorHandler: cfg.ErrorHandler,
	}, nil
}

// Consume starts consuming messages from Kafka
func (c *Consumer) Consume(ctx context.Context) error {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return err
			}
			c.errorHandler(err)
			continue
		}

		// 处理消息
		c.handler(message.Value)

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

// ConsumerGroup 管理多个消费者的消费者组
type ConsumerGroup struct {
	consumers []*Consumer
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// ConsumerGroupConfig 消费者组配置
type ConsumerGroupConfig struct {
	Brokers      []string
	GroupID      string
	Topics       []string
	NumConsumers int // 消费者数量
	MinBytes     int
	MaxBytes     int
	MaxWait      time.Duration
	StartOffset  int64
}

// NewConsumerGroup 创建一个新的消费者组
func NewConsumerGroup(cfg ConsumerGroupConfig, handler func(message []byte)) (*ConsumerGroup, error) {
	if cfg.NumConsumers <= 0 {
		cfg.NumConsumers = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	group := &ConsumerGroup{
		ctx:    ctx,
		cancel: cancel,
	}

	// 创建多个消费者
	for i := 0; i < cfg.NumConsumers; i++ {
		consumer, err := NewConsumer(ConsumerConfig{
			Brokers:     cfg.Brokers,
			GroupID:     cfg.GroupID,
			Topics:      cfg.Topics,
			MinBytes:    cfg.MinBytes,
			MaxBytes:    cfg.MaxBytes,
			MaxWait:     cfg.MaxWait,
			StartOffset: cfg.StartOffset,
		}, handler)
		if err != nil {
			group.Close()
			return nil, err
		}
		group.consumers = append(group.consumers, consumer)
	}

	return group, nil
}

// Start 启动所有消费者
func (g *ConsumerGroup) Start() error {
	for _, consumer := range g.consumers {
		g.wg.Add(1)
		go func(c *Consumer) {
			defer g.wg.Done()
			err := c.Consume(g.ctx)
			if err != nil && err != context.Canceled {
				// TODO: 处理错误，可以添加错误回调或日志
				return
			}
		}(consumer)
	}
	return nil
}

// Close 关闭消费者组
func (g *ConsumerGroup) Close() error {
	g.cancel()
	g.wg.Wait()

	var lastErr error
	for _, consumer := range g.consumers {
		if err := consumer.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
