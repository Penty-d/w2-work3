package repository

import (
	"context"
	"w2work3/internal/model"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository { return &TodoRepository{db: db} }

func (r *TodoRepository) CreateTodo(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, ids []uint, userid uint, todo *model.Todo, conds ...string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Todo{}).
		Where("id IN ? AND user_id = ?", ids, userid).
		Select(conds).
		Updates(todo)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TodoRepository) GetTodos(ctx context.Context, conds model.TodoQueryConditions) ([]model.Todo, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Todo{}).Where("user_id = ?", conds.UserID)
	if conds.Status != nil {
		db = db.Where("status = ?", *conds.Status)
	}
	if conds.Keyword != "" {
		kw := "%" + conds.Keyword + "%"
		db = db.Where("(title ILIKE ? OR content ILIKE ?)", kw, kw)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (conds.Page - 1) * conds.PageSize
	var todos []model.Todo
	db = db.Order("created_at DESC")
	if err := db.Limit(conds.PageSize).Offset(offset).Find(&todos).Error; err != nil {
		return nil, 0, err
	}
	return todos, total, nil
}

func (r *TodoRepository) DeleteTodosByID(ctx context.Context, userid uint, ids []uint) error {
	result := r.db.WithContext(ctx).Where("id IN ? AND user_id = ?", ids, userid).Delete(&model.Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TodoRepository) DeleteTodoByStatus(ctx context.Context, userid uint, status bool) (int64, error) {
	result := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userid, status).Delete(&model.Todo{})
	return result.RowsAffected, result.Error
}

func (r *TodoRepository) DeleteAllTodos(ctx context.Context, userid uint) (int64, error) {
	result := r.db.WithContext(ctx).Where("user_id = ?", userid).Delete(&model.Todo{})
	return result.RowsAffected, result.Error
}
