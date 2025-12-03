package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub    *services.WebSocketHub
	logger *logger.Logger
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler(hub *services.WebSocketHub, logger *logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		hub:    hub,
		logger: logger,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin validation
		return true
	},
}

// HandleConnection handles WebSocket connection upgrade and management
// GET /api/v1/ws
func (wh *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		wh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		wh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		wh.logger.Error("Failed to upgrade WebSocket connection", "error", err)
		http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	wh.logger.Info("WebSocket connection established", "tenant_id", tenantID, "user_id", userID)

	// Handle the client connection
	wh.hub.HandleClientConnection(conn, tenantID, userID)
}

// GetConnectionStats returns WebSocket connection statistics
// GET /api/v1/ws/stats
func (wh *WebSocketHandler) GetConnectionStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		wh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	totalClients := wh.hub.GetConnectedClientsCount()
	tenantClients := wh.hub.GetTenantClientsCount(tenantID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"total_connections": ` + string(rune(totalClients)) + `,
		"tenant_connections": ` + string(rune(tenantClients)) + `
	}`))
}
