package queue

import (
	"errors"
	"fmt"
	"sync"
	"taogin/core/utils"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMqClient struct {
	conn    *amqp.Connection
	Channel *amqp.Channel

	QueueName       string
	ExchangeName    string
	ExchangeType    string
	ExchangeDurable bool
	RoutingKey      string
	Durable         bool
	MqUrl           string

	notifyCloseConn    chan *amqp.Error
	notifyCloseChannel chan *amqp.Error
	notifyConfirm      chan amqp.Confirmation

	mustExitEvent bool
	forwait       chan bool

	count            int
	consumeFailedMap sync.Map // 消费失败记录

}

//初始参数
type InitParam struct {
	AmqpUrl         string `json:"amqp_url"`
	QueueName       string `json:"queue_name"`
	ExchangeName    string `json:"exchange_name"`
	ExchangeType    string `json:"exchange_type"`
	ExchangeDurable bool   `json:"exchange_durable"`
	RoutingKey      string `json:"routing_key"`
	Durable         bool   `json:"durable"`
}

const (
	// 单个消费者最多处理多少个消息后退出
	MaxConsumerCount = 1000000

	// 消费处理失败暂停的毫秒数量
	MaxConsumerFailedWaitTimes = 100

	// 单条记录消费失败重试次数
	MaxConsumerFailedRepeatCount = 3
)

func NewRabbitMqClient(param *InitParam) (*RabbitMqClient, error) {
	var err error

	rmq := &RabbitMqClient{
		MqUrl:           param.AmqpUrl,
		QueueName:       param.QueueName,
		ExchangeName:    param.ExchangeName,
		ExchangeType:    param.ExchangeType,
		ExchangeDurable: param.ExchangeDurable,
		RoutingKey:      param.RoutingKey,
		Durable:         param.Durable,
	}

	rmq.forwait = make(chan bool)
	rmq.notifyCloseConn = make(chan *amqp.Error)
	rmq.notifyCloseChannel = make(chan *amqp.Error)

	rmq.CheckParameter()

	rmq.conn, err = amqp.Dial(rmq.MqUrl)
	if nil != err {
		//log.Error("RabbitMqClient-err:failed to connect tp rabbitmq,", err)
		return nil, err
	}

	return rmq, nil
}

func (this *RabbitMqClient) ReConnect() (err error) {
	this.conn, err = amqp.Dial(this.MqUrl)
	if nil != err {
		//log.Error("RabbitMqClient-ReConnect-err:failed to connect tp rabbitmq,", err, this.MqUrl)
		return err
	}

	this.Channel = nil

	return nil
}

func (this *RabbitMqClient) CheckErr(err error, msg string) {
	if nil == err {
		return
	}

	//log.Error("RabbitMqClient-err,", msg, err)
	panic(err)
}

// 初始化通道并绑定
func (this *RabbitMqClient) NewChannelAndBind() error {
	if this.Channel == nil {
		this.NewChannel()
	}

	return this.Bind()
}

// 一个连接只有一个通道，避免tcp堵塞
func (this *RabbitMqClient) NewChannel() error {
	var err error

	this.Channel, err = this.conn.Channel()
	if nil != err {
		//log.Error("RabbitMq channel err,", err)
		return err
	}

	return nil
}

func (this *RabbitMqClient) Close() {
	this.Channel.Close()
	this.conn.Close()

	this.Channel = nil
}

func (this *RabbitMqClient) ConsumeNoAck(handler ConsumerHandler) error {
	if this.Channel == nil {
		if err := this.NewChannelAndBind(); nil != err {
			return err
		}
	}

	this.ConnectCheckTimer()

	msgs, err := this.Channel.Consume(this.QueueName, "", true /*true*/, false, false, false, nil)
	if nil != err {
		//log.Error("Rabbitmq create Channel.consume error ", err)
		return err
	}

	go func(msgs <-chan amqp.Delivery) {
		//defer utils.CatchException("ConsumeNoAck", true)

		tms := time.Tick(time.Second)

		for {
			select {
			case d := <-msgs:
				this.handleMsg(d, true, handler)

				if this.mustExitEvent {
					//log.Info("recv mustExitEvent")
					this.forwait <- true
					return
				}

				this.count++
				if MaxConsumerCount < this.count {
					//log.Info("handle msg count gt ", MaxConsumerCount, ",ready quit")
					this.forwait <- true
					return
				}
			case <-tms:
				if this.mustExitEvent {
					//log.Info("recv mustExitEvent-time")
					return
				}
			}
		}
	}(msgs)

	//log.Info(" [*] Waiting for messages... ")

	<-this.forwait

	return nil
}

// 需要应答的消费
func (this *RabbitMqClient) ConsumeMustAck(handler ConsumerHandler) error {
	if this.Channel == nil {
		if err := this.NewChannelAndBind(); nil != err {
			return err
		}
	}

	this.ConnectCheckTimer()

	err := this.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if nil != err {
		//log.Error("Rabbitmq channel.Qos error ", err)
		return err
	}

	msgs, err := this.Channel.Consume(this.QueueName, "", false /*autoAck*/, false, false, false, nil)
	if nil != err {
		//log.Error("Rabbitmq create Channel.consume error ", err)
		return err
	}

	go func(msgs <-chan amqp.Delivery) {
		//defer utils.CatchException("ConsumeMustAck", true)

		tms := time.Tick(time.Second)
		for {
			select {
			case d := <-msgs:
				this.handleMsg(d, false, handler)

				if this.mustExitEvent {
					//log.Info("recv mustExitEvent")
					this.forwait <- true
					return
				}

				this.count++
				if MaxConsumerCount < this.count {
					//log.Info("handle msg count gt ", MaxConsumerCount, ",ready quit")
					this.forwait <- true
					return
				}
			case <-tms:
				if this.mustExitEvent {
					//log.Info("recv mustExitEvent-time")
					return
				}
			}
		}
	}(msgs)

	//log.Info(" [*] Waiting for messages...")

	<-this.forwait

	return nil
}

// 自定义消费参数的消费
func (this *RabbitMqClient) ConsumeSelfDefine(queueName string, autoAck bool, handler ConsumerHandler) error {
	if this.Channel == nil {
		if err := this.NewChannelAndBind(); nil != err {
			return err
		}
	}

	this.ConnectCheckTimer()

	err := this.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if nil != err {
		//log.Error("Rabbitmq channel.Qos error ", err)
		return err
	}

	msgs, err := this.Channel.Consume(queueName, queueName, autoAck, false, false, false, nil)
	if nil != err {
		//log.Error("Rabbitmq create Channel.consume error ", err)
		return err
	}

	go func(msgs <-chan amqp.Delivery) {
		//defer utils.CatchException("ConsumeSelfDefine", true)

		tms := time.Tick(time.Second)
		for {
			select {
			case d := <-msgs:
				this.handleMsg(d, autoAck, handler)

				if this.mustExitEvent {
					//log.Info("recv mustExitEvent")
					this.forwait <- true
					return
				}

				this.count++
				if MaxConsumerCount < this.count {
					//log.Info("handle msg count gt ", MaxConsumerCount, ",ready quit")
					this.forwait <- true
					return
				}
			case <-tms:
				if this.mustExitEvent {
					//log.Info("recv mustExitEvent-time")
					return
				}
			}
		}
	}(msgs)

	//log.Info(" [*] Waiting for messages...")

	<-this.forwait
	//log.Info("Waiting for messages forwait ok")

	return nil
}

func (this *RabbitMqClient) handleMsg(d amqp.Delivery, autoAck bool, handler ConsumerHandler) {
	//log.Info(",no: ", this.count, "接受到消息: ", d, ", ", string(d.Body), ", ", this.QueueName, ", ", this.ExchangeName)
	if 0 == len(d.Body) {
		//log.Info("空消息")
		return
	}
	err := func() (err error) {
		/*defer utils.CatchExceptionFunc("RabbitMqClient.handleMsg", func() {
			err = errors.New("RabbitMqClient.call.handle-捕获异常")
		})*/
		err2 := handler(d.Body)
		return err2
	}()
	if nil != err {
		//log.Error("rabbitmq.consume.handle.error,", d.DeliveryTag, ", 消费内容=", d, ",", err)
		if !autoAck {
			md5val := utils.Md5(string(d.Body))
			var v interface{}
			var ok bool
			if v, ok = this.consumeFailedMap.Load(md5val); ok {
				//log.Debug("************失败消息记录", d.DeliveryTag, "v=", v)
				if v.(int) > MaxConsumerFailedRepeatCount {
					// 不再消费该消息
					//log.Error("************关闭该消息rabbitmq.consume.repeat.failed count max:", string(d.Body))
					this.consumeFailedMap.Delete(d.DeliveryTag)

					this.Channel.Ack(d.DeliveryTag, true) // 应答

					return
				}

				this.consumeFailedMap.Store(md5val, v.(int)+1)
			} else {
				//log.Debug("************LoadOrStore-nook：", d.DeliveryTag, " , ", v, ok)
				this.consumeFailedMap.Store(md5val, 1)
			}
			//log.Debug("************不应答处理：", d.DeliveryTag)
			this.Channel.Nack(d.DeliveryTag, true, true)

			time.Sleep(time.Duration(MaxConsumerFailedWaitTimes) * time.Millisecond)

			return
		}
	}

	//time.Sleep(time.Duration(500) * time.Millisecond)
	if !autoAck {
		this.Channel.Ack(d.DeliveryTag, true)
	}
	//log.Info( ",ack , wait 3 s, ", d.DeliveryTag)

	//time.Sleep(time.Duration(500) * time.Millisecond)

}

func (this *RabbitMqClient) ConnectCheckTimer() {
	c := this.conn.NotifyClose(this.notifyCloseConn)

	go func() {
		//defer utils.CatchException("ConnectCheckTimer", true)

		select {
		case e, ok := <-c:
			var msg string
			if ok {
				msg = e.Reason
			}
			//log.Info("connect is close by server", msg)
			fmt.Println(msg)
			this.mustExitEvent = true
		}
	}()

	cl := this.Channel.NotifyClose(this.notifyCloseChannel)
	go func() {
		//defer utils.CatchException("ConnectCheckTimer", true)

		//log.Info("wait..")
		select {
		case e, ok := <-cl:
			var msg string
			if ok {
				msg = e.Reason
			}
			//log.Info("channel is close by server", msg)
			fmt.Println(msg)
			this.mustExitEvent = true
		}
	}()
}

func (this *RabbitMqClient) Publish(exchange, routeKey string, data []byte) error {
	if nil == this.conn || this.conn.IsClosed() {
		//log.Debug("mq.conn.IsClosed-ReConnect")
		if err := this.ReConnect(); nil != err {
			//log.Error("mq.ReConnect-err", err)
			return err
		}
	}
	if nil == this.Channel {
		//log.Debug("mq.Channel=nil-NewChannel")
		if err := this.NewChannel(); nil != err {
			//log.Error("mq.NewChannel-err", err)
			return err
		}
	}
	//log.Debug("Channel.Publish-begin", string(data))
	err := this.Channel.Publish(
		exchange, // exchange
		routeKey, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: amqp.Transient, // 1=non-persistent/Transient, 2=Persistent
			Priority:     0})             // 0-9

	return err
}

