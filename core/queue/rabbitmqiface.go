package queue

type RabbitmqIface interface {
	// 初始化通道
	NewChannel() error
	// 绑定队列
	Bind() error
	// 初始化channel并绑定
	NewChannelAndBind() error
	// 不需要应答的消费，无需调用InitPublishMq
	ConsumeNoAck(handler ConsumerHandler) error
	// 需要应答的消费，无需调用InitPublishMq
	ConsumeMustAck(handler ConsumerHandler) error
	// 自定义消费参数的消费
	ConsumeSelfDefine(queueName string, autoAck bool, handler ConsumerHandler) error
	// 发送消息
	Publish(exchange, routeKey string, data []byte) error
	// 关闭链接
	Close()
	// 安全退出
	SecurityExit()
}
