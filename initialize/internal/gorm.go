package internal

import (
	"KubeGale/config"
	"KubeGale/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var Gorm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	var general config.GeneralDB
	general = global.KUBEGALE_CONFIG.Mysql.GeneralDB

	// 创建自定义日志记录器
	logWriter := NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags))

	// 配置日志级别
	logLevel := general.LogLevel()
	// 如果需要根据配置文件决定，可以取消注释下面的代码
	if general.LogMode == "debug" || general.LogMode == "Debug" {
		logLevel = logger.Info // 设置为Info以显示SQL语句 // 强制设置为 Error 级别，只记录错误
	}

	return &gorm.Config{
		Logger: logger.New(logWriter, logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logLevel,
			Colorful:      true,
			// 添加以下配置以记录所有SQL
			IgnoreRecordNotFoundError: true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}
