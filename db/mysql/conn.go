package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/BaiduNetDisk?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		println("连接mysql数据库失败,err=" + err.Error())
		os.Exit(1)
	}
}

//DBConn : 返回数据库连接
func DBConn() *sql.DB {
	return db
}
