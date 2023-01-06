// Package mlog 日志封装
package logprint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/niuniumart/gosdk/martlog"
	"github.com/niuniumart/gosdk/middleware/utils"
	"github.com/niuniumart/gosdk/requestid"
	"github.com/niuniumart/gosdk/tools"
	"io/ioutil"
	"strings"
	"time"
)

const MaxPrintBodyLen = 1024

// writer, use for log
type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

//Write write func
func (w bodyLogWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

const (
	REQUEST_KEY    = "Cloud-Trace-Id"
	REQUEST_MODULE = "Cloud-Module"
	REQUEST_ID_KEY = "Request-Id"
)

var ignoreReqLogUrlDic, ignoreRespLogUrl map[string]int

// init
func init() {
	ignoreReqLogUrlDic = make(map[string]int)
	ignoreRespLogUrl = make(map[string]int)
}

//RegisterIgnoreLogUrl func registerIgnoreLogUrl url
func RegisterIgnoreLogUrl(url string) {
	ignoreReqLogUrlDic[url] = 1
	ignoreRespLogUrl[url] = 1
}

//RegisterIgnoreReqLogUrl register ignore req log url/
func RegisterIgnoreReqLogUrl(url string) {
	ignoreReqLogUrlDic[url] = 1
}

//RegisterIgnoreRespLogUrl register ignore resp log url/
func RegisterIgnoreRespLogUrl(url string) {
	ignoreRespLogUrl[url] = 1
}

//InfoLog func infoLog
func InfoLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tools.InList(c.Request.URL.Path, utils.IgnorePaths) {
			return
		}
		beginTime := time.Now()
		// ***** 1. get request body ****** //
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close() //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// ***** 2. set requestID for goroutine ctx ****** //
		requestID := c.Request.Header.Get(REQUEST_KEY)
		requestModule := c.Request.Header.Get(REQUEST_MODULE)
		if requestID == "" {
			requestID = c.Request.Header.Get(REQUEST_ID_KEY)
		}
		if requestID == "" {
			var rg utils.RequestGetter
			err := json.Unmarshal(body, &rg)
			if err != nil {
				martlog.Warnf("Req Body json unmarshal requestID err %s", err.Error())
			}
			if requestModule == "" {
				requestModule = rg.Module
			}
		}
		utils.TotalCounterVec.WithLabelValues(requestModule, c.Request.URL.Path).Inc()
		if requestID == "" {
			requestID = fmt.Sprintf("%+v", uuid.New())
		}
		requestid.Set(requestID)
		defer requestid.Delete()
		if _, ok := ignoreReqLogUrlDic[c.Request.URL.Path]; !ok {
			martlog.Infof("Req Url: %s %+v,[Body]:%s; [Header]:%s", c.Request.Method, c.Request.URL,
				string(body), tools.GetFmtStr(c.Request.Header))
		}
		// ***** 3. set resp writer ****** //
		blw := utils.BodyLogWriter{BodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		// ***** 4. do Next ****** //
		c.Next()
		// ***** 5. log resp body ****** //
		strBody := strings.Trim(blw.BodyBuf.String(), "\n")
		if len(strBody) > MaxPrintBodyLen {
			strBody = strBody[:(MaxPrintBodyLen - 1)]
		}
		// ***** 6. judge logic error ****** //
		getterFactory := utils.GetRespGetterFactory()
		rspGetter := getterFactory()
		//var rspGetter utils.ResponseGetter
		json.Unmarshal(blw.BodyBuf.Bytes(), &rspGetter)
		if rspGetter.GetCode() != utils.REQUEST_SUCCESS {
			utils.ReqLogicErrorVec.WithLabelValues(requestModule, c.Request.URL.Path,
				fmt.Sprintf("%d", rspGetter.GetCode())).Inc()
		}
		if _, ok := ignoreRespLogUrl[c.Request.URL.Path]; !ok {
			martlog.Infof("Url: %+v, cost %v Resp Body %s", c.Request.URL,
				time.Since(beginTime), strBody)
		}
		duration := float64(time.Since(beginTime)) / float64(time.Second)
		martlog.Infof("ReqPath[%s]-Duration[%g]", c.Request.URL.Path, duration)
		utils.ReqDurationVec.WithLabelValues(requestModule, c.Request.URL.Path).Observe(duration)
	}
}

// return max(a, b)
func maxCode(a, b int) string {
	if a > b {
		return fmt.Sprintf("%d", a)
	}
	return fmt.Sprintf("%d", b)
}
