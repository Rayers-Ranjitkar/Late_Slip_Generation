package events

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

var clientManager = NewClientManager()

// WebSocketHandler handles WebSocket connections for both admin and student clients
func WebSocketHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	isAdmin := role == "admin"

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// Create new client
	client := NewClient(conn, userID, isAdmin, clientManager)

	// Register client with manager
	clientManager.Register(client)

	// Send initial connection message
	initialMessage := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: "CONNECTED",
		Data: gin.H{
			"message": "WebSocket connection established",
			"time":    time.Now(),
		},
	}

	client.SendJSON(initialMessage)
}