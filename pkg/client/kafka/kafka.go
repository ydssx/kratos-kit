package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

// SendMessage sends a message using the Producer.
//
// It takes a string `message` as a parameter and returns an error.
func (p *Producer) SendMessage(topic, message string) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})
	return err
}

type Consumer struct {
	consumer sarama.ConsumerGroup
	handler  sarama.ConsumerGroupHandler
	topics   []string
}

type ConsumerHandler struct {
	handle func(message []byte)
}

// NewConsumer 创建一个新的消费者。它接受 Kafka broker 地址列表、消费者组 ID、订阅的主题列表以及消息处理函数作为参数。
// 它会根据给定的配置创建一个 sarama 消费者,并用提供的处理函数封装一个 ConsumerHandler。
// 最后它会返回一个初始化好的 Consumer 结构体。如果创建消费者失败则返回错误。
func NewConsumer(brokers []string, groupID string, topics []string, handler func(message []byte)) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}
	consumerHandler := &ConsumerHandler{handle: handler}
	return &Consumer{
		consumer: consumer,
		handler:  consumerHandler,
		topics:   topics,
	}, nil
}

func (c *Consumer) Consume() error {
	err := c.consumer.Consume(context.Background(), c.topics, c.handler)
	if err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim consumes messages from a Kafka topic.
//
// It takes in a session from the sarama.ConsumerGroupSession and a claim from
// the sarama.ConsumerGroupClaim. The session is used to mark the consumption of
// messages and the claim is used to retrieve the messages from the Kafka topic.
//
// The function returns an error if there is an issue consuming the messages.
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handle(message.Value)
		session.MarkMessage(message, "")
		select {
		case <-session.Context().Done():
			return session.Context().Err()
		default:
			continue
		}
	}
	return nil
}
