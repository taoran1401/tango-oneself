package core

import (
	"github.com/go-redis/redis"
	"taogin/config/global"
)

type Redis struct {
}

func NewRedis() *Redis {
	return &Redis{}
}

func (this *Redis) RedisConn(addr string, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		global.LOG.Error("redis connect ping failed, err:" + err.Error())
	} else {
		global.LOG.Info("redis connect ping response:" + pong)
	}
	//程序关闭前关闭redis
	/*defer func() {
		fmt.Println("defer redis close")
		client.Close()
	}()*/
	return client
}

//连接多个
func (this *Redis) ConnList() map[string]*redis.Client {
	redisMap := make(map[string]*redis.Client)
	for _, info := range global.CONFIG.Redis {
		redisMap[info.AliasName] = this.RedisConn(info.RedisConf.Addr, info.RedisConf.Password, info.RedisConf.DB)
	}
	return redisMap
}
