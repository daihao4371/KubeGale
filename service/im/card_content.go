package im

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/im"
	"errors"
)

type CardContentService struct{}

var CardContentServiceApp = new(CardContentService)

// @function: CreateCardContent
// @description: 创建卡片内容配置
func (cardContentService *CardContentService) CreateCardContent(cardContent im.CardContentConfig) (err error) {
	// 检查关联的通知配置是否存在
	var count int64
	err = global.KUBEGALE_DB.Model(&im.NotificationConfig{}).Where("id = ?", cardContent.NotificationID).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("关联的通知配置不存在")
	}

	return global.KUBEGALE_DB.Create(&cardContent).Error
}

// @function: UpdateCardContent
// @description: 更新卡片内容配置
func (cardContentService *CardContentService) UpdateCardContent(cardContent im.CardContentConfig) (err error) {
	var oldContent im.CardContentConfig

	// 查询原有配置
	err = global.KUBEGALE_DB.Where("id = ?", cardContent.ID).First(&oldContent).Error
	if err != nil {
		return errors.New("卡片内容配置不存在")
	}

	// 如果修改了通知ID，检查新的通知配置是否存在
	if cardContent.NotificationID != oldContent.NotificationID {
		var count int64
		err = global.KUBEGALE_DB.Model(&im.NotificationConfig{}).Where("id = ?", cardContent.NotificationID).Count(&count).Error
		if err != nil {
			return err
		}
		if count == 0 {
			return errors.New("关联的通知配置不存在")
		}
	}

	// 保留原始的创建时间
	cardContent.CreatedAt = oldContent.CreatedAt

	// 使用Updates更新特定字段
	return global.KUBEGALE_DB.Model(&oldContent).Updates(map[string]interface{}{
		"notification_id":     cardContent.NotificationID,
		"alert_level":         cardContent.AlertLevel,
		"alert_name":          cardContent.AlertName,
		"notification_policy": cardContent.NotificationPolicy,
		"alert_content":       cardContent.AlertContent,
		"alert_time":          cardContent.AlertTime,
		"notified_users":      cardContent.NotifiedUsers,
		"last_similar_alert":  cardContent.LastSimilarAlert,
		"alert_handler":       cardContent.AlertHandler,
		"claim_alert":         cardContent.ClaimAlert,
		"resolve_alert":       cardContent.ResolveAlert,
		"mute_alert":          cardContent.MuteAlert,
		"unresolved_alert":    cardContent.UnresolvedAlert,
	}).Error
}

// @function: DeleteCardContent
// @description: 删除卡片内容配置
func (cardContentService *CardContentService) DeleteCardContent(id uint) (err error) {
	var cardContent im.CardContentConfig

	// 查询配置是否存在
	err = global.KUBEGALE_DB.Where("id = ?", id).First(&cardContent).Error
	if err != nil {
		return err
	}

	return global.KUBEGALE_DB.Delete(&cardContent).Error
}

// @function: GetCardContentList
// @description: 获取卡片内容配置列表
func (cardContentService *CardContentService) GetCardContentList(info request.PageInfo, notificationID uint) (list []im.CardContentConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.KUBEGALE_DB.Model(&im.CardContentConfig{})

	// 如果指定了通知ID，则只查询该通知的卡片内容
	if notificationID > 0 {
		db = db.Where("notification_id = ?", notificationID)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询卡片内容列表
	err = db.Order("id desc").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

// @function: GetCardContentById
// @description: 根据ID获取卡片内容配置
func (cardContentService *CardContentService) GetCardContentById(id uint) (cardContent im.CardContentConfig, err error) {
	err = global.KUBEGALE_DB.Where("id = ?", id).First(&cardContent).Error
	return cardContent, err
}
