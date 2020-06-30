package db

import "BaiduNetdisk_golang/db/mysql"

type UserFile struct {
	UserId     int64
	UserName   string
	FileSha1   string
	FileName   string
	FileSize   int64
	FilePath   string
	UpdateTime string
}

//OnUserFileUploadFinished:插入用户文件表
func OnUserFileUploadFinished(uid int64, username, filesha1, filename, filepath string, filesize int64) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user_file (`user_id`, `user_name`, `file_sha1`, `file_name`, `file_size`, `file_path`) values (?,?,?,?,?,?)")
	if err != nil {
		println("插入用户文件表失败1,err=" + err.Error())
		return false
	}
	_, err = stmt.Exec(uid, username, filesha1, filename, filesha1, filepath)
	if err != nil {
		println("插入用户文件表失败2,err=" + err.Error())
		return false
	}
	return true
}
