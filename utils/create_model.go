//生成model工具包
package utils

import (
	"fmt"
	"github.com/Gouplook/dzgin"
	"github.com/Gouplook/dzgin/logs"
	"github.com/Gouplook/dzgin/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"time"
)

func CreateModel(table string) {
	prefix := dzgin.AppConfig.String("db.prefix")
	ormers := orm.NewOrm()
	sql := fmt.Sprintf("select COLUMN_NAME from information_schema.COLUMNS where table_name ='%s'", prefix+table)
	var maps []orm.Params
	_, err := ormers.Raw(sql).Values(&maps)
	if err != nil {
		logs.Error("Sql:", sql, " Error,", err.Error())
	}
	fields := make([]string, 0)
	if len(maps) > 0 {
		for _, v := range maps {
			fields = append(fields, v["COLUMN_NAME"].(string))
		}
	} else {
		logs.Error("this table has not Field or table not exists.")
		os.Exit(1)
	}
	var tableF = table
	var tableA []string
	var rr string
	if strings.Contains(table, "_") {
		tableA = strings.Split(table, "_")
		for _, v := range tableA {
			rr += Capitalize(v)
		}
	} else {
		rr = Capitalize(table)
	}
	fileName := rr + "Model"
	f, err := os.OpenFile(dzgin.WorkPath+"/common/models/"+fileName+".go", os.O_CREATE, 0777)
	if err != nil {
		logs.Error("OpenFile Error:", err.Error())
		os.Exit(1)
	}
	defer f.Close()
	//头部及引入
	strContent := "//" + fileName + "\n" + "//" + time.Now().Format("2006-01-02 15:04:05") + "\n\n"
	strContent += `package models`
	strContent += "\n\n"
	strContent +=
		`import(
    "github.com/g/base/common/functions"
    "git.900sui.cn/kc/base/common/models/base"
    "git.900sui.cn/kc/dzgin/orm"
	)`

	//结构体
	strContent += "\n\n"
	strContent += "//表结构体\n"
	strContent += "type " + fileName + " struct {\n"
	strContent += "\tModel *base.Model\n"
	strContent += "\tField " + fileName + "Field\n"
	strContent += "}"
	strContent += "\n\n"
	//字段结构体
	strContent += "//表字段\n"
	strContent += "type " + fileName + "Field struct{\n"
	strContent += fmt.Sprintf("\tT_table\tstring\t`default:\"%s\"`\n", tableF)
	for _, v := range fields {
		strContent += fmt.Sprintf("\tF_%s\tstring\t`default:\"%s\"`\n", v, v)
	}
	strContent += "}"
	//Init
	strContent += "\n\n"
	tableIndexS := SubString(tableF, 0, 1)
	strContent += "//初始化\n"
	strContent += "func (" + tableIndexS + " *" + fileName + ") Init(ormer ...orm.Ormer) *" + fileName + "{\n"
	strContent += "\tfunctions.ReflectModel(&" + tableIndexS + ".Field)\n"
	strContent += "\t" + tableIndexS + ".Model = base.NewModel(" + tableIndexS + ".Field.T_table, ormer...)\n"
	strContent += "\treturn " + tableIndexS + "\n"
	strContent += "}"
	strContent += "\n\n"
	//Insert
	strContent += "//新增数据\n"
	strContent += fmt.Sprintf("func (%s *%s) Insert(data map[string]interface{}) (int){\n", tableIndexS, fileName)
	strContent += fmt.Sprintf("\tresult,_ := %s.Model.Data(data).Insert()\n", tableIndexS)
	strContent += "\treturn result\n"
	strContent += "}"

	n, err := f.WriteString(strContent)
	if err != nil {
		fmt.Println("WriteString Error,Error is", err.Error())
		return
	}
	logs.Info("create File:", fileName+".go success.", n)
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

//获取字符串首字母
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)
	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	// 返回子串
	return string(rs[begin:end])
}
