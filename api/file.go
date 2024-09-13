package api

import (
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/service"
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
			resp, err := l.UploadManyFiles(form)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			c.JSON(http.StatusOK, resp)
		}
	}
}