func (this *RabbitMqClient) CheckParameter() {
	if "" == this.MqUrl {
		panic(errors.New("MqUrl is not allow empty"))
	}
}

func (this *RabbitMqClient) Bind() error {
	//创建exchange(exchage可以将数据按照routeKey分发到对应的队列中去)
	//func (ch *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error
	//name:名称
	//kind:分发方式
	//durable:为true代表mq挂了自动恢复时，这个exchange还在
	//autoDelete:为true代表断开连接自动删除
	//internal:为true代表只有内部的exchange之间可以进行消息走动
	//noWait:是否等待返回
	//args:其他参数
	if err := this.Channel.ExchangeDeclare(
		this.ExchangeName,    // name of the exchange
		this.ExchangeType,    // type
		this.ExchangeDurable, // durable持久化
		false,                // delete when complete
		false,                // interna
		false,                // noWait
		nil,                  // arguments
	); nil != err {
		//log.Error("Rabbitmq Channel.ExchangeDeclare error ", err)
		return err
	}

	//通过 QueueDeclare 创建队列
	//func (ch *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error)
	//name:队列名，不填则会自动生成
	//durable: 为true代表mq挂了自动恢复时，这个队列还在
	//autoDelete: 为true代表断开连接就删除队列
	//exclusive: 为true代表只有当前连接可以使用该队列
	//noWait: 为true则不等待响应值返回
	//args: 其它附加参数
	queue, err := this.Channel.QueueDeclare(
		this.QueueName, // name of the queue
		this.Durable,   // durable 持久化
		false,          // delete when unused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	)
	if nil != err {
		//log.Error("Rabbitmq Channel.QueueDeclare error ", err)
		return err
	}

	//exchange要把数据分发到queue的前提条件是先创建queue,然后再将queue和exchange绑定
	if err := this.Channel.QueueBind(
		queue.Name,        // namethis.QueueName of the queue
		this.RoutingKey,   // bindingKey
		this.ExchangeName, // sourceExchange
		false,             // noWait
		nil,               // arguments
	); nil != err {
		//log.Error("Rabbitmq Channel.QueueBind error ", err)
		return err
	}

	return nil
}

func (this *RabbitMqClient) SecurityExit() {
	this.mustExitEvent = true
	//log.Info("exit-channel-----", this.QueueName, this.ExchangeName)
	//this.Channel.Cancel(this.QueueName, false)
	this.Channel.Close()
}
