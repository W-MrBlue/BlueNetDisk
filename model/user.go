package model

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModel struct {
	Id             int64      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	CreatedAt      *time.Time `gorm:"column:created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
	Username       string     `gorm:"column:username;unique"`
	PasswordDigest string     `gorm:"column:password_digest"`
	RootUuid       string     `gorm:"column:root_uuid"`
}

func (*UserModel) TableName() string {
	return "user"
}
func (u *UserModel) SetRootDir(uuid string) {
	u.RootUuid = uuid
}

func (u *UserModel) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

func (u *UserModel) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}
