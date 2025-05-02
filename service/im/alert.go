package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"KubeGale/model/im/response"
	messageIm "KubeGale/utils/im"
	"errors"
	"fmt"
	"time"
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

// @function: SendAlert
// @description: 发送告警通知
func (alertService *AlertService) SendAlert(notificationID uint, notificationType string, alertInfo AlertInfo) error {
	// 获取通知配置
	var cardContent im.CardContentConfig

	// 查询卡片内容配置
	global.KUBEGALE_DB.Where("notification_id = ?", notificationID).First(&cardContent)

	// 更新卡片内容
	cardContent.AlertLevel = alertInfo.Level
	cardContent.AlertName = alertInfo.Name
	cardContent.AlertContent = alertInfo.Content
	cardContent.AlertTime = alertInfo.Time
	// 确保告警时间不为零值
	if cardContent.AlertTime.IsZero() {
		cardContent.AlertTime = time.Now()
	}
	cardContent.NotifiedUsers = alertInfo.NotifyUsers
	cardContent.AlertHandler = alertInfo.Handler
	cardContent.ClaimAlert = false
	cardContent.ResolveAlert = false
	cardContent.MuteAlert = false
	cardContent.UnresolvedAlert = true

	// 保存卡片内容
	if cardContent.ID != 0 {
		global.KUBEGALE_DB.Save(&cardContent)
	} else {
		cardContent.NotificationID = notificationID
		global.KUBEGALE_DB.Create(&cardContent)
	}

	// 转换卡片内容为响应格式
	cardContentDetail := response.CardContentDetail{
		ID:                 cardContent.ID,
		NotificationID:     cardContent.NotificationID,
		AlertLevel:         cardContent.AlertLevel,
		AlertName:          cardContent.AlertName,
		NotificationPolicy: cardContent.NotificationPolicy,
		AlertContent:       cardContent.AlertContent,
		AlertTime:          cardContent.AlertTime,
		NotifiedUsers:      cardContent.NotifiedUsers,
		LastSimilarAlert:   cardContent.LastSimilarAlert,
		AlertHandler:       cardContent.AlertHandler,
		ClaimAlert:         cardContent.ClaimAlert,
		ResolveAlert:       cardContent.ResolveAlert,
		MuteAlert:          cardContent.MuteAlert,
		UnresolvedAlert:    cardContent.UnresolvedAlert,
	}

	// 根据通知类型发送告警
	switch notificationType {
	case im.NotificationTypeDingTalk:
		var dingTalk im.DingTalkConfig
		if err := global.KUBEGALE_DB.Where("id = ?", notificationID).First(&dingTalk).Error; err != nil {
			return err
		}
		config := response.NotificationDetailConfig{
			ID:                 dingTalk.ID,
			Name:               dingTalk.Name,
			Type:               dingTalk.Type,
			NotificationPolicy: dingTalk.NotificationPolicy,
			SendDailyStats:     dingTalk.SendDailyStats,
			CreatedAt:          dingTalk.CreatedAt,
			UpdatedAt:          dingTalk.UpdatedAt,
			SignatureKey:       dingTalk.SignatureKey,
			RobotURL:           dingTalk.RobotURL,
		}
		return messageIm.MessageServiceApp.SendDingTalkMessage(config, cardContentDetail, "")
	case im.NotificationTypeFeiShu:
		var feiShu im.FeiShuConfig
		if err := global.KUBEGALE_DB.Where("id = ?", notificationID).First(&feiShu).Error; err != nil {
			return err
		}
		config := response.NotificationDetailConfig{
			ID:                 feiShu.ID,
			Name:               feiShu.Name,
			Type:               feiShu.Type,
			NotificationPolicy: feiShu.NotificationPolicy,
			SendDailyStats:     feiShu.SendDailyStats,
			CreatedAt:          feiShu.CreatedAt,
			UpdatedAt:          feiShu.UpdatedAt,
			RobotURL:           feiShu.RobotURL,
		}
		return messageIm.MessageServiceApp.SendFeiShuMessage(config, cardContentDetail, "")
	default:
		return errors.New("不支持的通知类型")
	}
}

// @function: SendAlertToAll
// @description: 向所有通知配置发送告警
func (alertService *AlertService) SendAlertToAll(alertInfo AlertInfo) error {
	// 获取所有通知配置
	var dingTalkConfigs []im.DingTalkConfig
	var feiShuConfigs []im.FeiShuConfig

	// 查询钉钉配置
	if err := global.KUBEGALE_DB.Find(&dingTalkConfigs).Error; err != nil {
		return err
	}

	// 查询飞书配置
	if err := global.KUBEGALE_DB.Find(&feiShuConfigs).Error; err != nil {
		return err
	}

	// 发送钉钉通知
	for _, config := range dingTalkConfigs {
		err := alertService.SendAlert(config.ID, im.NotificationTypeDingTalk, alertInfo)
		if err != nil {
			global.KUBEGALE_LOG.Error(fmt.Sprintf("发送钉钉告警失败: %s", err.Error()))
		}
	}

	// 发送飞书通知
	for _, config := range feiShuConfigs {
		err := alertService.SendAlert(config.ID, im.NotificationTypeFeiShu, alertInfo)
		if err != nil {
			global.KUBEGALE_LOG.Error(fmt.Sprintf("发送飞书告警失败: %s", err.Error()))
		}
	}

	return nil
}
