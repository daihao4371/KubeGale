package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"KubeGale/model/im/request"
	"KubeGale/model/im/response"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type NotificationService struct{}

var NotificationServiceApp = new(NotificationService)

// @function: CreateDingTalk
// @description: 创建钉钉通知配置
func (notificationService *NotificationService) CreateDingTalk(req request.CreateDingTalkRequest) (dingTalk im.DingTalkConfig, err error) {
	// 检查名称是否已存在
	if !errors.Is(global.KUBEGALE_DB.Where("name = ?", req.Name).First(&im.NotificationConfig{}).Error, gorm.ErrRecordNotFound) {
		return dingTalk, errors.New("通知名称已存在")
	}

	// 创建钉钉通知配置
	dingTalk = im.DingTalkConfig{
		NotificationConfig: im.NotificationConfig{
			Name:               req.Name,
			Type:               im.NotificationTypeDingTalk,
			NotificationPolicy: req.NotificationPolicy,
			SendDailyStats:     req.SendDailyStats,
		},
		SignatureKey: req.SignatureKey,
		RobotURL:     req.RobotURL,
	}

	// 开启事务
	err = global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 创建钉钉配置
		if err := tx.Create(&dingTalk).Error; err != nil {
			return err
		}

		// 如果有卡片内容配置，则创建
		if req.CardContent.AlertLevel != "" {
			cardContent := req.CardContent
			cardContent.NotificationID = dingTalk.ID
			if err := tx.Create(&cardContent).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return dingTalk, err
}

// @function: CreateFeiShu
// @description: 创建飞书通知配置
func (notificationService *NotificationService) CreateFeiShu(req request.CreateFeiShuRequest) (feiShu im.FeiShuConfig, err error) {
	// 检查名称是否已存在
	if !errors.Is(global.KUBEGALE_DB.Where("name = ?", req.Name).First(&im.NotificationConfig{}).Error, gorm.ErrRecordNotFound) {
		return feiShu, errors.New("通知名称已存在")
	}

	// 创建飞书通知配置
	feiShu = im.FeiShuConfig{
		NotificationConfig: im.NotificationConfig{
			Name:               req.Name,
			Type:               im.NotificationTypeFeiShu,
			NotificationPolicy: req.NotificationPolicy,
			SendDailyStats:     req.SendDailyStats,
		},
		RobotURL: req.RobotURL,
	}

	// 开启事务
	err = global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 创建飞书配置
		if err := tx.Create(&feiShu).Error; err != nil {
			return err
		}

		// 如果有卡片内容配置，则创建
		if req.CardContent.AlertLevel != "" {
			cardContent := req.CardContent
			cardContent.NotificationID = feiShu.ID
			if err := tx.Create(&cardContent).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return feiShu, err
}

// @function: UpdateDingTalk
// @description: 更新钉钉通知配置
func (notificationService *NotificationService) UpdateDingTalk(req request.UpdateDingTalkRequest) (dingTalk im.DingTalkConfig, err error) {
	var oldConfig im.DingTalkConfig

	// 查询原有配置
	err = global.KUBEGALE_DB.Where("id = ?", req.ID).First(&oldConfig).Error
	if err != nil {
		return dingTalk, errors.New("通知配置不存在")
	}

	// 检查名称是否已被其他配置使用
	if req.Name != oldConfig.Name {
		var count int64
		global.KUBEGALE_DB.Model(&im.NotificationConfig{}).Where("name = ? AND id != ?", req.Name, req.ID).Count(&count)
		if count > 0 {
			return dingTalk, errors.New("通知名称已存在")
		}
	}

	// 更新钉钉通知配置
	updateMap := map[string]interface{}{
		"name":                req.Name,
		"notification_policy": req.NotificationPolicy,
		"send_daily_stats":    req.SendDailyStats,
		"signature_key":       req.SignatureKey,
		"robot_url":           req.RobotURL,
	}

	// 开启事务
	err = global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 更新钉钉配置
		if err := tx.Model(&oldConfig).Updates(updateMap).Error; err != nil {
			return err
		}

		// 如果有卡片内容配置，则更新或创建
		if req.CardContent.AlertLevel != "" {
			var cardContent im.CardContentConfig
			if err := tx.Where("notification_id = ?", req.ID).First(&cardContent).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 创建新的卡片内容
					cardContent = req.CardContent
					cardContent.NotificationID = req.ID
					if err := tx.Create(&cardContent).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				// 更新现有卡片内容
				if err := tx.Model(&cardContent).Updates(req.CardContent).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return dingTalk, err
	}

	// 获取更新后的配置
	err = global.KUBEGALE_DB.Where("id = ?", req.ID).First(&dingTalk).Error
	return dingTalk, err
}

// @function: UpdateFeiShu
// @description: 更新飞书通知配置
func (notificationService *NotificationService) UpdateFeiShu(req request.UpdateFeiShuRequest) (feiShu im.FeiShuConfig, err error) {
	var oldConfig im.FeiShuConfig

	// 查询原有配置
	err = global.KUBEGALE_DB.Where("id = ?", req.ID).First(&oldConfig).Error
	if err != nil {
		return feiShu, errors.New("通知配置不存在")
	}

	// 检查名称是否已被其他配置使用
	if req.Name != oldConfig.Name {
		var count int64
		global.KUBEGALE_DB.Model(&im.NotificationConfig{}).Where("name = ? AND id != ?", req.Name, req.ID).Count(&count)
		if count > 0 {
			return feiShu, errors.New("通知名称已存在")
		}
	}

	// 更新飞书通知配置
	updateMap := map[string]interface{}{
		"name":                req.Name,
		"notification_policy": req.NotificationPolicy,
		"send_daily_stats":    req.SendDailyStats,
		"robot_url":           req.RobotURL,
	}

	// 开启事务
	err = global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 更新飞书配置
		if err := tx.Model(&oldConfig).Updates(updateMap).Error; err != nil {
			return err
		}

		// 如果有卡片内容配置，则更新或创建
		if req.CardContent.AlertLevel != "" {
			var cardContent im.CardContentConfig
			if err := tx.Where("notification_id = ?", req.ID).First(&cardContent).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 创建新的卡片内容
					cardContent = req.CardContent
					cardContent.NotificationID = req.ID
					if err := tx.Create(&cardContent).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				// 更新现有卡片内容
				if err := tx.Model(&cardContent).Updates(req.CardContent).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return feiShu, err
	}

	// 获取更新后的配置
	err = global.KUBEGALE_DB.Where("id = ?", req.ID).First(&feiShu).Error
	return feiShu, err
}

// @function: DeleteNotification
// @description: 删除通知配置
func (notificationService *NotificationService) DeleteNotification(id uint, notificationType string) error {
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 删除卡片内容配置
		if err := tx.Where("notification_id = ?", id).Delete(&im.CardContentConfig{}).Error; err != nil {
			return err
		}

		// 根据通知类型删除对应的配置
		switch notificationType {
		case im.NotificationTypeDingTalk:
			if err := tx.Where("id = ?", id).Delete(&im.DingTalkConfig{}).Error; err != nil {
				return err
			}
		case im.NotificationTypeFeiShu:
			if err := tx.Where("id = ?", id).Delete(&im.FeiShuConfig{}).Error; err != nil {
				return err
			}
		default:
			return errors.New("不支持的通知类型")
		}

		return nil
	})
}

