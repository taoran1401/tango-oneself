package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"taogin/config/global"
)

func Viper() *viper.Viper {
	//创建实例
	viper := viper.New()
	//设置文件名，不需要后缀While parsing config: yaml: unmarshal errors
	viper.SetConfigName("config")
	//设置文件格式
	viper.SetConfigType("yaml")
	//查找配置文件路径
	viper.AddConfigPath("./")
	//读取
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	//监听配置文件
	viper.WatchConfig()
	//被监听配置文件修改后重新读取
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file name", e.Name)
		err = viper.Unmarshal(&global.CONFIG) //序列化
		if err != nil {
			fmt.Println(err)
		}
	})
	//解析配置文件
	err = viper.Unmarshal(&global.CONFIG) //序列化
	if err != nil {
		fmt.Println("解析配置文件失败", err)
	}
	return viper
}
