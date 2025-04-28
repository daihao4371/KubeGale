package im

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/model/im"
	"KubeGale/model/im/request"
	"KubeGale/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotificationApi struct{}

// CreateDingTalk
// @Summary 创建钉钉通知配置
func (n *NotificationApi) CreateDingTalk(c *gin.Context) {
	var req request.CreateDingTalkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, err = notificationService.CreateDingTalk(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

// CreateFeiShu
// @Summary 创建飞书通知配置
func (n *NotificationApi) CreateFeiShu(c *gin.Context) {
	var req request.CreateFeiShuRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, err = notificationService.CreateFeiShu(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

// UpdateDingTalk
// @Summary 更新钉钉通知配置
func (n *NotificationApi) UpdateDingTalk(c *gin.Context) {
	var req request.UpdateDingTalkRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, err = notificationService.UpdateDingTalk(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// UpdateFeiShu
// @Summary 更新飞书通知配置
func (n *NotificationApi) UpdateFeiShu(c *gin.Context) {
	var req request.UpdateFeiShuRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, err = notificationService.UpdateFeiShu(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteNotification
// @Summary 删除通知配置
func (n *NotificationApi) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	notificationType := c.Query("type")
	if notificationType == "" {
		response.FailWithMessage("通知类型不能为空", c)
		return
	}

	err = notificationService.DeleteNotification(uint(id), notificationType)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetNotificationList
// @Summary 获取通知配置列表
func (n *NotificationApi) GetNotificationList(c *gin.Context) {
	var params request.SearchNotificationParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(params.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := notificationService.GetNotificationList(params)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
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
		response.FailWithMessage("参数错误", c)
		return
	}

	notificationType := c.Query("type")
	if notificationType == "" {
		response.FailWithMessage("通知类型不能为空", c)
		return
	}

	notification, err := notificationService.GetNotificationById(uint(id), notificationType)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}

	response.OkWithData(notification, c)
}

// TestNotification
// @Summary 测试通知发送
func (n *NotificationApi) TestNotification(c *gin.Context) {
	var req request.TestNotificationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := notificationService.TestNotification(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("测试失败!", zap.Error(err))
		response.FailWithMessage("测试失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// CreateCardContent
// @Summary 创建卡片内容配置
func (n *NotificationApi) CreateCardContent(c *gin.Context) {
	var req im.CardContentConfig
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = cardContentService.CreateCardContent(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(req, "创建成功", c)
}

// UpdateCardContent
// @Summary 更新卡片内容配置
func (n *NotificationApi) UpdateCardContent(c *gin.Context) {
	var req im.CardContentConfig
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = cardContentService.UpdateCardContent(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(req, "更新成功", c)
}

// GetCardContentByNotificationId
// @Summary 根据通知ID获取卡片内容
func (n *NotificationApi) GetCardContentByNotificationId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("notification_id"), 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	cardContent, err := cardContentService.GetCardContentById(uint(id))
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}

	response.OkWithData(cardContent, c)
}
