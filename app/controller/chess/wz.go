package chess

import (
	"github.com/golang/protobuf/proto"
	"taogin/app/logic"
	"taogin/app/types/pb"
	"taogin/config/global"
	"taogin/core/ws"
	"taogin/protobuf/protopb"
)

type Wz struct {
}

func NewChess() *Wz {
	return &Wz{}
}

func (this *Wz) OnlineList(pbreq *pb.PbReq, cmd uint32, reqUid int64, data []byte) error {
	var (
		req  protopb.OnlineListReq
		resp protopb.OnlineListResp
	)

	if err := proto.Unmarshal(data, &req); err != nil {
		global.LOG.Error("proto 反序列化失败: ", err)
		return err
	}

	//调用logic
	err := logic.NewWz().OnlineList(&req, &resp)
	if err != nil {
		global.LOG.Error("逻辑异常：", err)
		return err
	}

	//发送消息
	ws.SendMsg(pbreq.UserId, cmd, reqUid, &resp)

	return nil
}

func (this *Wz) ApplyVs(pbreq *pb.PbReq, cmd uint32, reqUid int64, data []byte) error {
	var (
		req  protopb.ApplyVsReq
		resp protopb.ApplyVsResp
	)

	if err := proto.Unmarshal(data, &req); err != nil {
		global.LOG.Error("proto 反序列化失败: ", err)
		return err
	}

	//调用logic
	err := logic.NewWz().ApplyVs(&req, &resp, pbreq, cmd, reqUid)
	if err != nil {
		global.LOG.Error("逻辑异常：", err)
		return err
	}

	//响应
	ws.SendMsg(pbreq.UserId, cmd, reqUid, &resp)

	return nil
}

func (this *Wz) VsResp(pbreq *pb.PbReq, cmd uint32, reqUid int64, data []byte) error {
	var (
		req  protopb.VsRespReq
		resp protopb.VsRespResp
	)

	if err := proto.Unmarshal(data, &req); err != nil {
		global.LOG.Error("proto 反序列化失败: ", err)
		return err
	}

	//调用logic
	err := logic.NewWz().VsResp(&req, &resp)
	if err != nil {
		global.LOG.Error("逻辑异常：", err)
		return err
	}

	//发送消息
	ws.SendMsg(pbreq.UserId, cmd, reqUid, &resp)

	return nil
}

func (this *Wz) Chess(pbreq *pb.PbReq, cmd uint32, reqUid int64, data []byte) error {
	var (
		req  protopb.ChessReq
		resp protopb.ChessResp
	)

	if err := proto.Unmarshal(data, &req); err != nil {
		global.LOG.Error("proto 反序列化失败: ", err)
		return err
	}

	//调用logic
	err := logic.NewWz().Chess(&req, &resp)
	if err != nil {
		global.LOG.Error("逻辑异常：", err)
		return err
	}

	//发送消息
	ws.SendMsg(pbreq.UserId, cmd, reqUid, &resp)

	return nil
}

func (this *Wz) Surrender(pbreq *pb.PbReq, cmd uint32, reqUid int64, data []byte) error {
	var (
		req  protopb.SurrenderReq
		resp protopb.SurrenderResp
	)

	if err := proto.Unmarshal(data, &req); err != nil {
		global.LOG.Error("proto 反序列化失败: ", err)
		return err
	}

	//调用logic
	err := logic.NewWz().Surrender(&req, &resp)
	if err != nil {
		global.LOG.Error("逻辑异常：", err)
		return err
	}

	//发送消息
	ws.SendMsg(pbreq.UserId, cmd, reqUid, &resp)

	return nil
}
