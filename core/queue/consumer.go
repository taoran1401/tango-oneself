package queue

import (
	"errors"
	"sync"
	"sync/atomic"
	"taogin/config/load"
	"time"
)

//消费
type Consumer struct {
	conf        []*load.RabbitMQ //配置
	mqCliMaps   sync.Map         //mq客户端map
	handlerMaps sync.Map         //消费者方法map

	exitSignal chan bool //退出信号
	stopRun    bool
	quantity   int32
}

//消费方法(函数)
type ConsumerHandler func([]byte) error

func NewConsumer(conf []*load.RabbitMQ) *Consumer {
	return &Consumer{
		conf:       conf,
		exitSignal: make(chan bool),
	}
}

//启动消费 - 入口
func (this *Consumer) Start() {
	for _, mqconf := range this.conf {
		if _, ok := this.handlerMaps.Load(mqconf.QueueName); !ok {
			panic("请先定义消费的入口函数handler，请检查！")
		}
	}

	go func() {
		//异常打印日志
		//defer utils.CatchException("Start", true)
		//开启每个配置对应的消费者
		for i, _ := range this.conf {
			this.runOneQueueConsumer(this.conf[i])
		}
	}()
}

//启动消费 - 初始化参数、数据验证
func (this *Consumer) runOneQueueConsumer(config *load.RabbitMQ) {
	//用多少协程处理
	if config.RunQuantity <= 0 {
		return
	}

	//初始化配置
	param := &InitParam{
		AmqpUrl:         config.AmqpUrl,
		QueueName:       config.QueueName,
		ExchangeName:    config.ExchangeName,
		ExchangeType:    config.ExchangeType,
		ExchangeDurable: config.ExchangeDurable,
		RoutingKey:      config.RoutingKey,
		Durable:         config.Durable,
	}

	for i := 0; i < config.RunQuantity; i++ {
		go func(param *InitParam) {
			//defer utils.CatchException("startQueueConsumerWhile", true)
			//启动消费 - 获取mq客户端连接
			this.startQueueConsumerWhile(param)
		}(param)
	}
}

//启动消费 - 获取mq客户端连接启动
func (this *Consumer) startQueueConsumerWhile(param *InitParam) {
	//检查参数
	if err := this.checkParameter(param); err != nil {
		//log.Error("启动参数配置错误导致消息队列无法启动，请检查！NewRabbitMqClient-err:", err)
		return
	}

	atomic.AddInt32(&this.quantity, 1)
	for {
		//获取mq客户端
		mq, err := NewRabbitMqClient(param)
		if nil != err {
			//log.Error("消息队列无法启动，请检查！NewRabbitMqClient-err:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		//加入mq客户端map
		this.mqCliMaps.Store(mq, true)
		//启动消费
		err = this.startQueueConsumer(mq)
		if err != nil {
			time.Sleep(time.Second)
		}
		//删除mq客户端map
		this.mqCliMaps.Delete(mq)
	}
	atomic.AddInt32(&this.quantity, -1)
	this.exitSignal <- true
}

//启动消费
func (cr *Consumer) startQueueConsumer(mqCli *RabbitMqClient) error {
	//defer utils.CatchException("startQueueConsumer")
	defer mqCli.Close()

	//获取消费程序
	v, _ := cr.handlerMaps.Load(mqCli.QueueName)
	//断言是否消费者
	consumerH := v.(*ConsumerStru)

	//创建通道
	//写数据，读数据等操作,创建queue，exchange，都需要在channel的基础上进行
	//所以要先生成channel
	//通过dial.Channel(),得到channel和错误信息
	if err := mqCli.NewChannel(); nil != err {
		//log.Error("消息队列无法启动，请检查！NewChannel-err:", err)
		time.Sleep(3 * time.Second)
		return err
	}

	//创建exchange，queue以及建立绑定
	if err := mqCli.Bind(); nil != err {
		//log.Error("消息队列无法启动，请检查！Bind-err:", err)
		time.Sleep(3 * time.Second)
		return err
	}

	//自定义消费参数的消费
	if err := mqCli.ConsumeSelfDefine(consumerH.QueueName, consumerH.AutoAck, consumerH.Handler); nil != err {
		//log.Error("消息队列无法启动，请检查！Bind-err:", err)
		time.Sleep(3 * time.Second)
		return err
	}

	return nil
}

//检查参数
func (this *Consumer) checkParameter(param *InitParam) error {
	if "" == param.AmqpUrl {
		return errors.New("AmqpUrl参数非法")
	}
	if "" == param.QueueName {
		return errors.New("QueueName参数非法")
	}
	if "" == param.ExchangeName {
		return errors.New("ExchangeName参数非法")
	}
	return nil
}

//消费者
type ConsumerStru struct {
	QueueName string
	AutoAck   bool
	Handler   ConsumerHandler
}

//添加消费者
func (this *Consumer) AddConsumerHandler(queueName string, autoAck bool, handler ConsumerHandler) {
	this.handlerMaps.Store(queueName, &ConsumerStru{
		QueueName: queueName,
		AutoAck:   autoAck,
		Handler:   handler,
	})
}

func (this *Consumer) SecurityExit() {
	this.stopRun = true
	// 发送退出信号
	this.mqCliMaps.Range(func(key interface{}, value interface{}) bool {
		cli := key.(*RabbitMqClient)
		cli.SecurityExit()
		return true
	})
	count := 0
	tm := time.Tick(3 * time.Second)
	//定时循环检查通道是否关闭，超过3次后关闭exitSignal通道
	for {
		select {
		case <-this.exitSignal:
			if atomic.CompareAndSwapInt32(&this.quantity, 0, 0) {
				return
			}
		case <-tm:
			count++
			if 3 < count {
				atomic.StoreInt32(&this.quantity, 0)
				close(this.exitSignal)
				return
			}
			//log.Info("等待所有消费者退出。。。", count)
		}
	}
}
