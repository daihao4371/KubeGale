package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"KubeGale/model/im/request"
	"KubeGale/model/im/response"
	messageIm "KubeGale/utils/im"
	"errors"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

type NotificationService struct{}

var NotificationServiceApp = new(NotificationService)

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
			// 设置告警时间为当前时间，避免数据库插入空时间
			if cardContent.AlertTime.IsZero() {
				cardContent.AlertTime = time.Now()
			}
			if err := tx.Create(&cardContent).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return feiShu, err
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
func (notificationService *NotificationService) GetNotificationList(params request.SearchNotificationParams) (list []response.NotificationResponse, total int64, err error) {
	limit := params.PageSize
	offset := params.PageSize * (params.Page - 1)
	list = make([]response.NotificationResponse, 0)

	// 创建查询
	feiShuQuery := global.KUBEGALE_DB.Model(&im.FeiShuConfig{}).Select("id, name, 'feishu' as type, notification_policy, robot_url, created_at")

	// 添加搜索条件
	if params.Name != "" {
		feiShuQuery = feiShuQuery.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.Type != "" {
		if params.Type != im.NotificationTypeFeiShu {
			feiShuQuery = feiShuQuery.Where("1 = 0") // 不查询数据
		}
	}
	if params.NotificationPolicy != "" {
		feiShuQuery = feiShuQuery.Where("notification_policy = ?", params.NotificationPolicy)
	}

	// 计算总数
	var feiShuCount int64
	if err = feiShuQuery.Count(&feiShuCount).Error; err != nil {
		return nil, 0, err
	}
	total = feiShuCount

	// 构建排序
	order := "created_at"
	if params.OrderKey != "" {
		order = params.OrderKey
	}
	if params.Desc {
		order = order + " desc"
	} else {
		order = order + " asc"
	}

	// 查询飞书配置
	var feiShuConfigs []struct {
		ID                 uint      `gorm:"column:id"`
		Name               string    `gorm:"column:name"`
		Type               string    `gorm:"column:type"`
		NotificationPolicy string    `gorm:"column:notification_policy"`
		RobotURL           string    `gorm:"column:robot_url"`
		CreatedAt          time.Time `gorm:"column:created_at"`
	}
	if err = feiShuQuery.Order(order).Limit(limit).Offset(offset).Find(&feiShuConfigs).Error; err != nil {
		return nil, 0, err
	}

	// 合并数据
	for _, config := range feiShuConfigs {
		list = append(list, response.NotificationResponse{
			ID:                 config.ID,
			Name:               config.Name,
			Type:               config.Type,
			NotificationPolicy: config.NotificationPolicy,
			RobotURL:           config.RobotURL,
			CreatedAt:          config.CreatedAt,
		})
	}

	// 对合并后的数据进行排序
	sort.Slice(list, func(i, j int) bool {
		if params.Desc {
			return list[i].CreatedAt.After(list[j].CreatedAt)
		}
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})

	// 如果合并后的数据超过了页面大小，需要截取
	if len(list) > limit {
		list = list[:limit]
	}

	return list, total, nil
}

// @function: GetNotificationById
// @description: 根据ID获取通知配置
func (notificationService *NotificationService) GetNotificationById(id uint, notificationType string) (*response.NotificationDetailResponse, error) {
	var cardContent im.CardContentConfig
	var config response.NotificationDetailConfig

	// 查询卡片内容
	err := global.KUBEGALE_DB.Where("notification_id = ?", id).First(&cardContent).Error
	// 如果卡片内容不存在，这里不返回错误，而是使用空的卡片内容
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询卡片内容失败: %w", err)
	}

	// 根据类型查询对应的配置
	switch notificationType {
	case im.NotificationTypeFeiShu:
		var feiShu im.FeiShuConfig
		if err := global.KUBEGALE_DB.Where("id = ?", id).First(&feiShu).Error; err != nil {
			return nil, err
		}
		config = response.NotificationDetailConfig{
			ID:                 feiShu.ID,
			Name:               feiShu.Name,
			Type:               feiShu.Type,
			NotificationPolicy: feiShu.NotificationPolicy,
			SendDailyStats:     feiShu.SendDailyStats,
			CreatedAt:          feiShu.CreatedAt,
			UpdatedAt:          feiShu.UpdatedAt,
			RobotURL:           feiShu.RobotURL,
		}
	default:
		return nil, errors.New("不支持的通知类型")
	}

	return &response.NotificationDetailResponse{
		Config: config,
		CardContent: response.CardContentDetail{
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
		},
	}, nil
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
	case im.NotificationTypeFeiShu:
		err = messageIm.MessageServiceApp.SendFeiShuMessage(notification.Config, notification.CardContent, testMessage)
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
