package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"w2work3/internal/apperr"
	"w2work3/internal/model"

	"gorm.io/gorm"
)

type TodoService struct {
	todorepo TodoRepo
}

func NewTodoService(todorepo TodoRepo) *TodoService {
	return &TodoService{todorepo: todorepo}
}

func (s *TodoService) AddTodo(ctx context.Context, userid uint, title string, content string, startat time.Time, endat time.Time) (uint, error) {
	if userid == 0 {
		return 0, apperr.InvalidRequest("invalid request")
	}
	if strings.TrimSpace(title) == "" {
		return 0, apperr.InvalidRequest("invalid request")
	}
	if strings.TrimSpace(content) == "" {
		return 0, apperr.InvalidRequest("invalid request")
	}
	if startat.After(endat) {
		return 0, apperr.InvalidRequest("invalid request")
	}

	todo := &model.Todo{
		UserID:  userid,
		Title:   title,
		Content: content,
		StartAt: startat,
		EndAt:   endat,
		Views:   0,
	}
	if err := s.todorepo.CreateTodo(ctx, todo); err != nil {
		return 0, apperr.Internal(err)
	}
	return todo.ID, nil
}

func (s *TodoService) ListTodo(ctx context.Context, conds model.TodoQueryConditions) ([]model.Todo, int64, error) {
	if conds.UserID == 0 || conds.Page <= 0 || conds.PageSize <= 0 {
		return nil, 0, apperr.InvalidRequest("invalid request")
	}
	todos, total, err := s.todorepo.GetTodos(ctx, conds)
	if err != nil {
		return nil, 0, apperr.Internal(err)
	}
	return todos, total, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo *model.Todo, conds ...string) error {
	if err := validateConds(conds); err != nil {
		return err
	}
	if todo.ID == 0 || todo.UserID == 0 {
		return apperr.InvalidRequest("invalid request")
	}
	ids := []uint{todo.ID}
	if err := s.todorepo.UpdateTodo(ctx, ids, todo.UserID, todo, conds...); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("resource not found")
		}
		return apperr.Internal(err)
	}
	return nil
}

func (s *TodoService) UpdateTodosStatus(ctx context.Context, ids []uint, userid uint, status bool) error {
	if len(ids) == 0 || userid == 0 {
		return apperr.InvalidRequest("invalid request")
	}
	if err := s.todorepo.UpdateTodo(ctx, ids, userid, &model.Todo{Status: status}, "status"); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("resource not found")
		}
		return apperr.Internal(err)
	}
	return nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, userid uint, ids []uint) error {
	if len(ids) == 0 || userid == 0 {
		return apperr.InvalidRequest("invalid request")
	}
	if err := s.todorepo.DeleteTodosByID(ctx, userid, ids); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.NotFound("resource not found")
		}
		return apperr.Internal(err)
	}
	return nil
}

func (s *TodoService) DeleteTodoByStatus(ctx context.Context, userid uint, status bool) (int64, error) {
	if userid == 0 {
		return 0, apperr.InvalidRequest("invalid request")
	}
	total, err := s.todorepo.DeleteTodoByStatus(ctx, userid, status)
	if err != nil {
		return 0, apperr.Internal(err)
	}
	if total == 0 {
		return 0, apperr.NotFound("resource not found")
	}
	return total, nil
}

func (s *TodoService) DeleteAllTodos(ctx context.Context, userid uint) (int64, error) {
	if userid == 0 {
		return 0, apperr.InvalidRequest("invalid request")
	}
	total, err := s.todorepo.DeleteAllTodos(ctx, userid)
	if err != nil {
		return 0, apperr.Internal(err)
	}
	if total == 0 {
		return 0, apperr.NotFound("resource not found")
	}
	return total, nil
}

func validateConds(conds []string) error {
	allowed := map[string]struct{}{
		"title":    {},
		"content":  {},
		"status":   {},
		"start_at": {},
		"end_at":   {},
		"views":    {},
	}
	for _, cond := range conds {
		if _, ok := allowed[cond]; !ok {
			return apperr.InvalidRequest("invalid request")
		}
	}
	return nil
}
