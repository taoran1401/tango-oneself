package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"taogin/config/global"
)

type Server struct {
	onlineUserMap map[string]*Client
	userMapLock   sync.RWMutex
}

var server = Server{
	onlineUserMap: make(map[string]*Client),
}

func WsServer(ctx *gin.Context, w http.ResponseWriter, r *http.Request) {
	//http升级到websocket配置
	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//升级成websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		global.LOG.Error(err.Error())
		return
	}

	NewClient(ws, ctx)
}
