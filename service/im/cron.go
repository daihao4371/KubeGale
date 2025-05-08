package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"KubeGale/model/im/response"
	messageIm "KubeGale/utils/im"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type CronService struct{}

var CronServiceApp = new(CronService)

// SendDailyStats 发送每日统计
func (s *CronService) SendDailyStats() error {
	// 获取所有飞书配置
	var feiShuConfigs []im.FeiShuConfig
	if err := global.KUBEGALE_DB.Find(&feiShuConfigs).Error; err != nil {
		return err
	}

	// 获取今天的告警统计
	today := time.Now().Format("2006-01-02")
	var alertCount int64
	if err := global.KUBEGALE_DB.Model(&im.CardContentConfig{}).
		Where("DATE(alert_time) = ?", today).
		Count(&alertCount).Error; err != nil {
		return err
	}

	// 获取未解决的告警数量
	var unresolvedCount int64
	if err := global.KUBEGALE_DB.Model(&im.CardContentConfig{}).
		Where("unresolved_alert = ?", true).
		Count(&unresolvedCount).Error; err != nil {
		return err
	}

	// 构建统计消息
	statsMessage := fmt.Sprintf("每日告警统计\n"+
		"日期: %s\n"+
		"今日告警总数: %d\n"+
		"未解决告警数: %d", today, alertCount, unresolvedCount)

	// 发送统计消息到所有飞书配置
	for _, config := range feiShuConfigs {
		if config.SendDailyStats {
			// 构建通知配置详情
			notificationConfig := response.NotificationDetailConfig{
				ID:                 config.ID,
				Name:               config.Name,
				Type:               config.Type,
				NotificationPolicy: config.NotificationPolicy,
				SendDailyStats:     config.SendDailyStats,
				CreatedAt:          config.CreatedAt,
				UpdatedAt:          config.UpdatedAt,
				RobotURL:           config.RobotURL,
			}

			// 发送飞书消息
			err := messageIm.MessageServiceApp.SendFeiShuMessage(notificationConfig, response.CardContentDetail{}, statsMessage)
			if err != nil {
				global.KUBEGALE_LOG.Error("发送每日统计失败",
					zap.String("config", config.Name),
					zap.Error(err))
			}
		}
	}

	return nil
}
