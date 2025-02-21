package events

type EventType string

// Event interface
type Event interface {
	Type() EventType
	Data() interface{}
}

// SendMessageEventData is the data for SendMessageEvent.
type SendMessageEventData struct {
	DestinationAddr string
	Message         interface{} // 修改为 interface{}
}

// SendMessageEvent is an event that is triggered when a message needs to be sent.
type SendMessageEvent struct {
	EventData SendMessageEventData
}

func (e SendMessageEvent) Type() EventType {
	return "send_message" // 使用字符串常量
}

func (e SendMessageEvent) Data() interface{} {
	return e.EventData
}

// FileRequestEventData is the data for FileRequestEvent.
type FileRequestEventData struct {
	Filename        string
	DestinationAddr string
}

// FileRequestEvent is an event that is triggered when a file is requested.
type FileRequestEvent struct {
	EventData FileRequestEventData
}

func (e FileRequestEvent) Type() EventType {
	return "file_request"
}

func (e FileRequestEvent) Data() interface{} {
	return e.EventData
}
