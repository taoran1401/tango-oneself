package cache

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type Cache struct {
	Redis  *redis.Client
	Prefix string
	CacheInterface
}

func NewCache(redis *redis.Client, prefix string) *Cache {
	return &Cache{
		Redis:  redis,
		Prefix: prefix, //前缀
	}
}

func (this *Cache) Get(key string) string {
	value, err := this.Redis.Get(this.BuildKey(key)).Result()
	if err != nil {
		return ""
	}
	return value
}

func (this *Cache) Set(key string, value string, ttl time.Duration) bool {
	err := this.Redis.Set(this.BuildKey(key), value, ttl).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (this *Cache) Has(key string) bool {
	res := false
	result, err := this.Redis.Exists(this.BuildKey(key)).Result()
	if err != nil {
		log.Println(err)
		res = false
	}
	if result == 1 {
		res = true
	}
	return res
}

func (this *Cache) Delete(key string) bool {
	err := this.Redis.Del(this.BuildKey(key)).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (this *Cache) BuildKey(key string) string {
	return this.Prefix + key
}
