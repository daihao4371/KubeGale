package im

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/im"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CardContentService 卡片内容服务
type CardContentService struct{}

var CardContentServiceApp = new(CardContentService)

// @function: CreateCardContent
// @description: 创建卡片内容配置
func (cardContentService *CardContentService) CreateCardContent(req im.CardContentConfig) error {
	// 首先检查通用通知配置是否存在
	var notificationConfig im.NotificationConfig
	err := global.KUBEGALE_DB.Table("im_notification_configs").Where("id = ?", req.NotificationID).First(&notificationConfig).Error
	if err != nil {
		global.KUBEGALE_LOG.Error("查询通知配置失败",
			zap.Uint("notification_id", req.NotificationID),
			zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 尝试查询飞书配置表
			var feiShuConfig im.FeiShuConfig
			err = global.KUBEGALE_DB.Where("id = ?", req.NotificationID).First(&feiShuConfig).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("通知配置不存在")
				}
				return fmt.Errorf("查询飞书配置失败: %w", err)
			}
			// 找到了飞书配置，可以继续处理
			notificationConfig.ID = feiShuConfig.ID
			notificationConfig.Type = feiShuConfig.Type
		} else {
			return fmt.Errorf("查询通知配置失败: %w", err)
		}
	}

	// 设置告警时间为当前时间，如果未提供
	if req.AlertTime.IsZero() {
		req.AlertTime = time.Now()
	}

	// 创建卡片内容
	err = global.KUBEGALE_DB.Create(&req).Error
	if err != nil {
		global.KUBEGALE_LOG.Error("创建卡片内容失败",
			zap.Uint("notification_id", req.NotificationID),
			zap.Error(err))
		return fmt.Errorf("创建卡片内容失败: %w", err)
	}

	global.KUBEGALE_LOG.Info("创建卡片内容成功",
		zap.Uint("notification_id", req.NotificationID),
		zap.String("alert_name", req.AlertName))
	return nil
}

// @function: UpdateCardContent
// @description: 更新卡片内容配置
func (cardContentService *CardContentService) UpdateCardContent(req im.CardContentConfig) error {
	// 检查卡片内容是否存在
	var cardContent im.CardContentConfig
	if err := global.KUBEGALE_DB.Where("id = ?", req.ID).First(&cardContent).Error; err != nil {
		return errors.New("卡片内容不存在")
	}

	// 更新卡片内容
	return global.KUBEGALE_DB.Model(&cardContent).Updates(req).Error
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

// @function: GetCardContentByNotificationId
// @description: 根据通知ID获取卡片内容
func (s *CardContentService) GetCardContentByNotificationId(notificationId uint) (*im.CardContentConfig, error) {
	// 首先检查通知配置是否存在
	var notification im.FeiShuConfig
	if err := global.KUBEGALE_DB.Where("id = ?", notificationId).First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("通知配置不存在")
		}
		return nil, err
	}

	// 查询卡片内容
	var cardContent im.CardContentConfig
	err := global.KUBEGALE_DB.Where("notification_id = ?", notificationId).First(&cardContent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果卡片内容不存在，返回空的卡片内容
			return &im.CardContentConfig{
				NotificationID: notificationId,
			}, nil
		}
		return nil, err
	}
	return &cardContent, nil
}
