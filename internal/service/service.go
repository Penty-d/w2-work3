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

func (s *TodoService) 
