package main

import (
	"BlueNetDisk/config"
	"BlueNetDisk/consts"
	"BlueNetDisk/dao"
	"BlueNetDisk/pkg/utils"
	"BlueNetDisk/router"
)

func main() {
	config.InitConfig()
	dao.MysqlInit()
	utils.InitLog()
	consts.PathInit()
	r := router.NewRouter()
	_ = r.Run(config.Config.System.HttpPort)
}
