package util

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func InitLogs() {
	// 获取日志文件路径
	logPath := beego.AppConfig.String("logsfilename")
	if logPath == "" {
		logPath = "logs/app.log"
	}

	// 设置日志配置
	logConf := fmt.Sprintf(`{"filename":"%s"}`, logPath)
	err := logs.SetLogger(logs.AdapterFile, logConf)
	if err != nil {
		fmt.Sprintf("faild to init logs:%s", err)
		return
	}
	// 启用函数调用深度
	logs.EnableFuncCallDepth(true)
	logs.Info("Application started")
}
