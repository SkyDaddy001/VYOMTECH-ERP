package grpcservice

import (
	"database/sql"

	// "vyomtech-backend/api/pb/rbac" // TODO: Uncomment after running generate-proto.sh
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// RBACGRPCServer implements the RBAC gRPC service
// Note: Uncomment the pb/rbac import and implement methods after proto generation
type RBACGRPCServer struct {
	// rbac.UnimplementedRBACServiceServer // TODO: Uncomment after proto generation
	rbacService *services.RBACService
	db          *sql.DB
	logger      *logger.Logger
}

// NewRBACGRPCServer creates a new RBAC gRPC server
func NewRBACGRPCServer(rbacService *services.RBACService, db *sql.DB, logger *logger.Logger) *RBACGRPCServer {
	return &RBACGRPCServer{
		rbacService: rbacService,
		db:          db,
		logger:      logger,
	}
}

// CreateRole implements the CreateRole gRPC method
// TODO: Uncomment after proto generation
/*
func (s *RBACGRPCServer) CreateRole(ctx context.Context, req *rbac.CreateRoleRequest) (*rbac.CreateRoleResponse, error) {
	s.logger.Info("gRPC CreateRole called", "tenant_id", req.TenantId, "role_name", req.Name)

	// TODO: Implement role creation logic
	// This will call the existing RBACService.CreateRole method
	// and convert the response to protobuf format

	return &rbac.CreateRoleResponse{
		RoleId:  "role-123", // Generated ID
		Success: true,
		Message: "Role created successfully",
	}, nil
}

// GetRole implements the GetRole gRPC method
func (s *RBACGRPCServer) GetRole(ctx context.Context, req *rbac.GetRoleRequest) (*rbac.GetRoleResponse, error) {
	s.logger.Info("gRPC GetRole called", "tenant_id", req.TenantId, "role_id", req.RoleId)

	// TODO: Implement role retrieval logic
	// This will call the existing RBACService.GetRole method

	return &rbac.GetRoleResponse{
		Role: &rbac.Role{
			Id:        req.RoleId,
			TenantId:  req.TenantId,
			Name:      "Manager",
			IsActive:  true,
			CreatedAt: timestamppb.Now(),
		},
	}, nil
}

// ListRoles implements the ListRoles gRPC method
func (s *RBACGRPCServer) ListRoles(ctx context.Context, req *rbac.ListRolesRequest) (*rbac.ListRolesResponse, error) {
	s.logger.Info("gRPC ListRoles called", "tenant_id", req.TenantId)

	// TODO: Implement role listing logic
	// This will call the existing RBACService.ListRoles method

	return &rbac.ListRolesResponse{
		Roles: []*rbac.Role{},
		Total: 0,
	}, nil
}

// AssignPermissions implements the AssignPermissions gRPC method
func (s *RBACGRPCServer) AssignPermissions(ctx context.Context, req *rbac.AssignPermissionsRequest) (*rbac.AssignPermissionsResponse, error) {
	s.logger.Info("gRPC AssignPermissions called", "tenant_id", req.TenantId, "role_id", req.RoleId)

	// TODO: Implement permission assignment logic
	// This will call the existing RBACService methods

	return &rbac.AssignPermissionsResponse{
		Success:       true,
		Message:       "Permissions assigned successfully",
		AssignedCount: int32(len(req.PermissionIds)),
	}, nil
}

// VerifyPermission implements the VerifyPermission gRPC method
func (s *RBACGRPCServer) VerifyPermission(ctx context.Context, req *rbac.VerifyPermissionRequest) (*rbac.VerifyPermissionResponse, error) {
	s.logger.Info("gRPC VerifyPermission called", "tenant_id", req.TenantId, "user_id", req.UserId, "permission", req.PermissionCode)

	// Call the existing RBACService.VerifyPermission method
	err := s.rbacService.VerifyPermission(ctx, req.TenantId, req.UserId, req.PermissionCode)

	hasPermission := err == nil
	message := "Permission granted"
	if !hasPermission {
		message = "Permission denied"
	}

	return &rbac.VerifyPermissionResponse{
		HasPermission: hasPermission,
		Message:       message,
	}, nil
}

// GetUserPermissions implements the GetUserPermissions gRPC method
func (s *RBACGRPCServer) GetUserPermissions(ctx context.Context, req *rbac.VerifyPermissionRequest) (*rbac.ListPermissionsResponse, error) {
	s.logger.Info("gRPC GetUserPermissions called", "tenant_id", req.TenantId, "user_id", req.UserId)

	// TODO: Implement permission retrieval logic
	// This will call the existing RBACService.GetUserPermissions method

	return &rbac.ListPermissionsResponse{
		Permissions: []*rbac.Permission{},
		Total:       0,
	}, nil
}

// RegisterGRPCServices registers all gRPC services with the gRPC server
func RegisterGRPCServices(grpcServer *grpc.Server, rbacService *services.RBACService, db *sql.DB, logger *logger.Logger) {
	rbacGRPCServer := NewRBACGRPCServer(rbacService, db, logger)
	// rbac.RegisterRBACServiceServer(grpcServer, rbacGRPCServer) // TODO: Uncomment after proto generation

	logger.Info("gRPC services registered successfully - waiting for proto generation")
}
*/

// Placeholder function to prevent unused import warning
func init() {
	_ = timestamppb.Now
}
