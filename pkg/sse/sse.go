package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Client represents a connected SSE client
type Client struct {
	ID           string
	Messages     chan interface{}
	Disconnected chan struct{}
	writer       http.ResponseWriter
	flusher      http.Flusher
	mu           sync.Mutex
}

// Broker manages SSE clients and message broadcasting
type Broker struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
	events     chan interface{}
}

// NewBroker creates a new SSE broker
func NewBroker() *Broker {
	return &Broker{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		events:     make(chan interface{}, 100),
	}
}

// Start begins the broker's main loop
func (b *Broker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case client := <-b.register:
				b.mu.Lock()
				b.clients[client.ID] = client
				b.mu.Unlock()

			case client := <-b.unregister:
				b.mu.Lock()
				if _, ok := b.clients[client.ID]; ok {
					delete(b.clients, client.ID)
					close(client.Messages)
				}
				b.mu.Unlock()

			case event := <-b.events:
				b.broadcast(event)

			case <-ctx.Done():
				b.mu.Lock()
				for _, client := range b.clients {
					close(client.Messages)
				}
				b.clients = make(map[string]*Client)
				b.mu.Unlock()
				return
			}
		}
	}()
}

// broadcast sends a message to all connected clients
func (b *Broker) broadcast(event interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, client := range b.clients {
		select {
		case client.Messages <- event:
		default:
			// Client message buffer is full
		}
	}
}

// ServeHTTP handles the SSE endpoint
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := &Client{
		ID:           r.RemoteAddr,
		Messages:     make(chan interface{}, 100),
		Disconnected: make(chan struct{}),
		writer:       w,
		flusher:      flusher,
	}

	b.register <- client

	// Clean up when the client disconnects
	defer func() {
		b.unregister <- client
		close(client.Disconnected)
	}()

	// Keep the connection alive
	for {
		select {
		case msg, ok := <-client.Messages:
			if !ok {
				return
			}
			if err := client.Send(msg); err != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

// Send sends a message to the client
func (c *Client) Send(event interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.writer, "data: %s\n\n", data)
	c.flusher.Flush()
	return nil
}

// Publish sends an event to all connected clients
func (b *Broker) Publish(event interface{}) {
	b.events <- event
}