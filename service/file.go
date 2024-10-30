package service

import (
	"BlueNetDisk/consts"
	"BlueNetDisk/dao"
	"BlueNetDisk/model"
	"BlueNetDisk/pkg/ctl"
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/types"
	"context"
	"errors"
	"gorm.io/gorm"
	"mime/multipart"
	"path"
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
		fUUID, err := f.uploadFile(c, fileHeader, parentId)
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

// uploadFile uploads file to the given path,return err if any error occurs
// todo 数据库与存储写操作统一在DAO层进行，应该分离为DAO层和UTIL层，二者在SERVICE层合并而非在Dao层合并
func (*FileSrv) uploadFile(c context.Context, fileHeader *multipart.FileHeader, parentId string) (UUID string, err error) {

	u, err := ctl.GetUserInfo(c)
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}

	//同步文件信息到数据库
	f := dao.NewFileDao(c)
	//查询目标目录是否存在于用户文件表中
	parentDir, err := f.FindFileByUuid(u.Id, parentId)
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
	//创建软连接并更新用户文件表
	fileInfo, err := f.FindFileBySha1(shaStr)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		_ = f.IncreaseFileRef(fileInfo.UUID)
		_, err = f.CreateFile(fileHeader, u.Id, parentDir, shaStr, true, fileInfo.UUID)
		if err != nil {
			utils.Logrusobj.Error(err)
			return "", err
		}
		return fileInfo.UUID, nil
	}

	//写入文件池并更新用户文件表
	fUUID, err := f.CreateFile(fileHeader, u.Id, parentDir, shaStr, false, "")
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}

	return fUUID, err
}

// CreateDir 将在用户文件表和文件池下建立文件夹类型文件，但文件夹不会被实际写入文件池
func (*FileSrv) CreateDir(c context.Context, dirname string, parentUUID string) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return "", err
	}
	f := dao.NewFileDao(c)
	parentDir := &model.FileTreeModel{}
	//查询目标目录是否存在于用户文件表中
	if parentUUID == "0" {
		//创建假根文件
		parentDir = &model.FileTreeModel{
			FileName: "",
			FileAddr: "",
		}
	} else {
		parentDir, err = f.FindFileByUuid(u.Id, parentUUID)
		if err != nil {
			err = errors.New("failed to get the parent directory,it may not exit")
			return "", err
		}
	}

	//查找同级同名文件
	_, err = f.FindFileByNameAndParent(u.Id, dirname+".dir", parentUUID)
	//fmt.Println(err.Error())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("dir already exists")
		utils.Logrusobj.Error(err)
		return "", err
	}
	//创建文件夹
	fakeFileHeader := multipart.FileHeader{
		Filename: dirname + ".dir",
		Size:     consts.DirSize,
	}
	fUUID, err := f.CreateFile(&fakeFileHeader, u.Id, parentDir, "", false, "")
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}
	//如果是创建用户根目录,为service调用,返回fUUID
	//如果是创建其他目录，为api调用，返回resp
	if parentUUID == "0" {
		return fUUID, err
	} else {
		return ctl.RespSuccessWithData(fUUID), nil
	}
}

func (*FileSrv) GetRoot(c context.Context) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	f := dao.NewFileDao(c)
	root, err := f.GetRoot(u.Id)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	return ctl.RespSuccessWithData(root), nil
}

func (*FileSrv) ListFileByParentUUID(c context.Context, parentId string) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	f := dao.NewFileDao(c)

	_, err = f.FindFileByUuid(u.Id, parentId)
	if err != nil {
		err = errors.New("failed to get the directory,it may not exit")
		return "", err
	}

	files, total, err := f.ListFileByParent(u.Id, parentId)
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

func (*FileSrv) DeleteFile(c context.Context, filename string, parentId string) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	f := dao.NewFileDao(c)
	file, err := f.FindFileByNameAndParent(u.Id, filename, parentId)
	if err != nil {
		err = errors.New("failed to get the file,it may not exit")
		utils.Logrusobj.Error(err)
		return nil, err
	}
	err = f.DeleteFileTreeNodeWithTransaction(u.Id, file.UUID)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	return ctl.RespSuccess(), nil
}

func (*FileSrv) RenameFile(c context.Context, newFilename string, fileId, parentId string) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return "", err
	}
	f := dao.NewFileDao(c)
	//查找同级同名文件
	_, err = f.FindFileByNameAndParent(u.Id, newFilename, parentId)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("name already exists")
		utils.Logrusobj.Error(err)
		return "", err
	}

	err = f.RenameById(u.Id, fileId, parentId, newFilename)
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}
	return ctl.RespSuccess(), nil
}

func (*FileSrv) DownloadFile(c context.Context, filename string, parentId string) (downloadPath string, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}
	f := dao.NewFileDao(c)
	file, err := f.FindFileByNameAndParent(u.Id, filename, parentId)
	if err != nil {
		utils.Logrusobj.Error(err)
		return "", err
	}
	downloadPath = path.Join(consts.FilePoolPath, file.UUID+file.Ext)
	return downloadPath, nil
}
