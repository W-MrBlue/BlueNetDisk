package dao

import (
	"BlueNetDisk/model"
)

func migrate() {
	err := _db.Set("gorm:table_settings", "charset=utf8mb4&parseTime=True&loc=Local").AutoMigrate(&model.UserModel{}, &model.FileModel{}, &model.FileTreeModel{})
	if err != nil {
		panic(err)
	}
	return
}
