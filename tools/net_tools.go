package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/niuniumart/gosdk/response"
	"github.com/niuniumart/gosdk/seelog"
	"github.com/mitchellh/mapstructure"
)

var client *http.Client

// func init
func init() {
	client = &http.Client{Timeout: 60 * time.Second}
}

//HttpTrigger struct http trigger
type HttpTrigger struct {
	Method               string
	Url                  string
	From                 string
	To                   string
	DisableReqBodyPrint  bool
	DisableRespBodyPrint bool
	HeaderDic            map[string]string
	QueryDic             map[string]string
	Client               *http.Client
	QueryStr             string //如果不为空，则忽略QueryDic
}

//HttpOption func
type HttpOption func(*HttpTrigger)

//SendJsonRequest func sendJsonRequest
func (p *HttpTrigger) SendJsonRequest(client *http.Client, queryStrDic,
	headerDic map[string]string,
	body interface{}) ([]byte, error) {
	return SendRequest(client, queryStrDic, headerDic, p.Method, p.Url, body,
		p.From, p.To, p.DisableReqBodyPrint, p.DisableRespBodyPrint, p.QueryStr)
}

//Execute send and get response
func (p *HttpTrigger) Execute(body interface{}) (*response.RetData, error) {
	//不存在则设置默认json格式
	if p.HeaderDic == nil {
		p.HeaderDic = make(map[string]string)
	}
	if _, ok := p.HeaderDic["Content-Type"]; !ok {
		p.HeaderDic["Content-Type"] = "application/json;charset=utf-8"
	}

	respStr, err := p.SendJsonRequest(p.Client, p.QueryDic, p.HeaderDic, body)
	if err != nil {
		seelog.Errorf("SendJsonRequest error:%s", err.Error())
		return nil, response.RESP_HTTP_REQ_ERROR
	}

	var respData = &response.RetData{}
	err = json.Unmarshal(respStr, respData)
	if err != nil {
		seelog.Errorf("Unmarshal error:%s", err.Error())
		return nil, response.RESP_JSON_UNMARSHAL_ERROR
	}

	return respData, nil
}

//HttpGet JSON Get请求,需要服务端的响应体格式符合RetData
//resp必须是结构体的指针，字段是响应Data的内容
//响应报文JSON的key，需要和respData中的字段名相同(忽略大小写)，所以key最好是驼峰命名。否则需要在respData的字段标签中用指定名称，例如`mapstructure:"user_name"`
func (p *HttpTrigger) HttpGet(url string, retData interface{}) error {
	p.Method = http.MethodGet
	p.Url = url

	data, err := p.Execute(nil)
	if err != nil {
		return err
	}
	if data.RetCode != response.RESP_SUCC.RetCode {
		err := response.Build(data.RetCode, data.RetMsg)
		return err
	}

	if retData != nil {
		err = mapstructure.Decode(data.Data, retData)
		if err != nil {
			seelog.Errorf("mapstructure Decode error:%s", err.Error())
			return response.RESP_JSON_UNMARSHAL_ERROR
		}
	}
	return nil
}

//HttpPost JSON Post请求,需要服务端的响应体格式符合RetData
// body和resp必须是结构体的指针，resp字段是响应Data的内容
// 响应报文JSON的key，需要和respData中的字段名相同(忽略大小写)，
// 所以key最好是驼峰命名。否则需要在respData的字段标签中用指定名称，例如`mapstructure:"user_name"`
func (p *HttpTrigger) HttpPost(url string, body interface{}, retData interface{}) error {
	p.Method = http.MethodPost
	p.Url = url

	data, err := p.Execute(body)
	if err != nil {
		return err
	}
	if data.RetCode != response.RESP_SUCC.RetCode {
		err := response.Build(data.RetCode, data.RetMsg)
		return err
	}

	if retData != nil {
		err = mapstructure.Decode(data.Data, retData)
		if err != nil {
			seelog.Errorf("mapstructure Decode error:%s", err.Error())
			return response.RESP_JSON_UNMARSHAL_ERROR
		}
	}
	return nil
}

//SendRequest common func, send request
func SendRequest(client *http.Client, queryStrDic, headerDic map[string]string,
	method, reqUrl string, body interface{}, from, to string,
	disableReqBodyPrint, disableRespPrint bool, queryStr string) ([]byte, error) {
	// step 1: put query string into url
	baseUrl, err := url.Parse(reqUrl)
	if err != nil {
		seelog.Errorf("url.Parse err %s", err.Error())
		return nil, err
	}
	if queryStr != "" {
		//处理get请求直接透传的场景
		baseUrl.RawQuery = queryStr
	} else {
		params := url.Values{}
		for k, v := range queryStrDic {
			params.Add(k, v)
		}
		baseUrl.RawQuery = params.Encode()
	}
	reqUrl = baseUrl.String()
	seelog.Infof("call %s with reqUrl %s, method %s, from %s", to, reqUrl, method, from)
	if !disableReqBodyPrint {
		seelog.Infof("call %s with reqUrl %s, body is %v", to, reqUrl, body)
	}
	var reader io.Reader
	if method != http.MethodGet && body != nil {
		if bodyReader, ok := body.(io.Reader); ok {
			//处理post请求直接透传的场景
			reader = bodyReader
		} else {
			b, err := json.Marshal(body)
			if err != nil {
				seelog.Errorf("json.Marshal err %s", err.Error())
				return nil, err
			}
			seelog.Infof("post content %s", string(b))
			reqBody := bytes.NewBuffer(b)
			reader = reqBody
		}
	}

	req, err := http.NewRequest(method, reqUrl, reader)
	if err != nil {
		seelog.Errorf("http.NewRequest err %s, reqUrl %s", err.Error(), reqUrl)
		return nil, err
	}
	for k, v := range headerDic {
		//Host的请求头要特殊处理，否则设置不到header里去
		if k == "Host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}

	beginTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		seelog.Errorf("Client.Do err %s", err.Error())
		return nil, err
	}
	if resp == nil {
		errMsg := fmt.Sprintf("resp is nil, may be server can't reach")
		seelog.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		seelog.Errorf("ioutil.ReadAll err %s", err.Error())
		return nil, err
	}
	seelog.Infof("call %s with reqUrl %s success and cost %v", to, reqUrl, time.Since(beginTime))
	if !disableRespPrint {
		seelog.Infof("call %s with reqUrl %s resp %s", to, reqUrl, string(content))
	}
	defer resp.Body.Close()
	return content, nil
}

//GetByUrlWithoutParams 从url直接获取返回
//@Description: 从url直接获取返回
//@param url
//@return []byte
//@return error
func GetByUrlWithoutParams(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, err
}

// BuildRequestHeaderDic Description 把http请求头的数据，构造成HttpTrigger的HeaderDic。http请求头key对应的value列表中，只会取出第一个string
func BuildRequestHeaderDic(req *http.Request) map[string]string {
	retMap := make(map[string]string)

	for k, v := range req.Header {
		retMap[k] = v[0]
	}
	return retMap
}
