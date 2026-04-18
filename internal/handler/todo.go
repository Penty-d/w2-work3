package handler

import (
	"context"
	"encoding/json"
	"time"
	"w2work3/internal/constant"
	"w2work3/internal/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type TodoService interface {
	AddTodo(ctx context.Context, userid uint, title string, content string, startat time.Time, endat time.Time) (uint, error)
	ListTodo(ctx context.Context, conds model.TodoQueryConditions) ([]model.Todo, int64, error)
	UpdateTodo(ctx context.Context, todo *model.Todo, conds ...string) error
	UpdateTodosStatus(ctx context.Context, ids []uint, userid uint, status bool) error
	DeleteTodo(ctx context.Context, userid uint, ids []uint) error
	DeleteTodoByStatus(ctx context.Context, userid uint, status bool) (int64, error)
	DeleteAllTodos(ctx context.Context, userid uint) (int64, error)
}

type TodoHandler struct {
	todosvc TodoService
}

func NewTodoHandler(todosvc TodoService) *TodoHandler {
	return &TodoHandler{todosvc: todosvc}
}

type TodoRequest struct {
	IDs    []uint `json:"ids"`
	Status bool   `json:"status"`
}

type TodoCreateRequest struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	StartAt time.Time `json:"start_time"`
	EndAt   time.Time `json:"end_time"`
}

type TodoListRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pagesize"`
	Status   *bool  `json:"status,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
}

type TodoUpdateRequest struct {
	ID      uint       `json:"id"`
	Title   *string    `json:"title,omitempty"`
	Content *string    `json:"content,omitempty"`
	Status  *bool      `json:"status,omitempty"`
	StartAt *time.Time `json:"start_time,omitempty"`
	EndAt   *time.Time `json:"end_time,omitempty"`
	Views   *uint      `json:"view,omitempty"`
}

func (h *TodoHandler) AddTodo(ctx context.Context, c *app.RequestContext) {
	userid, ok := currentUserID(c)
	if !ok {
		return
	}
	var req TodoCreateRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	id, err := h.todosvc.AddTodo(ctx, userid, req.Title, req.Content, req.StartAt, req.EndAt)
	if err != nil {
		writeError(c, err)
		return
	}
	if id == 0 {
		c.JSON(consts.StatusInternalServerError, Response{Status: constant.StatusFailed, Msg: "failed to create todo"})
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success", Data: map[string]uint{"TodoID": id}})
}

func (h *TodoHandler) ListTodo(ctx context.Context, c *app.RequestContext) {
	userid, ok := currentUserID(c)
	if !ok {
		return
	}
	var req TodoListRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}

	conds := model.TodoQueryConditions{
		UserID:   userid,
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Keyword:  req.Keyword,
	}
	todos, total, err := h.todosvc.ListTodo(ctx, conds)
	if err != nil {
		writeError(c, err)
		return
	}
	if todos == nil {
		todos = make([]model.Todo, 0)
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
	userid, ok := currentUserID(c)
	if !ok {
		return
	}
	body := c.Request.Body()

	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	var req TodoUpdateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}

	todo := &model.Todo{ID: req.ID, UserID: userid}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Content != nil {
		todo.Content = *req.Content
	}
	if req.Status != nil {
		todo.Status = *req.Status
	}
	if req.StartAt != nil {
		todo.StartAt = *req.StartAt
	}
	if req.EndAt != nil {
		todo.EndAt = *req.EndAt
	}
	if req.Views != nil {
		todo.Views = *req.Views
	}

	conds := make([]string, 0, len(raw))
	for key := range raw {
		if key == "id" {
			continue
		}
		cond, ok := normalizeUpdateField(key)
		if !ok {
			c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
			return
		}
		conds = append(conds, cond)
	}
	if len(conds) == 0 {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}

	if err := h.todosvc.UpdateTodo(ctx, todo, conds...); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success"})
}

func (h *TodoHandler) UpdateTodosStatus(ctx context.Context, c *app.RequestContext) {
	userid, ok := currentUserID(c)
	if !ok {
		return
	}
	var req TodoRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	if err := h.todosvc.UpdateTodosStatus(ctx, req.IDs, userid, req.Status); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success"})
}

func (h *TodoHandler) DeleteTodo(ctx context.Context, c *app.RequestContext) {
	userid, ok := currentUserID(c)
	if !ok {
		return
	}
	var req TodoRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	if err := h.todosvc.DeleteTodo(ctx, userid, req.IDs); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success"})
}

func (h *TodoHandler) DeleteTodosByStatus(ctx context.Context, c *app.RequestContext) {
	userid, ok := currentUserID(c)
	if !ok {
		return
	}
	var req TodoRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	total, err := h.todosvc.DeleteTodoByStatus(ctx, userid, req.Status)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success", Data: map[string]int64{"total": total}})
}

func (h *TodoHandler) DeleteAllTodos(ctx context.Context, c *app.RequestContext) {
	userid, ok := currentUserID(c)
	if !ok {
		return
	}
	total, err := h.todosvc.DeleteAllTodos(ctx, userid)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success", Data: map[string]int64{"total": total}})
}

func normalizeUpdateField(key string) (string, bool) {
	switch key {
	case "title", "content", "status":
		return key, true
	case "start_time":
		return "start_at", true
	case "end_time":
		return "end_at", true
	case "view":
		return "views", true
	default:
		return "", false
	}
}

func currentUserID(c *app.RequestContext) (uint, bool) {
	uid, exists := c.Get("userid")
	userid, ok := uid.(uint)
	if !exists || !ok || userid == 0 {
		c.JSON(consts.StatusUnauthorized, Response{Status: constant.StatusUnauthorized, Msg: "unauthorized"})
		return 0, false
	}
	return userid, true
}