// @function: GetNotificationList
// @description: 获取通知配置列表
func (notificationService *NotificationService) GetNotificationList(params request.SearchNotificationParams) (list []interface{}, total int64, err error) {
	limit := params.PageSize
	offset := params.PageSize * (params.Page - 1)

	// 查询条件
	db := global.KUBEGALE_DB.Model(&im.NotificationConfig{})

	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Type != "" {
		db = db.Where("type = ?", params.Type)
	}

	if params.NotificationPolicy != "" {
		db = db.Where("notification_policy LIKE ?", "%"+params.NotificationPolicy+"%")
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 排序
	orderStr := "id desc"
	if params.OrderKey != "" {
		orderMap := map[string]bool{
			"id":                  true,
			"name":                true,
			"type":                true,
			"notification_policy": true,
			"created_at":          true,
			"updated_at":          true,
		}

		if !orderMap[params.OrderKey] {
			return nil, 0, fmt.Errorf("非法的排序字段: %v", params.OrderKey)
		}

		orderStr = params.OrderKey
		if params.Desc {
			orderStr += " desc"
		}
	}

	// 查询通知配置
	var configs []im.NotificationConfig
	err = db.Order(orderStr).Limit(limit).Offset(offset).Find(&configs).Error
	if err != nil {
		return nil, 0, err
	}

	// 根据类型获取详细配置
	for _, config := range configs {
		var item interface{}
		var cardContent im.CardContentConfig

		// 查询卡片内容
		global.KUBEGALE_DB.Where("notification_id = ?", config.ID).First(&cardContent)

		switch config.Type {
		case im.NotificationTypeDingTalk:
			var dingTalk im.DingTalkConfig
			if err := global.KUBEGALE_DB.Where("id = ?", config.ID).First(&dingTalk).Error; err != nil {
				continue
			}
			item = response.DingTalkResponse{
				Config:      dingTalk,
				CardContent: cardContent,
			}
		case im.NotificationTypeFeiShu:
			var feiShu im.FeiShuConfig
			if err := global.KUBEGALE_DB.Where("id = ?", config.ID).First(&feiShu).Error; err != nil {
				continue
			}
			item = response.FeiShuResponse{
				Config:      feiShu,
				CardContent: cardContent,
			}
		}

		if item != nil {
			list = append(list, item)
		}
	}

	return list, total, nil
}

// @function: GetNotificationById
// @description: 根据ID获取通知配置
func (notificationService *NotificationService) GetNotificationById(id uint, notificationType string) (interface{}, error) {
	var cardContent im.CardContentConfig
	
	// 查询卡片内容
	err := global.KUBEGALE_DB.Where("notification_id = ?", id).First(&cardContent).Error
	// 如果卡片内容不存在，这里不返回错误，而是使用空的卡片内容
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询卡片内容失败: %w", err)
	}
	
	switch notificationType {
	case im.NotificationTypeDingTalk:
		var dingTalk im.DingTalkConfig
		if err := global.KUBEGALE_DB.Where("id = ?", id).First(&dingTalk).Error; err != nil {
			return nil, err
		}
		return response.DingTalkResponse{
			Config:      dingTalk,
			CardContent: cardContent,
		}, nil
	case im.NotificationTypeFeiShu:
		var feiShu im.FeiShuConfig
		if err := global.KUBEGALE_DB.Where("id = ?", id).First(&feiShu).Error; err != nil {
			return nil, err
		}
		return response.FeiShuResponse{
			Config:      feiShu,
			CardContent: cardContent,
		}, nil
	default:
		return nil, errors.New("不支持的通知类型")
	}
}

