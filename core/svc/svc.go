package svc

import "taogin/config"

type ServiceContext struct {
	Conf config.Config
}

func NewServiceContext(conf *config.Config) *ServiceContext {
	return &ServiceContext{Conf: *conf}
}
