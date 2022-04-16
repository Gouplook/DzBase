package base

import (
	"errors"
	"fmt"
	"github.com/Gouplook/dzbase/common/functions"
	"github.com/Gouplook/dzgin"
	"github.com/Gouplook/dzgin/logs"
	"github.com/Gouplook/dzgin/orm"
	"reflect"
	"strconv"
	"strings"
)

//定义变量
var (
	T_PREFIX = dzgin.AppConfig.String("db.prefix")
)

const (
	MYSQL_CONDITION_OR = "OR"
)

//定义Model结构体
type Model struct {
	table   string
	o       orm.Ormer
	limit   []interface{}
	orderBy []string
	where   []WhereItem
	data    map[string]interface{}
	field   string
	sql     string
	lock    string
}

// where条件结构体
type WhereItem struct {
	Field string
	Value interface{}
}

//操作另一张表,表名不需要扩展
func (m *Model) Table(table string) *Model {
	m.table = T_PREFIX + table
	return m
}

//条件查询
// map[string]interface{"id":1,"name":[]interface{}{"in", []int{1,2}}} )
// []base.WhereItem{ {userinfo.Field.F_sex, 1}, {userinfo.Field.F_reg_channel, []interface{}{"in", []int{1,2}}} }
func (m *Model) Where(param interface{}) *Model {
	if where, ok := param.(map[string]interface{}); ok {
		if len(where) == 0 {
			return m
		}
		if m.where == nil {
			m.where = make([]WhereItem, 0)
		}
		for k, v := range where {
			m.where = append(m.where, WhereItem{k, v})
		}
	}
	if where, ok := param.([]WhereItem); ok {
		if len(where) == 0 {
			return m
		}
		if m.where == nil {
			m.where = make([]WhereItem, 0)
		}
		m.where = append(m.where, where...)
	}
	return m
}

//设置查询范围
//使用示例 Limit(10) limit(0,10)
func (m *Model) Limit(start interface{}, limit ...interface{}) *Model {
	if m.limit == nil {
		m.limit = make([]interface{}, 0)
	}
	if len(limit) == 0 {
		m.limit = append(m.limit, start)
	} else {
		m.limit = append(m.limit, start)
		m.limit = append(m.limit, limit[0])
	}
	return m
}

//设置排序
//使用示例 OrderBy("id asc","age desc")
func (m *Model) OrderBy(params ...string) *Model {
	if len(params) == 0 {
		return m
	}
	if m.orderBy == nil {
		m.orderBy = make([]string, len(params))
	}
	for k, v := range params {
		v = strings.ToLower(v)
		m.orderBy[k] = v
	}
	return m
}

//查询的字段
//使用示例 Field([]string{"name","age"})
func (m *Model) Field(param ...[]string) *Model {
	if len(param) == 0 {
		m.field = "*"
	} else {
		m.field = strings.TrimRight(strings.Join(param[0], ","), ",")
	}
	return m
}

//mysql for update
//只有在事务时才会生效，如果使用主键则锁行；否则锁表
func (m *Model) Lock(args bool) *Model {
	if args {
		m.lock = " for update"
	}
	return m
}

//存储新增、更新数据
//使用示例
//maps := make(map[string]interface{})
//maps["name"] = "lidazhao"
//maps["age"]  = 21
//map["level"] = []interface{}{"inc", 1} //自增1
//map["level"] = []interface{}{"dec", 1} //自减1
//map["level"] = []interface{}{"concat", "asdf"} //字符串连接
//Data(maps)
func (m *Model) Data(param map[string]interface{}) *Model {
	m.data = make(map[string]interface{})
	m.data = param
	return m
}

//新增数据
func (m *Model) Insert() (int, error) {
	//分析参数
	if len(m.data) == 0 {
		return 0, nil
	}
	var colsName, colsValue = "", ""
	param := []interface{}{}
	for i, v := range m.data {
		colsName += "`" + i + "`" + ","
		//如果为整型则转字符串类型
		colsValue += "?,"
		param = append(param, v)
	}
	colsName = strings.TrimRight(colsName, ",")
	colsValue = strings.TrimRight(colsValue, ",")
	// 组合数据写入SQL
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s);", m.table, colsName, colsValue)
	retData, err := m.o.Raw(sql, param...).Exec()
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", param))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
		return 0, err
	}
	LastId, err := retData.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(LastId), err
}

