package network

import (
	"log"

	"github.com/panjf2000/gnet"
)

// NetworkServer 定义网络服务的接口
type NetworkServer interface {
	Start() error
	Stop() error
	SendMessage(addr string, message []byte) error
	SetMessageHandler(handler func(string, []byte))
	SetConnectHandler(handler func(string))
	SetDisconnectHandler(handler func(string))
}

// Server wraps gnet.EventServer to manage network connections.
type Server struct {
	addr         string
	eventHandler *eventHandler
	// gnetServer   gnet.Server
	*gnet.Server
}

// NewServer creates a new Server instance.
func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		eventHandler: &eventHandler{
			messageHandler:    func(string, []byte) {}, // 默认空函数
			connectHandler:    func(string) {},         // 默认空函数
			disconnectHandler: func(string) {},         // 默认空函数
		},
	}
}

// Start starts the gnet server.
func (s *Server) Start() error {
	log.Printf("Starting network server on %s", s.addr)
	return gnet.Serve(s.eventHandler, "tcp://"+"localhost"+s.addr, gnet.WithMulticore(true))
}

// Stop stops the gnet server.
func (s *Server) Stop() error {
	log.Println("Stopping network server...")
	// return s.gnetServer.Shutdown()
	return nil
}

// SendMessage sends a message to a specific address.
func (s *Server) SendMessage(addr string, message []byte) error {
	// Implementation using gnet to send the message
	// (You'll need to maintain a map of connections to addresses)
	return nil // Placeholder
}

// SetMessageHandler 设置消息处理函数
func (s *Server) SetMessageHandler(handler func(string, []byte)) {
	s.eventHandler.messageHandler = handler
}

// SetConnectHandler 设置连接处理函数
func (s *Server) SetConnectHandler(handler func(string)) {
	s.eventHandler.connectHandler = handler
}

// SetDisconnectHandler 设置断开连接处理函数
func (s *Server) SetDisconnectHandler(handler func(string)) {
	s.eventHandler.disconnectHandler = handler
}

// eventHandler implements gnet.EventHandler interface.
type eventHandler struct {
	gnet.EventServer
	messageHandler    func(string, []byte)
	connectHandler    func(string)
	disconnectHandler func(string)
}

// OnInitComplete is called when the server is ready.
func (eh *eventHandler) OnInitComplete(server gnet.Server) (action gnet.Action) {
	log.Printf("Network server started on %s", server.Addr.String())
	return
}

// OnOpened is called when a new connection is opened.
func (eh *eventHandler) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	addr := c.RemoteAddr().String()
	log.Printf("Connection opened: %s", addr)
	eh.connectHandler(addr)
	return
}

// OnClosed is called when a connection is closed.
func (eh *eventHandler) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	addr := c.RemoteAddr().String()
	log.Printf("Connection closed: %s, error: %v", addr, err)
	eh.disconnectHandler(addr)
	return
}

// React is called when data is received.
func (eh *eventHandler) React(inputFrame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	addr := c.RemoteAddr().String()
	eh.messageHandler(addr, inputFrame)
	return
}
