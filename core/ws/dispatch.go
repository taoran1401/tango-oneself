package ws

import (
	"github.com/golang/protobuf/proto"
	"sync"
	"taogin/app/types/pb"
	"taogin/config/global"
	"taogin/protobuf/protopb"
	"time"
)

var (
	routeMaps map[uint32]DispatchFunc
	mux       sync.RWMutex
)

type DispatchFunc func(c *pb.PbReq, cmd uint32, reqUid int64, data []byte) error

func Dispatch(client *Client, pbreq *pb.PbReq) {
	message := pbreq.PbByte
	msg, err := decodeProto(message)
	if nil != err {
		//response.Error(ctx, atom.ERROR_CODE_PARAM, atom.GetMsgByCode(atom.ERROR_CODE_PARAM))
		return
	}
	if h, ok := routeMaps[msg.Cmd]; ok {
		if err := h(pbreq, msg.Cmd, msg.ReqUid, []byte(msg.Data)); nil != err {
			global.LOG.Error("cmd error:", err)
			//response.Error(ctx, atom.ERROR_CODE_EXCEPTION, atom.GetMsgByCode(atom.ERROR_CODE_EXCEPTION))
		}
	} else {
		global.LOG.Error("ws-cmd非法")
		//SendMsg(client, msg.Cmd， msg.)
		//response.Error(ctx, atom.ERROR_CODE_PARAM, atom.GetMsgByCode(atom.ERROR_CODE_PARAM))
	}
}

//proto反序列化
func decodeProto(bt []byte) (*protopb.Msg, error) {
	msg := &protopb.Msg{}
	if err := proto.Unmarshal(bt, msg); nil != err {
		return nil, err
	}
	return msg, nil
}

//proto序列化
func encodeProto(data proto.Message) ([]byte, error) {
	bt, err := proto.Marshal(data)
	if err != nil {
		global.LOG.Error("序列化proto.message失败！")
		return bt, nil
	}
	return bt, err
}

func AddWsHandleFunc(cmd uint32, fc DispatchFunc) {
	mux.Lock()
	defer mux.Unlock()

	if 0 == len(routeMaps) {
		routeMaps = make(map[uint32]DispatchFunc)
	}

	routeMaps[cmd] = fc
}

//发送消息给ws客户端
func SendMsg(userId uint64, cmd uint32, reqUnqid int64, data proto.Message) {
	//获取客户端
	client := GetClient(userId)

	//序列化需要发送的数据
	bt, err := proto.Marshal(data)
	if nil != err {
		global.LOG.Error("proto序列化失败:", err)
		return
	}

	respData := &protopb.Msg{
		Cmd:    cmd,
		ReqUid: reqUnqid,
		Tms:    time.Now().UnixNano(),
		Data:   bt,
	}

	respByte, err := encodeProto(respData)
	if err != nil {
		global.LOG.Error("proto序列化失败")
		return
	}

	//todo: 这里没有获取到对应client时会出现无效的指针引用问题

	//发送消息
	if client.sendChan != nil {
		client.sendChan <- respByte
	}
}
