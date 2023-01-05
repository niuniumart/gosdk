package gin

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/niuniumart/gosdk/seelog"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"time"
)

//GetRequest func get request
func GetRequest(uri string, param map[string]interface{}, router *gin.Engine) string {
	return baseRequest(uri, "POST", param, router)
}

// base request
func baseRequest(uri, method string, param map[string]interface{}, router *gin.Engine) string {

	jsonByte, _ := json.Marshal(param)
	uri = uri + ParseToStr(param)
	u, _ := url.Parse(uri)
	q := u.Query()
	u.RawQuery = q.Encode() //urlencode
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest(method, u.String(), bytes.NewReader(jsonByte))

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return string(body)
}

//ParseToStr parse map to str
func ParseToStr(mp map[string]interface{}) string {
	values := ""
	if len(mp) == 0 {
		return values
	}
	for key, val := range mp {
		values += "&" + key + "=" + interface2String(val)
	}
	temp := values[1:]
	values = "?" + temp
	return values
}

// interface to string
func interface2String(inter interface{}) string {
	result := ""
	switch inter.(type) {

	case string:
		result = inter.(string)
		break
	case int:
		result = strconv.Itoa(inter.(int))
		break
	case float64:
		strconv.FormatFloat(inter.(float64), 'f', -1, 64)
		break
	}
	return result
}

//DoRequest do request
func DoRequest(uri, method string, param map[string]interface{}) string {
	var reqBody []byte
	if method == http.MethodGet {
		uri = uri + ParseToStr(param)
		u, _ := url.Parse(uri)
		q := u.Query()
		u.RawQuery = q.Encode() //urlencode

		// 构造post请求，json数据以请求body的形式传递
	}
	if method == http.MethodPost {
		reqBody, _ = json.Marshal(param)
	}
	seelog.Infof("method %s, uri %s, reqBody %s", method, uri, string(reqBody))
	req, err := http.NewRequest(method, uri, bytes.NewReader(reqBody))
	if err != nil {
		seelog.Errorf("NewRequest err %s", err.Error())
		return err.Error()
	}

	// 初始化响应
	client := &http.Client{Timeout: 3000 * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		seelog.Errorf("NewRequest err %s", err.Error())
		return err.Error()
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
