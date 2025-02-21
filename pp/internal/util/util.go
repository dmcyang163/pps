package util

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a unique UUID.
func GenerateUUID() string {
	return uuid.New().String()
}
