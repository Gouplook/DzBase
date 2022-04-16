package components

import (
	"encoding/json"
	"fmt"
	"github.com/Gouplook/dzgin"
	"github.com/Gouplook/dzgin/logs"
	"time"
)

func InitLogger() (err error) {
	//打印环境变量
	logs.Info("Environment Variable:MSF_ENV:", dzgin.KcConfig.RunMode)
	//输出文件名和行号
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	//开启异步缓存
	logs.Async(1e3)

	fileName := time.Now().Format("20060102") + ".log"
	logConf := make(map[string]interface{})
	logConf = map[string]interface{}{
		"filename": "/opt/logs/" + dzgin.AppConfig.String("jaeger.serviceName") + "/" + fileName,
		"maxdays":  1,
	}

	confStr, err := json.Marshal(logConf)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(confStr))
	return
}
