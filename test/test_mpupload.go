package test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func multipartUpload(fileName, targetUrl string, chunkSize int) error {
	f, err := os.Open(fileName)
	if err != nil {
		println("文件打开错误,err="+err.Error())
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	for  {
		n, err := bfRd.Read(buf)
		if err != nil || n < 0 {
			break
		}
		index++

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			println("上传大小=%d",len(b))

			resp, err := http.Post(
				targetUrl + "?index=" + strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				println(strconv.Itoa(curIdx)+"上传失败,err="+err.Error())
			}
			body, err := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v %+v\n", string(body), err)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				println(err.Error())
			}
		}
	}

	for idx := 0; idx < index; idx++ {
		select {
		case res := <-ch:
			println("收到消息：%v", res)
		}
	}

	return nil
}

func main() {
	fileHash := "1ccf19994968035e3b3fb1b38370a8fe91975288"

	resp, err := http.PostForm(
		"http://localhost:8080/file/mpupload/init",
		url.Values{
			"uid": {"1"},
			"hash": {fileHash},
			"size": {"227129864"},
		})
	if err != nil {
		println("请求初始化分块上传接口失败,err="+err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("ReadAll失败,err="+err.Error())
		os.Exit(-1)
	}

	println("%v", body)
}