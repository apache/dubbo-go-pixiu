package redis

import (
	"github.com/go-redis/redis"
	"github.com/pantianying/dubbo-go-proxy/common/config"
	"github.com/pantianying/dubbo-go-proxy/common/logger"
	"github.com/pantianying/dubbo-go-proxy/common/util"
	"github.com/pantianying/dubbo-go-proxy/service"
	"time"
)

var rc *redis.Client

func init() {
	if config.Config.UseRedisMetadataCenter {
		rc = redis.NewClient(&redis.Options{
			Addr:        config.Config.RedisAddr,
			Password:    config.Config.RedisPassword,
			PoolSize:    100,
			MaxRetries:  3,
			IdleTimeout: time.Second * 10,
		})
		go func() {
			for {
				_, err := rc.Ping().Result()
				if err != nil {
					logger.Errorf("redis ping err:%+v", err)
				}
				time.Sleep(1 * time.Minute)
			}
		}()
	}

}

type redisMetaDataCenter struct {
}

func NewRedisMetaDataCenter() *redisMetaDataCenter {
	return &redisMetaDataCenter{}
}
func (r *redisMetaDataCenter) GetProviderMetaData(m service.MetadataIdentifier) *service.MetadataInfo {
	//todo 做缓存
	mInfo := &service.MetadataInfo{}
	cmd := rc.Get(getKey(m))
	if e := util.ParseJsonByStruct([]byte(cmd.String()), mInfo); e == nil {
		return mInfo
	}
	return nil
}

func getKey(m service.MetadataIdentifier) string {
	return m.ServiceInterface + ":provider:" + m.Application + ".metaData"
}
