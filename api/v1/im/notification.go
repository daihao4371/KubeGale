package im

import (
	"KubeGale/global"
	commonResponse "KubeGale/model/common/response"
	"KubeGale/model/im"
	"KubeGale/model/im/request"
	imResponse "KubeGale/model/im/response"
	"KubeGale/utils"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationApi struct{}

// CreateFeiShu
// @Summary 创建飞书通知配置
func (n *NotificationApi) CreateFeiShu(c *gin.Context) {
	var req request.CreateFeiShuRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		commonResponse.FailWithMessage(err.Error(), c)
		return
	}

	_, err = notificationService.CreateFeiShu(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		commonResponse.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithMessage("创建成功", c)
}

// UpdateFeiShu
// @Summary 更新飞书通知配置
func (n *NotificationApi) UpdateFeiShu(c *gin.Context) {
	var req request.UpdateFeiShuRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		commonResponse.FailWithMessage(err.Error(), c)
		return
	}

	_, err = notificationService.UpdateFeiShu(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		commonResponse.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithMessage("更新成功", c)
}

// DeleteNotification
// @Summary 删除通知配置
func (n *NotificationApi) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 32)
	if err != nil {
		commonResponse.FailWithMessage("参数错误", c)
		return
	}

	notificationType := c.Query("type")
	if notificationType == "" {
		commonResponse.FailWithMessage("通知类型不能为空", c)
		return
	}

	err = notificationService.DeleteNotification(uint(id), notificationType)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		commonResponse.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithMessage("删除成功", c)
}

// GetNotificationList
// @Summary 获取通知配置列表
func (n *NotificationApi) GetNotificationList(c *gin.Context) {
	var params request.SearchNotificationParams
	// 设置默认分页参数
	params.Page = 1
	params.PageSize = 10

	// 支持GET和POST两种方式获取参数
	if c.Request.Method == "POST" {
		if err := c.ShouldBindJSON(&params); err != nil {
			commonResponse.FailWithMessage(err.Error(), c)
			return
		}
	} else {
		// GET方式获取参数
		if err := c.ShouldBindQuery(&params); err != nil {
			commonResponse.FailWithMessage(err.Error(), c)
			return
		}
	}

	// 确保分页参数有效
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	err := utils.Verify(params.PageInfo, utils.PageInfoVerify)
	if err != nil {
		commonResponse.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := notificationService.GetNotificationList(params)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		commonResponse.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}

	// 返回统一的响应格式
	commonResponse.OkWithDetailed(commonResponse.PageResult{
		List:     list,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, "获取成功", c)
}

// GetNotificationById
// @Summary 根据ID获取通知配置
func (n *NotificationApi) GetNotificationById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 32)
	if err != nil {
		commonResponse.FailWithMessage("参数错误", c)
		return
	}

	notificationType := c.Query("type")
	if notificationType == "" {
		commonResponse.FailWithMessage("通知类型不能为空", c)
		return
	}

	notification, err := notificationService.GetNotificationById(uint(id), notificationType)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		commonResponse.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithData(notification, c)
}

