package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"multi-tenant-ai-callcenter/pkg/logger"
)

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type      string                 `json:"type"`     // notification_type
	EventID   string                 `json:"event_id"` // unique event identifier
	Timestamp time.Time              `json:"timestamp"`
	TenantID  string                 `json:"tenant_id"`
	UserID    int64                  `json:"user_id"`
	AgentID   int64                  `json:"agent_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
}

// WebSocketClient represents a connected WebSocket client
type WebSocketClient struct {
	ID       string
	TenantID string
	UserID   int64
	Conn     *websocket.Conn
	Send     chan *WebSocketMessage
}

// WebSocketHub manages WebSocket connections and broadcasts messages
type WebSocketHub struct {
	clients       map[*WebSocketClient]bool
	broadcast     chan *WebSocketMessage
	register      chan *WebSocketClient
	unregister    chan *WebSocketClient
	mu            sync.RWMutex
	logger        *logger.Logger
	tenantClients map[string]map[*WebSocketClient]bool // tenant-specific clients
}

// NewWebSocketHub creates a new WebSocketHub
func NewWebSocketHub(log *logger.Logger) *WebSocketHub {
	hub := &WebSocketHub{
		clients:       make(map[*WebSocketClient]bool),
		broadcast:     make(chan *WebSocketMessage, 256),
		register:      make(chan *WebSocketClient),
		unregister:    make(chan *WebSocketClient),
		logger:        log,
		tenantClients: make(map[string]map[*WebSocketClient]bool),
	}

	go hub.run()
	return hub
}

// run manages the hub's event loop
func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if h.tenantClients[client.TenantID] == nil {
				h.tenantClients[client.TenantID] = make(map[*WebSocketClient]bool)
			}
			h.tenantClients[client.TenantID][client] = true
			h.mu.Unlock()
			h.logger.Info("WebSocket client registered", "client_id", client.ID, "tenant_id", client.TenantID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.tenantClients[client.TenantID], client)
				close(client.Send)
			}
			h.mu.Unlock()
			h.logger.Info("WebSocket client unregistered", "client_id", client.ID, "tenant_id", client.TenantID)

		case message := <-h.broadcast:
			h.mu.RLock()
			// Broadcast to all clients in the tenant
			if tenantClients, ok := h.tenantClients[message.TenantID]; ok {
				for client := range tenantClients {
					select {
					case client.Send <- message:
					default:
						// Client's send channel is full, close it
						h.unregister <- client
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// RegisterClient registers a new WebSocket client
func (h *WebSocketHub) RegisterClient(client *WebSocketClient) {
	h.register <- client
}

// UnregisterClient unregisters a WebSocket client
func (h *WebSocketHub) UnregisterClient(client *WebSocketClient) {
	h.unregister <- client
}

// BroadcastAgentStatus broadcasts agent availability status updates
func (h *WebSocketHub) BroadcastAgentStatus(tenantID string, agentID int64, status string) {
	message := &WebSocketMessage{
		Type:      "agent_status_updated",
		EventID:   fmt.Sprintf("agent-status-%d-%d", agentID, time.Now().Unix()),
		Timestamp: time.Now(),
		TenantID:  tenantID,
		AgentID:   agentID,
		Data: map[string]interface{}{
			"agent_id": agentID,
			"status":   status,
		},
	}
	h.broadcast <- message
}

// NotifyAgentOfCall notifies an agent of an incoming call
func (h *WebSocketHub) NotifyAgentOfCall(tenantID string, agentID int64, callID int64, leadName string) {
	message := &WebSocketMessage{
		Type:      "incoming_call",
		EventID:   fmt.Sprintf("call-%d-%d", callID, time.Now().Unix()),
		Timestamp: time.Now(),
		TenantID:  tenantID,
		AgentID:   agentID,
		Data: map[string]interface{}{
			"call_id":   callID,
			"agent_id":  agentID,
			"lead_name": leadName,
		},
	}
	h.broadcast <- message
}

// BroadcastCampaignUpdate broadcasts campaign performance updates
func (h *WebSocketHub) BroadcastCampaignUpdate(tenantID string, campaignID int64, update map[string]interface{}) {
	update["campaign_id"] = campaignID
	message := &WebSocketMessage{
		Type:      "campaign_update",
		EventID:   fmt.Sprintf("campaign-%d-%d", campaignID, time.Now().Unix()),
		Timestamp: time.Now(),
		TenantID:  tenantID,
		Data:      update,
	}
	h.broadcast <- message
}

// BroadcastCallUpdate broadcasts call status updates
func (h *WebSocketHub) BroadcastCallUpdate(tenantID string, callID int64, status string, data map[string]interface{}) {
	data["call_id"] = callID
	data["status"] = status
	message := &WebSocketMessage{
		Type:      "call_status_updated",
		EventID:   fmt.Sprintf("call-update-%d-%d", callID, time.Now().Unix()),
		Timestamp: time.Now(),
		TenantID:  tenantID,
		Data:      data,
	}
	h.broadcast <- message
}

// BroadcastLeadUpdate broadcasts lead status updates
func (h *WebSocketHub) BroadcastLeadUpdate(tenantID string, leadID int64, status string, data map[string]interface{}) {
	data["lead_id"] = leadID
	data["status"] = status
	message := &WebSocketMessage{
		Type:      "lead_status_updated",
		EventID:   fmt.Sprintf("lead-update-%d-%d", leadID, time.Now().Unix()),
		Timestamp: time.Now(),
		TenantID:  tenantID,
		Data:      data,
	}
	h.broadcast <- message
}

// BroadcastGamificationEvent broadcasts gamification events (points awarded, badges earned)
func (h *WebSocketHub) BroadcastGamificationEvent(tenantID string, userID int64, eventType string, data map[string]interface{}) {
	data["event_type"] = eventType
	message := &WebSocketMessage{
		Type:      "gamification_event",
		EventID:   fmt.Sprintf("gamif-%d-%d", userID, time.Now().Unix()),
		Timestamp: time.Now(),
		TenantID:  tenantID,
		UserID:    userID,
		Data:      data,
	}
	h.broadcast <- message
}

// HandleClientConnection handles a new WebSocket client connection
func (h *WebSocketHub) HandleClientConnection(conn *websocket.Conn, tenantID string, userID int64) {
	client := &WebSocketClient{
		ID:       fmt.Sprintf("%s-%d-%d", tenantID, userID, time.Now().Unix()),
		TenantID: tenantID,
		UserID:   userID,
		Conn:     conn,
		Send:     make(chan *WebSocketMessage, 256),
	}

	h.RegisterClient(client)

	// Start reading messages
	go h.readMessages(client)
	// Start writing messages
	go h.writeMessages(client)
}

// readMessages reads incoming messages from a client
func (h *WebSocketHub) readMessages(client *WebSocketClient) {
	defer func() {
		h.UnregisterClient(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg WebSocketMessage
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("WebSocket error", "error", err)
			}
			break
		}
		// Process incoming message if needed
		h.logger.Debug("WebSocket message received", "client_id", client.ID, "type", msg.Type)
	}
}

// writeMessages writes outgoing messages to a client
func (h *WebSocketHub) writeMessages(client *WebSocketClient) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			json.NewEncoder(w).Encode(message)
			w.Close()

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// GetConnectedClientsCount returns the total number of connected clients
func (h *WebSocketHub) GetConnectedClientsCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetTenantClientsCount returns the number of connected clients for a tenant
func (h *WebSocketHub) GetTenantClientsCount(tenantID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.tenantClients[tenantID])
}
