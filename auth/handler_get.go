package auth

import (
	"encoding/json"
	"net/http"
)

type GetHandler struct {
	service *Service
}

func NewGetHandler(service *Service) *GetHandler {
	return &GetHandler{
		service: service,
	}
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	cookie, err := r.Cookie(h.service.CookieName())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetFromSession(ctx, cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var resp = struct {
		Username string `json:"username"`
	}{
		Username: user.Username,
	}

	json.NewEncoder(w).Encode(resp)
}
