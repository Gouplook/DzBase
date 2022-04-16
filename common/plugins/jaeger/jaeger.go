//Jaeger版本的链路追踪
package jaeger

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"io/ioutil"
)

//创建tracer
//@param service 当前服务名称
//@param Type 采集方式 有const, probabilistic, rateLimiting, remote 四种方式
//-- const 采样器始终对所有迹线做出相同的决定。它要么采样所有轨迹（sampler.param=1），要么不采样（sampler.param=0）。
//-- probabilistic 采样器做出随机采样决策，采样概率等于sampler.param属性值。例如，sampler.param=0.1大约有十分之一的轨迹将被采样
//-- rateLimiting 采样器使用漏斗速率限制器来确保以一定的恒定速率对轨迹进行采样。例如，sampler.param=2.0时，以每秒2条跟踪的速率对请求进行采样
//-- remote 采样器请咨询Jaeger代理以获取在当前服务中使用的适当采样策略。这允许从Jaeger后端的中央配置甚至动态地控制服务中的采样策略
//@param Param 配合Type使用的 采样概率值
//@param agentHost 代理地址:端口
func NewJaeger(service string, Type string, Param float64, agentHost string) (opentracing.Tracer, io.Closer, error) {
	if len(Type) == 0 {
		Type = "probabilistic"
		Param = 0.1
	}
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  Type,
			Param: Param,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentHost,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, nil, err
	}
	//设置为全局的
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}

//开启链路追踪
func OpenJaeger() (opentracing.Tracer, io.Closer, error) {
	isopen, err := kcgin.AppConfig.Bool("jaeger.open")
	if isopen == false {
		return nil, nil, err
	}
	serviceName := kcgin.AppConfig.String("jaeger.serviceName")
	if len(serviceName) == 0 {
		return nil, nil, errors.New("jaeger.serviceName is null")
	}
	jtype := kcgin.AppConfig.String("jaeger.jtype")
	param, err := kcgin.AppConfig.Float("jaeger.param")
	if err != nil {
		return nil, nil, err
	}
	agentHost := kcgin.AppConfig.String("jaeger.agentHost")
	if len(agentHost) == 0 {
		return nil, nil, errors.New("jaeger.agentHost is empty")
	}
	return NewJaeger(serviceName, jtype, param, agentHost)
}

type LogWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (lw *LogWriter) Write(b []byte) (int, error) {
	lw.Body.Write(b)
	return lw.ResponseWriter.Write(b)
}

func (lw *LogWriter) WriteString(s string) (int, error) {
	lw.Body.WriteString(s)
	return lw.ResponseWriter.WriteString(s)
}

//中间件
func SpanMiddle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span, cont, err := RpcxSpanWithContext(ctx.Request.Context(), fmt.Sprintf("请求地址：%s", ctx.Request.URL.Path), ctx.Request)
		if err == nil {
			defer func() {
				err2 := recover()
				if err2 != nil {
					span.SetTag("error", true)
					span.LogKV("错误信息", fmt.Sprint(err2))
				}
				span.Finish()
				if err2 != nil {
					panic(err2)
				}
			}()

			span.SetTag("请求方式", ctx.Request.Method)
			span.LogKV("get 参数", ctx.Request.URL.Query())
			method := ctx.Request.Method
			if method == "PUT" {
				raw, _ := ctx.GetRawData()
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
				span.LogKV("raw 参数", string(raw))
			} else {
				ctx.MultipartForm()
				span.LogKV("post 参数", ctx.Request.PostForm)
			}
			ctx.Request = ctx.Request.WithContext(cont)
		}

		responseBody := &LogWriter{
			ResponseWriter: ctx.Writer,
			Body:           bytes.NewBufferString(""),
		}
		ctx.Writer = responseBody
		ctx.Next()
		span.LogKV("返回数据：", responseBody.Body.String())
	}
}
