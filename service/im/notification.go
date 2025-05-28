package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"KubeGale/model/im/request"
	"KubeGale/model/im/response"
	messageIm "KubeGale/utils/im"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type NotificationService struct{}

var NotificationServiceApp = new(NotificationService)

// @function: CreateFeiShu
// @description: 创建飞书通知配置
func (notificationService *NotificationService) CreateFeiShu(req request.CreateFeiShuRequest) (*im.FeiShuConfig, error) {
	var createdFeiShuConfig im.FeiShuConfig

	err := global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create NotificationConfig
		notificationCfg := im.NotificationConfig{
			Name:               req.Name,
			Type:               im.NotificationTypeFeiShu, // Use constant
			NotificationPolicy: strings.Join(req.NotifyEvents, ","),
			SendDailyStats:     req.SendDailyStats,
		}
		if err := tx.Create(&notificationCfg).Error; err != nil {
			return fmt.Errorf("failed to create notification config: %w", err)
		}

		// 2. Create FeiShuConfig
		feiShuSpecificConfig := im.FeiShuConfig{
			NotificationConfigID: notificationCfg.ID, // Link to the created NotificationConfig
			RobotURL:             req.WebhookURL,
		}
		if err := tx.Create(&feiShuSpecificConfig).Error; err != nil {
			return fmt.Errorf("failed to create feishu config: %w", err)
		}

		// Store the ID for later retrieval
		createdFeiShuConfig.ID = feiShuSpecificConfig.ID
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Query and return the created FeiShuConfig with its associated NotificationConfig preloaded
	var resultFeiShuConfig im.FeiShuConfig
	if err := global.KUBEGALE_DB.Preload("NotificationConfig").First(&resultFeiShuConfig, createdFeiShuConfig.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch created feishu config with notification config: %w", err)
	}

	return &resultFeiShuConfig, nil
}

// @function: UpdateFeiShu
// @description: 更新飞书通知配置
func (notificationService *NotificationService) UpdateFeiShu(req request.UpdateFeiShuRequest) (im.FeiShuConfig, error) {
	var feiShuConfigToUpdate im.FeiShuConfig

	err := global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch existing FeiShuConfig with NotificationConfig
		if err := tx.Preload("NotificationConfig").First(&feiShuConfigToUpdate, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("飞书通知配置不存在")
			}
			return fmt.Errorf("failed to find feishu config: %w", err)
		}

		// 2. Name Uniqueness Check
		if req.Name != feiShuConfigToUpdate.NotificationConfig.Name {
			var count int64
			err := tx.Model(&im.NotificationConfig{}).
				Where("name = ? AND type = ? AND id != ?", req.Name, im.NotificationTypeFeiShu, feiShuConfigToUpdate.NotificationConfigID).
				Count(&count).Error
			if err != nil {
				return fmt.Errorf("failed to check name uniqueness: %w", err)
			}
			if count > 0 {
				return errors.New("通知名称已存在")
			}
		}

		// 3. Update NotificationConfig fields
		feiShuConfigToUpdate.NotificationConfig.Name = req.Name
		feiShuConfigToUpdate.NotificationConfig.NotificationPolicy = req.NotificationPolicy // Directly from UpdateFeiShuRequest
		feiShuConfigToUpdate.NotificationConfig.SendDailyStats = req.SendDailyStats

		if err := tx.Save(&feiShuConfigToUpdate.NotificationConfig).Error; err != nil {
			return fmt.Errorf("failed to save notification config: %w", err)
		}

		// 4. Update FeiShuConfig fields
		feiShuConfigToUpdate.RobotURL = req.RobotURL
		if err := tx.Save(&feiShuConfigToUpdate).Error; err != nil {
			return fmt.Errorf("failed to save feishu config: %w", err)
		}

		// 5. Card Content Handling (Simplified)
		if req.CardContent.AlertLevel != "" { // Check if CardContent is provided in the request
			var cardContent im.CardContentConfig
			// Attempt to find existing card content
			err := tx.Where("notification_id = ?", feiShuConfigToUpdate.NotificationConfigID).First(&cardContent).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Not found, create new if AlertLevel is present
					newCardContent := req.CardContent
					newCardContent.NotificationID = feiShuConfigToUpdate.NotificationConfigID // Set the foreign key
					if createErr := tx.Create(&newCardContent).Error; createErr != nil {
						return fmt.Errorf("failed to create card content: %w", createErr)
					}
				} else {
					// Other error finding card content
					return fmt.Errorf("failed to query card content: %w", err)
				}
			} else {
				// Found, update existing card content
				// Important: Ensure req.CardContent has an ID if you want to update a specific record by its ID.
				// GORM's Updates will update based on the primary key if present in req.CardContent,
				// or based on the Model(&cardContent) condition if not.
				// For simplicity, we assume Updates will correctly update the found cardContent.
				if updateErr := tx.Model(&cardContent).Updates(req.CardContent).Error; updateErr != nil {
					return fmt.Errorf("failed to update card content: %w", updateErr)
				}
			}
		}
		return nil
	})

	if err != nil {
		return im.FeiShuConfig{}, err // Return zero value of FeiShuConfig on error
	}

	// The feiShuConfigToUpdate already has NotificationConfig preloaded from the start of the transaction.
	return feiShuConfigToUpdate, nil
}

