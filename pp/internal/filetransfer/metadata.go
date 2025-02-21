package filetransfer

// Metadata represents file metadata.
type Metadata struct {
	FileID    string `json:"file_id"`
	Filename  string `json:"filename"`
	FileSize  int64  `json:"file_size"`
	ChunkSize int    `json:"chunk_size"`
}
