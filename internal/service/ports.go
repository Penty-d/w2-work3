package service

import (
	"context"
	"w2work3/internal/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
	GetUserByName(ctx context.Context, username string) (*model.User, error)
}

type TodoRepo interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
	UpdateTodo(ctx context.Context, ids []uint, userid uint, todo *model.Todo, conds ...string) error
	GetTodos(ctx context.Context, conds model.TodoQueryConditions) ([]model.Todo, int64, error)
	DeleteTodosByID(ctx context.Context, userid uint, ids []uint) error
	DeleteTodoByStatus(ctx context.Context, userid uint, status bool) (int64, error)
	DeleteAllTodos(ctx context.Context, userid uint) (int64, error)
}
