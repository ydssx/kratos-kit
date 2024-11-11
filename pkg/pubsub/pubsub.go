package pubsub

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type CloudEvent = cloudevents.Event

type EventHandler func(ctx context.Context, event *CloudEvent) error

type Publisher interface {
	PublishMessage(ctx context.Context, subject string, payload interface{}, opts ...Option) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler EventHandler, opts ...Subscription) error
	Unsubscribe(ctx context.Context, subject string) error
	UnsubscribeAll(ctx context.Context) error
}

type Consumer struct {
	Queue       string
	Type        SubscribeType
	Concurrency int
}

type SubscribeType int

const (
	SubscribeTypeQueue SubscribeType = iota
	SubscribeTypeTopic
)

type Subscription func(*Consumer)

// SubscribeToQueue returns a Subscription that sets the Consumer's Queue
// field to the provided queue string.
func SubscribeToQueue(queue string) Subscription {
	return func(c *Consumer) {
		c.Queue = queue
		c.Type = SubscribeTypeQueue
	}
}
