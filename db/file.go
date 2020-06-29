package db

import "BaiduNetdisk_golang/db/mysql"

func OnFileUploadFinish(fileHash, fileName, fileAddr string, fileSize int64) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil {
		println("插入文件数据失败,err=" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		println("数据库语句执行失败,err=" + err.Error())
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			println("存在相同的hash=", fileHash)
		}
		return true
	}
	return false
}
