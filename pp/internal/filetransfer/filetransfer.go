package filetransfer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	DefaultChunkSize = 1024 * 1024 // 1MB
)

// Manager manages file transfers.
type Manager struct {
	dataDir string
	mu      sync.Mutex
}

// NewManager creates a new Manager instance.
func NewManager(dataDir string) *Manager {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
	}
	return &Manager{dataDir: dataDir}
}

// GetMetadata retrieves file metadata.
func (m *Manager) GetMetadata(fileID string) (Metadata, error) {
	metadataPath := filepath.Join(m.dataDir, fileID+".metadata")
	file, err := os.Open(metadataPath)
	if err != nil {
		return Metadata{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var metadata Metadata
	err = decoder.Decode(&metadata)
	return metadata, err
}

// StoreMetadata stores file metadata.
func (m *Manager) StoreMetadata(metadata Metadata) error {
	metadataPath := filepath.Join(m.dataDir, metadata.FileID+".metadata")
	file, err := os.Create(metadataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(metadata)
	return err
}

// OpenFile opens a file for reading.
func (m *Manager) OpenFile(fileID string) (*os.File, error) {
	filePath := filepath.Join(m.dataDir, fileID)
	return os.Open(filePath)
}

// CreateFile creates a file for writing.
func (m *Manager) CreateFile(fileID string) (*os.File, error) {
	filePath := filepath.Join(m.dataDir, fileID)
	return os.Create(filePath)
}

// WriteChunk writes a chunk of data to a file.
func (m *Manager) WriteChunk(fileID string, chunk []byte, chunkIndex int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := filepath.Join(m.dataDir, fileID)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	offset := int64(chunkIndex) * int64(len(chunk))
	_, err = file.WriteAt(chunk, offset)
	return err
}

// SaveFile saves a file to the data directory.
func (m *Manager) SaveFile(filename string, reader io.Reader) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fileID := generateFileID(filename)
	filePath := filepath.Join(m.dataDir, fileID)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return "", err
	}

	// Create metadata
	metadata := Metadata{
		FileID:    fileID,
		Filename:  filename,
		FileSize:  getFileSize(filePath),
		ChunkSize: DefaultChunkSize,
	}
	if err := m.StoreMetadata(metadata); err != nil {
		os.Remove(filePath)
		return "", err
	}

	return fileID, nil
}

// generateFileID generates a unique file ID.
func generateFileID(filename string) string {
	// Implement a more robust ID generation strategy (e.g., UUID)
	return fmt.Sprintf("%s_%d", filename, os.Getpid())
}

// getFileSize returns the size of a file.
func getFileSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
