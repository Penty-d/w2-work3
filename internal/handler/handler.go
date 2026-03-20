package handler

import (
	"context"
	"encoding/json"
	"errors"
	"w2work3/internal/constant"
	"w2work3/internal/model"
	"w2work3/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(consts.StatusNotFound, Response{
				Status: constant.StatusNotFound,
				Msg:    "user not found",
			})
			return
		}
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
	err := h.authsvc.DeleteUser(ctx, req.UserName, req.PassWord)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(consts.StatusNotFound, Response{
				Status: constant.StatusNotFound,
				Msg:    "user not found",
			})
			return
		}
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

type TodoRequest struct {
	IDs    []uint `json:"ids"`
	Status bool   `json:"status"`
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
		return
	}
	var todo model.Todo
	if err := c.BindAndValidate(&todo); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	ID, err := h.todosvc.AddTodo(ctx, userid, todo.Title, todo.Content, todo.StartTime, todo.EndTime)
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
		return
	}
	var conds model.TodoQueryConditions
	if err := c.BindAndValidate(&conds); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	conds.UserID = userid
	todos, total, err := h.todosvc.ListTodo(ctx, conds)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
		})
		return
	}
	if total == 0 {
		c.JSON(consts.StatusNotFound, Response{
			Status: constant.StatusNotFound,
			Msg:    "todo not found",
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data: map[string]any{
			"items": todos,
			"total": total,
		},
	})
}

func (h *TodoHandler) UpdateTodo(ctx context.Context, c *app.RequestContext) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
		})
		return
	}
	body := c.Request.Body()
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	var todo model.Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	todo.UserID = userid
	conds := make([]string, 0, len(raw))
	for c := range raw {
		conds = append(conds, c)
	}
	err := h.todosvc.UpdateTodo(ctx, &todo, conds...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(consts.StatusNotFound, Response{
				Status: constant.StatusNotFound,
				Msg:    "todo not found",
			})
			return
		}
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

func (h *TodoHandler) UpdateTodosStatus(ctx context.Context, c *app.RequestContext) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
		})
		return
	}
	var req TodoRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	err := h.todosvc.UpdateTodosStatus(ctx, req.IDs, userid, req.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(consts.StatusNotFound, Response{
				Status: constant.StatusNotFound,
				Msg:    "todo not found",
			})
			return
		}
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

func (h *TodoHandler) DeleteTodo(ctx context.Context, c *app.RequestContext) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
		})
		return
	}
	var req TodoRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	if err := h.todosvc.DeleteTodo(ctx, userid, req.IDs); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(consts.StatusNotFound, Response{
				Status: constant.StatusNotFound,
				Msg:    "todo not found",
			})
			return
		}
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

func (h *TodoHandler) DeleteTodosByStatus(ctx context.Context, c *app.RequestContext) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
		})
		return
	}
	var req TodoRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: constant.StatusInvalidRequest,
			Msg:    "invalid request",
		})
		return
	}
	total, err := h.todosvc.DeleteTodoByStatus(ctx, userid, req.Status)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
		})
		return
	}
	if total == 0 {
		c.JSON(consts.StatusNotFound, Response{
			Status: constant.StatusNotFound,
			Msg:    "todo not found",
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data:   map[string]int64{"total": total},
	})
}

func (h *TodoHandler) DeleteAllTodos(ctx context.Context, c *app.RequestContext) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: constant.StatusUnauthorized,
			Msg:    "unauthorized",
		})
		return
	}
	total, err := h.todosvc.DeleteAllTodos(ctx, userid)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: constant.StatusFailed,
			Msg:    err.Error(),
		})
		return
	}
	if total == 0 {
		c.JSON(consts.StatusNotFound, Response{
			Status: constant.StatusNotFound,
			Msg:    "todo not found",
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data:   map[string]int64{"total": total},
	})
}
