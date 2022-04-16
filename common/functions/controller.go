package functions

import (
	"github.com/Gouplook/dzgin"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	dzgin.Controller
	Input  *input
	Output *output
	Public public
}

type public struct {
	Utoken  string  // 登录验证token
	Channel int     // 渠道
	Device  int     // 设备
	Version string  // 版本
	Cid     int     // 城市id
	Lng     float64 // 经度
	Lat     float64 // 维度
}

func (k *Controller) Init(ctx *gin.Context, method string) {
	k.Controller.Init(ctx, method)

	k.Input = MakeInput(k)
	k.Output = MakeOutput(k)

	k.Public.Utoken = k.Input.Header("utoken").String()
	k.Public.Channel = k.Input.Header("channel").Int()
	k.Public.Device = k.Input.Header("device").Int()
	k.Public.Version = k.Input.Header("version").String()
	k.Public.Cid = k.Input.Header("cid").Int()
	k.Public.Lng = k.Input.Header("lng").Float64()
	k.Public.Lat = k.Input.Header("lat").Float64()
}
