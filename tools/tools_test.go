package tools

import (
	"encoding/base64"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	"io/ioutil"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTimeTool(t *testing.T) {
	t1 := "2019-01-08 13:50:30"
	tt, err := GetTimeFromStr(t1)
	fmt.Println("format 2019-01-08 13:50:30: ", tt, err)

	t2 := "2019-01-08"
	tt, err = GetTimeFromDayStr(t2)
	fmt.Println("format 2019-01-08: ", tt, err)

	t3 := "2019/01/08 13:50:30"
	tt, err = GetTimeFromSpritStr(t3)
	fmt.Println("format 2019/01/08 13:50:30: ", tt, err)
}

func AppendAry(ary *[]int) {
	*ary = append(*ary, 9)
}

func TestArray(t *testing.T) {
	//	var ary = []int{5,6,7}
	ary := make([]int, 0)
	ary = append(ary, 8)
	fmt.Println(ary)
	AppendAry(&ary)
	fmt.Println(ary)
}

func TestQr(t *testing.T) {
	//err := qrcode.WriteFile("http://www.baidu.com/",qrcode.Medium,1024,"./blog_qrcode1024.png")
	content, err := qrcode.Encode("http://www.baidu.com/", qrcode.Medium, 512)
	fmt.Println(err)
	err = ioutil.WriteFile("./blog_qrcode512.png", content, 0644)
	fmt.Println(err)
}

func TestIdCardReplace(t *testing.T) {
	idCard := "510681111302266117"
	s := IdCardWithoutMonthDay(idCard)
	fmt.Println(s)
}

func TestIdCardMask(t *testing.T) {
	fmt.Println(IdCardMask("4413811199210315555"))
	fmt.Println(IdCardMask("123"))
	fmt.Println(IdCardMask("12"))
}

func TestMailMask(t *testing.T) {
	fmt.Println(MailMask("su@21cn.com"))
	fmt.Println(MailMask("ray@qq.com"))
	fmt.Println(MailMask("ray123@qq.com"))
}

func TestBase64(t *testing.T) {
	str := "abcd"
	b := base64.StdEncoding.EncodeToString([]byte(str))
	n, err := base64.StdEncoding.DecodeString(b)
	fmt.Println(err)
	fmt.Println(string(n))

}

func putMap(res gin.H) {
	res["b"] = 3
}

func TestMap(t *testing.T) {
	var res = gin.H{}
	res["a"] = 2
	putMap(res)
	fmt.Println(res)
}

func TestSlice(t *testing.T) {
	randNumStr := fmt.Sprintf("%d", time.Now().Unix())
	fmt.Println(randNumStr)
	agentCode := randNumStr[len(randNumStr)-8:]
	fmt.Println(agentCode)
}

func TestDicToQueryString(t *testing.T) {
	baseUrl, err := url.Parse("http://www.baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	var queryStrDic = map[string]string{
		"searchId": "12345",
	}
	params := url.Values{}
	for k, v := range queryStrDic {
		params.Add(k, v)
	}
	baseUrl.RawQuery = params.Encode()
	fmt.Println(baseUrl)

}

func TestTimeCost(t *testing.T) {
	crt := time.Now()
	time.Sleep(2 * time.Second)
	d := time.Since(crt)
	dd := d.Milliseconds()
	fmt.Println(dd)
	if int(dd) > 2000 {
		fmt.Println("warning")
	}
	fmt.Printf("%v", time.Since(crt))
}

func TestRandomInt(t *testing.T) {

	rand.Seed(time.Now().Unix())
	for i := 0; i < 50; i++ {
		pos := rand.Int63n(2)
		//		pos := rand.Intn(2)
		if pos > 0 {
			fmt.Println(true)
		}
		fmt.Println(pos)

	}
	return
	for i := 0; i < 10; i++ {
		base := 2
		fmt.Println("base ", base)
		fmt.Println("add ", 1)
		rand.Seed(time.Now().Unix())
		pos := rand.Intn(base) + 1
		fmt.Println("pos ", pos)
	}

}

func TestSwitch(t *testing.T) {
	var stage string
	stage = "stage_1"
	switch stage {
	case "stage_1":
		fmt.Println("stage 1")
		stage = "stage_3"
		fallthrough
	case "stage_2":
		fmt.Println("stage 2")
		fallthrough
	case "stage_3":
		fmt.Println("stage 3")
	}
}

func TestSm3(t *testing.T) {
	a := SM3("fabricgosdk")
	fmt.Println("a is ", a, len(a))

}

type T struct {
	name string `json:"dciId"`
	age  int
}

func (p *T) NewT(name string, age int) {

}

func TestStrIndex(t *testing.T) {
	str := "1231321_cat"
	idx := strings.Index(str, "_")
	fmt.Println(idx)
	fmt.Println(str[:idx])
	fmt.Println(str[idx+1:])
}

func TestStructToMap(t *testing.T) {
	var tt = &T{
		name: "cag",
		age:  32,
	}

	m := structs.Map(tt)
	fmt.Println(GetFmtStr(m))

}
func TestGetByUrlWithoutParams(t *testing.T) {
	content, err := GetByUrlWithoutParams("http://www.baidu.com")
	fmt.Println(content)
	fmt.Println(err)
}

func TestRuneString(t *testing.T) {
	dciName := "一二三四五六七八九十一二三四五六七八九ab靠112"
	nameRune := []rune(dciName)
	if len(nameRune) > 20 {
		nameRune = nameRune[:20]
		dciName = string(nameRune)
		dciName = dciName + "..."
	}
	fmt.Println(dciName)
}

func TestMsStampToS(t *testing.T) {
	var a int64
	a = 1616728416000
	a = a / 1000
	tm := time.Unix(a, 0)
	fmt.Println(tm)
}

func TestSlice2(t *testing.T) {
	a := "cataaaaabccc"
	c := Min(len(a), 10)
	b := a[:c]
	fmt.Println(b)
}

func TestMd5(t *testing.T) {
	a := GetMd5([]byte("hahaha"))
	fmt.Println(a, len(a))
}

func TestHashToInt(t *testing.T) {
	hash := "2761a234a2d13584f069ecdb47e315029a828792249172a4b514cff82b226d00"
	hash = hash[len(hash)-2:]
	fmt.Println(hash)
	index, err := strconv.ParseInt(hash, 16, 32)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("table index: %d\n", index)
}

func TestHashStrAndModToInt(t *testing.T) {
	fmt.Println(HashStrAndConvertToInt("haha", 3))
}

/*
// for support gorm batch insert
func BatchInsert(db *gorm.DB, table string, data interface{}) error {
	sql := genBatchInsertSql(table, data)
	e := db.Debug().Exec(sql).Error
	if e != nil {
		return e
	}
	return nil
}

func genBatchInsertSql(table string, data interface{}) string {
	var fields, values string
	j := 0
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)

		for i := 0; i < s.Len(); i++ {
			f, v := genFieldSql(table, s.Index(i).Interface())
			if i == 0 {
				fields = f
			}
			if j == 0 {
				values += v
				j++
				continue
			}
			values += "," + v
		}
	}
	str := fields + " values " + values
	return str
}


func genFieldSql(table string, intf interface{}) (string, string) {
	var str string
	str = fmt.Sprintf("insert into %s(", table)
	value := "("
	rt := reflect.TypeOf(intf)
	rv := reflect.ValueOf(intf)
	for i, j := 0, rt.NumField(); i < j; i++ {
		rtf := rt.Field(i)
		rvf := rv.Field(i)
		switch rtf.Type.Name() {
		case "uint64":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(uint64))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(uint64))
		case "uint32":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(uint32))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(uint32))
		case "uint16":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(uint16))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(uint16))
		case "uint8":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(uint8))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(uint8))
		case "int":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(int))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(int))
		case "int8":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(int8))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(int8))
		case "int16":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(int16))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(int16))
		case "int32":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(int32))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(int32))
		case "int64":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += fmt.Sprintf("%d", rvf.Interface().(int64))
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += "," + fmt.Sprintf("%d", rvf.Interface().(int64))
		case "string":
			if i == 0 {
				str += fieldFormat(rtf.Name)
				value += "'" + rvf.String() + "'"
				continue
			}
			str += "," + fieldFormat(rtf.Name)
			value += ",'" + rvf.String() + "'"
		}
	}
	str += ")"
	value += ")"
	return str, value
}

func fieldFormat(filed string) string {
	var list []byte
	n := len(filed) * 2
	list = make([]byte, 0, n)
	for i := 0; i < len(filed); i++ {
		if unicode.IsUpper(rune(filed[i])) && i != 0 {
			list = append(list, '_')
		}
		list = append(list, filed[i])
	}
	str := string(list)
	str = strings.ToLower(str)
	return str
}*/

func getRuneLen(str string) int {
	return len([]rune(str))
}

func TestDciName(t *testing.T) {
	nameList := []string{
		"张三啊啊啊", "李四呀呀呀", "王二汪汪汪", "赵五嗷嗷嗷", "刘六mia",
	}
	str := getNameStr(nameList, maxRighterShowNum, maxRighterMsgLength)
	fmt.Println(str)
}

const (
	maxRighterMsgLength = 26
	maxRighterShowNum   = 5
)

func getNameStr(nameList []string, maxNum, maxLength int) string {
	var nameMsg string
	for i, name := range nameList {
		if i == maxNum {
			break
		}
		if i == 0 {
			tmpMsg := name
			if getRuneLen(nameMsg)+getRuneLen(tmpMsg) > maxLength {
				nameMsg = tmpMsg
				tmpNameMsg := []rune(tmpMsg)
				tmpNameMsg = tmpNameMsg[:maxLength]
				nameMsg = string(tmpNameMsg) + "..."
				break
			}
			nameMsg += tmpMsg
			continue
		}
		tmpMsg := "、" + name
		if getRuneLen(nameMsg)+getRuneLen(tmpMsg) > maxLength {
			nameMsg += " ..."
			break
		}
		nameMsg += tmpMsg
	}
	return nameMsg
}

func TestPower(t *testing.T) {
	a := 1 << 2
	fmt.Println(a)
}

func TestRandDuration(t *testing.T) {
	for i := 0; i < 100; i++ {

		internelTime := 1 * time.Second
		step := RandNum(500)
		internelTime += time.Duration(step) * time.Millisecond
		fmt.Println(internelTime)
	}
	fmt.Printf("%v\n", true)
}

func RandNum(num int64) int64 {
	step := rand.Int63n(num) + int64(1)
	flag := rand.Int63n(2)
	if flag == 0 {
		return -step
	}
	return step
}

func TestUuid(t *testing.T) {
	requestIDStr := fmt.Sprintf("%+v", uuid.New())
	fmt.Println(requestIDStr)
}
