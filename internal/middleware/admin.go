package middleware

import (
	"net/http"

	"vyomtech-backend/pkg/logger"
)

// AdminMiddleware validates that the user has admin role
func AdminMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract role from context (set by auth middleware)
			role, ok := r.Context().Value("role").(string)
			if !ok {
				log.Warn("No role found in context")
				http.Error(w, "forbidden: admin role required", http.StatusForbidden)
				return
			}

			// Check if user is admin
			if role != "admin" && role != "super_admin" {
				userID, _ := r.Context().Value("userID").(int64)
				log.Warn("Non-admin user attempted admin operation", "userID", userID, "role", role)
				http.Error(w, "forbidden: admin role required", http.StatusForbidden)
				return
			}

			// User is admin, proceed
			next.ServeHTTP(w, r)
		})
	}
}
