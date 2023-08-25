package tool

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func GenMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
