// Package gormcli 用于管理mysql连接
package gormcli

import (
	"fmt"
	"strings"
	// load mysql enum
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	DefaultMaxConn      = 12000
	DefaultMaxIdleConn  = 7000
	DefaultIdleTimeout  = 10 // use time.Second
	DefaultReadTimeout  = 10
	DefaultWriteTimeout = 10
)

// GormLogger Gorm用来打日志的结构体
type GormLogger struct {
}

var (
	Factory GormFactory
)

// GormFactory 用来生成Gorm指针的工厂
type GormFactory struct {
	MaxIdleConn  int
	MaxConn      int
	IdleTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

const (
	GormDuplicateErrKey = "Duplicate entry"
)

// judge err dup
func IsDupErr(err error) bool {
	return strings.Contains(err.Error(), GormDuplicateErrKey)
}

//	url like this :  "127.0.0.1:3306"
func (p *GormFactory) CreateTBassGorm(user, pwd, url, database string) (*gorm.DB, error) {
	if p.ReadTimeout == 0 {
		p.ReadTimeout = DefaultReadTimeout
	}
	if p.WriteTimeout == 0 {
		p.WriteTimeout = DefaultWriteTimeout
	}
	auth := user + ":" + pwd
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s&readTimeout=%ds&writeTimeout=%ds",
			auth, url, database, p.ReadTimeout, p.WriteTimeout))
	if err != nil {
		return db, err
	}
	logger := &GormLogger{}
	db.LogMode(true)
	db.SetLogger(logger)
	maxIdleConn := p.MaxIdleConn
	if maxIdleConn == 0 {
		maxIdleConn = DefaultMaxIdleConn
	}
	maxConn := p.MaxConn
	if maxConn == 0 {
		maxConn = DefaultMaxConn
	}
	idleTimeout := p.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = DefaultIdleTimeout
	}
	db.DB().SetMaxIdleConns(maxIdleConn)
	db.DB().SetMaxOpenConns(maxConn)
	db.DB().SetConnMaxLifetime(time.Duration(idleTimeout) * time.Second)
	return db, err
}

// 实现logger的print函数
func (logger *GormLogger) Print(values ...interface{}) {
	var (
		level = values[0]
	)
	if level == "sql" {
		martlog.Infof("%+v %s \"\"", values, level)
	} else {
		martlog.Infof("%+v", values)
	}
}
