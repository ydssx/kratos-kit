package sse

// Event represents an SSE event
type Event struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	ID      string      `json:"id,omitempty"`
	Retry   int         `json:"retry,omitempty"`
}

// NewEvent creates a new SSE event
func NewEvent(eventType string, data interface{}) Event {
	return Event{
		Type: eventType,
		Data: data,
	}
} 