// @function: CreateDingTalk
// @description: 创建钉钉通知配置
func (notificationService *NotificationService) CreateDingTalk(req request.CreateDingTalkRequest) (*im.DingTalkConfig, error) {
	var dingTalkConfig im.DingTalkConfig
	err := global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		notificationCfg := im.NotificationConfig{
			Name:               req.Name,
			Type:               im.NotificationTypeDingTalk,
			NotificationPolicy: strings.Join(req.NotifyEvents, ","),
			SendDailyStats:     req.SendDailyStats,
		}
		if err := tx.Create(&notificationCfg).Error; err != nil {
			return fmt.Errorf("failed to create notification config: %w", err)
		}

		newDingTalkConfig := im.DingTalkConfig{
			NotificationConfigID: notificationCfg.ID,
			WebhookURL:           req.WebhookURL,
			Secret:               req.Secret,
		}
		if err := tx.Create(&newDingTalkConfig).Error; err != nil {
			return fmt.Errorf("failed to create dingtalk config: %w", err)
		}
		dingTalkConfig = newDingTalkConfig // Assign to outer scope variable
		return nil
	})

	if err != nil {
		return nil, err
	}
	// Preload NotificationConfig for the response
	err = global.KUBEGALE_DB.Preload("NotificationConfig").First(&dingTalkConfig, dingTalkConfig.ID).Error
	return &dingTalkConfig, err
}

