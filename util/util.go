package util

import (
	"crypto/sha256"
	"io"
	"os"
)

func fileHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		println("读取文件失败")
		return ""
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {

	}
	sum := hash.Sum(nil)
	return string(sum)
}
