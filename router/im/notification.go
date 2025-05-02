package im

import (
	"KubeGale/middleware"

	"github.com/gin-gonic/gin"
)

type NotificationRouter struct{}

func (r *NotificationRouter) InitNotificationRouter(Router *gin.RouterGroup) {
	notificationRouter := Router.Group("notification").Use(middleware.OperationRecord())
	notificationRouterWithoutRecord := Router.Group("notification")

	// 需要记录操作的接口
	{
		notificationRouter.POST("createDingTalk", notificationApi.CreateDingTalk)           // 创建钉钉通知
		notificationRouter.POST("createFeiShu", notificationApi.CreateFeiShu)               // 创建飞书通知
		notificationRouter.PUT("updateDingTalk", notificationApi.UpdateDingTalk)            // 更新钉钉通知
		notificationRouter.PUT("updateFeiShu", notificationApi.UpdateFeiShu)                // 更新飞书通知
		notificationRouter.DELETE("deleteNotification", notificationApi.DeleteNotification) // 删除通知配置
		notificationRouter.POST("testNotification", notificationApi.TestNotification)       // 测试通知发送
		notificationRouter.POST("createCardContent", notificationApi.CreateCardContent)     // 创建卡片内容
		notificationRouter.PUT("updateCardContent", notificationApi.UpdateCardContent)      // 更新卡片内容
	}

	// 不需要记录操作的查询接口
	{
		notificationRouterWithoutRecord.POST("getNotificationList", notificationApi.GetNotificationList)      // 获取通知配置列表(POST方式)
		notificationRouterWithoutRecord.GET("getNotificationList", notificationApi.GetNotificationList)       // 获取通知配置列表(GET方式)
		notificationRouterWithoutRecord.GET("getNotificationById", notificationApi.GetNotificationById)       // 根据ID获取通知配置
		notificationRouterWithoutRecord.GET("getCardContent", notificationApi.GetCardContentByNotificationId) // 根据通知ID获取卡片内容
	}
}
