package context

import (
	"context"

	"github.com/jhampac/picha/model"
)

// WithUser is a wrapper for a custom context object
func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, "user", user)
}
