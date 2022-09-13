package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

//生成密码
func CreatePassword(password string, salt string) string {
	if salt == "" {
		salt = RandString(5)
	}
	passwordSalt := []byte(password + salt)
	sha1Obj := sha1.New()
	sha1Obj.Write(passwordSalt)
	sha1String := hex.EncodeToString(sha1Obj.Sum(nil))
	return Md5(sha1String)
}

//比较密码
func ComparePassword(encrypted_password string, password string, salt string) bool {
	if encrypted_password != CreatePassword(password, salt) {
		return true
	}
	return false
}
