package node

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"pp/internal/config"
	"pp/internal/filetransfer"
	"pp/internal/message"
	"pp/internal/network"
	"pp/internal/peer"
)

// Node represents a peer in the P2P network.
type Node struct {
	config              *config.Config
	serverAddr          string
	peerManager         *peer.Manager
	MessageRouter       *message.Router
	networkServer       network.NetworkServer // 使用接口
	fileTransferManager *filetransfer.Manager
	shutdownCh          chan struct{}
	wg                  sync.WaitGroup
}

// NewNode creates a new Node instance.
func NewNode(cfg *config.Config, networkServer network.NetworkServer) (*Node, error) { // 传入接口
	node := &Node{
		config:              cfg,
		serverAddr:          ":" + strconv.Itoa(cfg.Port),
		peerManager:         peer.NewManager(cfg.MaxPeers),
		MessageRouter:       message.NewRouter(),
		networkServer:       networkServer, // 使用传入的接口
		shutdownCh:          make(chan struct{}),
		fileTransferManager: filetransfer.NewManager(cfg.DataDir),
	}

	//node.networkServer = network.NewServer(node.serverAddr, node.handleIncomingMessage, node.peerConnected, node.peerDisconnected) // 移除
	networkServer.SetMessageHandler(node.handleIncomingMessage) // 设置消息处理函数
	networkServer.SetConnectHandler(node.peerConnected)         // 设置连接处理函数
	networkServer.SetDisconnectHandler(node.peerDisconnected)   // 设置断开连接处理函数

	return node, nil
}

// Start starts the node's network server and background tasks.
func (n *Node) Start() error {
	log.Printf("Starting node on %s", n.serverAddr)

	// Start the network server
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		if err := n.networkServer.Start(); err != nil { // 移除 gnet.ErrServerClosed 检查
			log.Printf("Error starting network server: %v", err)
		}
	}()

	// Start peer discovery
	n.wg.Add(1)
	go n.discoverPeers()

	// Start pinging peers
	n.wg.Add(1)
	go n.pingPeers()

	return nil
}

// Shutdown gracefully shuts down the node.
func (n *Node) Shutdown() {
	log.Println("Shutting down node...")
	close(n.shutdownCh) // Signal goroutines to stop

	if err := n.networkServer.Stop(); err != nil {
		log.Printf("Error stopping network server: %v", err)
	}

	n.wg.Wait() // Wait for all goroutines to finish
	log.Println("Node shutdown complete.")
}

// handleIncomingMessage processes incoming messages.
func (n *Node) handleIncomingMessage(addr string, data []byte) { // 修改参数类型
	msg, err := message.Deserialize(data)
	if err != nil {
		log.Printf("Error deserializing message: %v", err)
		return
	}

	log.Printf("Received message from %s: Type=%s, Data=%v", msg.Sender, msg.Type, msg.Data)
	if handler, ok := n.MessageRouter.GetHandler(msg.Type); ok {
		handler.Handle(addr, msg)
	} else {
		log.Printf("No handler found for message type: %s", msg.Type)
	}
}

// peerConnected is called when a new peer connects.
func (n *Node) peerConnected(addr string) {
	n.peerManager.AddPeer(addr)
	log.Printf("Peer connected: %s", addr)
}

// peerDisconnected is called when a peer disconnects.
func (n *Node) peerDisconnected(addr string) {
	n.peerManager.RemovePeer(addr)
	log.Printf("Peer disconnected: %s", addr)
}

// discoverPeers attempts to connect to seed nodes and discover new peers.
func (n *Node) discoverPeers() {
	defer n.wg.Done()

	// Connect to seed nodes
	for _, seedNode := range n.config.SeedNodes {
		if err := n.connectToPeer(seedNode); err != nil {
			log.Printf("Error connecting to seed node %s: %v", seedNode, err)
		}
	}

	// Periodically try to discover more peers (optional)
	ticker := time.NewTicker(30 * time.Second) // Adjust interval as needed
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Implement logic to request peer lists from known peers
			// and connect to new peers.  This is a more advanced topic
			// and depends on your desired network topology.
		case <-n.shutdownCh:
			return
		}
	}
}

// pingPeers periodically pings all known peers.
func (n *Node) pingPeers() {
	defer n.wg.Done()

	ticker := time.NewTicker(time.Duration(n.config.PingInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			peers := n.peerManager.GetPeers()
			for _, peerAddr := range peers {
				go func(addr string) {
					pingMsg := message.Message{Type: "ping", Data: "ping", Sender: n.serverAddr}
					if err := n.SendMessage(addr, pingMsg); err != nil {
						log.Printf("Error sending ping to %s: %v", addr, err)
						n.peerManager.RemovePeer(addr)
					}
				}(peerAddr)
			}
		case <-n.shutdownCh:
			return
		}
	}
}

// connectToPeer attempts to establish a connection with a peer.
func (n *Node) connectToPeer(peerAddr string) error {
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send a "new_peer" message to introduce ourselves
	newPeerMsg := message.Message{Type: "new_peer", Data: n.serverAddr, Sender: n.serverAddr}
	if err := n.SendMessage(peerAddr, newPeerMsg); err != nil {
		return err
	}

	n.peerManager.AddPeer(peerAddr) // Add to peer list immediately
	return nil
}

// sendMessage sends a message to a specific peer.
func (n *Node) SendMessage(addr string, msg message.Message) error {
	msgBytes, err := message.Serialize(msg)
	if err != nil {
		return err
	}

	return n.networkServer.SendMessage(addr, msgBytes)
}

// broadcastMessage sends a message to all known peers.
func (n *Node) broadcastMessage(msg message.Message) {
	peers := n.peerManager.GetPeers()
	for _, peerAddr := range peers {
		go func(addr string) {
			if err := n.SendMessage(addr, msg); err != nil {
				log.Printf("Error sending message to %s: %v", addr, err)
			}
		}(peerAddr)
	}
}

// SendChatMessage sends a chat message to all peers.
func (n *Node) SendChatMessage(text string) error {
	msg := message.Message{Type: "chat", Data: text, Sender: n.serverAddr}
	n.broadcastMessage(msg)
	return nil
}

// SendFileRequest sends a file request message to a specific peer.
func (n *Node) SendFileRequest(peerAddr string, fileID string) error {
	msg := message.Message{Type: "file_request", Data: fileID, Sender: n.serverAddr}
	return n.SendMessage(peerAddr, msg)
}

// SendFileChunk sends a file chunk message to a specific peer.
func (n *Node) SendFileChunk(peerAddr string, fileID string, chunk []byte, chunkIndex int) error {
	data := map[string]interface{}{
		"file_id":     fileID,
		"chunk_index": chunkIndex,
		"chunk_data":  chunk,
	}
	msg := message.Message{Type: "file_chunk", Data: data, Sender: n.serverAddr}
	return n.SendMessage(peerAddr, msg)
}

// SendFileMetadata sends file metadata to a specific peer.
func (n *Node) SendFileMetadata(peerAddr string, metadata filetransfer.Metadata) error {
	msg := message.Message{Type: "file_metadata", Data: metadata, Sender: n.serverAddr}
	return n.SendMessage(peerAddr, msg)
}
