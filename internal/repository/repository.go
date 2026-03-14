package repository

import (
	"context"
	"w2work3/internal/config"
	"w2work3/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open((cfg.DB.DSN())), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	/*
		if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
			return nil, err
		} 自动建表，但不需要
	*/
	return db, nil
}

// user

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *UserRepository) GetUserByName(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("user_name = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// todo

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository { return &TodoRepository{db: db} }

func (r *TodoRepository) CreateTodo(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *TodoRepository) UpdateTodoStatus(ctx context.Context, id uint, status bool) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Update("status", status).Error
}
