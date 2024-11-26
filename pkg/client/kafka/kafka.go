package kafka

import (
	"context"
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
	reader  *kafka.Reader
	handler func(message []byte)
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	Brokers     []string
	GroupID     string
	Topics      []string
	MinBytes    int
	MaxBytes    int
	MaxWait     time.Duration
	StartOffset int64
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
		reader:  reader,
		handler: handler,
	}, nil
}

// Consume starts consuming messages from Kafka
func (c *Consumer) Consume(ctx context.Context) error {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
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
