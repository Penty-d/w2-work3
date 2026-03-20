package middleware

import (
	"context"
	"strings"

	"w2work3/internal/handler"
	"w2work3/internal/utils/jwt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func JWTAuth(secret string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		auth := string(c.GetHeader("Authorization"))
		if auth == "" {
			c.JSON(consts.StatusUnauthorized, handler.Response{
				Status: consts.StatusUnauthorized,
				Msg:    "missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			c.JSON(consts.StatusUnauthorized, handler.Response{
				Status: consts.StatusUnauthorized,
				Msg:    "invalid authorization header",
			})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(secret, parts[1])
		if err != nil {
			c.JSON(consts.StatusUnauthorized, handler.Response{
				Status: consts.StatusUnauthorized,
				Msg:    "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("userid", claims.UserID)
		c.Set("username", claims.UserName)

		c.Next(ctx)
	}
}