//批量添加
func (m *Model) InsertAll(data []map[string]interface{}) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	var keys []string
	var colsName, colsValue = "", ""
	for i, _ := range data[0] {
		colsName += "`" + i + "`" + ","
		keys = append(keys, i)
		//如果为整型则转字符串类型
	}
	colsName = strings.TrimRight(colsName, ",")
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES ", m.table, colsName)
	values := []interface{}{}
	for _, v := range data {
		colsValue += "("
		for _, k := range keys {
			colsValue += "?,"
			values = append(values, v[k])
		}
		colsValue = strings.TrimRight(colsValue, ",")
		colsValue += "),"
	}
	colsValue = strings.TrimRight(colsValue, ",")
	sql = fmt.Sprintf("%s %s;", sql, colsValue)
	retData, err := m.o.Raw(sql, values...).Exec()
	m.sql = fmt.Sprintf("%s - `%s`", sql, functions.Implode("`, `", values))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
		return 0, err
	}
	LastId, err := retData.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(LastId), err
}

//查询多条数据
func (m *Model) Select() []map[string]interface{} {
	var field string
	if m.field == "" {
		field = "*"
	} else {
		field = m.field
	}
	where, param := m.whereString()
	var orderBy string
	if len(m.orderBy) > 0 {
		for _, v := range m.orderBy {
			orderBy += v + ","
		}
		orderBy = " ORDER BY " + strings.TrimRight(orderBy, ",")
	}
	var limit string
	if len(m.limit) > 0 {
		if len(m.limit) == 1 {
			limit = strconv.Itoa(m.limit[0].(int))
		} else {
			limit = strconv.Itoa(m.limit[0].(int)) + "," + strconv.Itoa(m.limit[1].(int))
		}
		limit = " LIMIT " + limit
	}
	sql := fmt.Sprintf("SELECT %s FROM %s%s%s%s", field, m.table, where, orderBy, limit)
	var res []orm.Params
	_, err := m.o.Raw(sql, param...).Values(&res)
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", param))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
	}
	var maps = make([]map[string]interface{}, len(res))
	if len(res) > 0 {
		for i, v := range res {
			maps[i] = v
		}
	}
	return maps
}

//查询单条数据
func (m *Model) Find() map[string]interface{} {
	var field string
	if m.field == "" {
		field = "*"
	} else {
		field = m.field
	}
	where, param := m.whereString()
	sql := ""
	if m.lock == "" {
		sql = fmt.Sprintf("SELECT %s FROM %s%s LIMIT 1", field, m.table, where)
	} else {
		sql = fmt.Sprintf("SELECT %s FROM %s%s LIMIT 1 %s", field, m.table, where, m.lock)
	}
	var res []orm.Params
	_, err := m.o.Raw(sql, param...).Values(&res)
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", param))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
	}
	if len(res) == 0 {
		return make(map[string]interface{})
	}
	return res[0]
}

//更新
func (m *Model) Update() (int, error) {
	//分析参数
	if len(m.data) == 0 {
		return 0, nil
	}
	var updateStr string
	field := []interface{}{}
	for i, v := range m.data {
		if val, ok := v.([]interface{}); ok {
			if len(val) == 2 {
				m.arrayData(i, val[0].(string), val[1], &updateStr, &field)
			}
		} else {
			updateStr += i + "=?,"
			field = append(field, v)
		}
	}
	updateStr = strings.TrimRight(updateStr, ",")
	where, param := m.whereString()

	if where == "" {
		return 0, errors.New("must have where")
	}

	sql := fmt.Sprintf("UPDATE %s SET %s%s", m.table, updateStr, where)
	param = append(field, param...)
	sqlSource, err := m.o.Raw(sql, param...).Exec()
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", param))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
		return 0, err
	}
	num, _ := sqlSource.RowsAffected()
	return int(num), err
}

