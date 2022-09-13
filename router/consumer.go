package router

import (
	consumer_handle "taogin/app/amqp/consumer"
	"taogin/core/queue"
)

func InitRouteConsumer(consumer *queue.Consumer) {
	//参数说明： 队列名称, 是否自动应答, 消费逻辑函数
	consumer.AddConsumerHandler("dev-common-test1", false, consumer_handle.TestConsumer)
}
