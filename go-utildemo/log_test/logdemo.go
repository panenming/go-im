package main

import (
	"errors"

	"github.com/panenming/go-im/libs/logger"
)

func main() {
	//初始化
	logger.Initialize("./log", "LoginServer")

	//设置选项
	logger.SetConsole(true)
	logger.SetLevel(logger.INFO)

	//单一输出
	logger.Debug("I'm debug log!")
	logger.Info("I'm info log!")
	logger.Warn("I'm warn log!")
	logger.Error("I'm error log!")

	//格式化输出
	logger.Debugf("I'm %s log! ", "debug")
	logger.Infof("I'm %s log!", "info")
	logger.Warnf("I'm %s log!", "warn")
	logger.Errorf("I'm %s log!", "error")

	//行输出
	logger.Debugln("I'm", "debug", "log!")
	logger.Infoln("I'm", "info", "log!")
	logger.Warnln("I'm", "warn", "log!")
	logger.Errorln("I'm", "error", "log!")

	//异常捕获
	err := errors.New("程序崩溃！")
	defer logger.CatchException()
	panic(err) //此panic会被logger.CatchException()捕获，并保存到exception目录
}
