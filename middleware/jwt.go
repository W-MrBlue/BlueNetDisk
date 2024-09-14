package middleware

import (
	"BlueNetDisk/pkg/ctl"
	"BlueNetDisk/pkg/e"
	"BlueNetDisk/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.Success
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			code = http.StatusUnauthorized
			c.JSON(code, gin.H{
				"code": code,
				"msg":  "Token is empty",
			})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			code = http.StatusUnauthorized
			c.JSON(code, gin.H{
				"code": code,
				"msg":  "Invalid token",
			})
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(),
			&ctl.UserInfo{Id: claims.Id, Username: claims.Username}))
		c.Next()
	}
}
