package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

type DeleteHandler struct {
	service *Service
}

func NewDeleteHandler(service *Service) *DeleteHandler {
	return &DeleteHandler{
		service: service,
	}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.service.DeleteUser(ctx)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("User deleted successfully"); err != nil {
		log.ErrorContext(ctx, "failed to decode response to json", "error", err)
	}
}
