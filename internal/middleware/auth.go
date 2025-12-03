package middleware

import (
	"context"
	"net/http"
	"strings"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

const (
	UserIDKey   = "user_id"
	TenantIDKey = "tenant_id"
	RoleKey     = "role"
)

// AuthMiddleware validates JWT token in Authorization header
func AuthMiddleware(authService *services.AuthService, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token
			user, err := authService.ValidateToken(token)
			if err != nil {
				log.Warn("Invalid token", "error", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user info to request context
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			ctx = context.WithValue(ctx, TenantIDKey, user.TenantID)
			ctx = context.WithValue(ctx, RoleKey, user.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TenantIsolationMiddleware ensures tenant-specific access control
func TenantIsolationMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID, ok := r.Context().Value(TenantIDKey).(string)
			if !ok || tenantID == "" {
				http.Error(w, "Missing tenant context", http.StatusForbidden)
				return
			}

			// Tenant isolation is enforced at the database query level
			// This middleware just ensures tenant context is present
			log.WithTenant(tenantID).Debug("Processing request")

			next.ServeHTTP(w, r)
		})
	}
}

// RoleBasedAccessMiddleware enforces role-based access control
func RoleBasedAccessMiddleware(requiredRoles []string, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "Missing role context", http.StatusForbidden)
				return
			}

			// Check if user role is in required roles
			allowed := false
			for _, required := range requiredRoles {
				if userRole == required {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers for all requests
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Tenant-ID, X-User-Role, X-User-ID")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ErrorRecoveryMiddleware recovers from panics
func ErrorRecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("Panic recovered", "error", err, "path", r.URL.Path)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestLoggingMiddleware logs all incoming requests
func RequestLoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Incoming request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}
}
