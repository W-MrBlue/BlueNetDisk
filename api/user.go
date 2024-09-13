package api

import (
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/service"
	"BlueNetDisk/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserRegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Logrusobj.Infoln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		} else {
			l := service.GetUserSrv()
			resp, err := l.UserRegister(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

func UserLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserLoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Logrusobj.Infoln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		} else {
			l := service.GetUserSrv()
			resp, err := l.UserLogin(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}
