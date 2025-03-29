package middleware

import (
	"KubeGale/global"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// SQLLogger 自定义SQL日志记录器
type SQLLogger struct {
	logger.Interface
}

// NewSQLLogger 创建自定义SQL日志记录器
func NewSQLLogger() *SQLLogger {
	return &SQLLogger{
		Interface: global.KUBEGALE_DB.Logger,
	}
}

// Trace 记录SQL执行情况
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	
	// 记录SQL语句
	global.KUBEGALE_LOG.Info("SQL执行",
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
		zap.Error(err),
	)
	
	// 调用原始日志记录器
	l.Interface.Trace(ctx, begin, fc, err)
}

// SQLLogMiddleware 记录SQL日志的中间件
func SQLLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在请求处理前，替换为带有SQL日志的DB实例
		if global.KUBEGALE_CONFIG.Mysql.LogMode == "debug" || global.KUBEGALE_CONFIG.Mysql.LogMode == "info" {
			// 创建一个新的DB实例，使用自定义日志记录器
			sqlLogger := NewSQLLogger()
			db := global.KUBEGALE_DB.Session(&gorm.Session{
				Logger: sqlLogger,
			})
			
			// 临时保存原始DB
			originalDB := global.KUBEGALE_DB
			
			// 替换全局DB为带日志的DB
			global.KUBEGALE_DB = db
			
			// 请求结束后恢复原始DB
			c.Next()
			
			global.KUBEGALE_DB = originalDB
		} else {
			c.Next()
		}
	}
}