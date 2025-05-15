package im

import "KubeGale/service"

type ApiGroup struct {
	NotificationApi
}

var (
	notificationService = service.ServiceGroupApp.ImServiceGroup.NotificationService
	cardContentService  = service.ServiceGroupApp.ImServiceGroup.CardContentService
)