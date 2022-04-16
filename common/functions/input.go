package functions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strconv"
	"strings"
)

type Data string
func (d Data) Form(data interface{}) Data {
	return Data(fmt.Sprint(data))
}
func (d Data) String() string {
	return string(d)
}
func (d Data) Int() int {
	data, _ := strconv.Atoi(string(d))
	return data
}
func (d Data) Int64() int64 {
	data, _ := strconv.ParseInt(string(d), 10, 64)
	return data
}
func (d Data) Uint64() uint64 {
	data, _ := strconv.ParseUint(string(d), 10, 64)
	return data
}
func (d Data) Float64() float64 {
	data, _ := strconv.ParseFloat(string(d), 64)
	return data
}
func (d Data) Bool() bool {
	data, _ := strconv.ParseBool(string(d))
	return data
}
func (d Data) StringArray(sep... string) []string {
	sp := ","
	if len(sep) > 0 {
		sp = sep[0]
	}
	return strings.Split(strings.Trim(string(d), "[] "), sp)
}
func (d Data) IntArray(sep... string) []int {
	sp := ","
	if len(sep) > 0 {
		sp = sep[0]
	}
	data := strings.Split(strings.Trim(string(d), "[] "), sp)
	var output []int
	for _,v := range data {
		output = append(output, Data(v).Int())
	}
	return output
}

type input struct {
	Ctx *gin.Context
	encry bool
}

func MakeInput(k *Controller) *input {
	i := &input{Ctx:k.Ctx}
	i.init()
	return i
}

// 获取Get参数
func (i *input)Get(key string, defaultValue... string) Data{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetQuery(key); ok {
		return Data(i.getSrcStr(value))
	}
	if defaultValue != nil {
		return Data(defaultValue[0])
	}
	return ""
}

// 获取GetArray参数
func (i *input)GetArray(key string) []string{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetQueryArray(key); ok {
		for k,v := range value {
			value[k] = i.getSrcStr(v)
		}
		return value
	}
	return nil
}

// 获取GetMap参数
func (i *input)GetMap(key string) map[string]string{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetQueryMap(key); ok {
		res := make(map[string]string)
		for k,v := range value {
			res[k] = i.getSrcStr(v)
		}
		return res
	}
	return nil
}

// 获取Post参数
func (i *input)Post(key string, defaultValue... string) Data{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetPostForm(key); ok {
		return Data(i.getSrcStr(value))
	}
	if defaultValue != nil {
		return Data(defaultValue[0])
	}
	return ""
}

// 获取PostArray参数
func (i *input)PostArray(key string) []string{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetPostFormArray(key); ok {
		for k,v := range value {
			value[k] = i.getSrcStr(v)
		}
		return value
	}
	return nil
}

// 获取PostMap参数
func (i *input)PostMap(key string) map[string]string{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetPostFormMap(key); ok {
		res := make(map[string]string)
		for k,v := range value {
			res[k] = i.getSrcStr(v)
		}
		return res
	}
	return nil
}

// 获取GetPost参数
func (i *input)GetPost(key string, defaultValue... string) Data{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetQuery(key); ok {
		return Data(i.getSrcStr(value))
	}
	if value, ok := i.Ctx.GetPostForm(key); ok {
		return Data(i.getSrcStr(value))
	}
	if defaultValue != nil {
		return Data(defaultValue[0])
	}
	return ""
}

// 获取GetPostArray参数
func (i *input)GetPostArray(key string) []string{
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetQueryArray(key); ok {
		for k,v := range value {
			value[k] = i.getSrcStr(v)
		}
		return value
	}
	if value, ok := i.Ctx.GetPostFormArray(key); ok {
		for k,v := range value {
			value[k] = i.getSrcStr(v)
		}
		return value
	}
	return nil
}

// 获取GetPostMap参数
func (i *input)GetPostMap(key string) (res map[string]string){
	key = i.getSrcStr(key)
	if value, ok := i.Ctx.GetQueryMap(key); ok {
		res = make(map[string]string)
		for k,v := range value {
			res[k] = i.getSrcStr(v)
		}
	}
	if value, ok := i.Ctx.GetPostFormMap(key); ok {
		if res == nil{
			res = make(map[string]string)
		}
		for k,v := range value {
			res[k] = i.getSrcStr(v)
		}
	}
	return
}

// 判断是否加密
func (i *input)IsEncry() bool{
	return i.encry
}

// 初始化
func (i *input)init(){
	i.encry = false
}

// 获取字符串原始数据
func (i *input)getSrcStr(str string) string {
	if i.IsEncry() {
		return DecodeStr(str)
	}
	return str
}

func (i *input) Cookie(name string) (string, error) {
	return i.Ctx.Cookie(name)
}

func (i *input) File(name string) (*multipart.FileHeader, error) {
	return i.Ctx.FormFile(name)
}

func (i *input) MulFile() (*multipart.Form, error) {
	return i.Ctx.MultipartForm()
}

func (i *input) Header(key string) Data {
	return Data(i.Ctx.GetHeader(key))
}
