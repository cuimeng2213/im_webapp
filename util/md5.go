package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	ret := h.Sum(nil)
	return hex.EncodeToString(ret)
}
func ValidatePasswd(plainpwd, salt, passwd string) bool {
	fmt.Printf("input=[%s] %s \n", Md5Encode(plainpwd+salt), passwd)
	return Md5Encode(plainpwd+salt) == passwd
}
func MakePasswd(passwd, salt string) string {

	return Md5Encode(passwd + salt)
}
