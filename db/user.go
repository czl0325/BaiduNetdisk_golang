package db

import (
	"BaiduNetdisk_golang/db/mysql"
	"BaiduNetdisk_golang/util"
	"database/sql"
	"errors"
)

const pwd_salt = "#890"
type TableUser struct {
	Id sql.NullInt64
	UserName sql.NullString
	Password sql.NullString
	Phone sql.NullString
}

func OnLoginHandle(username, password string) (*TableUser, error) {
	stmt, err := mysql.DBConn().Prepare("select id, user_name, user_pwd, phone from tbl_user where user_name=? and status=1 limit 1")
	if err != nil {
		println("用户登录语句执行失败1,err="+err.Error())
		return nil, err
	}
	defer stmt.Close()
	tableUser := TableUser{}
	err = stmt.QueryRow(username).Scan(&tableUser.Id, &tableUser.UserName, &tableUser.Password, &tableUser.Phone)
	if err != nil {
		println("用户登录语句执行失败1,err="+err.Error())
		return nil, err
	}
	encPassword := util.Sha1([]byte(password + pwd_salt))
	if encPassword != tableUser.Password.String {
		return nil, errors.New("密码错误")
	}
	return &tableUser, nil
}

func OnSignUpHandle(username, password string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		println("注册用户数据库语句创建失败,err=" + err.Error())
		return false
	}
	defer stmt.Close()

	encPassword := util.Sha1([]byte(password + pwd_salt))

	ret, err := stmt.Exec(username, encPassword)
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

func OnGetUserHandle(id int64) (*TableUser, error) {
	stmt, err := mysql.DBConn().Prepare("select id, user_name, phone from tbl_user where id=? and status=1 limit 1")
	if err != nil {
		println("获取用户语句执行失败1,err="+err.Error())
		return nil, err
	}
	defer stmt.Close()
	tableUser := TableUser{}
	err = stmt.QueryRow(id).Scan(&tableUser.Id, &tableUser.UserName, &tableUser.Phone)
	if err != nil {
		println("获取用户语句执行失败2,err="+err.Error())
		return nil, err
	}
	return &tableUser, nil
}


