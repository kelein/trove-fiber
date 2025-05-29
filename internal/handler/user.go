package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber"

	v1 "github.com/kelein/trove-fiber/internal/api/v1"
	"github.com/kelein/trove-fiber/internal/service"
)

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

// Register godoc
// @Summary 用户注册
// @Schemes
// @Description 目前只支持邮箱登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response
// @Router /register [post]
func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	req := new(v1.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return h.Failed(ctx, http.StatusBadRequest, err)
	}
	if err := h.userService.Register(ctx.Context(), req); err != nil {
		return h.Failed(ctx, http.StatusBadRequest, err)
	}
	return h.Succeed(ctx, nil)
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var req v1.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.Failed(ctx, http.StatusBadRequest, err)
	}

	token, err := h.userService.Login(ctx.Context(), &req)
	if err != nil {
		return h.Failed(ctx, http.StatusUnauthorized, err)
	}
	return h.Succeed(ctx, token)
}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetProfileResponse
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	userID := h.ParseUserID(ctx)
	if userID == "" {
		err := errors.New("unauthorized - userID empty")
		return h.Failed(ctx, http.StatusUnauthorized, err)
	}

	user, err := h.userService.GetProfile(ctx.Context(), userID)
	if err != nil {
		return h.Failed(ctx, http.StatusBadRequest, err)
	}
	return h.Succeed(ctx, user)
}

// UpdateProfile godoc
// @Summary 修改用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user [put]
func (h *UserHandler) UpdateProfile(ctx *fiber.Ctx) error {
	userID := h.ParseUserID(ctx)
	var req v1.UpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.Failed(ctx, http.StatusBadRequest, err)
	}
	if err := h.userService.UpdateProfile(ctx.Context(), userID, &req); err != nil {
		return h.Failed(ctx, http.StatusBadRequest, err)
	}
	return h.Succeed(ctx, nil)
}
