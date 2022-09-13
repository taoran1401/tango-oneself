package cache

type CacheInterface interface {
	//获取
	Get(key string)
	//设置
	Set(key string, value string, ttl int) bool
	//判断是否存在
	Has(key string) bool
	//清除指定key的缓存
	Delete(key string) bool
}
