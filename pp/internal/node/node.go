package node

import (
	"log"
	"strconv"
	"sync"

	"pp/internal/config"
	"pp/internal/events"
	"pp/internal/filetransfer"
	"pp/internal/message"
	"pp/internal/network"
	"pp/internal/peer"
)

// Node represents a peer in the P2P network.
type Node struct {
	config              *config.Config
	ServerAddr          string
	PeerManager         *peer.Manager
	MessageRouter       *message.Router
	networkServer       network.NetworkServer // 使用接口
	FileTransferManager *filetransfer.Manager
	shutdownCh          chan struct{}
	wg                  sync.WaitGroup
	EventManager        *events.EventManager // 添加事件管理器
}

// NewNode creates a new Node instance.
func NewNode(cfg *config.Config, networkServer network.NetworkServer) (*Node, error) { // 传入接口
	node := &Node{
		config:              cfg,
		ServerAddr:          ":" + strconv.Itoa(cfg.Port),
		PeerManager:         peer.NewManager(cfg.MaxPeers),
		networkServer:       networkServer, // 使用传入的接口
		shutdownCh:          make(chan struct{}),
		FileTransferManager: filetransfer.NewManager(cfg.DataDir),
		EventManager:        events.NewEventManager(), // 初始化事件管理器
	}

	node.MessageRouter = message.NewRouter(node.EventManager)

	networkServer.SetMessageHandler(node.handleIncomingMessage) // 设置消息处理函数
	networkServer.SetConnectHandler(node.peerConnected)         // 设置连接处理函数
	networkServer.SetDisconnectHandler(node.peerDisconnected)   // 设置断开连接处理函数

	//  不再在这里注册 handlers，而是在 main.go 中注册

	// 订阅 SendMessageEvent
	node.EventManager.Subscribe("send_message", node.handleSendMessageEvent)

	// 订阅 FileRequestEvent
	node.EventManager.Subscribe("file_request", node.handleFileRequestEvent)

	return node, nil
}

// handleSendMessageEvent handles SendMessageEvent.
func (n *Node) handleSendMessageEvent(event events.Event) { // 修改参数类型
	sendMessageEvent, ok := event.(events.SendMessageEvent)
	if !ok {
		log.Printf("Invalid event type: %T", event)
		return
	}

	data := sendMessageEvent.Data().(events.SendMessageEventData)
	msg, ok := data.Message.(message.Message) // 类型断言
	if !ok {
		log.Printf("Invalid message type: %T", data.Message)
		return
	}
	err := n.SendMessage(data.DestinationAddr, msg)
	if err != nil {
		log.Printf("Error sending message to %s: %v", data.DestinationAddr, err)
	}
}

// handleFileRequestEvent handles FileRequestEvent.
func (n *Node) handleFileRequestEvent(event events.Event) {
	fileRequestEvent, ok := event.(events.FileRequestEvent)
	if !ok {
		log.Printf("Invalid event type: %T", event)
		return
	}

	data := fileRequestEvent.Data().(events.FileRequestEventData)
	err := n.sendFile(data.DestinationAddr, data.Filename)
	if err != nil {
		log.Printf("Error sending file %s to %s: %v", data.Filename, data.DestinationAddr, err)
	}
}

// handleIncomingMessage handles incoming messages from the network.
func (n *Node) handleIncomingMessage(addr string, data []byte) {
	msg, err := message.Deserialize(data)
	if err != nil {
		log.Printf("Error deserializing message: %v", err)
		return
	}

	handler, ok := n.MessageRouter.GetHandler(msg.Type)
	if !ok {
		log.Printf("No handler found for message type: %s", msg.Type)
		return
	}

	handler.Handle(addr, msg)
}

// peerConnected is called when a new peer connects to the node.
func (n *Node) peerConnected(addr string) {
	log.Printf("New peer connected: %s", addr)
	n.PeerManager.AddPeer(addr)
}

// peerDisconnected is called when a peer disconnects from the node.
func (n *Node) peerDisconnected(addr string) {
	log.Printf("Peer disconnected: %s", addr)
	n.PeerManager.RemovePeer(addr)
}

// Start starts the node.
func (n *Node) Start() error {
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		if err := n.networkServer.Start(); err != nil {
			log.Printf("Network server failed: %v", err)
			n.Shutdown()
		}
	}()

	log.Printf("Node started on %s", n.ServerAddr)
	return nil
}

// Shutdown shuts down the node.
func (n *Node) Shutdown() {
	log.Println("Shutting down node...")
	close(n.shutdownCh)
	n.networkServer.Stop()
	n.wg.Wait()
	log.Println("Node shutdown complete.")
}

// sendMessage sends a message to a specific peer.
func (n *Node) SendMessage(addr string, msg message.Message) error {
	msgBytes, err := message.Serialize(msg)
	if err != nil {
		return err
	}

	return n.networkServer.SendMessage(addr, msgBytes)
}

// sendFile sends a file to a specific peer.
func (n *Node) sendFile(destinationAddr string, filename string) error {
	return n.FileTransferManager.SendFile(destinationAddr, filename)
}
