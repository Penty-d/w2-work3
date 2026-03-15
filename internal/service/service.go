package service

import (
	"context"
	"errors"
	"strings"
	"w2work3/internal/model"
	"w2work3/internal/repository"
	"w2work3/internal/utils"

	"gorm.io/gorm"
)

//user

type AuthService struct {
	userrepo repository.UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}
func (r *AuthService) SignupUser(ctx context.Context, username string, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return errors.New("invalid username")
	}
	if len(password) < 6 {
		return errors.New("password too short")
	}
	if len(password) > 72 {
		return errors.New("password too long")
	}
	_, err := r.userrepo.GetUserByName(ctx, username)
	if err == nil {
		return errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hashed, err := utils.Hash(password)
	if err != nil {
		return err
	}
	return r.userrepo.CreateUser(ctx, &model.User{
		UserName:     username,
		PasswordHash: hashed,
	})

}
func (r *AuthService) LoginUser(ctx context.Context, username string, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return errors.New("invalid username")
	}
	if len(password) < 6 || len(password) > 72 {
		return errors.New("invalid password")
	}

}
