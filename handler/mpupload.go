package handler

import (
	"BaiduNetdisk_golang/cache/redis"
	"BaiduNetdisk_golang/db"
	"BaiduNetdisk_golang/meta"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadId   string
	ChunkSize  int //分块大小
	ChunkCount int //分多少块
}

// 初始化分块上传
func initialMultipartUploadHandle(w http.ResponseWriter, r *http.Request) {
	//1. 解析用户请求参数
	uid, _ := r.Cookie("uid")
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

	r.ParseForm()

	fileHash := r.Form.Get("hash")
	fileSize, _ := strconv.Atoi(r.Form.Get("size"))

	// 2.获取一个redis连接池
	conn := redis.RedisPool().Get()
	defer conn.Close()

	// 3.生成分块上传的初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash:   fileHash,
		FileSize:   fileSize,
		UploadId:   user.UserName + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, //5MB
		ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))),
	}

	conn.Do("HSET", "MP_"+uploadInfo.UploadId, "FileHash", uploadInfo.FileHash)
	conn.Do("HSET", "MP_"+uploadInfo.UploadId, "fileSize", uploadInfo.FileSize)
	conn.Do("HSET", "MP_"+uploadInfo.UploadId, "ChunkCount", uploadInfo.ChunkCount)

	res := db.BaseResponse{
		Code:    0,
		Message: "成功",
		Data:    uploadInfo,
	}

	w.Header().Set("content-type", "text/json")
	data, _ := json.Marshal(res)
	w.Write(data)
}

// 上传文件分块
func UploadPartHandle(w http.ResponseWriter, r *http.Request)  {
	res := db.BaseResponse{
		Code:    0,
		Message: "成功",
		Data:    nil,
	}

	defer func() {
		w.Header().Set("content-type", "text/json")
		data, _ := json.Marshal(res)
		w.Write(data)
	}()

	//1. 解析用户请求参数
	//uid, _ := r.Cookie("uid")
	//if uid == nil {
	//	println("获取cookie的uid失败")
	//	res.Code = 500
	//	res.Message = "查无此用户"
	//	return
	//}
	//
	//uid2, _ := strconv.ParseInt(uid.Value, 10, 64)
	//user, err := meta.GetUserMetaByIdDB(uid2)
	//if err != nil {
	//	println("上传文件uid错误，无此用户")
	//	res.Code = 500
	//	res.Message = "查无此用户"
	//	return
	//}

	r.ParseForm()

	uploadId := r.Form.Get("upload_id")
	chunkIndex := r.Form.Get("index")

	// 2.获取redis连接池
	conn := redis.RedisPool().Get()
	defer conn.Close()

	// 3.获取文件句柄，用于存储分块内容
	fd, err := os.Create("./temp/data/" + uploadId + "/" + chunkIndex)
	if err != nil {
		println("创建文件分块错误,err=" + err.Error())
		res.Code = 500
		res.Message = "创建文件分块错误"
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4.更新redis缓存
	conn.Do("HSET", "MP_" + uploadId, "idx_" + chunkIndex, 1)

	// 5.返回数据给前端
	
}

// 通知上传合并接口
func CompleteUploadHandle(w http.ResponseWriter, r *http.Request)  {
	res := db.BaseResponse{
		Code:    0,
		Message: "成功",
		Data:    nil,
	}

	defer func() {
		w.Header().Set("content-type", "text/json")
		data, _ := json.Marshal(res)
		w.Write(data)
	}()
	//1. 解析用户请求参数
	uid, _ := r.Cookie("uid")
	if uid == nil {
		println("获取cookie的uid失败")
		res.Code = 500
		res.Message = "查无此用户"
		return
	}

	uid2, _ := strconv.ParseInt(uid.Value, 10, 64)
	user, err := meta.GetUserMetaByIdDB(uid2)
	if err != nil {
		println("上传文件uid错误，无此用户")
		res.Code = 500
		res.Message = "查无此用户"
		return
	}

	r.ParseForm()

	upload_id := r.Form.Get("upload_id")
	fileHash := r.Form.Get("hash")
	fileSize := r.Form.Get("size")
	fileName := r.Form.Get("name")



	// 2.获取redis连接池
	conn := redis.RedisPool().Get()
	defer conn.Close()

	// 3.查询是否所有分块都上传完成
	data, err := conn.Do("HGETALL", "MP_" + upload_id)
	if err != nil {
		println("redis查询失败,err="+err.Error())
		res.Code = 500
		res.Message = "redis查询失败"
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i+=2 {

	}

	// 4.分块合并

	// 5.更新唯一文件表和用户文件表

	// 6.返回结果给前端
}