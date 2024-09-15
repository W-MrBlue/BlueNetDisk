package service

import (
	"BlueNetDisk/consts"
	"BlueNetDisk/dao"
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

		//为用户建立root文件夹
		//此处需要查询用户信息写入db后自动产生的id，以保存到context中进行跟踪
		user, err = u.FindUserByName(user.Username)
		if err != nil {
			utils.Logrusobj.Error(err)
			return nil, err
		}
		f := GetFileSrv()
		c = context.WithValue(c, consts.UserKey, &ctl.UserInfo{Id: user.Id, Username: user.Username})
		fUUID, err := f.CreateDir(c, user.Username, "0")
		if err != nil {
			utils.Logrusobj.Error(err)
			return nil, err
		}

		if fUUID, ok := fUUID.(string); ok {
			user.SetRootDir(fUUID)
		} else {
			err = errors.New("wrong return type")
			utils.Logrusobj.Error(err)
			return nil, err
		}

		return ctl.RespSuccessWithData(fUUID), nil
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

	token, err := utils.GenerateToken(user.Id, user.Username)
	if err != nil {
		utils.Logrusobj.Error(err)
		return nil, err
	}

	userResp := types.TokenData{
		UserInfo: types.UserInfoResp{
			Id:       user.Id,
			Username: user.Username,
		},
		AccessToken: token,
	}

	return ctl.RespSuccessWithData(userResp), nil
}
