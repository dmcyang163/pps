package message

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
