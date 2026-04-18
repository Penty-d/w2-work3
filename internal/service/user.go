package service

import (
	"context"
	"errors"
	"strings"
	"w2work3/internal/apperr"
	"w2work3/internal/model"
	jwtutil "w2work3/internal/utils/jwt"
	pass "w2work3/internal/utils/password"

	"gorm.io/gorm"
)

type AuthService struct {
	userrepo    UserRepo
	jwtsecret   string
	jwtexpirehr int
}

func NewAuthService(userrepo UserRepo, jwtsecret string, jwtexpirehr int) *AuthService {
	return &AuthService{userrepo: userrepo, jwtsecret: jwtsecret, jwtexpirehr: jwtexpirehr}
}

func (s *AuthService) SignupUser(ctx context.Context, username string, password string) (uint, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return 0, apperr.InvalidRequest("invalid request")
	}
	if len(password) < 6 {
		return 0, apperr.InvalidRequest("invalid request")
	}
	if len(password) > 72 {
		return 0, apperr.InvalidRequest("invalid request")
	}

	_, err := s.userrepo.GetUserByName(ctx, username)
	if err == nil {
		return 0, apperr.Conflict("resource already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, apperr.Internal(err)
	}

	hashed, err := pass.Hash(password)
	if err != nil {
		return 0, apperr.Internal(err)
	}
	user := &model.User{UserName: username, PasswordHash: hashed}
	if err := s.userrepo.CreateUser(ctx, user); err != nil {
		return 0, apperr.Internal(err)
	}
	return user.ID, nil
}

func (s *AuthService) LoginUser(ctx context.Context, username string, password string) (string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return "", apperr.InvalidRequest("invalid request")
	}
	if len(password) < 6 || len(password) > 72 {
		return "", apperr.InvalidRequest("invalid request")
	}

	user, err := s.userrepo.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperr.Unauthorized("invalid username or password")
		}
		return "", apperr.Internal(err)
	}
	if !pass.Verify(user.PasswordHash, password) {
		return "", apperr.Unauthorized("invalid username or password")
	}

	token, err := jwtutil.GenerateToken(s.jwtsecret, s.jwtexpirehr, user.ID, user.UserName)
	if err != nil {
		return "", apperr.Internal(err)
	}
	return token, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, username string, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || len(username) < 4 {
		return apperr.InvalidRequest("invalid request")
	}
	if len(password) < 6 || len(password) > 72 {
		return apperr.InvalidRequest("invalid request")
	}

	user, err := s.userrepo.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Unauthorized("invalid username or password")
		}
		return apperr.Internal(err)
	}
	if !pass.Verify(user.PasswordHash, password) {
		return apperr.Unauthorized("invalid username or password")
	}
	if err := s.userrepo.DeleteUser(ctx, user.ID); err != nil {
		return apperr.Internal(err)
	}
	return nil
}
