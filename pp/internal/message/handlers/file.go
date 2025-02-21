package handlers

import (
	"log"
	"pp/internal/filetransfer"
	"pp/internal/message"
	// "pp/internal/node" //不再需要node
)

// FileRequestHandler handles file request messages.
type FileRequestHandler struct {
	//node                *node.Node //不再需要node
	fileTransferManager *filetransfer.Manager
	serverAddr          string //需要serverAddr
	sendMessageFunc     func(addr string, msg message.Message) error
}

// NewFileRequestHandler creates a new FileRequestHandler instance.
func NewFileRequestHandler(fileTransferManager *filetransfer.Manager, serverAddr string, sendMessageFunc func(addr string, msg message.Message) error) *FileRequestHandler {
	return &FileRequestHandler{fileTransferManager: fileTransferManager, serverAddr: serverAddr, sendMessageFunc: sendMessageFunc}
}

// Handle processes a file request message.
func (h *FileRequestHandler) Handle(senderAddr string, msg message.Message) {
	fileID, ok := msg.Data.(string)
	if !ok {
		log.Printf("Invalid file request data from %s", senderAddr)
		return
	}

	// Check if we have the file
	metadata, err := h.fileTransferManager.GetMetadata(fileID)
	if err != nil {
		log.Printf("File not found: %s", fileID)
		return
	}

	// Send file metadata to the requester
	msg = message.Message{Type: "file_metadata", Data: metadata, Sender: h.serverAddr}
	if err := h.sendMessageFunc(senderAddr, msg); err != nil {
		log.Printf("Error sending file metadata to %s: %v", senderAddr, err)
	}

	// Start sending file chunks
	go h.sendFileChunks(senderAddr, fileID, metadata.ChunkSize)
}

// sendFileChunks sends file chunks to the requester.
func (h *FileRequestHandler) sendFileChunks(senderAddr string, fileID string, chunkSize int) {
	file, err := h.fileTransferManager.OpenFile(fileID)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)
	chunkIndex := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file
			}
			log.Printf("Error reading file: %v", err)
			return
		}

		chunk := buffer[:bytesRead]
		data := map[string]interface{}{
			"file_id":     fileID,
			"chunk_index": chunkIndex,
			"chunk_data":  chunk,
		}
		msg := message.Message{Type: "file_chunk", Data: data, Sender: h.serverAddr}
		if err := h.sendMessageFunc(senderAddr, msg); err != nil {
			log.Printf("Error sending file chunk to %s: %v", senderAddr, err)
			return
		}

		chunkIndex++
	}

	log.Printf("File %s sent to %s", fileID, senderAddr)
}

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
