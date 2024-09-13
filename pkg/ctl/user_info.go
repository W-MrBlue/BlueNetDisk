package ctl

import (
	"context"
	"errors"
)

var userKey = "userKey"

type UserInfo struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

func GetUserInfo(c context.Context) (*UserInfo, error) {
	user, ok := FromContext(c)
	if !ok {
		return nil, errors.New("no user info found")
	}
	return user, nil
}

func NewContext(c context.Context, u *UserInfo) context.Context {
	return context.WithValue(c, userKey, u)
}

func FromContext(c context.Context) (*UserInfo, bool) {
	u, ok := c.Value(userKey).(*UserInfo)
	return u, ok
}
