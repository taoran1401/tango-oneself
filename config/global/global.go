package global

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"taogin/config"
	"taogin/core/cache"
	"taogin/core/email_utils"
	"taogin/core/jwt"
	"taogin/core/sms"
)

var (
	CONFIG config.Config
	VIPER  *viper.Viper
	//LOG     *log.Logger
	LOG     *zap.SugaredLogger
	DB      map[string]*gorm.DB
	MONGODB map[string]*mongo.Client
	REDIS   map[string]*redis.Client
	JWT     *jwt.JWT
	CACHE   *cache.Cache
	EMAIL   *email_utils.Email
	SMS     *sms.Sms
	//CRON    *cron2.Cron
	//QUEUE   map[string]*queue.RabbitMQ
)
