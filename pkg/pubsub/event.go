package pubsub

import (
	"context"
	"time"

	"github.com/ydssx/kratos-kit/pkg/middleware/kratos"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

type (
	eventSource struct{}
	eventType   struct{}
)

func EventSourceFromCtx(ctx context.Context) string {
	t, ok := ctx.Value(eventSource{}).(string)
	if ok {
		return t
	}
	return "api"
}

func EventTypeFromCtx(ctx context.Context) string {
	t, ok := ctx.Value(eventType{}).(string)
	if ok {
		return t
	}
	return "null"
}

type Event struct {
	contentType ContentType
	source      string
	eventType   string
	traceID     string
}

type ContentType string

const (
	ApplicationJSON ContentType = cloudevents.ApplicationJSON
	ApplicationXML  ContentType = cloudevents.ApplicationXML
	TextPlain       ContentType = cloudevents.TextPlain
)

func getTraceIDFromCtx(ctx context.Context) string {
	return kratos.TraceIDFromContext(ctx)
}

func newDefaultEvent(ctx context.Context) *Event {
	return &Event{
		contentType: ApplicationJSON,
		source:      EventSourceFromCtx(ctx),
		eventType:   EventTypeFromCtx(ctx),
		traceID:     getTraceIDFromCtx(ctx),
	}
}

type Option func(*Event)

func WithContentType(t ContentType) Option {
	return func(e *Event) { e.contentType = t }
}

func WithSource(s string) Option {
	return func(e *Event) { e.source = s }
}

func WithEventType(t string) Option {
	return func(e *Event) { e.eventType = t }
}

func NewEvent(ctx context.Context, payload interface{}, opts ...Option) (event.Event, error) {
	ev := newDefaultEvent(ctx)
	for _, v := range opts {
		v(ev)
	}

	e := cloudevents.NewEvent()
	e.SetID(uuid.New().String())
	e.SetType(ev.eventType)
	e.SetTime(time.Now().Local())
	e.SetSource(ev.source)
	e.SetExtension("traceid", ev.traceID)
	err := e.SetData(string(ev.contentType), payload)

	return e, err
}
