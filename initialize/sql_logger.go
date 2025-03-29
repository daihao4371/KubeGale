package initialize

import (
	"KubeGale/global"
	"context"
	"go.uber.org/zap"
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
