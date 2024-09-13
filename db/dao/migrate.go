package dao

import (
	"BlueNetDisk/model"
)

func migrate() {
	err := _db.Set("gorm:table_settings", "charset=utf8&parseTime=True&loc=Local").AutoMigrate(&model.UserModel{})
	if err != nil {
		panic(err)
	}
	return
}
