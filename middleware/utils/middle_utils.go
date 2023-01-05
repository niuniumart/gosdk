// Package utils gin中间件封装
package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	REQUEST_SUCCESS = 0
)

//RequestGetter 获取包体通用结构
type RequestGetter struct {
	RequestID string `json:"RequestId" form:"RequestId"`
	Token     string `json:"token" form:"token"`
	Module    string `json:"module" form:"module"`
}

//DefaultRespGetter 默认返回结构获取器
type DefaultRespGetter struct {
	Code int    `json:"retCode"`
	Msg  string `json:"retMsg"`
}

//GetCode 获取Code属性
func (p *DefaultRespGetter) GetCode() int {
	return p.Code
}

var defaultRespGetterFactory = func() RespGetter {
	return new(DefaultRespGetter)
}

// 初始化
func init() {
	respGetterFactory = defaultRespGetterFactory
}

//SetRespGetterFactory 设置返回器工厂
func SetRespGetterFactory(factory RespGetterFactory) {
	respGetterFactory = factory
}

//GetRespGetterFactory 获取返回器工厂
func GetRespGetterFactory() RespGetterFactory {
	return respGetterFactory
}

//RespGetterFactory
type RespGetterFactory func() RespGetter

var respGetterFactory RespGetterFactory

//RespGetter 返回处理器
type RespGetter interface {
	GetCode() int
}

//BodyLogWriter 日志打印器
type BodyLogWriter struct {
	gin.ResponseWriter
	BodyBuf *bytes.Buffer
}

//Write implement body log writer
func (w BodyLogWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.BodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

/**
UrlMetrics
URL_CONFIG_UPDATE_NOTIFY
UrlHeartBeat
*/
const (
	UrlMetrics   = "/metrics"
	UrlHeartBeat = "/heartbeat"
)

//IgnorePaths
var IgnorePaths []string

func init() {
	IgnorePaths = []string{
		UrlMetrics,
		UrlHeartBeat,
	}
}

/**
TotalCounterVec
ReqDurationVec
ReqLogicErrorVec
ReqSystemErrorVec
*/
var (
	TotalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Total number of HTTP requests made",
		},
		[]string{"module", "operation"},
	)
	ReqDurationVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "request_latency",
		Help: "record request latency",
	}, []string{"module", "operation"})
	ReqLogicErrorVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_error_count",
		Help: "Total request error count of the host",
	}, []string{"module", "operation", "code"})
	ReqSystemErrorVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_error",
		Help: "Total error count of the request",
	}, []string{"module", "operation"})
)

func init() {
	prometheus.MustRegister(
		TotalCounterVec,
		ReqDurationVec,
		ReqLogicErrorVec,
		ReqSystemErrorVec,
	)
}
