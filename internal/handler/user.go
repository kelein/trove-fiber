package handler

import "github.com/kelein/trove-fiber/internal/service"

// UserHandler handles user related requests
type UserHandler struct {
	*BaseHandler
	userService service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(handler *BaseHandler, userService service.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: handler,
		userService: userService,
	}
}