// @function: UpdateDingTalk
// @description: 更新钉钉通知配置
func (notificationService *NotificationService) UpdateDingTalk(req request.UpdateDingTalkRequest) (*im.DingTalkConfig, error) {
	var updatedDingTalkConfig im.DingTalkConfig
	err := global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		var existingDingTalkConfig im.DingTalkConfig
		if err := tx.Preload("NotificationConfig").First(&existingDingTalkConfig, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("dingtalk configuration not found")
			}
			return fmt.Errorf("failed to fetch dingtalk config: %w", err)
		}

		// Name Uniqueness Check (if name is being changed)
		if req.Name != "" && req.Name != existingDingTalkConfig.NotificationConfig.Name {
			var count int64
			if err := tx.Model(&im.NotificationConfig{}).
				Where("name = ? AND type = ? AND id != ?", req.Name, im.NotificationTypeDingTalk, existingDingTalkConfig.NotificationConfigID).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check name uniqueness: %w", err)
			}
			if count > 0 {
				return errors.New("notification name already exists for this type")
			}
			existingDingTalkConfig.NotificationConfig.Name = req.Name
		}

		if req.NotificationPolicy != "" { // Assuming NotificationPolicy is a string in request
			existingDingTalkConfig.NotificationConfig.NotificationPolicy = req.NotificationPolicy
		}
		if req.SendDailyStats != nil {
			existingDingTalkConfig.NotificationConfig.SendDailyStats = *req.SendDailyStats
		}
		if err := tx.Save(&existingDingTalkConfig.NotificationConfig).Error; err != nil {
			return fmt.Errorf("failed to save notification config: %w", err)
		}

		if req.WebhookURL != "" {
			existingDingTalkConfig.WebhookURL = req.WebhookURL
		}
		if req.Secret != nil { // If *string is nil, means not provided; if points to empty string, means clear
			existingDingTalkConfig.Secret = *req.Secret
		}
		if err := tx.Save(&existingDingTalkConfig).Error; err != nil {
			return fmt.Errorf("failed to save dingtalk config: %w", err)
		}

		// Simplified Card Content Handling (as in UpdateFeiShu)
		if req.CardContent.AlertLevel != "" { // Check if CardContent is provided
			var cardContent im.CardContentConfig
			err := tx.Where("notification_id = ?", existingDingTalkConfig.NotificationConfigID).First(&cardContent).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) { // Create new
					newCardContent := req.CardContent
					newCardContent.NotificationID = existingDingTalkConfig.NotificationConfigID
					if createErr := tx.Create(&newCardContent).Error; createErr != nil {
						return fmt.Errorf("failed to create card content: %w", createErr)
					}
				} else { // Other error fetching card content
					return fmt.Errorf("failed to fetch card content: %w", err)
				}
			} else { // Update existing
				updateData := req.CardContent
				updateData.ID = cardContent.ID                         // Ensure we update the correct record
				updateData.NotificationID = cardContent.NotificationID // Keep original link
				if updateErr := tx.Model(&cardContent).Updates(updateData).Error; updateErr != nil {
					return fmt.Errorf("failed to update card content: %w", updateErr)
				}
			}
		}
		updatedDingTalkConfig = existingDingTalkConfig // Assign to outer scope
		return nil
	})
	return &updatedDingTalkConfig, err // updatedDingTalkConfig will have NotificationConfig preloaded
}

// @function: DeleteNotification
// @description: 删除通知配置 (id is NotificationConfigID)
func (notificationService *NotificationService) DeleteNotification(id uint, notificationType string) error { // id is NotificationConfigID
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch NotificationConfig to confirm existence and get its actual Type
		var notificationConfig im.NotificationConfig
		if err := tx.First(&notificationConfig, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("notification config not found")
			}
			return fmt.Errorf("failed to fetch notification config: %w", err)
		}

		// 2. Delete associated CardContentConfig
		if err := tx.Where("notification_id = ?", id).Delete(&im.CardContentConfig{}).Error; err != nil {
			// Allow not found for card content, but fail on other errors
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to delete card content config: %w", err)
			}
		}

		// 3. Based on notificationConfig.Type, delete specific config
		// The notificationType param can be used as a hint but notificationConfig.Type is authoritative
		actualType := notificationConfig.Type
		if notificationType != "" && notificationType != actualType {
			// Optional: Log a warning if provided type mismatches stored type, but proceed with stored type.
			// global.KUBEGALE_LOG.Warn(fmt.Sprintf("Provided type %s mismatches stored type %s for ID %d", notificationType, actualType, id))
		}

		switch actualType {
		case im.NotificationTypeFeiShu:
			if err := tx.Where("notification_config_id = ?", id).Delete(&im.FeiShuConfig{}).Error; err != nil {
				// Allow not found, but fail on other errors
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("failed to delete feishu config: %w", err)
				}
			}
		case im.NotificationTypeDingTalk:
			if err := tx.Where("notification_config_id = ?", id).Delete(&im.DingTalkConfig{}).Error; err != nil {
				// Allow not found, but fail on other errors
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("failed to delete dingtalk config: %w", err)
				}
			}
		default:
			return fmt.Errorf("unsupported notification type for deletion: %s", actualType)
		}

		// 4. Delete the NotificationConfig itself
		if err := tx.Delete(&notificationConfig).Error; err != nil { // Delete by instance
			return fmt.Errorf("failed to delete notification config: %w", err)
		}

		return nil
	})
}

