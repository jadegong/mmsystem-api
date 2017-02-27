package g

import (
	"gopkg.in/redis.v3"
	"github.com/maemual/go-cache"
	"github.com/Sirupsen/logrus"
)

var (
	Redis *redis.Client
	Cache *cache.Cache
)

func InitCache() {
	if Conf.Cache == CACHE_REDIS {
		Redis = redis.NewClient(&redis.Options{
			Addr:     Conf.RedisUrl,
			Password: Conf.RedisPassword,
			DB:       0,
			PoolSize: 10,
		})

		_, err := Redis.Ping().Result()
		if err != nil {
			panic(err)
		}
	} else {
		Cache = cache.New(0, 120)
	}
	logrus.Infof("use [%s] as cache engine.", Conf.Cache)
}
