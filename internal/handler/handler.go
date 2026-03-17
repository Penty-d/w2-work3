package handler

import (
	"context"
	"time"
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
		})
		return
	}
	if _, err := h.authsvc.SignupUser(ctx, req.UserName, req.PassWord); err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
	})
}

func (h *AuthHandler) LoginUser(ctx context.Context, c *app.RequestContext) {
	var req AuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	token, err := h.authsvc.LoginUser(ctx, req.UserName, req.PassWord)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
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
		})
		return
	}
	if err := h.authsvc.DeleteUser(ctx, req.UserName, req.PassWord); err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
	})
}

//todo

type TodoHandler struct {
	todosvc *service.TodoService
}

type Todo struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	StartAt time.Time `json:"startat"`
	EndAt   time.Time `json:"endat"`
}

type QueryConditions struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pagesize"`
	Status   bool   `json:"status,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
}

func NewTodoHandler(todosvc *service.TodoService) *TodoHandler {
	return &TodoHandler{todosvc: todosvc}
}

func (h *TodoHandler) AddTodo(ctx context.Context, c *app.RequestContext) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
		})
	}
	var todo Todo
	if err := c.BindAndValidate(&todo); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	ID, err := h.todosvc.AddTodo(ctx, userid, todo.Title, todo.Content, todo.StartAt, todo.EndAt)
	if err != nil || ID == 0 {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data:   map[string]uint{"TodoID": ID},
	})
}

func (h *TodoHandler) ListTodo(ctx context.Context, c *app.RequestContext) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
		})
	}
	var conds QueryConditions
	if err := c.BindAndValidate(&conds); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}

}
