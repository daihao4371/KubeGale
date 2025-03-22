package core

import (
	"KubeGale/global"
	"KubeGale/initialize"
	"KubeGale/middleware"
	"KubeGale/utils"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct{}
type server interface {
	ListenAndServe() error
}

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "KubeGale server",
		Long: "The KubeGale server controller is a daemon that embeds the core control loops.",
		Run: func(cmd *cobra.Command, args []string) {
			s := Server{}
			// 加载所有配置文件
			s.Init()
		},
	}
	return cmd
}

func (s Server) Init() {

	global.KUBEGALE_VP = Viper()            // 初始化viper配置
	global.KUBEGALE_LOG = Zap()             // 初始化zap日志
	zap.ReplaceGlobals(global.KUBEGALE_LOG) // 替换zap
	global.KUBEGALE_DB = initialize.Gorm()  // 初始化gorm数据库连接
	if global.KUBEGALE_DB != nil {
		initialize.RegisterTables() //初始化表结构
		db, _ := global.KUBEGALE_DB.DB()
		defer db.Close()
	}
	if global.KUBEGALE_CONFIG.System.UseMultipoint || global.KUBEGALE_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
		initialize.RedisList()
	}

	// 获取Redis客户端并传递给JWT处理器
	var redisClient redis.Cmdable
	if global.KUBEGALE_CONFIG.System.UseRedis {
		redisClient = global.KUBEGALE_REDIS
	}
	handler := utils.NewJWTHandler(redisClient)
	enforcer := middleware.InitCasbin(global.KUBEGALE_DB)
	initialize.InitMiddlewares(handler, global.KUBEGALE_LOG, enforcer)

	// 初始化服务器组件
	RunServer()
}
func RunServer() {
	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.KUBEGALE_CONFIG.System.Addr)
	s := initServer(address, Router)

	global.KUBEGALE_LOG.Info("server run success on ", zap.String("address", address))

	global.KUBEGALE_LOG.Error(s.ListenAndServe().Error())
}
