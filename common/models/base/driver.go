//配置基础数据库引擎
//配置基础数据库参数
package base

import (
	"fmt"
	"github.com/Gouplook/dzgin"
	"github.com/Gouplook/dzgin/logs"
	"github.com/Gouplook/dzgin/orm"

	//"git.900sui.cn/kc/dzgin"
	//"git.900sui.cn/kc/dzgin/logs"
	//"git.900sui.cn/kc/dzgin/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var (
	err error
)

//初始化驱动
func init() {
	logs.Info("Init driver.go mysql start")
	//设置驱动数据库连接参数
	dataSource := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s", dzgin.AppConfig.String("db.user"), dzgin.AppConfig.String("db.pwd"),
		dzgin.AppConfig.String("db.host"), dzgin.AppConfig.String("db.port"), dzgin.AppConfig.String("db.name"), dzgin.AppConfig.String("db.charset"))
	//打印连接数据库参数
	logs.Info("DatabaseDriverConnect String:", dataSource)
	maxIdle, _ := strconv.Atoi(dzgin.AppConfig.DefaultString("db.maxidle", "10"))
	maxConn, _ := strconv.Atoi(dzgin.AppConfig.DefaultString("db.maxconn", "0"))
	maxTime := dzgin.AppConfig.DefaultInt("db.maxlifetime", 10800)
	logs.Info("connMaxLifeTime(s)：", maxTime)
	//设置注册数据库
	if err == nil {
		err = orm.RegisterDataBase("default", dzgin.AppConfig.String("db.type"), dataSource, maxIdle, maxConn, maxTime)
	}
}
