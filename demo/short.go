package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
)

var chars = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")

// 1) 哈希实现
func hashShortUrl(url string) {
	hex := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	resUrl := make([]string, 4)
	for i := 0; i < 4; i++ {
		val, _ := strconv.ParseInt(hex[i*8:i*8+8], 16, 0)
		lHexLong := val & 0x3fffffff
		outChars := ""
		for j := 0; j < 6; j++ {
			outChars += chars[0x0000003D&lHexLong]
			lHexLong >>= 5
		}
		resUrl[i] = outChars
	}
	fmt.Println(resUrl)
}

// 2) 自增长算法
func autoShortUrl(id uint64) string {
	return GetString62(Encode62(id))
}

func Encode62(id uint64) []uint64 {
	tempE := []uint64{}

	for id > 0 {
		tempE = append(tempE, id%62)
		id /= 62
	}
	return tempE
}

func GetString62(indexA []uint64) string {
	res := ""

	for _, val := range indexA {
		res += chars[val]
	}
	return reverseString(res)
}

// 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func main() {
	hashShortUrl("https://www.baidu.com")

	//生成 - 入库
	//访问 - 中转请求,记录相关信息 -
	fmt.Println(autoShortUrl(12345678978945612345))
}
