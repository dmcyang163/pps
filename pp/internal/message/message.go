package message

import (
	"encoding/json"
)

// Message represents a P2P message.
type Message struct {
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
	Sender string      `json:"sender"`
}

// Serialize serializes a Message to JSON.
func Serialize(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

// Deserialize deserializes a Message from JSON.
func Deserialize(data []byte) (Message, error) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	return msg, err
}
