package tools

import (
	"testing"
)

type EmailQueryReq struct {
	Email string `json:"email" binding:"required"`
}

type UserInfoQueryData struct {
	Email  string `json:"email"`
	AppId  string `json:"appId"`
	UserId int    `json:"userId"`
}

func TestHttpTrigger_HttpGet(t *testing.T) {
	var resp UserInfoQueryData
	trigger := HttpTrigger{
		Client: client,
	}
	//如果要加入GET后缀，就加入下面的设置
	m := make(map[string]string)
	m["test"] = "444"
	trigger.QueryDic = m

	err := trigger.HttpGet("http://127.0.0.1:30303/ump/user_info/query/email", &resp)
	if err != nil {
		return
	}
}

func TestHttpTrigger_HttpPost(t *testing.T) {
	var body EmailQueryReq
	var resp UserInfoQueryData
	trigger := HttpTrigger{
		Client: client,
	}
	//如果要自定义请求头，就加入下面的设置
	m := make(map[string]string)
	m["test"] = "444"
	trigger.HeaderDic = m

	err := trigger.HttpPost("http://127.0.0.1:30303/ump/user_info/change/email", &body, &resp)
	if err != nil {
		return
	}
}
