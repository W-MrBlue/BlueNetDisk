package dao

import (
	"BlueNetDisk/model"
	"context"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(c context.Context) *UserDao {
	if c == nil {
		c = context.Background()
	}
	return &UserDao{NewDbClient(c)}
}

func (u *UserDao) CreateUser(user *model.UserModel) error {
	return u.db.Create(user).Error
}

func (u *UserDao) FindUserByName(username string) (user *model.UserModel, err error) {
	err = u.db.Model(&model.UserModel{}).Where("username=?", username).
		First(&user).Error
	return
}
