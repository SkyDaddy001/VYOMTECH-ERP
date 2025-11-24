package handlers

import (
    "encoding/json"
    "net/http"

    "multi-tenant-ai-callcenter/internal/services"
)

type PasswordResetHandler struct {
    resetService *services.PasswordResetService
}

func NewPasswordResetHandler(resetService *services.PasswordResetService) *PasswordResetHandler {
    return &PasswordResetHandler{
        resetService: resetService,
    }
}

func (h *PasswordResetHandler) RequestReset(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if err := h.resetService.RequestPasswordReset(req.Email); err != nil {
        http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Reset email sent"})
}

func (h *PasswordResetHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Token       string `json:"token"`
        NewPassword string `json:"new_password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if err := h.resetService.ResetPassword(req.Token, req.NewPassword); err != nil {
        http.Error(w, "Failed to reset password", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successful"})
}
