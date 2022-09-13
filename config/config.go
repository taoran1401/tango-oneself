package config

import "taogin/config/load"

type Config struct {
	App      load.App
	Zap      load.Zap
	Jwt      load.JWT
	Mysql    []load.Mysql     `mapstructure:"Mysql" json:"Mysql" yaml:"Mysql"`
	MongoDB  []load.MongoDB   `mapstructure:"MongoDB" json:"MongoDB" yaml:"MongoDB"`
	Redis    []load.Redis     `mapstructure:"Redis" json:"Redis" yaml:"Redis"`
	RabbitMQ []*load.RabbitMQ `mapstructure:"RabbitMQ" json:"RabbitMQ" yaml:"RabbitMQ"`
	Email    load.Email
	Crontab  load.Crontab
	Cache    load.Cache
	Sms      load.Sms
}
