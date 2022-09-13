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
	var action uint32

	//服务器地址
	url := "ws://localhost:8088/ws/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6NTQyNzMsImlzcyI6InJhbmJsb2dzIiwiZXhwIjoxNjYzMjExMjQ5LCJuYmYiOjE2NjI2MDY0NDksImlhdCI6MTY2MjYwNjQ0OX0.q-3ZCDPOHa1ol3AAVT5YWr3fiZXVMnadUEIjYyMTyRo"
	header := make(map[string][]string)
	header["debug"] = []string{
		"dev-debug",
	}
	ws, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal(err)
	}

	go WriteMsgTest(ws, action)

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

		ResponseTest(pbmsg)

		go WriteMsgTest(ws, action)
	}
}

func OnlineList() (msg []byte) {
	//msgpb
	msgpb := &protopb.Msg{
		Cmd:    196608,
		ReqUid: time.Now().UnixNano(),
		Tms:    time.Now().UnixNano(),
	}
	//pb转换[]byte
	msg, _ = proto.Marshal(msgpb)
	return msg
}

func ApplyVs() (msg []byte) {
	var userId uint64
	userId = 54272
	vspb := &protopb.ApplyVsReq{
		TargetUserId: userId,
	}

	vsb, err := proto.Marshal(vspb)
	if err != nil {
		log.Fatal(err)
	}

	//msgpb
	msgpb := &protopb.Msg{
		Cmd:    262144,
		ReqUid: time.Now().UnixNano(),
		Tms:    time.Now().UnixNano(),
		Data:   vsb,
	}

	//pb转换[]byte
	msg, _ = proto.Marshal(msgpb)
	return msg
}

func RouterTest(action uint32) []byte {
	var msg []byte
	if action == (uint32(protopb.CmdBase_CmdBaseOnlineList) << 16) {
		msg = OnlineList()
	} else if action == (uint32(protopb.CmdBase_CmdBaseApplyVs) << 16) {
		msg = ApplyVs()
	}
	return msg
}

func WriteMsgTest(ws *websocket.Conn, action uint32) {
	//param
	fmt.Println("请输入: ")
	_, err := fmt.Scanln(&action)
	if err != nil {
		log.Fatal(err)
	}

	//router
	msg := RouterTest(action)

	//write message
	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 2)
}

func ResponseTest(pbmsg *protopb.Msg) {
	fmt.Println("response data ----")
	fmt.Println(pbmsg)
	if pbmsg.Cmd == (uint32(protopb.CmdBase_CmdBaseOnlineList) << 16) {
		resp := protopb.OnlineListResp{}
		proto.Unmarshal(pbmsg.Data, &resp)
		fmt.Println(">>> online list:")
		for _, user := range resp.UserList {
			fmt.Printf("id: %d; nickname: %s;\n", user.Id, user.Nickname)
		}
	} else if pbmsg.Cmd == (uint32(protopb.CmdBase_CmdBaseApplyVs) << 16) {

	}
	fmt.Println("response data ----")
}
