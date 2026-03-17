package handler

import (
	"context"
	"w2work3/internal/constant"
	"w2work3/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data,omitempty"`
}

//user

type AuthHandler struct {
	authsvc *service.AuthService
}

func NewAuthHandler(authsvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authsvc: authsvc}
}

type AuthRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func (h *AuthHandler) SignupUser(ctx context.Context, c *app.RequestContext) {
	var req AuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
			Data:   nil,
		})
		return
	}
	if _, err := h.authsvc.SignupUser(ctx, req.UserName, req.PassWord); err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
			Data:   nil,
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data:   nil,
	})
}

func (h *AuthHandler) LoginUser(ctx context.Context, c *app.RequestContext) {
	var req AuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
			Data:   nil,
		})
		return
	}
	token, err := h.authsvc.LoginUser(ctx, req.UserName, req.PassWord)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
			Data:   nil,
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data:   map[string]string{"Token": token},
	})
}

func (h *AuthHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	var req AuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
			Data:   nil,
		})
		return
	}
	if err := h.authsvc.DeleteUser(ctx, req.UserName, req.PassWord); err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
			Data:   nil,
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data:   nil,
	})
}

//todo

type TodoHandler struct {
	todosvc *service.TodoService
}

func NewTodoHandler(todosvc *service.TodoService) *TodoHandler {
	return &TodoHandler{todosvc: todosvc}
}

func AddTodo(ctx context.Context, c *app.RequestContext) {
	userid, exists := c.Get("userid")
	if !exists || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
			Data:   nil,
		})
	}

	c.BindAndValidate()
}
