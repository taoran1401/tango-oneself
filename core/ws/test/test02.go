package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"taogin/protobuf/protopb"

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
	//msgpb
	msgpb := &protopb.Msg{
		Cmd:    196608,
		ReqUid: 12389123,
		Tms:    123123123123,
	}
	//pb转换[]byte
	msg, _ := proto.Marshal(msgpb)

	//服务器地址
	url := "ws://localhost:8088/ws/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NTQyNzMsImlzcyI6InJhbmJsb2dzIiwiZXhwIjoxNjYyNjE5NTQxLCJuYmYiOjE2NjIwMTQ3NDEsImlhdCI6MTY2MjAxNDc0MX0.wo2Fy9l7u8LYJItm8zPwIqvKqbyOs1R0_ChXnBiGYS0"
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
		//消息反序列化
		pbmsg := &protopb.Msg{}
		err = proto.Unmarshal(data, pbmsg)
		if err != nil {
			fmt.Println("proto.Unmarshal failed ")
		}
		fmt.Println(pbmsg)
	}
}
