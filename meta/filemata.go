package meta

import "BaiduNetdisk_golang/db"

type FileMeta struct {
	FileSha1   string
	FileName   string
	FileSize   int64
	Location   string
	CreateTime string
}

func UpdateFileMetaDB(meta FileMeta) bool {
	return db.OnFileUploadFinish(meta.FileSha1, meta.FileName, meta.Location, meta.FileSize)
}

//GetFileMetaDB:从mysql获取文件元信息
func GetFileMetaDB(fileHash string) (FileMeta, error) {
	tFile, err := db.GetFileMeta(fileHash)
	if err != nil {
		return FileMeta{}, err
	}
	fMeta := FileMeta{
		FileSha1: tFile.FileHash,
		FileName: tFile.FileName.String,
		FileSize: tFile.FileSize.Int64,
		Location: tFile.FileAddr.String,
	}
	return fMeta, nil
}
