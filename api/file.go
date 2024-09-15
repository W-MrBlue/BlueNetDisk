package api

import (
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/service"
	"BlueNetDisk/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
)

// UploadFileHandler handles multipart/form-data.
func UploadFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			utils.Logrusobj.Error(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
		l := service.GetFileSrv()
		resp, err := l.UploadManyFiles(c.Request.Context(), form)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func ListFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ListFileReq types.ListFileReq
		if err := c.ShouldBindJSON(&ListFileReq); err != nil {
			utils.Logrusobj.Error(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
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

func MkdirHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var MkdirReq types.MkdirReq
		if err := c.ShouldBindJSON(&MkdirReq); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
		l := service.GetFileSrv()
		//这里的返回值不规范，后续还会修改
		resp, err := l.CreateDir(c.Request.Context(), MkdirReq.DirName, MkdirReq.ParentId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func DeleteFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var DeleteFileReq types.DeleteFileReq
		if err := c.ShouldBindJSON(&DeleteFileReq); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
		l := service.GetFileSrv()
		resp, err := l.DeleteFile(c.Request.Context(), DeleteFileReq.FileName, DeleteFileReq.ParentId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetFileRootHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := service.GetFileSrv()
		resp, err := l.GetRoot(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func RenameFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var RenameFileReq types.RenameFileReq
		if err := c.ShouldBindJSON(&RenameFileReq); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
		l := service.GetFileSrv()
		resp, err := l.RenameFile(c.Request.Context(), RenameFileReq.NewFileName, RenameFileReq.FileId, RenameFileReq.ParentId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		}
		c.JSON(http.StatusOK, resp)
	}
}

func DownloadFileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var DownloadFileReq types.DownloadFileReq
		if err := c.ShouldBindJSON(&DownloadFileReq); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		l := service.GetFileSrv()
		downloadPath, err := l.DownloadFile(c.Request.Context(), DownloadFileReq.FileName, DownloadFileReq.ParentId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		fmt.Println("Got Download Path: ", downloadPath)
		if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"Status": 404,
				"Msg":    "File not found",
			})
			return
		}
		encodedFileName := url.QueryEscape(DownloadFileReq.FileName)
		c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFileName)
		c.Header("Content-Type", "application/octet-stream")
		c.File(downloadPath)
	}
}
