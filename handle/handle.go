package handle

import (
	"io/ioutil"
	"net/http"
)

func UploadHandle(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("../static/html/file_upload.html")
		if err != nil {
			println("file_upload.html读取失败!")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(data)
	} else {

	}
}