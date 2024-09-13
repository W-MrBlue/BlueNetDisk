package types

type UserRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfoResp struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenData struct {
	UserInfo    UserInfoResp `json:"user_info"`
	AccessToken string       `json:"access_token"`
}
