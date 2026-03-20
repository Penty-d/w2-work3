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
	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		return nil, err
	}
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
	if len(conds.Keywords) > 0 {
		for _, kw := range conds.Keywords {
			like := "%" + kw + "%"
			db = db.Where("(title ILIKE ? OR content ILIKE ?)", like, like)
		}
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (conds.Page - 1) * conds.PageSize
	var rtodos []model.Todo
	//补一个按时间排序
	db = db.Order("created_at DESC")
	if err := db.Limit(conds.PageSize).Offset(offset).Find(&rtodos).Error; err != nil {
		return nil, 0, err
	}
	return rtodos, total, nil
}

func (r *TodoRepository) DeleteTodosByID(ctx context.Context, userid uint, ids []uint) error { //删单个
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
	db := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userid, status)
	result := db.Delete(&model.Todo{})
	return result.RowsAffected, result.Error
}

func (r *TodoRepository) DeleteAllTodos(ctx context.Context, userid uint) (int64, error) {
	db := r.db.WithContext(ctx).Where("user_id = ?", userid)
	result := db.Delete(&model.Todo{})
	return result.RowsAffected, result.Error
}
