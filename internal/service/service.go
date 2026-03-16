package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"w2work3/internal/model"
	"w2work3/internal/repository"
	jwtutil "w2work3/internal/utils/jwt"
	pass "w2work3/internal/utils/password"

	"gorm.io/gorm"
)

//user

type AuthService struct {
	userrepo    repository.UserRepository
	jwtsecret   string
	jwtexpirehr int
}

func NewAuthService(userrepo repository.UserRepository, jwtsecret string, jwtexpirehr int) *AuthService {
	return &AuthService{userrepo: userrepo, jwtsecret: jwtsecret, jwtexpirehr: jwtexpirehr}
}
func (s *AuthService) SignupUser(ctx context.Context, username string, password string) (uint, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return 0, errors.New("invalid username")
	}
	if len(password) < 6 {
		return 0, errors.New("password too short")
	}
	if len(password) > 72 {
		return 0, errors.New("password too long")
	}
	_, err := s.userrepo.GetUserByName(ctx, username)
	if err == nil {
		return 0, errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	hashed, err := pass.Hash(password)
	if err != nil {
		return 0, err
	}
	user := &model.User{
		UserName:     username,
		PasswordHash: hashed,
	}
	if err := s.userrepo.CreateUser(ctx, user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *AuthService) LoginUser(ctx context.Context, username string, password string) (string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return "", errors.New("invalid username")
	}
	if len(password) < 6 || len(password) > 72 {
		return "", errors.New("invalid password")
	}
	user, err := s.userrepo.GetUserByName(ctx, username)
	if err != nil {
		return "", err
	}
	if !pass.Verify(user.PasswordHash, password) {
		return "", errors.New("wrong password")
	}
	return jwtutil.GenerateToken(s.jwtsecret, s.jwtexpirehr, user.ID, user.UserName)
}
func (s *AuthService) DeleteUser(ctx context.Context, username string, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return errors.New("invalid username")
	}
	if len(password) < 6 || len(password) > 72 {
		return errors.New("invalid password")
	}
	user, err := s.userrepo.GetUserByName(ctx, username)
	if err != nil {
		return err
	}
	if !pass.Verify(user.PasswordHash, password) {
		return errors.New("wrong password")
	}
	return s.userrepo.DeleteUser(ctx, user.ID)
}

//todo

type TodoService struct {
	todorepo repository.TodoRepository
}

func NewTodoService(todorepo repository.TodoRepository) TodoService {
	return TodoService{todorepo: todorepo}
}

func (s *TodoService) AddTodo(ctx context.Context, userid uint, title string, content string, startat time.Time, endat time.Time) (uint, error) {
	todo := &model.Todo{
		UserID:  userid,
		Title:   title,
		Content: content,
		StartAt: startat,
		EndAt:   endat,
		Views:   0,
	}
	if err := s.todorepo.CreateTodo(ctx, todo); err != nil {
		return 0, err
	}
	return todo.ID, nil
}

func (s *TodoService) ListTodo(ctx context.Context, lmt repository.TodoQueryLimit) ([]model.Todo, int64, error) {
	if lmt.UserID == 0 {
		return nil, 0, errors.New("invalid user")
	}
	if lmt.Page <= 0 {
		lmt.Page = 1
	}
	if lmt.PageSize <= 0 {
		lmt.PageSize = 1
	}
	var todos []model.Todo
	var total int64
	todos, total, err := s.todorepo.GetTodos(ctx, lmt)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, errors.New("todo not found")
	}
	return todos, total, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, ids []uint, userid uint, todo *model.Todo, conds ...string) error {
	if err := validateConds(conds); err != nil {
		return err
	}
	if len(ids) == 0 {
		return errors.New("invalid id")
	}
	if userid == 0 {
		return errors.New("invalid userid")
	}
	return s.todorepo.UpdateTodo(ctx, ids, userid, todo, conds...)
}

func (s *TodoService) DeleteTodo(ctx context.Context, userid uint, id uint) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	if userid == 0 {
		return errors.New("invalid userid")
	}
	return s.todorepo.DeleteTodoByID(ctx, userid, id)
}

func (s *TodoService) DeleteTodoByStatus(ctx context.Context, userid uint, status bool) (int64, error) {
	if userid == 0 {
		return 0, errors.New("invalid userid")
	}
	return s.todorepo.DeleteTodoByStatus(ctx, userid, status)
}

func (s *TodoService) DeleteAllTodos(ctx context.Context, userid uint) (int64, error) {
	if userid == 0 {
		return 0, errors.New("invalid userid")
	}
	return s.todorepo.DeleteAllTodos(ctx, userid)
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
			return errors.New("invalid cond: " + cond)
		}
	}
	return nil
}
