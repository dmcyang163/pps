package handlers

import (
	"fmt"
	"log"
	"pp/internal/message"
)

// ChatHandler handles chat messages.
type ChatHandler struct {
	// node *node.Node  //不再需要node
}

// NewChatHandler creates a new ChatHandler instance.
func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

// Handle processes a chat message.
func (h *ChatHandler) Handle(senderAddr string, msg message.Message) {
	chatText, ok := msg.Data.(string)
	if !ok {
		log.Printf("Invalid chat message data from %s", senderAddr)
		return
	}

	fmt.Printf("[%s]: %s\n", senderAddr, chatText)
}
