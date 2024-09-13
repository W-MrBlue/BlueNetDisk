package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logrusobj *logrus.Logger

func InitLog() {
	if Logrusobj != nil {
		src, _ := setOutPutFile()
		Logrusobj.Out = src
		return
	}
	logger := logrus.New()
	src, _ := setOutPutFile()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Logrusobj = logger
}

func setOutPutFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/log/"
		fmt.Println(err)
	}

	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(logFilePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	logFileName := now.Format("2006-01-02") + ".log"
	filePath := logFilePath + logFileName
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
	}
	src, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend|os.ModePerm)
	if err != nil {
		return nil, err
	}
	return src, nil
}
