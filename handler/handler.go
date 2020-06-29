package handler

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandle(w http.ResponseWriter, r *http.Request)  {
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

		newFile, err := os.Create("./temp/" + header.Filename)
		if err != nil {
			println("服务端创建文件失败,err=" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			println("文件写入失败,err=" + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

func UploadSuccessHandle(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "文件上传成功!")
}