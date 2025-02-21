package peer

import (
	"sync"
)

// Manager manages the list of peers.
type Manager struct {
	maxPeers int
	peers    sync.Map // map[string]bool
	mu       sync.Mutex
}

// NewManager creates a new Manager instance.
func NewManager(maxPeers int) *Manager {
	return &Manager{
		maxPeers: maxPeers,
		peers:    sync.Map{},
	}
}

// AddPeer adds a peer to the list.
func (m *Manager) AddPeer(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.peers.Load(addr); ok {
		return
	}

	// Enforce max peer limit (optional)
	count := 0
	m.peers.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count >= m.maxPeers {
		// Implement peer eviction strategy (e.g., remove least recently seen)
		return
	}

	m.peers.Store(addr, true)
}

// RemovePeer removes a peer from the list.
func (m *Manager) RemovePeer(addr string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.peers.Delete(addr)
}

// GetPeers returns a list of all peers.
func (m *Manager) GetPeers() []string {
	peers := []string{}
	m.peers.Range(func(key, value interface{}) bool {
		peers = append(peers, key.(string))
		return true
	})
	return peers
}
