package middleware

import "github.com/jhampac/picha/model"

type RequireUser struct {
	model.UserService
}