//条件更新 - 批量更新
func (m *Model) UpdateCase(data []map[string]interface{}, field_key string) (int, error) {
	f := map[string]string{}
	fp := map[string][]interface{}{}
	where := []interface{}{}
	whereStr := ""

	for _, val := range data {
		if _, ok := val[field_key]; !ok {
			continue
		}
		where = append(where, val[field_key])
		whereStr += "?,"
		for k, v := range val {
			if k == field_key {
				continue
			}
			if v1, ok := v.([]interface{}); ok {
				if len(v1) == 2 {
					str := ""
					m.arrayData(k, v1[0].(string), v1[1], &str, &[]interface{}{})
					str = strings.TrimRight(str, ",")
					str = str[len(k+"="):]
					f[k] += " when ? then " + str
					fp[k] = append(fp[k], val[field_key], v1[1])
				}
			} else {
				f[k] += " when ? then ? "
				fp[k] = append(fp[k], val[field_key], v)
			}
		}
	}
	updateStr := ""
	param := []interface{}{}
	for k, v := range f {
		updateStr += fmt.Sprintf("%s=case %s%s else %s end,", k, field_key, v, k)
		param = append(param, fp[k]...)
	}

	param = append(param, where...)

	updateStr = strings.TrimRight(updateStr, ",")
	where1, param1 := m.whereString()
	param = append(param, param1...)
	where1 = strings.Replace(where1, "WHERE", " and ", 1)

	whereStr = strings.TrimRight(whereStr, ",")
	updateStr = strings.TrimRight(updateStr, ",")
	sql := fmt.Sprintf("UPDATE %s SET %s where %s in (%s) %s", m.table, updateStr, field_key, whereStr, where1)
	sqlSource, err := m.o.Raw(sql, param...).Exec()

	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", param))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
		return 0, err
	}
	num, _ := sqlSource.RowsAffected()
	return int(num), err
}

//物理删除
func (m *Model) Delete() (int, error) {
	where, param := m.whereString()
	if where == "" {
		return 0, nil
	}
	sql := fmt.Sprintf("DELETE FROM %s%s", m.table, where)
	sqlSource, err := m.o.Raw(sql, param...).Exec()
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", param))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
		return 0, nil
	}
	num, _ := sqlSource.RowsAffected()
	return int(num), err
}

//统计数据
//使用示例
//Count() 或 Count("id") //id为字段名
func (m *Model) Count(param ...string) int {
	co := "*"
	if len(param) != 0 {
		co = param[0]
	}
	where, pargm := m.whereString()
	sql := fmt.Sprintf("SELECT COUNT(%s) FROM %s%s", co, m.table, where)
	var maps []orm.Params
	_, err := m.o.Raw(sql, pargm...).Values(&maps)
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", pargm))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
	}
	num, _ := strconv.Atoi(maps[0]["COUNT("+co+")"].(string))
	return num
}

//聚合函数-sum
func (m *Model) Sum(field string) float64 {
	where, pargm := m.whereString()
	sql := fmt.Sprintf("SELECT SUM(%s) FROM %s%s", field, m.table, where)
	var maps []orm.Params
	_, err := m.o.Raw(sql, pargm...).Values(&maps)
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", pargm))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
	}
	var num float64
	if maps[0]["SUM("+field+")"] != nil {
		num, _ = strconv.ParseFloat(maps[0]["SUM("+field+")"].(string), 64)
	} else {
		num = 0
	}
	return num
}

//聚合函数-Avg
func (m *Model) Avg(field string) float64 {
	where, pargm := m.whereString()
	sql := fmt.Sprintf("SELECT AVG(%s) FROM %s%s", field, m.table, where)
	var maps []orm.Params
	_, err := m.o.Raw(sql, pargm...).Values(&maps)
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", pargm))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
	}
	var num float64
	if maps[0]["AVG("+field+")"] != nil {
		num, _ = strconv.ParseFloat(maps[0]["AVG("+field+")"].(string), 64)
	} else {
		num = 0
	}
	return num
}

//聚合函数-Min
func (m *Model) Min(field string) float64 {
	where, pargm := m.whereString()
	sql := fmt.Sprintf("SELECT MIN(%s) FROM %s%s", field, m.table, where)
	var maps []orm.Params
	_, err := m.o.Raw(sql, pargm...).Values(&maps)
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", pargm))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
	}
	var num float64
	if maps[0]["MIN("+field+")"] != nil {
		num, _ = strconv.ParseFloat(maps[0]["MIN("+field+")"].(string), 64)
	} else {
		num = 0
	}
	return num
}

//聚合函数-max
func (m *Model) Max(field string) float64 {
	where, pargm := m.whereString()
	sql := fmt.Sprintf("SELECT MAX(%s) FROM %s%s", field, m.table, where)
	var maps []orm.Params
	_, err := m.o.Raw(sql, pargm...).Values(&maps)
	m.sql = fmt.Sprintf("%s-`%s`", sql, functions.Implode("`, `", pargm))
	if err != nil {
		if dzgin.KcConfig.RunMode != dzgin.PROD {
			logs.Error("Sql:", sql, " Error,", err.Error())
		}
	}
	var num float64
	if maps[0]["MAX("+field+")"] != nil {
		num, _ = strconv.ParseFloat(maps[0]["MAX("+field+")"].(string), 64)
	} else {
		num = 0
	}
	return num
}

