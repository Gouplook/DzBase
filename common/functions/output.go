package functions

import (
	"bytes"
	"encoding/json"
	"github.com/Gouplook/dzbase/lang"
	"github.com/Gouplook/dzgin"
	"github.com/gin-gonic/gin"
)

type output struct {
	Ctx        *gin.Context
	encry      bool
	data       map[string]interface{}
	controller *Controller

	templates string
}

//正确输出
func (o *output) Success(params ...interface{}) {
	var msg string
	var odata interface{}
	if len(params) > 0 {
		odata = params[0]
		if len(params) > 1 {
			if _, ok := params[1].(string); ok {
				msg = params[1].(string)
			}
		}
	} else {
		odata = make([]interface{}, 0)
	}
	maps := map[string]interface{}{
		"error":    0,
		"errorMsg": msg,
		"data":     odata,
	}
	var content []byte
	o.Ctx.Header("Content-Type", "application/json; charset=utf-8")
	content, _ = json.Marshal(maps)
	jsonS := string(content)
	if jsonCallback := o.Ctx.Query("jsonCallback"); jsonCallback != "" {
		jsonS = jsonCallback + "(" + jsonS + ")"
	}
	res := StringsToJSON(jsonS)
	if o.encry {
		res = EncodeStr(res)
	}
	o.Ctx.String(200, res)
}

//错误输出
func (o *output) Error(errCode interface{}, errorMsg ...string) {
	odata := make([]interface{}, 0)
	var errMsg = ""
	if len(errorMsg) > 0 {
		errMsg = errorMsg[0]
	} else {
		errMsg = lang.GetLang(errCode.(string))

	}

	maps := map[string]interface{}{
		"error":    errCode,
		"errorMsg": errMsg,
		"data":     odata,
	}
	var content []byte
	o.Ctx.Header("Content-Type", "application/json; charset=utf-8")
	content, _ = json.Marshal(maps)
	jsonS := string(content)
	if jsonCallback := o.Ctx.Query("jsonCallback"); jsonCallback != "" {
		jsonS = jsonCallback + "(" + jsonS + ")"
	}
	res := StringsToJSON(jsonS)
	if o.encry {
		res = EncodeStr(res)
	}
	o.Ctx.String(200, res)
}

func MakeOutput(k *Controller) *output {
	o := &output{Ctx: k.Ctx, controller: k}
	o.init()
	return o
}

// 初始化
func (o *output) init() {
	o.encry = false
}

// 判断是否加密
func (o *output) IsEncry() bool {
	return o.encry
}

func (o *output) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	o.Ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (o *output) Header(key string, value string) {
	o.Ctx.Header(key, value)
}

func (o *output) SetDefaultTmpl(templates string) {
	o.templates = templates
}

func (o *output) getDefaultTmpl() string {
	if o.templates != "" {
		return o.templates
	}
	return dzgin.AppConfig.DefaultString("tmpl.defaultview", "templates")
}

// 设置模板数据
func (o *output) Assign(key string, value interface{}) {
	if o.data == nil {
		o.data = map[string]interface{}{
			key: value,
		}
	} else {
		o.data[key] = value
	}
}

// 渲染模板
func (o *output) Html(tmplname string) {
	if dzgin.KcConfig.RunMode != dzgin.PROD {
		_ = dzgin.BuildTemplate(o.getDefaultTmpl())
	}
	var buf bytes.Buffer
	_ = dzgin.ExecuteViewPathTemplate(&buf, tmplname, o.getDefaultTmpl(), o.data)
	o.Header("Content-Type", "text/html; charset=utf-8")
	o.Ctx.String(200, string(buf.Bytes()))
}
