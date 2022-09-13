package config

//服务配置，自动获取该文件中的方法执行
type ServiceConfig struct {
}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{}
}
