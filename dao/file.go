package dao

import (
	"BlueNetDisk/consts"
	"BlueNetDisk/model"
	"BlueNetDisk/pkg/utils"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"path"
)

type FileDao struct {
	db *gorm.DB
}

func NewFileDao(c context.Context) *FileDao {
	if c == nil {
		c = context.Background()
	}
	return &FileDao{NewDbClient(c)}
}

// CreateFile 将传入文件头解析成文件，写信息到文件池表，用户文件树表，如果flag为True则只写用户文件树表
func (f *FileDao) CreateFile(fileHeader *multipart.FileHeader, uid int64, parentId string, shaStr string, OnlyUpDateFileTree bool) (UUID string, err error) {
	fUUID := uuid.New().String()
	fileInfo := &model.FileModel{
		UUID:     fUUID,
		Sha1:     shaStr,
		Filename: fileHeader.Filename,
		Filesize: fileHeader.Size,
		Fileaddr: consts.FilePoolPath,
		Status:   consts.Available,
		Ext:      path.Ext(fileHeader.Filename),
		Ref:      1,
	}
	nodeInfo := &model.FileTreeModel{
		Uid:        uid,
		UUID:       fUUID,
		ParentUUID: parentId,
		FileName:   fileHeader.Filename,
		Filesize:   fileHeader.Size,
		Status:     consts.Available,
		Ext:        path.Ext(fileHeader.Filename),
	}
	//原子事务操作
	//写入文件信息到文件池表
	begin := f.db.Begin()
	//如果是正在创建软连接,则不写信息到文件池表
	if !OnlyUpDateFileTree {
		err := begin.Model(&model.FileModel{}).Create(fileInfo).Error
		if err != nil {
			begin.Rollback()
			return "", err
		}
	}
	//写入文件信息到用户文件树表
	err = begin.Model(&model.FileTreeModel{}).Create(nodeInfo).Error
	if err != nil {
		begin.Rollback()
		return "", err
	}
	//如果是文件夹或正在创建软连接,则不写文件到文件池
	if fileInfo.Ext != ".dir" || OnlyUpDateFileTree {
		err = utils.Writer(fileHeader, consts.FilePoolPath, fileInfo.UUID)
		if err != nil {
			begin.Rollback()
			return "", err
		}
	}
	begin.Commit()
	return fUUID, err
}

// FindFileBySha1 在文件池中按sha1值查询某一文件
func (f *FileDao) FindFileBySha1(shaStr string) (fileInfo *model.FileModel, err error) {
	err = f.db.Model(&model.FileModel{}).Where("sha1 = ?", shaStr).First(&fileInfo).Error
	return
}

// FindFileByNameAndParent 在用户文件表中查询某一级的文件
func (f *FileDao) FindFileByNameAndParent(uid int64, name string, parentId string) (fileInfo *model.FileTreeModel, err error) {
	err = f.db.Model(&model.FileTreeModel{}).Where("uid = ? AND file_name = ? AND parent_uuid = ?", uid, name, parentId).First(&fileInfo).Error
	return
}

func (f *FileDao) FindFileByUuid(uid int64, uuid string) (fileInfo *model.FileTreeModel, err error) {
	err = f.db.Model(&model.FileTreeModel{}).Where("uuid = ? AND uid = ?", uuid, uid).First(&fileInfo).Error
	return
}

func (f *FileDao) ListFileByParent(uid int64, parentId string) (fileInfos []*model.FileTreeModel, total int64, err error) {
	err = f.db.Model(&model.FileTreeModel{}).Where("uid = ? AND parent_uuid = ?", uid, parentId).Find(&fileInfos).Count(&total).Error
	return
}

func (f *FileDao) UpdateFileRef(uuid string) error {
	file := &model.FileModel{}
	return f.db.Model(&model.FileModel{}).Where("uuid = ?", uuid).First(file).Update("ref", file.Ref+1).Error
}
