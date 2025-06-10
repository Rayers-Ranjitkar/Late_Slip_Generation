package events

import (
	"sync"
)

type ClientManager struct {
	adminClients   map[*Client]bool
	studentClients map[string]*Client
	mu             sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		adminClients:   make(map[*Client]bool),
		studentClients: make(map[string]*Client),
	}
}

// Register adds a new client to the manager
func (cm *ClientManager) Register(client *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client.IsAdmin {
		cm.adminClients[client] = true
	} else {
		cm.studentClients[client.UserID] = client
	}

	// Start client goroutines
	go client.ReadPump()
	go client.WritePump()
}

// Unregister removes a client from the manager
func (cm *ClientManager) Unregister(client *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client.IsAdmin {
		if _, ok := cm.adminClients[client]; ok {
			delete(cm.adminClients, client)
			close(client.Send)
		}
	} else {
		if existingClient, ok := cm.studentClients[client.UserID]; ok && existingClient == client {
			delete(cm.studentClients, client.UserID)
			close(client.Send)
		}
	}
}