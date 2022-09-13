package pb

import (
	"github.com/gin-gonic/gin"
	"taogin/app/types"
)

type PbReq struct {
	User   *types.UserBase `json:"user"`    // 用户信息
	UserId uint64          `json:"uid"`     // 用户id
	PbByte []byte          `json:"pb_byte"` // 上传的protobuf的byte
	Ctx    *gin.Context
}

type PbResp struct {
	Code uint32     `json:"code"` // 0成功，其他失败
	Msg  string     `json:"msg"`  // 失败字符描述
	Data PbRespData `json:"data"` // 返回给app的内容
}

type PbRespData struct {
	Cmd       uint32 `json:"cmd"`        // 返回给app的cmd
	OriginCmd uint32 `json:"origin_cmd"` // 返回给app的cmd
	ReqUnqid  int64  `json:"req_unqid"`  // 返回给app的序列号
	PbByte    []byte `json:"pb_byte"`    // 返回给app的protobuf
}
