package dao

import (
	"BlueNetDisk/consts"
	"BlueNetDisk/model"
	"BlueNetDisk/pkg/utils"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
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

// CreateFile 将传入文件头解析成文件，写信息到文件池表，用户文件树表.
// 如果flag为True则正在创建软连接，只写用户文件树表，需要传入链接目标的UUID
func (f *FileDao) CreateFile(fileHeader *multipart.FileHeader, uid int64, parentDir *model.FileTreeModel, shaStr string, OnlyUpDateFileTree bool, linkUUID string) (UUID string, err error) {
	fUUID := uuid.New().String()
	fileInfo := model.NewFile(fileHeader, fUUID, shaStr)
	nodeInfo := model.NewFileNode(fileHeader, uid, parentDir, fUUID)
	//原子事务操作
	//写入文件信息到文件池表
	// 开启事务
	tx := f.db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// 插入文件树信息
	if OnlyUpDateFileTree {
		nodeInfo.UUID = linkUUID
	}
	if err := tx.Model(&model.FileTreeModel{}).Create(nodeInfo).Error; err != nil {
		return "", err
	}

	// 如果不是仅更新文件树，则插入文件池信息以及写文件池存储
	if !OnlyUpDateFileTree {
		if err := tx.Model(&model.FileModel{}).Create(fileInfo).Error; err != nil {
			return "", err
		}

		if fileInfo.Ext != ".dir" {
			if err := utils.WriteFile(fileHeader, consts.FilePoolPath, fileInfo.UUID); err != nil {
				return "", err
			}
		}
	}
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
	// 先获取总数
	err = f.db.Model(&model.FileTreeModel{}).
		Where("uid = ? AND parent_uuid = ?", uid, parentId).
		Count(&total).Error
	if err != nil {
		return
	}

	// 再查询实际的数据
	err = f.db.Model(&model.FileTreeModel{}).
		Where("uid = ? AND parent_uuid = ?", uid, parentId).
		Find(&fileInfos).Error
	return
}

func (f *FileDao) IncreaseFileRef(uuid string) error {
	file := &model.FileModel{}
	return f.db.Model(&model.FileModel{}).Where("uuid = ?", uuid).First(file).Update("ref", file.Ref+1).Error
}

func deleteFileWithRef(tx *gorm.DB, uuid string) error {

	file := &model.FileModel{}
	err := tx.Model(&model.FileModel{}).Where("uuid = ?", uuid).First(&file).Error

	if err != nil {
		return err
	}

	// 根据 ref 值决定是更新还是删除
	if file.Ref > 1 {
		// 如果 ref > 1，执行 ref - 1 的更新操作
		err = tx.Model(&model.FileModel{}).Where("uuid = ?", uuid).Update("ref", file.Ref-1).Error
		if err != nil {
			return err
		}
	} else {
		// 如果 ref <= 1，执行删除操作
		err = tx.Model(&model.FileModel{}).Where("uuid = ?", uuid).Delete(&model.FileModel{}).Error
		if err != nil {
			return err
		}
		err = utils.RemoveFile(uuid + file.Ext)
	}

	return nil
}

// DeleteFileTreeNodeWithTransaction 使用事务递归删除节点及其所有子节点
func (f *FileDao) DeleteFileTreeNodeWithTransaction(uid int64, uuid string) error {
	// 开始事务
	tx := f.db.Begin()

	// 递归删除节点及其所有子节点
	err := f.deleteFileTreeNodeWithTx(tx, uid, uuid)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// deleteFileTreeNodeWithTx 递归删除节点及其所有子节点，带事务处理
func (f *FileDao) deleteFileTreeNodeWithTx(tx *gorm.DB, uid int64, uuid string) error {
	// 查找所有子节点
	var children []*model.FileTreeModel
	err := tx.Where("parent_uuid = ? AND uid = ?", uuid, uid).Find(&children).Error
	if err != nil {
		return err
	}

	// 递归删除每个子节点
	for _, child := range children {
		err = f.deleteFileTreeNodeWithTx(tx, uid, child.UUID)
		if err != nil {
			return err
		}
	}

	// 删除当前节点,并修改其对应物理文件的REF
	err = deleteFileWithRef(tx, uuid)
	if err != nil {
		return err
	}
	err = tx.Where("uuid = ? AND uid = ?", uuid, uid).Delete(&model.FileTreeModel{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetRoot 返回用户的根目录UUID
func (f *FileDao) GetRoot(uid int64) (UUID string, err error) {
	fileInfo := &model.FileTreeModel{}
	err = f.db.Model(&model.FileTreeModel{}).Where("uid = ? AND parent_uuid = 0", uid).First(fileInfo).Error
	if err != nil {
		return "", err
	}
	return fileInfo.UUID, nil
}

func (f *FileDao) RenameById(uid int64, uuid, parentId, newName string) error {
	return f.db.Model(&model.FileTreeModel{}).Where("uid = ? AND uuid = ? AND parent_uuid = ?", uid, uuid, parentId).Update("file_name", newName).Error
}
