package gormcli

import (
	"fmt"
	"testing"
	"time"
)

func TestGormCliHttp(t *testing.T) {
	url := "127.0.0.1"
	user := "root"
	pwd := ""
	dbName := "niuniumart"
	Factory.WriteTimeout = 2
	Factory.ReadTimeout = 2
	MysqlCli, err := Factory.CreateTBassGorm(user, pwd, url, dbName)
	if err != nil {
		fmt.Println(err)
		return
	}
	type A struct {
		EvidStatus int
	}
	fmt.Println(time.Now())
	var a A
	err = MysqlCli.Table(" t_ump_evidence_202103").
		Limit(1).
		Offset(300000).
		Find(&a).Error
	fmt.Println(time.Now())
	fmt.Println(a)
	fmt.Println(err)

	fmt.Println("second")
	err = MysqlCli.Table(" t_ump_evidence_202103").
		Limit(1).
		Offset(1).
		Find(&a).Error
	fmt.Println(time.Now())
	fmt.Println(a)
	fmt.Println(err)
}
