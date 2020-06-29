package db

import (
	"BaiduNetdisk_golang/db/mysql"
	"database/sql"
)

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

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}
func GetFileMeta(fileHash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil {
		println("查询文件语句执行失败1,err=" + err.Error())
		return nil, err
	}
	defer stmt.Close()
	tFile := TableFile{}
	stmt.QueryRow(fileHash).Scan(&tFile.FileHash, &tFile.FileName, &tFile.FileSize, &tFile.FileAddr)
	if err != nil {
		println("查询文件语句执行失败2,err=" + err.Error())
		return nil, err
	}
	return &tFile, nil
}