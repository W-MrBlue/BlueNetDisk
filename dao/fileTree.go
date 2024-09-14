package dao

import (
	"BlueNetDisk/model"
	"gorm.io/gorm"
)

type FileTreeDao struct {
	db *gorm.DB
}

func (f *FileTreeDao) CreateFileNode(node *model.FileTreeModel) error {
	return f.db.Create(node).Error
}
