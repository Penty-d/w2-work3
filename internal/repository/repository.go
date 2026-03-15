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
		Model(&model.User{}).
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

type TodoQueryLimit struct {
	UserID   uint
	Page     int
	PageSize int
	Status   *bool
	Keyword  string
}

func NewTodoRepository(db *gorm.DB) *TodoRepository { return &TodoRepository{db: db} }

func (r *TodoRepository) CreateTodo(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, todo *model.Todo, fields ...string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Where("id = ? AND user_id = ?", todo.ID, todo.UserID).
		Select(fields).
		Updates(todo)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TodoRepository) GetTodos(ctx context.Context, lmt TodoQueryLimit) ([]model.Todo, int64, error) {
	if lmt.Page <= 0 {
		lmt.Page = 1
	}
	if lmt.PageSize <= 0 {
		lmt.PageSize = 1
	}

	db := r.db.WithContext(ctx).Model(&model.Todo{}).Where("user_id = ?", lmt.UserID)
	if lmt.Status != nil {
		db = db.Where("status = ?", *lmt.Status)
	}
	if lmt.Keyword != "" {
		kw := "%" + lmt.Keyword + "%"
		db = db.Where("(title ILIKE ? OR content ILIKE ?)", kw, kw)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (lmt.Page - 1) * lmt.PageSize
	var rtodos []model.Todo
	//补一个按时间排序
	db = db.Order("created_at DESC")
	if err := db.Limit(lmt.PageSize).Offset(offset).Find(&rtodos).Error; err != nil {
		return nil, 0, err
	}
	return rtodos, total, nil
}

func (r *TodoRepository) DeleteTodoByID(ctx context.Context, userid uint, id uint) error { //删单个
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userid).Delete(&model.Todo{}).Error
}

func (r *TodoRepository) DeleteTodoByStatus(ctx context.Context, userid uint, status *bool) (int64, error) { //偷个懒，status留空为全删
	db := r.db.WithContext(ctx).Where("user_id = ?", userid)
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	result := db.Delete(&model.Todo{})
	return result.RowsAffected, result.Error
}
