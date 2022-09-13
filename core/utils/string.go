package utils

import (
	"math/rand"
	"time"
)

//生成随机字符串
func RandString(length int) string {
	str := "0123456789abcdefghigklmnopqrstuvwxyz"
	strList := []byte(str)
	result := []byte{}
	i := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i < length {
		new := strList[r.Intn(len(strList))]
		result = append(result, new)
		i = i + 1
	}
	return string(result)
}
