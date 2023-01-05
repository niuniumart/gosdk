// Package websdk for web frame
package gin

import (
	"context"
	"fmt"
	"github.com/niuniumart/gosdk/middleware/logprint"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/niuniumart/gosdk/middleware/cors"

	// 加入pprof功能
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	recoverSdk "github.com/niuniumart/gosdk/middleware/recover"
	"github.com/niuniumart/gosdk/middleware/utils"
	"github.com/niuniumart/gosdk/seelog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	debug.SetGCPercent(1000)
}

//CreateTBaasGin create tbaas gin instance
func CreateTBaasGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.Use(recoverSdk.PanicRecover())
	engine.Use(logprint.InfoLog())
	engine.Use(cors.Cors())
	engine.GET(utils.UrlMetrics, gin.WrapH(promhttp.Handler()))
	engine.GET(utils.UrlHeartBeat, HeartBeat)
	return engine
}

//RunByPort run with port
func RunByPort(engine *gin.Engine, port int) error {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				seelog.Errorf("In pprof PanicRecover,Error:%s", err)
			}
		}()
		err := http.ListenAndServe(fmt.Sprintf(":%d", port+3), nil) //开启一个http服务
		if err != nil {
			seelog.Errorf("ListenAndServe: ", err)
			return
		}
	}()
	return Run(engine, fmt.Sprintf("%d", port))
}

//Run run web sever
//param engine: instance of gin.Engine
//param port: format as :port, for example :31112
func Run(engine *gin.Engine, port string) error {
	var runPort string
	if port[0] == ':' {
		runPort = port
	} else {
		runPort = fmt.Sprintf(":%s", port)
	}

	return engine.Run(runPort)
}

//RunWithGraceShutDown run with grace shutdown
func RunWithGraceShutDown(engine *gin.Engine, port string, timeout int) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			seelog.Errorf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of input seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	seelog.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		seelog.Errorf("Server Shutdown err:%s", err)
	}
	// catching ctx.Done(). timeout of input seconds.
	select {
	case <-ctx.Done():
		seelog.Infof("Reach timeout of %d seconds.", timeout)
	}
	seelog.Infof("Server exiting")
}

//SetRespGetterFactory set resp getter factory
func SetRespGetterFactory(factory utils.RespGetterFactory) {
	utils.SetRespGetterFactory(factory)
}

//HeartBeat heart beat
func HeartBeat(c *gin.Context) {
	c.String(http.StatusOK, "SUCCESS")
}
