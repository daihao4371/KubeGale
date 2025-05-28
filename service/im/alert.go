package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"KubeGale/model/im/response"
	messageIm "KubeGale/utils/im"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type AlertService struct{}

var AlertServiceApp = new(AlertService)

// AlertInfo 告警信息结构体
type AlertInfo struct {
	Level       string    // 告警等级
	Name        string    // 告警名称
	Content     string    // 告警内容
	Time        time.Time // 告警时间
	NotifyUsers string    // 通知人
	Handler     string    // 处理人
}

// SendAlert 发送告警通知
func (s *AlertService) SendAlert(alert *im.CardContentConfig) error {
	// 查询通知配置
	var notification im.NotificationConfig
	err := global.KUBEGALE_DB.Where("id = ?", alert.NotificationID).First(&notification).Error
	if err != nil {
		return fmt.Errorf("查询通知配置失败: %w", err)
	}

	// 根据通知类型发送消息
	switch notification.Type {
	case im.NotificationTypeFeiShu:
		var feiShuConfig im.FeiShuConfig
		if err := global.KUBEGALE_DB.Where("notification_config_id = ?", notification.ID).First(&feiShuConfig).Error; err != nil {
			return fmt.Errorf("查询飞书配置失败: %w", err)
		}
		// 构建通知配置详情
		config := response.NotificationDetailConfig{
			ID:                 notification.ID,
			Name:               notification.Name,
			Type:               notification.Type,
			NotificationPolicy: notification.NotificationPolicy,
			SendDailyStats:     notification.SendDailyStats,
			CreatedAt:          notification.CreatedAt,
			UpdatedAt:          notification.UpdatedAt,
			RobotURL:           feiShuConfig.RobotURL,
		}

		// 构建卡片内容详情
		cardContent := response.CardContentDetail{
			ID:                 alert.ID,
			NotificationID:     alert.NotificationID,
			AlertLevel:         alert.AlertLevel,
			AlertName:          alert.AlertName,
			NotificationPolicy: alert.NotificationPolicy,
			AlertContent:       alert.AlertContent,
			AlertTime:          alert.AlertTime,
			NotifiedUsers:      alert.NotifiedUsers,
			LastSimilarAlert:   alert.LastSimilarAlert,
			AlertHandler:       alert.AlertHandler,
			ClaimAlert:         alert.ClaimAlert,
			ResolveAlert:       alert.ResolveAlert,
			MuteAlert:          alert.MuteAlert,
			UnresolvedAlert:    alert.UnresolvedAlert,
		}

		// 发送飞书消息
		err = messageIm.MessageServiceApp.SendFeiShuMessage(config, cardContent, alert.AlertContent)
		if err != nil {
			return fmt.Errorf("发送飞书消息失败: %w", err)
		}
	default:
		return errors.New("不支持的通知类型")
	}

	return nil
}

// SendAlertToAll 向所有通知配置发送告警
func (s *AlertService) SendAlertToAll(alertInfo AlertInfo) error {
	// 获取所有通知配置
	var notifications []im.NotificationConfig

	// 查询所有通知配置
	if err := global.KUBEGALE_DB.Find(&notifications).Error; err != nil {
		return err
	}

	// 发送通知
	for _, notification := range notifications {
		if notification.Type != im.NotificationTypeFeiShu {
			continue
		}

		var feiShuConfig im.FeiShuConfig
		if err := global.KUBEGALE_DB.Where("notification_config_id = ?", notification.ID).First(&feiShuConfig).Error; err != nil {
			global.KUBEGALE_LOG.Error("查询飞书配置失败",
				zap.String("config", notification.Name),
				zap.Error(err))
			continue
		}

		// 创建卡片内容
		cardContent := &im.CardContentConfig{
			NotificationID:  notification.ID,
			AlertLevel:      alertInfo.Level,
			AlertName:       alertInfo.Name,
			AlertContent:    alertInfo.Content,
			AlertTime:       alertInfo.Time,
			NotifiedUsers:   alertInfo.NotifyUsers,
			AlertHandler:    alertInfo.Handler,
			ClaimAlert:      false,
			ResolveAlert:    false,
			MuteAlert:       false,
			UnresolvedAlert: true,
		}

		// 确保告警时间不为零值
		if cardContent.AlertTime.IsZero() {
			cardContent.AlertTime = time.Now()
		}

		// 保存卡片内容
		if err := global.KUBEGALE_DB.Create(cardContent).Error; err != nil {
			global.KUBEGALE_LOG.Error("创建卡片内容失败",
				zap.String("config", notification.Name),
				zap.Error(err))
			continue
		}

		// 发送告警
		if err := s.SendAlert(cardContent); err != nil {
			global.KUBEGALE_LOG.Error("发送飞书告警失败",
				zap.String("config", notification.Name),
				zap.Error(err))
		}
	}

	return nil
}

// GetAlertList 获取告警列表
func (s *AlertService) GetAlertList(notificationID uint) ([]im.CardContentConfig, error) {
	var alerts []im.CardContentConfig
	err := global.KUBEGALE_DB.Where("notification_id = ?", notificationID).Find(&alerts).Error
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

// GetAlertById 根据ID获取告警
func (s *AlertService) GetAlertById(id uint) (*im.CardContentConfig, error) {
	var alert im.CardContentConfig
	err := global.KUBEGALE_DB.Where("id = ?", id).First(&alert).Error
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateAlert 更新告警
func (s *AlertService) UpdateAlert(alert *im.CardContentConfig) error {
	return global.KUBEGALE_DB.Model(alert).Updates(alert).Error
}

// DeleteAlert 删除告警
func (s *AlertService) DeleteAlert(id uint) error {
	return global.KUBEGALE_DB.Delete(&im.CardContentConfig{}, id).Error
}
