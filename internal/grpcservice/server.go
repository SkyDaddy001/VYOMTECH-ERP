package grpcservice

import (
	"database/sql"
	"net"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server wraps the gRPC server with configuration
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	address    string
	logger     *logger.Logger
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(address string, rbacService *services.RBACService, db *sql.DB, logger *logger.Logger) (*Server, error) {
	// Create gRPC server with options
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(50 * 1024 * 1024), // 50MB max receive message size
		grpc.MaxSendMsgSize(50 * 1024 * 1024), // 50MB max send message size
	}

	grpcServer := grpc.NewServer(opts...)

	// Register services (TODO: Uncomment after proto generation and RegisterGRPCServices implementation)
	// RegisterGRPCServices(grpcServer, rbacService, db, logger)

	// Enable reflection for debugging and testing
	reflection.Register(grpcServer)

	logger.Info("gRPC server created", "address", address)

	return &Server{
		grpcServer: grpcServer,
		address:    address,
		logger:     logger,
	}, nil
}

// Start starts the gRPC server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.logger.Error("Failed to listen on gRPC port", "error", err, "address", s.address)
		return err
	}

	s.listener = listener
	s.logger.Info("gRPC server starting", "address", s.address)

	// Start server in blocking mode
	return s.grpcServer.Serve(listener)
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server")
	s.grpcServer.GracefulStop()
	s.logger.Info("gRPC server stopped")
}

// GetGRPCServer returns the underlying grpc.Server
func (s *Server) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

// GetListener returns the network listener
func (s *Server) GetListener() net.Listener {
	return s.listener
}
