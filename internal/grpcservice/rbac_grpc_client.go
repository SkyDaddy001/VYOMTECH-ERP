package grpcservice

import (

	// "vyomtech-backend/api/pb/rbac" // TODO: Uncomment after running generate-proto.sh
	"vyomtech-backend/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RBACGRPCClient wraps the gRPC client for RBAC service
type RBACGRPCClient struct {
	// client rbac.RBACServiceClient // TODO: Uncomment after proto generation
	conn   *grpc.ClientConn
	logger *logger.Logger
}

// NewRBACGRPCClient creates a new RBAC gRPC client
func NewRBACGRPCClient(address string, logger *logger.Logger) (*RBACGRPCClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to create gRPC connection", "error", err, "address", address)
		return nil, err
	}

	// client := rbac.NewRBACServiceClient(conn) // TODO: Uncomment after proto generation

	return &RBACGRPCClient{
		// client: client, // TODO: Uncomment after proto generation
		conn:   conn,
		logger: logger,
	}, nil
}

// CreateRole calls the gRPC CreateRole method
// TODO: Uncomment after proto generation
/*
func (c *RBACGRPCClient) CreateRole(ctx context.Context, req *rbac.CreateRoleRequest) (*rbac.CreateRoleResponse, error) {
	c.logger.Info("gRPC call: CreateRole", "tenant_id", req.TenantId, "role_name", req.Name)
	return c.client.CreateRole(ctx, req)
}

// GetRole calls the gRPC GetRole method
func (c *RBACGRPCClient) GetRole(ctx context.Context, req *rbac.GetRoleRequest) (*rbac.GetRoleResponse, error) {
	c.logger.Info("gRPC call: GetRole", "tenant_id", req.TenantId, "role_id", req.RoleId)
	return c.client.GetRole(ctx, req)
}

// ListRoles calls the gRPC ListRoles method
func (c *RBACGRPCClient) ListRoles(ctx context.Context, req *rbac.ListRolesRequest) (*rbac.ListRolesResponse, error) {
	c.logger.Info("gRPC call: ListRoles", "tenant_id", req.TenantId)
	return c.client.ListRoles(ctx, req)
}

// AssignPermissions calls the gRPC AssignPermissions method
func (c *RBACGRPCClient) AssignPermissions(ctx context.Context, req *rbac.AssignPermissionsRequest) (*rbac.AssignPermissionsResponse, error) {
	c.logger.Info("gRPC call: AssignPermissions", "tenant_id", req.TenantId, "role_id", req.RoleId)
	return c.client.AssignPermissions(ctx, req)
}

// VerifyPermission calls the gRPC VerifyPermission method
func (c *RBACGRPCClient) VerifyPermission(ctx context.Context, req *rbac.VerifyPermissionRequest) (*rbac.VerifyPermissionResponse, error) {
	c.logger.Info("gRPC call: VerifyPermission", "tenant_id", req.TenantId, "user_id", req.UserId)
	return c.client.VerifyPermission(ctx, req)
}

// GetUserPermissions calls the gRPC GetUserPermissions method
func (c *RBACGRPCClient) GetUserPermissions(ctx context.Context, req *rbac.VerifyPermissionRequest) (*rbac.ListPermissionsResponse, error) {
	c.logger.Info("gRPC call: GetUserPermissions", "tenant_id", req.TenantId, "user_id", req.UserId)
	return c.client.GetUserPermissions(ctx, req)
}
*/

// Close closes the gRPC connection
func (c *RBACGRPCClient) Close() error {
	c.logger.Info("Closing gRPC connection to RBAC service")
	return c.conn.Close()
}
