package context

import (
	"context"

	"github.com/jhampac/picha/model"
)

type privateContextKey string

const (
	userKey privateContextKey = "user"
)

// WithUser is a wrapper for a custom context object; this guarantees that the value we get back will always be a user
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User retrieves a user that was attached to the context
func User(ctx context.Context) *model.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*model.User); ok {
			return user
		}
	}
	return nil
}
