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
	"strconv"
	"time"
)

const COOKIE_MAX_AGE = time.Hour * 24 / time.Second

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
		uid, _ := r.Cookie("uid")
		println("打印=%v", uid)
		if uid == nil {
			println("注册用户才能上传文件")
			w.WriteHeader(http.StatusInternalServerError)
			http.Redirect(w, r, "/file/upload", http.StatusFound)
			return
		}

		uid2, _ := strconv.ParseInt(uid.Value, 10, 64)
		user, err := meta.GetUserMetaByIdDB(uid2)
		if err != nil {
			println("上传文件uid错误，无此用户")
			http.Redirect(w, r, "/file/upload", http.StatusFound)
			return
		}

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

		ret := db.OnUserFileUploadFinished(user.Id, user.UserName, fileMeta.FileSha1, fileMeta.FileName, fileMeta.Location, fileMeta.FileSize)
		if ret == true {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/file/upload/success", http.StatusFound)
		}
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

func UserLoginHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/html/user_login.html")
		if err != nil {
			println("user_login.html读取失败!err=" + err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(data)
	} else {
		r.ParseForm()

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if username == "" || password == "" {
			http.Redirect(w, r, "/user/login", http.StatusFound)
		}

		user, err := meta.GetUserMetaDB(username, password)
		if err != nil {
			http.Redirect(w, r, "/user/login", http.StatusFound)
		} else {
			maxAge := int(COOKIE_MAX_AGE)
			cookie := &http.Cookie{
				Name:     "uid",
				Value:    strconv.FormatInt(user.Id, 10),
				Path:     "/",
				HttpOnly: false,
				MaxAge:   maxAge,
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func UserSignUpHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/html/user_signup.html")
		if err != nil {
			println("user_signup.html读取失败!err=" + err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(data)
	} else {
		r.ParseForm()

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if username == "" || password == "" {
			http.Redirect(w, r, "/user/signup", http.StatusFound)
		}

		ret := db.OnSignUpHandle(username, password)
		if ret == false {
			http.Redirect(w, r, "/user/signup", http.StatusFound)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func UserInfoHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res := db.BaseResponse{
		Code:    0,
		Message: "成功",
		Data:    nil,
	}

	defer func() {
		w.Header().Set("content-type","text/json")
		data, _ := json.Marshal(res)
		w.Write(data)
	}()

	id := r.Form.Get("id")
	if id == "" {
		res.Code = 500
		res.Message = "缺少必要参数"
		return
	}
	id2, _ := strconv.ParseInt(id, 10, 64)
	user, err := meta.GetUserMetaByIdDB(id2)
	if err != nil {
		res.Code = 500
		res.Message = err.Error()
		return
	}
	res.Data = user
}

func OnHomeHandle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./static/html/index.html")
	if err != nil {
		println("index.html读取失败!err=" + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(data)
}

func FileQueryHandle(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	res := db.BaseResponse{
		Code:    0,
		Message: "成功",
		Data:    nil,
	}

	defer func() {
		w.Header().Set("content-type","text/json")
		data, _ := json.Marshal(res)
		w.Write(data)
	}()

	limit, err := strconv.ParseInt(r.Form.Get("limit"), 10, 64)
	if err != nil {
		limit = 20
	}
	uid, err := strconv.ParseInt(r.Form.Get("uid"), 10, 64)
	if err != nil {
		println("获取uid失败,err="+err.Error())
		res.Code = 500
		res.Message = "获取uid失败,err="+err.Error()
		return
	}
	userFiles, err := db.QueryUserFileMetas(uid, limit)
	if err != nil {
		println("查询用户文件失败,err="+err.Error())
		res.Code = 500
		res.Message = "查询用户文件失败,err="+err.Error()
		return
	}
	res.Data = userFiles
}