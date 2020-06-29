package db

import "BaiduNetdisk_golang/db/mysql"

func OnSignUpHandle(username, password string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user `user_name`,`user_pwd` values (?,?)")
	if err != nil {
		println("注册用户数据库语句创建失败,err=" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		println("注册用户数据库语句执行失败,err=" + err.Error())
		return false
	}
	if r, err := ret.RowsAffected(); err == nil {
		if r <= 0 {
			println("存在相同用户,username=" + username)
			return false
		}
		return true
	}
	return false
}

