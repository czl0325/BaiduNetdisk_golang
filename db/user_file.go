package db

import "BaiduNetdisk_golang/db/mysql"

type UserFile struct {
	Id         int64
	UserId     int64
	UserName   string
	FileSha1   string
	FileName   string
	FileSize   int64
	FilePath   string
	CreateTime string
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

func QueryUserFileMetas(uid, limit int64) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare("select id, user_id, user_name, file_sha1, file_size, file_name, file_path, create_time from tbl_user_file where user_id=? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(uid, limit)
	if err != nil {
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		uFile := UserFile{}
		err = rows.Scan(&uFile.Id, &uFile.UserId, &uFile.UserName, &uFile.FileSha1, &uFile.FileSize, &uFile.FileName, &uFile.FilePath, &uFile.CreateTime)
		if err != nil {
			println("查询文件失败,err=" + err.Error())
			break
		}
		userFiles = append(userFiles, uFile)
	}
	return userFiles, nil
}
