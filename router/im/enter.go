package im

import api "KubeGale/api/v1"

type RouterGroup struct {
	NotificationRouter
}

var (
	notificationApi = api.ApiGroupApp.ImApiGroup.NotificationApi
)
