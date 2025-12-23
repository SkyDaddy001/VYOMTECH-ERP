package grpcservice

import (
	"context"
	"net/http"

	"vyomtech-backend/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// GRPCGatewayMux creates a new gRPC Gateway multiplexer
// This allows REST clients to communicate with gRPC services
func GRPCGatewayMux(ctx context.Context, grpcServerAddr string, logger *logger.Logger) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()

	// Create gRPC client connection
	// Note: conn is not directly needed here but kept for future handler registrations
	_, err := grpc.DialContext(
		ctx,
		grpcServerAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Error("Failed to dial gRPC server", "error", err, "address", grpcServerAddr)
		return nil, err
	}

	// Register RBAC service handlers
	// TODO: Add handler registrations as gRPC services are implemented
	// Example: rbac.RegisterRBACServiceHandler(ctx, mux, conn)

	logger.Info("gRPC Gateway mux created", "grpc_address", grpcServerAddr)

	return mux, nil
}

// GRPCGatewayHandler creates an HTTP handler that bridges REST to gRPC
func GRPCGatewayHandler(mux *runtime.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add common headers
		w.Header().Set("Content-Type", "application/json")

		// Route to gRPC Gateway
		mux.ServeHTTP(w, r)
	})
}

// ServiceRegistry holds references to all registered gRPC services
type ServiceRegistry struct {
	RBACService  interface{} // rbac.RBACServiceClient
	AuditService interface{} // audit.AuditServiceClient
	logger       *logger.Logger
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(logger *logger.Logger) *ServiceRegistry {
	return &ServiceRegistry{
		logger: logger,
	}
}

// RegisterService registers a service in the registry
func (sr *ServiceRegistry) RegisterService(serviceName string, service interface{}) {
	sr.logger.Info("Registering gRPC service", "service", serviceName)

	switch serviceName {
	case "rbac":
		sr.RBACService = service
	case "audit":
		sr.AuditService = service
	default:
		sr.logger.Warn("Unknown service registered", "service", serviceName)
	}
}

// GetService retrieves a registered service
func (sr *ServiceRegistry) GetService(serviceName string) interface{} {
	switch serviceName {
	case "rbac":
		return sr.RBACService
	case "audit":
		return sr.AuditService
	default:
		return nil
	}
}
