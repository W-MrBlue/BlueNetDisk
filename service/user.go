package service

import (
	"BlueNetDisk/db/dao"
	"BlueNetDisk/model"
	"BlueNetDisk/pkg/ctl"
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/types"
	"context"
	"errors"
	"gorm.io/gorm"
	"sync"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct{}

// GetUserSrv returns an instance of UserService
func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (*UserSrv) UserRegister(c context.Context, req *types.UserRegisterReq) (resp interface{}, err error) {
	u := dao.NewUserDao(c)
	user, err := u.FindUserByName(req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &model.UserModel{
			Username: req.Username,
		}
		if err = user.SetPassword(req.Password); err != nil {
			utils.Logrusobj.Error(err)
			return
		}
		if err = u.CreateUser(user); err != nil {
			utils.Logrusobj.Error(err)
			return nil, err
		}
		return ctl.RespSuccess(), nil
	}
	err = errors.New("user already exists")
	utils.Logrusobj.Error(err)
	return nil, err
}

func (*UserSrv) UserLogin(c context.Context, req *types.UserLoginReq) (resp interface{}, err error) {
	u := dao.NewUserDao(c)
	user, err := u.FindUserByName(req.Username)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}
	if ok := user.CheckPassword(req.Password); !ok {
		err = errors.New("invalid password")
		utils.Logrusobj.Error(err)
		return nil, err
	}

	//todo :jwt delivery

	return ctl.RespSuccess(), nil
}