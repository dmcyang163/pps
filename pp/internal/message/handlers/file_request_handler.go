package handlers

import (
	"log"
	"pp/internal/events"
	"pp/internal/filetransfer"
	"pp/internal/message"
)

// FileRequestHandler handles file request messages.
type FileRequestHandler struct {
	fileTransferManager *filetransfer.Manager
	serverAddr          string
	eventManager        *events.EventManager
}

// NewFileRequestHandler creates a new FileRequestHandler instance.
func NewFileRequestHandler(fileTransferManager *filetransfer.Manager, serverAddr string, eventManager *events.EventManager) *FileRequestHandler {
	return &FileRequestHandler{
		fileTransferManager: fileTransferManager,
		serverAddr:          serverAddr,
		eventManager:        eventManager,
	}
}

// Handle processes a file request message.
func (h *FileRequestHandler) Handle(senderAddr string, msg message.Message) {
	filename, ok := msg.Data.(string)
	if !ok {
		log.Printf("Invalid file request data: %T", msg.Data)
		return
	}

	log.Printf("Received file request for %s from %s", filename, senderAddr)

	// 触发一个事件，通知 Node 发送文件
	eventData := events.FileRequestEventData{
		Filename:        filename,
		DestinationAddr: senderAddr,
	}
	h.eventManager.Publish(events.FileRequestEvent{EventData: eventData})
}
