package api

import (
	"BlueNetDisk/pkg/ctl"
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/service"
	"BlueNetDisk/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadFileHandler handles multipart/form-data.
func UploadFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			utils.Logrusobj.Error(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		} else {
			l := service.GetFileSrv()
			resp, err := l.UploadManyFiles(c.Request.Context(), form)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

func ListFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ListFileReq types.ListFileReq
		if err := c.ShouldBindJSON(&ListFileReq); err != nil {
			utils.Logrusobj.Error(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		} else {
			fmt.Println("received ListFileReq : ", ListFileReq)
			l := service.GetFileSrv()
			resp, err := l.ListFileByParentUUID(c.Request.Context(), ListFileReq.ParentId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}

func MkdirHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var MkdirReq types.MkdirReq
		if err := c.ShouldBindJSON(&MkdirReq); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		} else {
			l := service.GetFileSrv()
			//这里的返回值不规范，后续还会修改
			resp, err := l.CreateDir(c.Request.Context(), MkdirReq.DirName, MkdirReq.ParentId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusOK, ctl.RespSuccessWithData(resp))
		}
	}
}
