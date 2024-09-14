package service

import (
	"BlueNetDisk/consts"
	"BlueNetDisk/dao"
	"BlueNetDisk/pkg/ctl"
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/types"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
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
func (f *FileSrv) UploadManyFiles(c context.Context, form *multipart.Form) (resp interface{}, err error) {
	var total int64 = 0
	if form.File["files"] == nil {
		err := errors.New("no file given")
		return total, err
	}

	fileHeaders := form.File["files"]
	parentId := form.Value["parentDir"][0]
	respList := make([]*types.UploadFileResp, 0)
	for _, fileHeader := range fileHeaders {
		fUUID, err := f.UploadFile(c, fileHeader, parentId)
		if err != nil {
			respList = append(respList, &types.UploadFileResp{
				FileName: fileHeader.Filename,
				FileSize: fileHeader.Size,
				Info:     err.Error(),
			})
			total += 1
		} else {
			respList = append(respList, &types.UploadFileResp{
				FileName: fileHeader.Filename,
				FileSize: fileHeader.Size,
				Info:     fUUID,
			})
		}
	}
	return ctl.RespList(respList, total), nil
}

// UploadFile uploads file to the given path,return err if any error occurs
func (*FileSrv) UploadFile(c context.Context, fileHeader *multipart.FileHeader, parentId string) (UUID string, err error) {

	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return "", err
	}

	//同步文件信息到数据库
	f := dao.NewFileDao(c)
	//查询目标目录是否存在于用户文件表中
	_, err = f.FindFileByUuid(u.Id, parentId)
	if err != nil {
		err = errors.New("failed to get the parent directory,it may not exit")
		return "", err
	}
	//查询同级目录下是否存在同名文件
	_, err = f.FindFileByNameAndParent(u.Id, fileHeader.Filename, parentId)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		err := errors.New("file already exists")
		utils.Logrusobj.Error(err)
		return "", err
	}

	//查询文件池中是否存在相同文件,如果存在则创建软连接，仅更新用户文件表
	shaStr := utils.Sha1(fileHeader)
	if shaStr == "" {
		err := errors.New("failed to generate sha1")
		utils.Logrusobj.Error(err)
		return "", err
	}
	fileInfo, err := f.FindFileBySha1(shaStr)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		_ = f.UpdateFileRef(fileInfo.UUID)
		_, err = f.CreateFile(fileHeader, u.Id, parentId, shaStr, true)
		if err != nil {
			utils.Logrusobj.Error(err)
			return "", err
		}
		return fileInfo.UUID, nil
	}

	fUUID, err := f.CreateFile(fileHeader, u.Id, parentId, shaStr, false)
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}
	return fUUID, err
}

// CreateDir 将在用户文件表和文件池下建立文件夹类型文件，但文件夹不会被实际写入文件池
func (*FileSrv) CreateDir(c context.Context, dirname string, parentId string) (UUID string, err error) {
	u, err := ctl.GetUserInfo(c)
	f := dao.NewFileDao(c)
	fUUID := uuid.New().String()
	if err != nil {
		return
	}
	//查找同级同名文件
	_, err = f.FindFileByNameAndParent(u.Id, dirname, parentId)
	fmt.Println(err.Error())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		err := errors.New("dir already exists")
		utils.Logrusobj.Error(err)
		return "", err
	}
	//创建文件夹
	fakeFileHeader := multipart.FileHeader{
		Filename: dirname + ".dir",
		Size:     consts.DirSize,
	}
	_, err = f.CreateFile(&fakeFileHeader, u.Id, parentId, "", false)
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}
	return fUUID, err
}

func (*FileSrv) ListFileByParentUUID(c context.Context, parentId string) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	fmt.Println("Going to select where uid = ", u.Id, "AND parentId = ", parentId)
	files, total, err := dao.NewFileDao(c).ListFileByParent(u.Id, parentId)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	fRespList := make([]*types.ListFileResp, 0)
	for _, file := range files {
		fRespList = append(fRespList, &types.ListFileResp{
			UUID:      file.UUID,
			FileName:  file.FileName,
			FileSize:  file.Filesize,
			Ext:       file.Ext,
			CreatedAt: file.CreatedAt,
			UpdatedAt: file.UpdatedAt,
		})
	}
	return ctl.RespList(fRespList, total), nil
}
