package model

import (
	"time"
)

type FileModel struct {
	Id        int64      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	UUID      string     `gorm:"column:uuid;size:36;unique"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	Sha1      string     `gorm:"column:sha1"`
	Filename  string     `gorm:"column:filename"`
	Filesize  int64      `gorm:"column:file_size"`
	Fileaddr  string     `gorm:"column:file_addr"`
	Status    int        `gorm:"column:status"`
	Ext       string     `gorm:"column:ext"`
	Ref       int        `gorm:"column:ref"`
}

func (*FileModel) TableName() string {
	return "file"
}