// @function: GetNotificationList
// @description: 获取通知配置列表
func (notificationService *NotificationService) GetNotificationList(params request.SearchNotificationParams) (list []response.NotificationResponse, total int64, err error) {
	limit := params.PageSize
	offset := params.PageSize * (params.Page - 1)
	db := global.KUBEGALE_DB.Model(&im.NotificationConfig{})

	// Apply filters
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.Type != "" {
		db = db.Where("type = ?", params.Type)
	}
	if params.NotificationPolicy != "" {
		// Assuming NotificationPolicy is an exact match, adjust if partial match is needed
		db = db.Where("notification_policy = ?", params.NotificationPolicy)
	}

	// Get total count
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count notification configs: %w", err)
	}

	// Apply sorting and pagination
	order := "created_at" // Default sort key
	if params.OrderKey != "" {
		order = params.OrderKey
	}
	if params.Desc {
		order += " desc"
	} else {
		order += " asc"
	}
	db = db.Order(order).Limit(limit).Offset(offset)

	var baseConfigs []im.NotificationConfig
	if err = db.Find(&baseConfigs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch notification configs: %w", err)
	}

	list = make([]response.NotificationResponse, 0, len(baseConfigs))
	for _, baseConfig := range baseConfigs {
		respItem := response.NotificationResponse{
			ID:                 baseConfig.ID,
			Name:               baseConfig.Name,
			Type:               baseConfig.Type,
			NotificationPolicy: baseConfig.NotificationPolicy,
			CreatedAt:          baseConfig.CreatedAt,
			// RobotURL will be populated below based on type
		}

		switch baseConfig.Type {
		case im.NotificationTypeFeiShu:
			var feiShuConfig im.FeiShuConfig
			err := global.KUBEGALE_DB.Where("notification_config_id = ?", baseConfig.ID).First(&feiShuConfig).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					global.KUBEGALE_LOG.Warn(fmt.Sprintf("FeiShuConfig not found for NotificationConfigID: %d", baseConfig.ID))
				} else {
					return nil, 0, fmt.Errorf("failed to fetch feishu config for ID %d: %w", baseConfig.ID, err)
				}
			} else {
				respItem.RobotURL = feiShuConfig.RobotURL
			}
		case im.NotificationTypeDingTalk:
			var dingTalkConfig im.DingTalkConfig
			if errDb := global.KUBEGALE_DB.Where("notification_config_id = ?", baseConfig.ID).First(&dingTalkConfig).Error; errDb == nil {
				respItem.RobotURL = dingTalkConfig.WebhookURL // Reusing RobotURL for WebhookURL
			} else if !errors.Is(errDb, gorm.ErrRecordNotFound) {
				global.KUBEGALE_LOG.Warn(fmt.Sprintf("DingTalkConfig not found for NotificationConfigID: %d, error: %v", baseConfig.ID, errDb))
			}
		default:
			global.KUBEGALE_LOG.Warn(fmt.Sprintf("Unhandled notification type '%s' in GetNotificationList for ID %d", baseConfig.Type, baseConfig.ID))
		}
		list = append(list, respItem)
	}

	// Sorting is already handled by the database query if OrderKey is standard
	// If custom sorting logic is needed post-fetch (e.g. combining fields from different types),
	// it would be applied here. The current sort is by NotificationConfig fields.

	return list, total, nil
}

