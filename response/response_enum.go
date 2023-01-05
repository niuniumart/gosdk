// Package response base info
package response

var (
	RESP_SUCC *RetBase
	RESP_FAIL *RetBase

	RESP_PARAMS_INVALID        *RetBase
	RESP_HTTP_REQ_ERROR        *RetBase
	RESP_BUSINESS_REQ_ERROR    *RetBase
	RESP_STRUCT_COPY_REQ_ERROR *RetBase

	RESP_DB_ERROR                  *RetBase
	RESP_DB_SELECT_ERROR           *RetBase
	RESP_DB_UPDATE_ERROR           *RetBase
	RESP_DB_INSERT_ERROR           *RetBase
	RESP_DB_DELETE_ERROR           *RetBase
	RESP_DB_RECORD_NOT_FOUND_ERROR *RetBase
	RESP_DB_RECORD_EXIST_ERROR     *RetBase

	RESP_JSON_MARSHAL_ERROR   *RetBase
	RESP_JSON_UNMARSHAL_ERROR *RetBase

	RESP_DECRYPT_ERROR *RetBase
	RESP_ENCRYPT_ERROR *RetBase

	RESP_REDIS_GET_ERROR    *RetBase
	RESP_REDIS_SET_ERROR    *RetBase
	RESP_REDIS_TTL_ERROR    *RetBase
	RESP_REDIS_EXPIRE_ERROR *RetBase
	RESP_REDIS_DELETE_ERROR *RetBase
	RESP_REDIS_SCRIPT_ERROR *RetBase

	RESP_NFT_PUBKEY_FORMAT_ERROR      *RetBase
	RESP_NFT_SIGN_FORMAT_ERROR        *RetBase
	RESP_NFT_ADDRESS_FORMAT_ERROR     *RetBase
	RESP_NFT_SIGN_VERIFY_FORMAT_ERROR *RetBase
)

//0～1000为公共响应码，为各服务共用。其余的响应码分段，各服务自行维护
func init() {
	RESP_SUCC = Build(0, "请求成功")
	RESP_FAIL = Build(0xFFFF, "请求失败")

	RESP_PARAMS_INVALID = Build(10, "请求参数无效")
	RESP_HTTP_REQ_ERROR = Build(11, "Http请求失败")
	RESP_BUSINESS_REQ_ERROR = Build(12, "Http请求业务错误")
	RESP_STRUCT_COPY_REQ_ERROR = Build(13, "结构体同名字段拷贝失败")

	RESP_DB_ERROR = Build(30, "数据库异常")
	RESP_DB_SELECT_ERROR = Build(31, "数据库读失败")
	RESP_DB_UPDATE_ERROR = Build(32, "数据库更新失败")
	RESP_DB_INSERT_ERROR = Build(33, "数据库插入失败")
	RESP_DB_DELETE_ERROR = Build(34, "数据库删除失败")
	RESP_DB_RECORD_NOT_FOUND_ERROR = Build(35, "记录不存在")
	RESP_DB_RECORD_EXIST_ERROR = Build(36, "记录已存在")

	RESP_JSON_MARSHAL_ERROR = Build(38, "json序列化失败")
	RESP_JSON_UNMARSHAL_ERROR = Build(39, "json反序列化失败")

	RESP_DECRYPT_ERROR = Build(40, "解密失败")
	RESP_ENCRYPT_ERROR = Build(41, "加密失败")

	RESP_REDIS_GET_ERROR = Build(45, "redis执行GET失败")
	RESP_REDIS_SET_ERROR = Build(46, "redis执行SET失败")
	RESP_REDIS_TTL_ERROR = Build(47, "redis执行TTL失败")
	RESP_REDIS_EXPIRE_ERROR = Build(48, "redis执行EXPIRE失败")
	RESP_REDIS_DELETE_ERROR = Build(49, "redis执行DEL失败")
	RESP_REDIS_SCRIPT_ERROR = Build(50, "redis脚本错误")

	RESP_NFT_PUBKEY_FORMAT_ERROR = Build(55, "NFT-公钥格式错误")
	RESP_NFT_SIGN_FORMAT_ERROR = Build(56, "NFT-签名格式错误")
	RESP_NFT_ADDRESS_FORMAT_ERROR = Build(57, "NFT-地址格式错误")
	RESP_NFT_SIGN_VERIFY_FORMAT_ERROR = Build(58, "NFT-验签错误")
}

// Build func for build ret struct
func Build(retCode int, retMsg string) *RetBase {
	return &RetBase{
		RetCode: retCode,
		RetMsg:  retMsg,
	}
}
