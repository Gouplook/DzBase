package utils

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"strings"
)

// 输入日志到文件

var Logerr *log.Logger

type DzLog struct {
}

//写日志
func (k *DzLog) PrintLog(logInfo string) {
	Logerr.Println(logInfo)
}

func (k *DzLog) GetInstance(filePath string, logLevel ...string) *DzLog {
	logL := "Info"
	if len(logLevel) > 0 {
		logL = logLevel[0]
	}
	file, _ := rotatelogs.New(filePath)
	Logerr = log.New(io.MultiWriter(file, os.Stderr),
		strings.ToUpper(logL)+"：",
		log.Ldate|log.Ltime|log.Lshortfile)

	return k
}
