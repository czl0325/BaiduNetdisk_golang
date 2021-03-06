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
	"path"
	"strconv"
	"strings"
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
func InitialMultipartUploadHandle(w http.ResponseWriter, r *http.Request) {
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
	r.ParseForm()

	uid, _ := strconv.Atoi(r.Form.Get("uid"))
	user, err := meta.GetUserMetaByIdDB(int64(uid))
	if err != nil {
		println("上传文件uid错误，无此用户")
		res.Code = 500
		res.Message = "查无此用户"
		return
	}
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
	r.ParseForm()

	uploadId := r.Form.Get("upload_id")
	chunkIndex := r.Form.Get("index")

	// 2.获取redis连接池
	conn := redis.RedisPool().Get()
	defer conn.Close()

	// 3.获取文件句柄，用于存储分块内容
	fPath := "./temp/data/" + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(fPath), 0744)
	fd, err := os.Create(fPath)
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

	uploadId := r.Form.Get("upload_id")
	fileHash := r.Form.Get("hash")
	fileSize, _ := strconv.Atoi(r.Form.Get("size"))
	fileName := r.Form.Get("name")



	// 2.获取redis连接池
	conn := redis.RedisPool().Get()
	defer conn.Close()

	// 3.查询是否所有分块都上传完成
	data, err := conn.Do("HGETALL", "MP_" + uploadId)
	if err != nil {
		println("redis查询失败,err="+err.Error())
		res.Code = 500
		res.Message = "redis查询失败"
		return
	}
	totalCount := 0
	chunkCount := 0
	dic := data.(map[string]string)

	for key := range dic{
		if key == "ChunkCount" {
			totalCount, _ = strconv.Atoi(dic[key])
		} else if strings.HasPrefix(key, "idx_")  {
			v, _ := strconv.Atoi(dic[key])
			if v == 1 {
				chunkCount += 1
			}
		}
	}

	if totalCount != chunkCount {
		res.Code = 500
		res.Message = "非法请求"
		return
	}

	// 4.分块合并

	// 5.更新唯一文件表和用户文件表
	db.OnFileUploadFinish(fileHash, fileName, "", int64(fileSize))
	db.OnUserFileUploadFinished(user.Id, user.UserName, fileHash, fileName, "", int64(fileSize))

	// 6.返回结果给前端
}