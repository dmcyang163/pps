package events

import (
	"github.com/asaskevich/EventBus"
)

// EventManager manages events.
type EventManager struct {
	bus EventBus.Bus
}

// NewEventManager creates a new EventManager instance.
func NewEventManager() *EventManager {
	return &EventManager{
		bus: EventBus.New(),
	}
}

// Subscribe subscribes a handler to an event type.
func (em *EventManager) Subscribe(eventType EventType, handler func(Event)) {
	em.bus.Subscribe(string(eventType), handler)
}

// Unsubscribe unsubscribes a handler from an event type.
func (em *EventManager) Unsubscribe(eventType EventType, handler func(Event)) {
	em.bus.Unsubscribe(string(eventType), handler)
}

// Publish publishes an event.
func (em *EventManager) Publish(event Event) {
	em.bus.Publish(string(event.Type()), event)
}
