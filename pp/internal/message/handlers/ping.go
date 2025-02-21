package handlers

import (
	"log"
	"pp/internal/events"
	"pp/internal/message"
)

// PingHandler handles ping messages.
type PingHandler struct {
	eventManager *events.EventManager
	serverAddr   string // 需要 serverAddr
}

// NewPingHandler creates a new PingHandler instance.
func NewPingHandler(eventManager *events.EventManager, serverAddr string) *PingHandler {
	return &PingHandler{
		eventManager: eventManager,
		serverAddr:   serverAddr,
	}
}

// Handle processes a ping message.
func (h *PingHandler) Handle(senderAddr string, msg message.Message) {
	log.Printf("Received ping from %s", senderAddr)

	// 触发一个事件，通知 Node 发送 pong 消息
	eventData := events.SendMessageEventData{
		DestinationAddr: senderAddr,
		Message: message.Message{
			Type:   "pong",
			Data:   "pong",
			Sender: h.serverAddr,
		},
	}
	h.eventManager.Publish(events.SendMessageEvent{EventData: eventData})
}
