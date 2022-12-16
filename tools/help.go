package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
