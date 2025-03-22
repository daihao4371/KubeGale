package initialize

import (
	"KubeGale/config"
	"KubeGale/global"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func initRedisClient(redisCfg config.Redis) (redis.UniversalClient, error) {
	var client redis.UniversalClient
	// 使用集群模式
	if redisCfg.UseCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.ClusterAddrs,
			Password: redisCfg.Password,
		})
	} else {
		// 使用单例模式
		client = redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		})
	}
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.KUBEGALE_LOG.Error("redis connect ping failed, err:", zap.String("name", redisCfg.Name), zap.Error(err))
		return nil, err
	}

	global.KUBEGALE_LOG.Info("redis connect ping response:", zap.String("name", redisCfg.Name), zap.String("pong", pong))
	return client, nil
}

func Redis() redis.Cmdable {
	redisClient, err := initRedisClient(global.KUBEGALE_CONFIG.Redis)
	if err != nil {
		panic(err)
	}
	global.KUBEGALE_REDIS = redisClient
	return redisClient
}

func RedisList() {
	redisMap := make(map[string]redis.UniversalClient)

	for _, redisCfg := range global.KUBEGALE_CONFIG.RedisList {
		client, err := initRedisClient(redisCfg)
		if err != nil {
			panic(err)
		}
		redisMap[redisCfg.Name] = client
	}

	global.KUBEGALE_REDISList = redisMap
}
