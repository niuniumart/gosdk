package tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/niuniumart/gosdk/requestid"
	"hash/fnv"
	"strconv"
	"strings"
	"time"
)

var doubleNameDic map[string]bool

// func init
func init() {
	doubleNameDic = map[string]bool{
		"欧阳": true, "太史": true, "端木": true, "上官": true, "司马": true,
		"东方": true, "独孤": true, "南宫": true, "万俟": true, "闻人": true,
		"夏侯": true, "诸葛": true, "尉迟": true, "公羊": true, "赫连": true,
		"澹台": true, "皇甫": true, "宗政": true, "濮阳": true, "公冶": true,
		"太叔": true, "申屠": true, "公孙": true, "慕容": true, "仲孙": true,
		"钟离": true, "长孙": true, "宇文": true, "司徒": true, "鲜于": true,
		"司空": true, "闾丘": true, "子车": true, "亓官": true, "司寇": true,
		"巫马": true, "公西": true, "颛孙": true, "壤驷": true, "公良": true,
		"漆雕": true, "乐正": true, "宰父": true, "谷梁": true, "拓跋": true,
		"夹谷": true, "轩辕": true, "令狐": true, "段干": true, "百里": true,
		"呼延": true, "东郭": true, "南门": true, "羊舌": true, "微生": true,
		"公户": true, "公玉": true, "公仪": true, "梁丘": true, "公仲": true,
		"公上": true, "公门": true, "公山": true, "公坚": true, "左丘": true,
		"公伯": true, "西门": true, "公祖": true, "第五": true, "公乘": true,
		"贯丘": true, "公皙": true, "南荣": true, "东里": true, "东宫": true,
		"仲长": true, "子书": true, "子桑": true, "即墨": true, "达奚": true,
		"褚师": true,
	}
}

//GetFmtStr func
func GetFmtStr(data interface{}) string {
	resp, _ := json.Marshal(data)
	respStr := string(resp)
	if respStr == "" {
		respStr = fmt.Sprintf("%+v", data)
	}
	return respStr
}

//InList func
func InList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

//IdCardWithoutMonthDay id card cut month and day
func IdCardWithoutMonthDay(idCard string) string {
	idCardNew := idCard[:10] + "****" + idCard[14:]
	return idCardNew
}

//IdCardMask Description: 只留前3后2，隐藏13位 idCard 身份证号 return:  string 隐藏后的身份证号
func IdCardMask(idCard string) string {
	if len(idCard) < 3 {
		return idCard
	} else if len(idCard) < 17 {
		return idCard[:3] + "**************"
	} else {
		return idCard[:3] + "**************" + idCard[len(idCard)-2:]
	}
}

// PersonalNameMask @Description: 只保留姓氏，复姓两字，name 姓名 return: string 隐藏后的姓名
func PersonalNameMask(name string) string {
	chName := []rune(name)
	if len(chName) < 3 {
		return string(chName[:1]) + "**"
	} else {
		// 大于三字判断是否包含复姓
		if doubleNameDic[string(chName[:2])] {
			return string(chName[:2]) + "**"
		} else {
			return string(chName[:1]) + "**"
		}
	}
}

//MailMask Description: @前面小于3位全显示+3位*号，大于3位只取前3位+3位*号 mailAddress 邮件地址
func MailMask(mailAddress string) string {
	splitIndex := strings.Index(mailAddress, "@")
	stringLength := len(mailAddress)
	// 如果存在字符@则对前缀处理，否则对全字符串处理
	if splitIndex > 0 {
		stringLength = splitIndex
	}
	if stringLength > 3 {
		return mailAddress[:3] + "***" + mailAddress[stringLength:]
	} else {
		return mailAddress[:stringLength] + "***" + mailAddress[stringLength:]
	}
}

// CharacterStringLength 字符串字符长度，即支持中英文混用 		i"原字符串"
func CharacterStringLength(str string) int {
	temp := []rune(str)
	return len(temp)
}

// CharacterSubString 截取字符串，下标为字符计数，即支持中英文混用 str:"原字符串" begin:"起始位置" end:"终止位置+1"
func CharacterSubString(str string, begin int, end int) string {
	if begin < 0 || begin >= end {
		return ""
	}

	temp := []rune(str)
	length := len(temp)
	if end > length {
		return ""
	}
	return string(temp[begin:end])
}

//GetBusinessNumber  func getBusinessNumber
func GetBusinessNumber(bType int) string {
	return fmt.Sprintf("%d%s%d", time.Now().UnixNano()/1000, GetRandNo(8), bType)
}

// GetRandNo  func getRandNo
func GetRandNo(length int) string {
	s := fmt.Sprintf("%.8d", requestid.Goid())
	return s[len(s)-8:]
}

// GetMd5 func getMd5
func GetMd5(b []byte) string {
	//给哈希算法添加数据
	res := md5.Sum(b)                    //返回值：[Size]byte 数组
	result := hex.EncodeToString(res[:]) //对应的参数为：切片，需要将数组转换为切片。
	return result
}

// ModStr sum32 str, and mod
func ModStr(s string, len uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))

	return h.Sum32() % len
}

// HashStrAndConvertToInt func HashStrAndConvertToInt
func HashStrAndConvertToInt(str string, lenth int) int {
	hashStr := GetMd5([]byte(str))
	hashStr = hashStr[len(hashStr)-lenth:]
	index, err := strconv.ParseInt(hashStr, 16, 32)
	if err != nil {
		martlog.Errorf("parse int error %s", err.Error())
		return 0
	}
	return int(index) + 1
}
