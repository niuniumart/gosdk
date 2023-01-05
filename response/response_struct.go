package response

import "strconv"

//RetBase struct
type RetBase struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

//Error func error
func (r *RetBase) Error() string {
	return strconv.Itoa(r.RetCode) + "-" + r.RetMsg
}

//RetData struct ret data
type RetData struct {
	RetBase
	Data interface{} `json:"data"`
}

//NewRetData 构建一般响应，通常情况下，非成功时data都会为nil
func NewRetData(err error, data interface{}) *RetData {
	var resp *RetBase
	resp, ok := err.(*RetBase)
	if !ok {
		resp = RESP_FAIL
	}

	return &RetData{
		RetBase: *resp,
		Data:    data,
	}
}

//BuildSuccResp 构建成功响应
func BuildSuccResp(data interface{}) *RetData {
	return NewRetData(RESP_SUCC, data)
}

//BuildFailResp 构建异常/未知的错误响应，通常情况下，非成功时data都会为nil
func BuildFailResp(data interface{}) *RetData {
	return NewRetData(RESP_FAIL, data)
}
