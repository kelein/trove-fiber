package handler

import (
	"github.com/gofiber/fiber"
)

// Response Message Constants
const (
	MessageSuccess = "success"
	MessageFailed  = "failed"
)

// ServerResponse stands for server response
type ServerResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	ReqID   string `json:"reqid,omitempty"`
	Message string `json:"message"`
}

// BaseHandler for common handler methods
type BaseHandler struct{}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// Succeed handles successful responses
func (h *BaseHandler) Succeed(ctx *fiber.Ctx, data any) error {
	return ctx.JSON(ServerResponse{
		Data:    data,
		Code:    fiber.StatusOK,
		Message: MessageSuccess,
	})
}

// Failed handles failed responses
func (h *BaseHandler) Failed(ctx *fiber.Ctx, code int, err error) error {
	return ctx.JSON(ServerResponse{
		Code:    code,
		Data:    map[string]any{},
		Message: err.Error(),
	})
}

// ParseUserID parses the user ID from the context claims
func (h *BaseHandler) ParseUserID(ctx *fiber.Ctx) string {
	// 	v, exists := ctx.Get("claims")
	// 	if !exists {
	// 		return ""
	// 	}
	// 	return v.(*jwt.MyCustomClaims).UserId
	return ctx.Get("claims", "")
}