// @function: TestNotification
// @description: 测试通知发送
func (notificationService *NotificationService) TestNotification(req request.TestNotificationRequest) (response.TestNotificationResponse, error) {
	// 根据ID和类型获取通知配置
	notification, err := notificationService.GetNotificationById(req.ID, req.Type)
	if err != nil {
		return response.TestNotificationResponse{
			Success: false,
			Message: "获取通知配置失败: " + err.Error(),
		}, err
	}

	// 准备测试消息
	testMessage := "这是一条测试消息，请忽略。"
	if req.Message != "" {
		testMessage = req.Message
	}

	// 根据通知类型发送测试消息
	switch req.Type {
	case im.NotificationTypeDingTalk:
		dingTalkResp := notification.(response.DingTalkResponse)
		err = MessageServiceApp.SendDingTalkMessage(dingTalkResp.Config, dingTalkResp.CardContent, testMessage)
	case im.NotificationTypeFeiShu:
		feiShuResp := notification.(response.FeiShuResponse)
		err = MessageServiceApp.SendFeiShuMessage(feiShuResp.Config, feiShuResp.CardContent, testMessage)
	default:
		return response.TestNotificationResponse{
			Success: false,
			Message: "不支持的通知类型",
		}, errors.New("不支持的通知类型")
	}

	if err != nil {
		return response.TestNotificationResponse{
			Success: false,
			Message: "发送测试消息失败: " + err.Error(),
		}, err
	}

	return response.TestNotificationResponse{
		Success: true,
		Message: "测试消息发送成功",
	}, nil
}
