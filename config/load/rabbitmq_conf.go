package load

type RabbitMQ struct {
	AmqpUrl         string `mapstructure:"AmqpUrl" json:"AmqpUrl" yaml:"AmqpUrl"`
	QueueName       string `mapstructure:"QueueName" json:"QueueName" yaml:"QueueName"`
	ExchangeName    string `mapstructure:"ExchangeName" json:"ExchangeName" yaml:"ExchangeName"`
	ExchangeType    string `mapstructure:"ExchangeType" json:"ExchangeType" yaml:"ExchangeType"`
	ExchangeDurable bool   `mapstructure:"ExchangeDurable" json:"ExchangeDurable" yaml:"ExchangeDurable"`
	RoutingKey      string `mapstructure:"RoutingKey" json:"RoutingKey" yaml:"RoutingKey"`
	Durable         bool   `mapstructure:"Durable" json:"Durable" yaml:"Durable"`
	RunQuantity     int    `mapstructure:"RunQuantity" json:"RunQuantity" yaml:"RunQuantity"`
}
