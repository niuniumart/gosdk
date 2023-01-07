// Package recover middileware
package recover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/niuniumart/gosdk/martlog"
	"github.com/niuniumart/gosdk/middleware/utils"
	"github.com/niuniumart/gosdk/response"
	"io/ioutil"
	"net/http"
	"runtime"
	"runtime/debug"
)

// body writer
type bodyWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

//Write func write
func (w bodyWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

//PanicRecover func panicRecover
func PanicRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				martlog.Errorf("In PanicRecover,Error:%s", err)
				var rg utils.RequestGetter
				body, _ := ioutil.ReadAll(c.Request.Body)
				err := json.Unmarshal(body, &rg)
				if err != nil {
					martlog.Warnf("Req Body json unmarshal requestID err %s", err.Error())
				}
				utils.ReqSystemErrorVec.WithLabelValues(rg.Module, c.Request.URL.Path).Inc()
				//打印调用栈信息
				debug.PrintStack()
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])
				martlog.Errorf("panic stack info %s\n", stackInfo)
				/*blw := bodyWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
				c.Writer = blw*/
				c.JSON(http.StatusOK, *response.BuildFailResp(nil))
				return
			}
		}()
		c.Next()
	}
}
