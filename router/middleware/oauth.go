package middleware

import (
	"net/http"
	"time"

	"github.com/hohice/gin-web/pkg/jwt"
	"github.com/hohice/gin-web/router/ex"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int

		code = ex.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = ex.INVALID_PARAMS
		} else {
			claims, err := jwt.ParseToken(token)
			if err != nil {
				code = ex.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = ex.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != ex.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  ex.GetMsg(code),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
