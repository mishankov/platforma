package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mishankov/platforma/log"
)

type ChangePasswordHandler struct {
	service *Service
}

func NewChangePasswordHandler(service *Service) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		service: service,
	}
}

func (h *ChangePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.service.ChangePassword(r.Context(), req.CurrentPassword, req.NewPassword)
	log.DebugContext(r.Context(), "error from change password", "error", err)

	if err != nil {
		if errors.Is(err, ErrCurrentPasswordIncorrect) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if errors.Is(err, ErrInvalidPassword) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("Password changed successfully"); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
