package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	social_api_ws "taogin/core/ws/test/protopb"

	//"golang.org/x/net/websocket"
	"log"
	"time"
)

/**
websocket调试工具
- 连接，断开连接
- 支持protobuf传递(需自行安装protobuf转换工具)
- 设置头信息
- 发送消息
- 响应并打印
*/
func main() {
	//接口参数pb
	p := &social_api_ws.EveluateReq{
		Score:    3670016,
		AnchorId: 54079,
		Tags:     "7-8-9",
		OrderId:  123123123,
	}
	b, _ := proto.Marshal(p)

	//msgpb
	msgpb := &social_api_ws.Msg{
		Cmd:    3670016,
		ReqUid: 12389123,
		Tms:    123123123123,
		Data:   b,
	}
	//pb转换[]byte
	msg, _ := proto.Marshal(msgpb)

	//服务器地址
	url := "ws://localhost:7002/api/v1/ws/941591a21ef651193539cb1eb1e603c0549f041d66203c1661754427138805022"
	//social-login-tk-57927221de533fbf96dd524e3699f93a5e5746c335fcf11661506597400177858
	//16011000125_54093_1661506597400184014_10111
	header := make(map[string][]string)
	header["debug"] = []string{
		"dev-debug",
	}
	ws, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 2)
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("receive: ", data)
	}
}

func BytesCombine(pBytes ...[]byte) []byte {
	length := len(pBytes)
	s := make([][]byte, length)
	for index := 0; index < length; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