//事务开始
func (m *Model) Begin() *Model {
	err := m.o.Begin()
	if err != nil {
		logs.Error("Begin Error", err.Error())
	}
	return m
}

//事务提交
func (m *Model) Commit() *Model {
	err := m.o.Commit()
	if err != nil {
		logs.Error("Commit Error", err.Error())
	}
	return m
}

//事务回滚
func (m *Model) RollBack() *Model {
	err := m.o.Rollback()
	if err != nil {
		logs.Error("RollBack Error", err.Error())
	}
	return m
}

//打印sql语句
//使用示例
//GetLastSql() 返回sql语句 GetLastSql(true) 打印控制台
func (m *Model) GetLastSql(param ...bool) string {
	var isPrint bool
	if len(param) != 0 {
		isPrint = true
	}
	if isPrint {
		fmt.Println("sql:", m.sql)
		return ""
	}
	return m.sql
}

//组织where字符串
func (m *Model) whereString() (string, []interface{}) {
	param := []interface{}{}
	var where string = ""
	hasOr := 0
	if len(m.where) != 0 {
		for wi, v := range m.where {
			nextIsOr := false
			if (wi + 1) < len(m.where) {
				if strings.ToUpper(m.where[wi+1].Field) == MYSQL_CONDITION_OR {
					nextIsOr = true
				}
			}
			if k, ok := v.Value.([]interface{}); ok {
				if len(k) == 2 {
					m.arrayWhere(v.Field, k[0].(string), k[1], &where, &param, nextIsOr)
				}
			} else if k, ok := v.Value.([]string); ok {
				if len(k) == 2 {
					m.arrayWhere(v.Field, k[0], k[1], &where, &param, nextIsOr)
				}
			} else if strings.ToUpper(v.Field) == MYSQL_CONDITION_OR {
				hasOr = 1
				where += " ) OR  ( "
			} else {
				if nextIsOr == true {
					where += v.Field + "=? "
				} else {
					where += v.Field + "=? AND "
				}

				param = append(param, v.Value)
			}
		}

		if hasOr == 1 {
			where = " WHERE ( " + strings.TrimRight(where, " AND") + " ) "
		} else {
			where = " WHERE " + strings.TrimRight(where, " AND")
		}
	}
	m.where = make([]WhereItem, 0)
	return where, param
}

func (m *Model) arrayWhere(name string, condition string, v interface{}, where *string, param *[]interface{}, nextIsOr bool) {
	switch strings.ToLower(condition) {
	case "eq":
		condition = "="
		break
	case "neq":
		condition = "<>"
		break
	case "gt":
		condition = ">"
		break
	case "egt":
		condition = ">="
		break
	case "lt":
		condition = "<"
		break
	case "elt":
		condition = "<="
		break
	case "between":
		if nextIsOr == true {
			*where += name + " " + condition + " ? and ?  "
		} else {
			*where += name + " " + condition + " ? and ? AND "
		}
		*param = append(*param, v)
		return
	}

	ars := "?"
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		for i := 1; i < reflect.ValueOf(v).Len(); i++ {
			ars += ",?"
		}
	}
	if nextIsOr == true {
		*where += name + " " + condition + "(" + ars + ") "
	} else {
		*where += name + " " + condition + "(" + ars + ") AND "
	}
	*param = append(*param, v)
}

func (m *Model) arrayData(name string, condition string, v interface{}, field *string, param *[]interface{}) {
	switch strings.ToLower(condition) {
	case "inc":
		*field += name + "=" + name + "+?,"
		break
	case "dec":
		*field += name + "=" + name + "-?,"
		break
	case "concat":
		*field += name + "=concat(" + name + ",?),"
		break
	}
	*param = append(*param, v)
}

func (m *Model) GetOrmer() orm.Ormer {
	return m.o
}

//实例化Model引用
//@param string table 表名称
//使用示例
//NewModel("student")
func NewModel(table string, ormer ...orm.Ormer) *Model {
	var ormers orm.Ormer
	if len(ormer) > 0 {
		ormers = ormer[0]
	} else {
		ormers = orm.NewOrm()
	}
	return &Model{
		table: T_PREFIX + table,
		o:     ormers,
	}
}
