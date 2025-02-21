package message

import (
	"encoding/json"
)

// Message represents a P2P message.
type Message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	Sender string      `json:"sender"`
}

// Serialize serializes a Message to JSON.
func Serialize(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

// Deserialize deserializes a Message from JSON.
func Deserialize(data []byte) (Message, error) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// Router handles message routing to different handlers.
type Router struct {
	handlers map[string]Handler
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]Handler),
	}
}

// RegisterHandler registers a message handler for a specific type.
func (r *Router) RegisterHandler(msgType string, handler Handler) {
	r.handlers[msgType] = handler
}

// GetHandler returns the handler for a specific message type.
func (r *Router) GetHandler(msgType string) (Handler, bool) {
	handler, ok := r.handlers[msgType]
	return handler, ok
}

// Handler interface for message handlers.
type Handler interface {
	Handle(senderAddr string, msg Message)
}
