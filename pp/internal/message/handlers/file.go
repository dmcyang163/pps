package handlers

import (
	"log"
	"pp/internal/filetransfer"
	"pp/internal/message"
	// "pp/internal/node" //不再需要node
)

// FileChunkHandler handles file chunk messages.
type FileChunkHandler struct {
	fileTransferManager *filetransfer.Manager
}

// NewFileChunkHandler creates a new FileChunkHandler instance.
func NewFileChunkHandler(fileTransferManager *filetransfer.Manager) *FileChunkHandler {
	return &FileChunkHandler{fileTransferManager: fileTransferManager}
}

// Handle processes a file chunk message.
func (h *FileChunkHandler) Handle(senderAddr string, msg message.Message) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid file chunk data from %s", senderAddr)
		return
	}

	fileID, ok := data["file_id"].(string)
	if !ok {
		log.Printf("Invalid file_id in chunk from %s", senderAddr)
		return
	}

	chunkIndexFloat, ok := data["chunk_index"].(float64)
	if !ok {
		log.Printf("Invalid chunk_index in chunk from %s", senderAddr)
		return
	}
	chunkIndex := int(chunkIndexFloat)

	chunkData, ok := data["chunk_data"].([]byte)
	if !ok {
		log.Printf("Invalid chunk_data in chunk from %s", senderAddr)
		return
	}

	if err := h.fileTransferManager.WriteChunk(fileID, chunkData, chunkIndex); err != nil {
		log.Printf("Error writing chunk to file: %v", err)
	}
}

// FileMetadataHandler handles file metadata messages.
type FileMetadataHandler struct {
	fileTransferManager *filetransfer.Manager
}

// NewFileMetadataHandler creates a new FileMetadataHandler instance.
func NewFileMetadataHandler(fileTransferManager *filetransfer.Manager) *FileMetadataHandler {
	return &FileMetadataHandler{fileTransferManager: fileTransferManager}
}

// Handle processes a file metadata message.
func (h *FileMetadataHandler) Handle(senderAddr string, msg message.Message) {
	metadata, ok := msg.Data.(filetransfer.Metadata)
	if !ok {
		log.Printf("Invalid file metadata from %s", senderAddr)
		return
	}

	if err := h.fileTransferManager.StoreMetadata(metadata); err != nil {
		log.Printf("Error storing file metadata: %v", err)
	}
}
