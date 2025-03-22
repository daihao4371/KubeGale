package global

import (
	"github.com/redis/go-redis/v9"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/qmgo"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"KubeGale/config"

	"github.com/spf13/viper"
	"gorm.io/gorm"
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
