package model

import (
	"time"
)

type FileModel struct {
	Id        int64      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Uid       int64      `gorm:"column:uid"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	Sha1      string     `gorm:"column:sha1;unique"`
	Filename  string     `gorm:"column:filename"`
	Filesize  int64      `gorm:"column:file_size"`
	Fileaddr  string     `gorm:"column:file_addr"`
	status    int        `gorm:"column:status"`
}

func (*FileModel) TableName() string {
	return "file"
}
