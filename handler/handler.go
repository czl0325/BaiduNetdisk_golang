package handler

import (
	"BaiduNetdisk_golang/db"
	"BaiduNetdisk_golang/meta"
	"BaiduNetdisk_golang/util"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/html/file_upload.html")
		if err != nil {
			println("file_upload.html读取失败!err=" + err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(data)
	} else {
		file, header, err := r.FormFile("file")
		if err != nil {
			println("读取文件失败,err=" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		path := "./temp/" + header.Filename
		newFile, err := os.Create(path)
		if err != nil {
			println("服务端创建文件失败,err=" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer newFile.Close()

		var fileMeta = meta.FileMeta{
			FileName: header.Filename,
			Location: path,
		}

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			println("文件写入失败,err=" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMetaDB(fileMeta)

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

func UploadSuccessHandle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "文件上传成功!")
}

func GetFileMetaHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := db.BaseResponse{
		Code:    0,
		Message: "成功",
		Data:    nil,
	}

	fileHash := r.Form.Get("hash")
	fMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		res.Code = 404
		res.Message = err.Error()
		data, _ := json.Marshal(res)
		w.Write(data)
		return
	}
	res.Data = fMeta
	data, _ := json.Marshal(res)
	w.Write(data)
}
