package model

import (
	"BlueNetDisk/consts"
	"mime/multipart"
	"path"
	"strings"
	"time"
)

type FileTreeModel struct {
	Id         int64      `gorm:"primary_key;AUTO_INCREMENT"`
	Uid        int64      `gorm:"column:uid"`
	UUID       string     `gorm:"column:uuid"`
	FileName   string     `gorm:"column:file_name"`
	Filesize   int64      `gorm:"column:file_size"`
	FileAddr   string     `gorm:"column:file_addr"`
	ParentUUID string     `gorm:"column:parent_uuid"`
	CreatedAt  *time.Time `gorm:"column:created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
	Status     int        `gorm:"column:status"`
	Ext        string     `gorm:"column:ext"`
}

func (*FileTreeModel) TableName() string {
	return "file_tree"
}

func NewFileNode(fileHeader *multipart.FileHeader, uid int64, parentDir *FileTreeModel, fUUID string) *FileTreeModel {
	return &FileTreeModel{
		UUID:       fUUID,
		Uid:        uid,
		ParentUUID: parentDir.UUID,
		FileName:   fileHeader.Filename,
		Filesize:   fileHeader.Size,
		FileAddr:   path.Join(parentDir.FileAddr, strings.TrimSuffix(parentDir.FileName, ".dir")),
		Status:     consts.Available,
		Ext:        path.Ext(fileHeader.Filename),
	}
}
