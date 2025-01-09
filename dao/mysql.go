package dao

import (
	sqlx "database/sql"
	"fmt"
	"time"
)

type MysqlConfig struct {
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Pwd           string `json:"Password"`
	Database      string `json:"database"`
	MaxConn       int    `json:"maxConn"`
	MaxIdle       int    `json:"maxIdle"`
	MaxConnSecond int64  `json:"maxConnSecond"`
}

func OpenMysql(conf *MysqlConfig) (db *sqlx.DB, err error) {
	if conf == nil {
		panic("invalid empty mysql config")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", conf.User, conf.Pwd, conf.IP, conf.Port, conf.Database)

	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		panic("打开数据库失败,err:%v\n", err)
	}

	//尝试连接数据库，Ping方法可检查数据源名称是否合法,账号密码是否正确。
	err = db.Ping()
	if err != nil {
		panic("连接数据库失败,err:%v\n", err)
		return
	}

	db.SetMaxOpenConns(conf.MaxConn) // 设置最大连接数
	db.SetMaxIdleConns(conf.MaxIdle) // 设置最大连接数
	db.SetConnMaxLifetime(time.Duration(conf.MaxConnSecond * int64(time.Second)))
	return
}
