package auth

import (
	"context"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(UserContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}