// @function: GetNotificationById
// @description: 根据ID获取通知配置 (id is NotificationConfigID)
func (notificationService *NotificationService) GetNotificationById(id uint, notificationType string) (*response.NotificationDetailResponse, error) {
	// 1. Fetch NotificationConfig by id
	var fetchedNotificationConfig im.NotificationConfig
	if err := global.KUBEGALE_DB.First(&fetchedNotificationConfig, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification config not found")
		}
		return nil, fmt.Errorf("failed to fetch notification config: %w", err)
	}

	// Optional: Validate notificationType if provided, but use fetchedNotificationConfig.Type as source of truth.
	if notificationType != "" && fetchedNotificationConfig.Type != notificationType {
		// Log warning or return error if strict type matching is required
		// global.KUBEGALE_LOG.Warn(fmt.Sprintf("Requested type '%s' does not match stored type '%s' for ID %d", notificationType, fetchedNotificationConfig.Type, id))
		// For now, proceed with the actual type from DB.
	}

	// 2. Fetch associated CardContentConfig
	var cardContent im.CardContentConfig
	var cardContentDetail response.CardContentDetail // Prepare for response struct
	err := global.KUBEGALE_DB.Where("notification_id = ?", id).First(&cardContent).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to query card content: %w", err)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) { // If found, populate detail
		cardContentDetail = response.CardContentDetail{
			ID:                 cardContent.ID,
			NotificationID:     cardContent.NotificationID,
			AlertLevel:         cardContent.AlertLevel,
			AlertName:          cardContent.AlertName,
			NotificationPolicy: cardContent.NotificationPolicy, // Assuming this is the card's own policy if distinct
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
	}

	// 3. Initialize detailConfig and populate common fields
	detailConfig := response.NotificationDetailConfig{
		ID:                 fetchedNotificationConfig.ID,
		Name:               fetchedNotificationConfig.Name,
		Type:               fetchedNotificationConfig.Type,
		NotificationPolicy: fetchedNotificationConfig.NotificationPolicy,
		SendDailyStats:     fetchedNotificationConfig.SendDailyStats,
		CreatedAt:          fetchedNotificationConfig.CreatedAt,
		UpdatedAt:          fetchedNotificationConfig.UpdatedAt,
	}

	// 4. Based on fetchedNotificationConfig.Type, fetch specific config and populate specific fields
	switch fetchedNotificationConfig.Type {
	case im.NotificationTypeFeiShu:
		var feiShuConfig im.FeiShuConfig
		if err := global.KUBEGALE_DB.Where("notification_config_id = ?", id).First(&feiShuConfig).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("feishu config not found for notification id %d: %w", id, err)
			}
			return nil, fmt.Errorf("failed to fetch feishu config: %w", err)
		}
		detailConfig.RobotURL = feiShuConfig.RobotURL // Populate FeiShu specific field
	case im.NotificationTypeDingTalk:
		var dingTalkConfig im.DingTalkConfig
		if err := global.KUBEGALE_DB.Where("notification_config_id = ?", id).First(&dingTalkConfig).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("dingtalk config not found for notification id %d: %w", id, err)
			}
			return nil, fmt.Errorf("failed to fetch dingtalk config details: %w", err)
		}
		detailConfig.RobotURL = dingTalkConfig.WebhookURL // Reusing RobotURL for WebhookURL
		if dingTalkConfig.Secret != "" {                  // Only include secret if it's not empty
			sec := dingTalkConfig.Secret // Local variable for pointer
			detailConfig.Secret = &sec
		}
	default:
		return nil, fmt.Errorf("unsupported notification type: %s", fetchedNotificationConfig.Type)
	}

	// 5. Construct and return response
	return &response.NotificationDetailResponse{
		Config:      detailConfig,
		CardContent: cardContentDetail,
	}, nil
}