// TestNotification
// @Summary 测试通知发送
func (n *NotificationApi) TestNotification(c *gin.Context) {
	var req request.TestNotificationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		commonResponse.FailWithMessage(err.Error(), c)
		return
	}

	result, err := notificationService.TestNotification(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("测试失败!", zap.Error(err))
		commonResponse.FailWithMessage("测试失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithData(result, c)
}

// CreateDingTalk
// @Tags Notification
// @Summary 创建钉钉通知配置
// @Produce application/json
// @Param data body request.CreateDingTalkRequest true "钉钉通知配置信息"
// @Success 200 {object} commonResponse.Response{data=im.DingTalkConfig,msg=string} "创建成功"
// @Router /notification/createDingTalk [post]
func (n *NotificationApi) CreateDingTalk(c *gin.Context) {
	var req request.CreateDingTalkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		commonResponse.FailWithMessage(err.Error(), c)
		return
	}
	
	createdDingTalkConfig, err := notificationService.CreateDingTalk(req) // Assuming notificationService is available as in other handlers
	if err != nil {
		global.KUBEGALE_LOG.Error("创建钉钉配置失败!", zap.Error(err))
		commonResponse.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	commonResponse.OkWithData(createdDingTalkConfig, c) // Return created data
}

// UpdateDingTalk
// @Tags Notification
// @Summary 更新钉钉通知配置
// @Produce application/json
// @Param data body request.UpdateDingTalkRequest true "钉钉通知配置更新信息"
// @Success 200 {object} commonResponse.Response{data=im.DingTalkConfig,msg=string} "更新成功"
// @Router /notification/updateDingTalk [put]
func (n *NotificationApi) UpdateDingTalk(c *gin.Context) {
	var req request.UpdateDingTalkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		commonResponse.FailWithMessage(err.Error(), c)
		return
	}

	updatedDingTalkConfig, err := notificationService.UpdateDingTalk(req) // Assuming notificationService is available
	if err != nil {
		global.KUBEGALE_LOG.Error("更新钉钉配置失败!", zap.Error(err))
		commonResponse.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	commonResponse.OkWithData(updatedDingTalkConfig, c) // Return updated data
}

// CreateCardContent
// @Summary 创建卡片内容配置
func (n *NotificationApi) CreateCardContent(c *gin.Context) {
	var req request.CreateCardContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		commonResponse.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 将通知用户数组转换为JSON字符串
	var notifiedUsersStr string
	if len(req.NotifiedUsers) > 0 {
		notifiedUsersBytes, err := json.Marshal(req.NotifiedUsers)
		if err != nil {
			global.KUBEGALE_LOG.Error("序列化通知用户失败", zap.Error(err))
			commonResponse.FailWithMessage("创建失败: 序列化通知用户失败", c)
			return
		}
		notifiedUsersStr = string(notifiedUsersBytes)
	}

	// 将请求转换为卡片内容配置
	cardContent := im.CardContentConfig{
		NotificationID:     req.NotificationID,
		AlertLevel:         req.AlertLevel,
		AlertName:          req.AlertName,
		NotificationPolicy: req.NotificationPolicy,
		AlertContent:       req.AlertContent,
		AlertTime:          req.AlertTime,
		NotifiedUsers:      notifiedUsersStr,
		LastSimilarAlert:   req.LastSimilarAlert,
		AlertHandler:       req.AlertHandler,
		ClaimAlert:         req.ClaimAlert,
		ResolveAlert:       req.ResolveAlert,
		MuteAlert:          req.MuteAlert,
		UnresolvedAlert:    req.UnresolvedAlert,
	}

	err := cardContentService.CreateCardContent(cardContent)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建卡片内容失败", zap.Error(err))
		commonResponse.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithMessage("创建成功", c)
}

// UpdateCardContent
// @Summary 更新卡片内容配置
func (n *NotificationApi) UpdateCardContent(c *gin.Context) {
	var req im.CardContentConfig
	err := c.ShouldBindJSON(&req)
	if err != nil {
		commonResponse.FailWithMessage(err.Error(), c)
		return
	}

	err = cardContentService.UpdateCardContent(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		commonResponse.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithDetailed(req, "更新成功", c)
}

// GetCardContentByNotificationId
// @Summary 根据通知ID获取卡片内容
func (n *NotificationApi) GetCardContentByNotificationId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("notification_id"), 10, 32)
	if err != nil {
		commonResponse.FailWithMessage("参数错误", c)
		return
	}

	// 先查询通知配置的类型
	var notificationType string
	var notificationConfig im.NotificationConfig
	err = global.KUBEGALE_DB.Where("id = ?", id).First(&notificationConfig).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果通知配置不存在，返回空的卡片内容
			commonResponse.OkWithData(imResponse.NotificationDetailResponse{
				Config:      imResponse.NotificationDetailConfig{},
				CardContent: imResponse.CardContentDetail{},
			}, c)
			return
		}
		global.KUBEGALE_LOG.Error("获取通知配置失败!", zap.Error(err))
		commonResponse.FailWithMessage("获取通知配置失败: "+err.Error(), c)
		return
	}
	notificationType = notificationConfig.Type

	// 获取通知配置和卡片内容
	result, err := notificationService.GetNotificationById(uint(id), notificationType)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取通知配置失败!", zap.Error(err))
		commonResponse.FailWithMessage("获取通知配置失败: "+err.Error(), c)
		return
	}

	commonResponse.OkWithData(result, c)
}
