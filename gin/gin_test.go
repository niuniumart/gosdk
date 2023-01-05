package gin

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niuniumart/gosdk/middleware/utils"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateTBaasGin(t *testing.T) {
	convey.Convey("TestCreateTBaasGin", t, func() {
		seelog.Infof("just test log", 5, 666)
		engine := CreateTBaasGin()
		engine.POST("/reverse", Reverse)
		engine.GET("/panic", MustPanic)
		engine.GET("/logicerr", LogicError)
		go func() {
			engine.Run(":30001")
			time.Sleep(20 * time.Second)
		}()
		var url = "http://127.0.0.1:30001/reverse"
		param := make(map[string]interface{})
		param["abc"] = "xxx"
		//		requestid.Set("heiheiheihei")
		resp := DoRequest(url, http.MethodPost, param)
		fmt.Printf("resp %s\n", resp)
	})
}

type RespGetter struct {
	Code int    `json:"retCode"`
	Msg  string `json:"retMsg"`
}

func (p *RespGetter) GetCode() int {
	return p.Code
}

func TestCreateTBaasGinServer(t *testing.T) {
	convey.Convey("TestCreateTBaasGinServer", t, func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		var fac utils.RespGetterFactory = func() utils.RespGetter {
			return new(RespGetter)
		}
		SetRespGetterFactory(fac)
		engine := CreateTBaasGin()
		engine.GET("/ping", Pong)
		engine.POST("/reverse", Reverse)
		err := Run(engine, ":31112")
		fmt.Println(err)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar\n"))
}

type ConfigTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Pong(c *gin.Context) {
	fmt.Println("pong")
	c.JSON(http.StatusOK, "pong")
}

func Reverse(c *gin.Context) {
	c.JSON(http.StatusOK, "into")
}

func MustPanic(c *gin.Context) {
	panic(nil)
}

type Resp struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

func LogicError(c *gin.Context) {
	var resp Resp
	resp.RetCode = 10005
	c.JSON(http.StatusOK, &resp)
}

func TestPanic(t *testing.T) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				seelog.Errorf("In PanicRecover,Error:%s", err)
				//打印调用栈信息
			}
		}()
		var dic map[string]int
		dic["aaa"] = 1
		fmt.Printf("dic %v\n", dic)
	}()
}
