package service

import (
	"BlueNetDisk/pkg/ctl"
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/types"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"sync"
)

var FileSrvIns *FileSrv
var FileSrvOnce sync.Once

type FileSrv struct{}

// GetFileSrv returns an instance of FileService
func GetFileSrv() *FileSrv {
	FileSrvOnce.Do(func() {
		FileSrvIns = &FileSrv{}
	})
	return FileSrvIns
}

// UploadManyFiles calls UploadFile for each file.
// Resp includes how many uploads have failed and specified failure information
func (f *FileSrv) UploadManyFiles(form *multipart.Form) (resp interface{}, err error) {
	var total int64 = 0
	if form.File["files"] == nil {
		err := errors.New("no file given")
		return total, err
	}

	files := form.File["files"]
	dirPath := form.Value["filepath"][0]
	respList := make([]*types.UploadFileResp, 0)
	for _, file := range files {
		err = f.UploadFile(file, dirPath)
		if err != nil {
			respList = append(respList, &types.UploadFileResp{
				FileName: file.Filename,
				FileSize: file.Size,
				Info:     err.Error(),
			})
			total += 1
		}
	}
	return ctl.RespList(respList, total), nil
}

// UploadFile uploads file to the given path,return err if any error occurs
func (*FileSrv) UploadFile(file *multipart.FileHeader, dirPath string) error {
	src, err := file.Open()
	if err != nil {
		utils.Logrusobj.Error(err)
		return err
	}
	filePath := dirPath + "/" + file.Filename
	defer src.Close()
	out, err := os.Create(filePath)
	if err != nil {
		utils.Logrusobj.Error(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return nil
}
