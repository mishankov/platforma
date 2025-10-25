package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/mishankov/platforma/log"
)

type userService interface {
	GetFromSession(ctx context.Context, sessionId string) (*User, error)
	CookieName() string
}

type AuthenticationMiddleware struct {
	userService userService
}

func NewAuthenticationMiddleware(userService userService) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{userService: userService}
}

func (m *AuthenticationMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(m.userService.CookieName())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := m.userService.GetFromSession(r.Context(), cookie.Value)
		if errors.Is(err, ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "failed to get user", http.StatusInternalServerError)
			return
		}

		ctxWithUserId := context.WithValue(r.Context(), log.UserIdKey, user.ID)
		ctxWithUser := context.WithValue(ctxWithUserId, UserContextKey, user)
		requestWithUser := r.WithContext(ctxWithUser)

		next.ServeHTTP(w, requestWithUser)
	})
}
