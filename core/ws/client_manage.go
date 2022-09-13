package ws

import (
	"sync"
)

type ClientManage struct {
	Clients map[string]*Client
	mutex   *sync.Mutex
}

var ClientManages = ClientManage{
	Clients: make(map[string]*Client),
	mutex:   new(sync.Mutex),
}

//获取客户端
func GetClient(userId uint64) *Client {
	client := ClientManages.Clients[string(userId)]
	return client
}
