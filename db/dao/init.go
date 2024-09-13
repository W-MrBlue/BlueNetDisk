package dao

import (
	"BlueNetDisk/config"
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

var _db *gorm.DB

func MysqlInit() {
	conf := config.Config.MySql["default"]
	dsn := strings.Join([]string{conf.UserName, ":", conf.Password, "@tcp(", conf.DbHost, ":", conf.DbPort, ")/",
		conf.DbName, "?charset=", conf.Charset, "&parseTime=True&loc=Local"}, "")

	var ormLogger = logger.Default
	if gin.Mode() == "Debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	_db = db

	migrate()
}

func NewDbClient(c context.Context) *gorm.DB {
	db := _db
	return db.WithContext(c)
}
