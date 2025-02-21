package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pp/internal/config"
	"pp/internal/message/handlers"
	"pp/internal/network"
	"pp/internal/node"
	"syscall"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create network server
	networkServer := network.NewServer(":" + fmt.Sprintf("%d", cfg.Port))

	// Create node
	node, err := node.NewNode(cfg, networkServer)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	//  创建 handlers，并将 eventManager 和其他依赖项传递给它们
	pingHandler := handlers.NewPingHandler(node.EventManager, node.ServerAddr)
	chatHandler := handlers.NewChatHandler()
	fileRequestHandler := handlers.NewFileRequestHandler(node.FileTransferManager, node.ServerAddr, node.EventManager) // 这里的 sendMessage 需要修改，通过事件触发
	fileChunkHandler := handlers.NewFileChunkHandler(node.FileTransferManager)
	fileMetadataHandler := handlers.NewFileMetadataHandler(node.FileTransferManager)

	//  注册 handlers 到 router
	node.MessageRouter.RegisterHandler("ping", pingHandler)
	node.MessageRouter.RegisterHandler("chat", chatHandler)
	node.MessageRouter.RegisterHandler("file_request", fileRequestHandler)
	node.MessageRouter.RegisterHandler("file_chunk", fileChunkHandler)
	node.MessageRouter.RegisterHandler("file_metadata", fileMetadataHandler)

	// Start the node
	if err := node.Start(); err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}

	// Handle shutdown signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh // Wait for shutdown signal
	node.Shutdown()
}
