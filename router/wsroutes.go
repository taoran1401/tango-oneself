package router

import (
	"taogin/app/controller/chess"
	"taogin/core/ws"
	"taogin/protobuf/protopb"
)

func WsInitRoute() {
	//online list
	ws.AddWsHandleFunc((uint32(protopb.CmdBase_CmdBaseOnlineList) << 16), chess.NewChess().OnlineList) //cmd: 196608
	ws.AddWsHandleFunc((uint32(protopb.CmdBase_CmdBaseApplyVs) << 16), chess.NewChess().ApplyVs)       //262144
	ws.AddWsHandleFunc((uint32(protopb.CmdBase_CmdBaseVsResp) << 16), chess.NewChess().VsResp)         //327680
	ws.AddWsHandleFunc((uint32(protopb.CmdBase_CmdBaseChess) << 16), chess.NewChess().Chess)           //393216
	ws.AddWsHandleFunc((uint32(protopb.CmdBase_CmdBaseSurrender) << 16), chess.NewChess().Surrender)   //458752
}
