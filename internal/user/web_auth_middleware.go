package user

import (
	"fmt"
	"marketplace_server/internal/servers/web/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationKey = "Authorization"
	UserIDKey        = "username"
)

type AuthMiddleware struct {
	UserApp UserAppInterface
}

func NewAuthMiddleware(userApp UserAppInterface) *AuthMiddleware {
	return &AuthMiddleware{
		UserApp: userApp,
	}
}

func (a *AuthMiddleware) Auth(c *gin.Context) {
	// 获取 token
	token := c.GetHeader(AuthorizationKey)
	if token == "" {
		response.Err(c, http.StatusUnauthorized, "token is empty")
		c.Abort()
		return
	}

	// token認證失敗
	authInfo, err := a.UserApp.GetAuthInfo(token)
	if err != nil {
		response.Err(c, http.StatusUnauthorized, fmt.Sprintf("token auth fail msg:%s", err.Error()))
		c.Abort()
		return
	}

	// 保存用户信息
	c.Set(UserIDKey, authInfo.UserID)
}
