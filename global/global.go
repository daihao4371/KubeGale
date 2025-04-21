package global

import (
	"KubeGale/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"sync"
)

var (
	KUBEGALE_DB        *gorm.DB
	KUBEGALE_DBList    map[string]*gorm.DB
	KUBEGALE_REDIS     redis.UniversalClient
	KUBEGALE_REDISList map[string]redis.UniversalClient
	KUBEGALE_CONFIG    config.Server
	KUBEGALE_VP        *viper.Viper
	// KUBEGALE_LOG    *oplogging.Logger
	KUBEGALE_LOG                 *zap.Logger
	KUBEGALE_Concurrency_Control = &singleflight.Group{}
	KUBEGALE_ROUTERS             gin.RoutesInfo
	KUBEGALE_ACTIVE_DBNAME       *string
	BlackCache                   local_cache.Cache
	lock                         sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return KUBEGALE_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := KUBEGALE_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

func GetRedis(name string) redis.UniversalClient {
	redis, ok := KUBEGALE_REDISList[name]
	if !ok || redis == nil {
		panic(fmt.Sprintf("redis `%s` no init", name))
	}
	return redis
}
