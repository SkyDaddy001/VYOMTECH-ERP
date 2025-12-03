package middleware

import (
	"net/http"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// PermissionMiddleware checks if user has required permission
func PermissionMiddleware(rbacService *services.RBACService, requiredPermission string, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(int64)
			if !ok {
				log.Warn("User ID not found in context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tenantID, ok := r.Context().Value(TenantIDKey).(string)
			if !ok {
				log.Warn("Tenant ID not found in context")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Check permission
			hasPermission, err := rbacService.HasPermission(r.Context(), tenantID, userID, requiredPermission)
			if err != nil {
				log.Error("Failed to check permission", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				log.Warn("Permission denied", "user_id", userID, "permission", requiredPermission)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// PermissionBasedAccessMiddleware restricts endpoint access to specific roles via RBAC service
func PermissionBasedAccessMiddleware(rbacService *services.RBACService, allowedRoles []string, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(int64)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tenantID, ok := r.Context().Value(TenantIDKey).(string)
			if !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Get user roles
			userRoles, err := rbacService.GetUserRoles(r.Context(), tenantID, userID)
			if err != nil {
				log.Error("Failed to get user roles", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Check if user has one of the allowed roles
			hasRole := false
			for _, userRole := range userRoles {
				for _, allowedRole := range allowedRoles {
					if userRole.Name == allowedRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				log.Warn("Access denied due to insufficient role", "user_id", userID, "allowed_roles", allowedRoles)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AuditMiddleware logs all requests for audit trail
func AuditMiddleware(auditService *services.AuditService, log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userID, _ := ctx.Value(UserIDKey).(int64)
			tenantID, _ := ctx.Value(TenantIDKey).(string)

			// Capture response status
			wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call next handler
			next.ServeHTTP(wrappedWriter, r)

			// Log after response is written
			if userID > 0 && tenantID != "" {
				status := "success"
				if wrappedWriter.statusCode >= 400 {
					status = "failure"
				}

				_ = auditService.LogUserAction(ctx, tenantID, userID,
					r.Method,
					r.URL.Path,
					map[string]interface{}{
						"method": r.Method,
						"path":   r.URL.Path,
						"query":  r.URL.RawQuery,
						"status": wrappedWriter.statusCode,
					},
					getClientIP(r),
					r.UserAgent(),
					status,
				)
			}
		})
	}
}

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Enable XSS protection
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Clickjacking protection
			w.Header().Set("X-Frame-Options", "DENY")

			// Referrer policy
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Content Security Policy
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

			// HSTS
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware prevents abuse with rate limiting
func RateLimitMiddleware(requestsPerMinute int, log *logger.Logger) func(http.Handler) http.Handler {
	// In production, use Redis or similar for distributed rate limiting
	ipRequests := make(map[string][]int64)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			currentTime := int64(0) // Placeholder for actual implementation

			// Simplified rate limiting - in production use a proper library
			// This is a placeholder implementation
			ipRequests[clientIP] = append(ipRequests[clientIP], currentTime)

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// getClientIP extracts the client IP address from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (load balancer)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// DataMaskingMiddleware for sensitive data in logs (apply in production)
func DataMaskingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mask sensitive parameters in URL query
			maskSensitiveQuery(r)

			next.ServeHTTP(w, r)
		})
	}
}

// maskSensitiveQuery masks sensitive query parameters
func maskSensitiveQuery(r *http.Request) {
	sensitiveParams := []string{"password", "token", "key", "secret", "api_key", "ssn", "credit_card"}
	query := r.URL.Query()

	for _, param := range sensitiveParams {
		if query.Get(param) != "" {
			query.Set(param, "***MASKED***")
		}
	}

	r.URL.RawQuery = query.Encode()
}

// EnforceHTTPSMiddleware redirects HTTP to HTTPS in production
func EnforceHTTPSMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if already HTTPS or in development
			if r.Header.Get("X-Forwarded-Proto") == "http" {
				http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