// @function: TestNotification
// @description: 测试通知发送 (req.ID is NotificationConfigID)
func (notificationService *NotificationService) TestNotification(req request.TestNotificationRequest) (response.TestNotificationResponse, error) {
	// 1. Fetch NotificationConfig by req.ID
	var fetchedNotificationConfig im.NotificationConfig
	if err := global.KUBEGALE_DB.First(&fetchedNotificationConfig, req.ID).Error; err != nil {
		msg := "获取通知配置失败: " + err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = fmt.Sprintf("通知配置ID %d 不存在", req.ID)
		}
		return response.TestNotificationResponse{Success: false, Message: msg}, err
	}

	// 2. Confirm req.Type matches fetchedNotificationConfig.Type
	if req.Type != fetchedNotificationConfig.Type {
		msg := fmt.Sprintf("请求的类型 '%s'与存储的类型 '%s' 不匹配", req.Type, fetchedNotificationConfig.Type)
		return response.TestNotificationResponse{Success: false, Message: msg}, errors.New(msg)
	}

	// 3. Fetch associated CardContentConfig
	var cardContent im.CardContentConfig
	var cardContentDetail response.CardContentDetail // For sender
	err := global.KUBEGALE_DB.Where("notification_id = ?", req.ID).First(&cardContent).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		msg := "获取卡片内容配置失败: " + err.Error()
		return response.TestNotificationResponse{Success: false, Message: msg}, err
	}
	// If cardContent is found, populate cardContentDetail (similar to GetNotificationById)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		cardContentDetail = response.CardContentDetail{
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
	}

	testMessage := "这是一条测试消息，请忽略。"
	if req.Message != "" {
		testMessage = req.Message
	}

	var sendErr error
	// 4. Based on fetchedNotificationConfig.Type, fetch specific config and send message
	switch fetchedNotificationConfig.Type {
	case im.NotificationTypeFeiShu:
		var feiShuConfig im.FeiShuConfig
		if err := global.KUBEGALE_DB.Where("notification_config_id = ?", req.ID).First(&feiShuConfig).Error; err != nil {
			msg := "获取飞书具体配置失败: " + err.Error()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				msg = fmt.Sprintf("飞书具体配置未找到，通知配置ID: %d", req.ID)
			}
			return response.TestNotificationResponse{Success: false, Message: msg}, err
		}
		// Prepare NotificationDetailConfig for the sender
		notificationDetailCfg := response.NotificationDetailConfig{
			ID:                 fetchedNotificationConfig.ID,
			Name:               fetchedNotificationConfig.Name,
			Type:               fetchedNotificationConfig.Type,
			NotificationPolicy: fetchedNotificationConfig.NotificationPolicy,
			SendDailyStats:     fetchedNotificationConfig.SendDailyStats,
			CreatedAt:          fetchedNotificationConfig.CreatedAt,
			UpdatedAt:          fetchedNotificationConfig.UpdatedAt,
			RobotURL:           feiShuConfig.RobotURL,
		}
		// Call sender
		sendErr = messageIm.MessageServiceApp.SendFeiShuMessage(notificationDetailCfg, cardContentDetail, testMessage)
	case im.NotificationTypeDingTalk:
		var dingTalkConfig im.DingTalkConfig
		if errDb := global.KUBEGALE_DB.Where("notification_config_id = ?", req.ID).First(&dingTalkConfig).Error; errDb != nil {
			msg := "获取钉钉配置失败: " + errDb.Error()
			if errors.Is(errDb, gorm.ErrRecordNotFound) {
				msg = fmt.Sprintf("钉钉具体配置未找到，通知配置ID: %d", req.ID)
			}
			return response.TestNotificationResponse{Success: false, Message: msg}, errDb
		}
		// Prepare NotificationDetailConfig for DingTalk (including WebhookURL and Secret)
		var detailConfigForSender response.NotificationDetailConfig
		detailConfigForSender.ID = fetchedNotificationConfig.ID
		detailConfigForSender.Name = fetchedNotificationConfig.Name
		detailConfigForSender.Type = fetchedNotificationConfig.Type
		detailConfigForSender.NotificationPolicy = fetchedNotificationConfig.NotificationPolicy
		detailConfigForSender.SendDailyStats = fetchedNotificationConfig.SendDailyStats
		detailConfigForSender.CreatedAt = fetchedNotificationConfig.CreatedAt
		detailConfigForSender.UpdatedAt = fetchedNotificationConfig.UpdatedAt
		detailConfigForSender.RobotURL = dingTalkConfig.WebhookURL
		if dingTalkConfig.Secret != "" {
			sec := dingTalkConfig.Secret
			detailConfigForSender.Secret = &sec
		}

		// Call the actual DingTalk message sending utility
		sendErr = messageIm.MessageServiceApp.SendDingTalkMessage(detailConfigForSender, cardContentDetail, testMessage)
	default:
		msg := fmt.Sprintf("不支持的通知类型 '%s' 进行测试发送", fetchedNotificationConfig.Type)
		return response.TestNotificationResponse{Success: false, Message: msg}, errors.New(msg)
	}

	if sendErr != nil {
		return response.TestNotificationResponse{Success: false, Message: "发送测试消息失败: " + sendErr.Error()}, sendErr
	}

	return response.TestNotificationResponse{Success: true, Message: "测试消息发送成功"}, nil
}
