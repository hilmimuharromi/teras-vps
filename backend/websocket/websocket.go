package websocket

import (
	"fmt"
	"log"
	"sync"
	"time"

	"teras-vps/backend/proxmox"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

// Hub maintains active WebSocket connections
type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

// Register adds a new client
func (h *Hub) Register(client *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
	log.Printf("Client connected. Total clients: %d", len(h.clients))
}

// Unregister removes a client
func (h *Hub) Unregister(client *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		client.Close()
		log.Printf("Client disconnected. Total clients: %d", len(h.clients))
	}
}

// Broadcast sends a message to all clients
func (h *Hub) Broadcast(message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if err := client.WriteJSON(message); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			h.Unregister(client)
		}
	}
}

// VMStats represents VM statistics
type VMStats struct {
	VMID      int                    `json:"vmid"`
	Username  string                 `json:"username"`
	Status    string                 `json:"status"`
	CPU       float64                `json:"cpu"`
	Memory    int64                  `json:"memory"`
	MemoryMax int64                  `json:"memory_max"`
	Disk      int64                  `json:"disk"`
	DiskMax   int64                  `json:"disk_max"`
	Network   map[string]interface{} `json:"network"`
	Timestamp int64                  `json:"timestamp"`
}

// StatsService handles VM statistics polling
type StatsService struct {
	hub      *Hub
	proxmox  *proxmox.Client
	stopChan  chan struct{}
}

// NewStatsService creates a new stats service
func NewStatsService(hub *Hub, proxmoxClient *proxmox.Client) *StatsService {
	return &StatsService{
		hub:     hub,
		proxmox: proxmoxClient,
		stopChan: make(chan struct{}),
	}
}

// Start begins polling VM stats
func (s *StatsService) Start() {
	ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds

	go func() {
		for {
			select {
			case <-s.stopChan:
				ticker.Stop()
				return
			case <-ticker.C:
				s.pollVMStats()
			}
		}
	}()

	log.Println("Stats service started")
}

// Stop stops the stats service
func (s *StatsService) Stop() {
	close(s.stopChan)
	log.Println("Stats service stopped")
}

// pollVMStats polls all VMs for statistics
func (s *StatsService) pollVMStats() {
	// TODO: Get all VMs from database and poll their stats
	// For now, send placeholder data

	stats := VMStats{
		VMID:      101,
		Username:  "demo",
		Status:    "running",
		CPU:       12.5,
		Memory:    512 * 1024 * 1024,
		MemoryMax: 1024 * 1024 * 1024,
		Disk:      8 * 1024 * 1024 * 1024,
		DiskMax:   20 * 1024 * 1024 * 1024,
		Network: map[string]interface{}{
			"in":  1250 * 1024 * 1024,
			"out": 3200 * 1024 * 1024,
		},
		Timestamp: time.Now().Unix(),
	}

	s.hub.Broadcast(stats)
}

// SetupWebSocket sets up WebSocket routes
func SetupWebSocket(app *fiber.App, hub *Hub, proxmoxClient *proxmox.Client) {
	// WebSocket route for VM stats
	app.Get("/ws/vm-stats", websocket.New(func(c *websocket.Conn) {
		// Register client
		hub.Register(c)
		defer hub.Unregister(c)

		// Send initial stats
		stats := VMStats{
			VMID:      101,
			Username:  "demo",
			Status:    "running",
			CPU:       12.5,
			Memory:    512 * 1024 * 1024,
			MemoryMax: 1024 * 1024 * 1024,
			Disk:      8 * 1024 * 1024 * 1024,
			DiskMax:   20 * 1024 * 1024 * 1024,
			Network: map[string]interface{}{
				"in":  1250 * 1024 * 1024,
				"out": 3200 * 1024 * 1024,
			},
			Timestamp: time.Now().Unix(),
		}

		c.WriteJSON(stats)

		// Wait for disconnect
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				break
			}

			log.Printf("Received: type=%d, message=%s", messageType, message)
		}
	}))

	// Start stats service
	statsService := NewStatsService(hub, proxmoxClient)
	statsService.Start()

	log.Println("WebSocket server started on /ws/vm-stats")
}
