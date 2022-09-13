package queue

import (
	"taogin/config/load"
)

type ProducerIface interface {
	Init() error
	GetMq() RabbitmqIface
}

type Producer struct {
	conf  *load.RabbitMQ
	MqCli *RabbitMqClient
}

func NewProducer(c *load.RabbitMQ) *Producer {
	return &Producer{
		conf: c,
	}
}

func (this *Producer) Init() error {
	var err error

	c := this.conf
	this.MqCli, err = NewRabbitMqClient(&InitParam{
		AmqpUrl:         c.AmqpUrl,
		QueueName:       c.QueueName,
		ExchangeName:    c.ExchangeName,
		ExchangeType:    c.ExchangeType,
		ExchangeDurable: c.ExchangeDurable,
		RoutingKey:      c.RoutingKey,
		Durable:         c.Durable,
	})
	if nil != err {
		return err
	}

	return nil
}

func (p *Producer) GetMq() RabbitmqIface {
	return p.MqCli
}
