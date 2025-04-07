package global

import (
	"KubeGale/config"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"
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
	KUBEGALE_MONGO     *qmgo.QmgoClient
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
