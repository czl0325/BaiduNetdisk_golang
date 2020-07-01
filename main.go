package main

import (
	"BaiduNetdisk_golang/handler"
	"net/http"
)

func main() {
	// 启动静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 接口
	http.HandleFunc("/", handler.OnHomeHandle)

	http.HandleFunc("/file/upload", handler.UploadHandle)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandle)
	http.HandleFunc("/file/get", handler.GetFileMetaHandle)

	// 用户
	http.HandleFunc("/user/signup", handler.UserSignUpHandle)
	http.HandleFunc("/user/login", handler.UserLoginHandle)
	http.HandleFunc("/user/get", handler.UserInfoHandle)
	http.HandleFunc("/user/files", handler.FileQueryHandle)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		println("服务器启动失败!err=" + err.Error())
	}
}
