package core

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"taogin/config/global"
)

type MongoDB struct {
}

func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

func (this *MongoDB) MongoDB(uri string) *mongo.Client {
	var err error
	clientOptions := options.Client().ApplyURI(uri)
	// 连接到MongoDB
	m, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = m.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("初始化mongodb失败：", err)
		panic(err)
	}
	return m
}

//连接多个
func (this *MongoDB) ConnList() map[string]*mongo.Client {
	mongoMap := make(map[string]*mongo.Client)
	for _, info := range global.CONFIG.MongoDB {
		mongoMap[info.AliasName] = this.MongoDB(info.MongoDBConf.Uri)
	}
	return mongoMap
}
