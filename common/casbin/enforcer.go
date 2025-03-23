package casbin

import (
	"KubeGale/global"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"sync"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

// GetEnforcer 获取 Casbin 强制器的单例实例
func GetEnforcer() (*casbin.Enforcer, error) {
	var err error
	once.Do(func() {
		// 使用 GORM 适配器
		adapter, adapterErr := gormadapter.NewAdapterByDB(global.KUBEGALE_DB)
		if adapterErr != nil {
			global.KUBEGALE_LOG.Error("创建Casbin适配器失败", zap.Error(adapterErr))
			err = fmt.Errorf("创建Casbin适配器失败: %w", adapterErr)
			return
		}

		// 从配置文件和适配器创建强制器
		e, enforcerErr := casbin.NewEnforcer("model.conf", adapter)
		if enforcerErr != nil {
			global.KUBEGALE_LOG.Error("创建Casbin强制器失败", zap.Error(enforcerErr))
			err = fmt.Errorf("创建Casbin强制器失败: %w", enforcerErr)
			return
		}

		// 加载策略
		if loadErr := e.LoadPolicy(); loadErr != nil {
			global.KUBEGALE_LOG.Error("加载Casbin策略失败", zap.Error(loadErr))
			err = fmt.Errorf("加载Casbin策略失败: %w", loadErr)
			return
		}

		enforcer = e
		global.KUBEGALE_LOG.Info("Casbin强制器初始化成功")
	})

	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

// GetMethodName 将HTTP方法数字转换为字符串
func GetMethodName(method int) string {
	switch method {
	case 1:
		return "GET"
	case 2:
		return "POST"
	case 3:
		return "PUT"
	case 4:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}