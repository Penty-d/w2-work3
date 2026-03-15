package service

import (
	"context"
	"errors"
	"strings"
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
	hashed, err := pass.Hash(password)
	if err != nil {
		return err
	}
	return r.userrepo.CreateUser(ctx, &model.User{
		UserName:     username,
		PasswordHash: hashed,
	})

}
func (r *AuthService) LoginUser(ctx context.Context, username string, password string) (string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return "", errors.New("invalid username")
	}
	if len(password) < 6 || len(password) > 72 {
		return "", errors.New("invalid password")
	}
	user, err := r.userrepo.GetUserByName(ctx, username)
	if err != nil {
		return "", err
	}
	if !pass.Verify(user.PasswordHash, password) {
		return "", errors.New("wrong password")
	}
	return jwtutil.GenerateToken(r.jwtsecret, r.jwtexpirehr, user.ID, user.UserName)
}
func (r *AuthService) DeleteUser(ctx context.Context, username string, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return errors.New("invalid username")
	}
	if len(password) < 6 || len(password) > 72 {
		return errors.New("invalid password")
	}
	user, err := r.userrepo.GetUserByName(ctx, username)
	if err != nil {
		return err
	}
	if !pass.Verify(user.PasswordHash, password) {
		return errors.New("wrong password")
	}
	return r.userrepo.DeleteUser(ctx, user.ID)
}

//todo
