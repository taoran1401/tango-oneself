package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"taogin/app/types"
	"taogin/app/types/pb"
	"taogin/config/global"
	"time"
)

type Client struct {
	wsConn   *websocket.Conn
	addr     string
	sendChan chan []byte
	//用户信息
	User   *types.UserBase
	UserId uint64
	//前一次心跳时间
	HeartbeatTime int64
}

func NewClient(ws *websocket.Conn, cxt *gin.Context) {
	userInfo, ok := cxt.Get("user_info")
	if !ok {
		global.LOG.Error("用户信息异常")
		return
	}

	UserBase, ok := userInfo.(*types.UserBase)
	if !ok {
		global.LOG.Error("用户信息异常")
		return
	}

	var clt = &Client{
		wsConn:   ws,
		addr:     ws.RemoteAddr().String(),
		sendChan: make(chan []byte),
		User:     UserBase,
		UserId:   UserBase.Id,
	}

	//添加user信息到redis
	clt.AddUserCache()

	//绑定客户端
	ClientManages.Clients[string(clt.UserId)] = clt

	//clt.online()
	go clt.sendMessage()
	go clt.recvMessage(cxt)
}

//添加客户端
func (this *Client) addOnlineUserMap(client *Client) {
	server.userMapLock.Lock()
	server.onlineUserMap[client.addr] = client
	server.userMapLock.Unlock()
}

//删除客户端
func (this *Client) deleteOnlineUserMap(client *Client) {
	server.userMapLock.Lock()
	delete(server.onlineUserMap, client.addr)
	server.userMapLock.Unlock()
}

//发送
func (this *Client) sendMessage() {
	//defer this.offline()
	for {
		buf := <-this.sendChan
		err := this.wsConn.WriteMessage(websocket.TextMessage, buf)
		if err != nil {
			global.LOG.Error(err.Error())
			return
		}
	}
}

//接收
func (this *Client) recvMessage(ctx *gin.Context) {
	//defer this.offline()
	for {
		_, message, err := this.wsConn.ReadMessage()
		if err != nil {
			global.LOG.Error(err.Error())
			return
		}
		//this.sendChan <- message

		req := &pb.PbReq{
			User:   this.User,
			UserId: this.UserId,
			PbByte: message,
			Ctx:    ctx,
		}
		//分发
		Dispatch(this, req)
		global.LOG.Info(string(message) + "接收消息：")
	}
}

//上线
func (this *Client) online() {
	this.addOnlineUserMap(this)

	global.LOG.Info(this.addr + "上线了")
}

//下线
func (this *Client) offline() {
	this.deleteOnlineUserMap(this)

	global.LOG.Info(this.addr + "下线了")
}

func (this *Client) AddUserCache() {
	wzUserKey := fmt.Sprintf("%s%d", "wz:user:", this.UserId)
	wzData := types.Wz{
		User:   this.User,
		Room:   0,
		Status: 0,
	}
	wzUserValue, err := json.Marshal(wzData)
	if err != nil {
		fmt.Println(err)
		global.LOG.Error("json序列化失败")
		return
	}
	ttl := 3600 * 24 * time.Second
	global.CACHE.Set(wzUserKey, string(wzUserValue), ttl)
}
