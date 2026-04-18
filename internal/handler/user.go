package handler

import (
	"context"
	"errors"
	"log"
	"w2work3/internal/apperr"
	"w2work3/internal/constant"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data,omitempty"`
}

func writeError(c *app.RequestContext, err error) {
	if err != nil {
		log.Printf("request error: %v", err)
	}
	httpCode, bizCode, msg := mapError(err)
	c.JSON(httpCode, Response{
		Status: bizCode,
		Msg:    msg,
	})
}

func mapError(err error) (int, int, string) {
	if err == nil {
		return consts.StatusInternalServerError, constant.StatusFailed, "internal server error"
	}
	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case apperr.CodeInvalidRequest:
			return consts.StatusBadRequest, constant.StatusInvalidRequest, appErr.Message
		case apperr.CodeUnauthorized:
			return consts.StatusUnauthorized, constant.StatusUnauthorized, appErr.Message
		case apperr.CodeNotFound:
			return consts.StatusNotFound, constant.StatusNotFound, appErr.Message
		case apperr.CodeConflict:
			return consts.StatusConflict, constant.StatusFailed, appErr.Message
		default:
			return consts.StatusInternalServerError, constant.StatusFailed, "internal server error"
		}
	}
	return consts.StatusInternalServerError, constant.StatusFailed, "internal server error"
}

type AuthService interface {
	SignupUser(ctx context.Context, username string, password string) (uint, error)
	LoginUser(ctx context.Context, username string, password string) (string, error)
	DeleteUser(ctx context.Context, username string, password string) error
}

type AuthHandler struct {
	authsvc AuthService
}

func NewAuthHandler(authsvc AuthService) *AuthHandler {
	return &AuthHandler{authsvc: authsvc}
}

type AuthRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func (h *AuthHandler) SignupUser(ctx context.Context, c *app.RequestContext) {
	var req AuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	if _, err := h.authsvc.SignupUser(ctx, req.UserName, req.PassWord); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success"})
}

func (h *AuthHandler) LoginUser(ctx context.Context, c *app.RequestContext) {
	var req AuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	token, err := h.authsvc.LoginUser(ctx, req.UserName, req.PassWord)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{
		Status: constant.StatusOK,
		Msg:    "success",
		Data:   map[string]string{"token": token},
	})
}

func (h *AuthHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	var req AuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{Status: constant.StatusInvalidRequest, Msg: "invalid request"})
		return
	}
	if err := h.authsvc.DeleteUser(ctx, req.UserName, req.PassWord); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, Response{Status: constant.StatusOK, Msg: "success"})
}
