package Ingress

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sIngressRouter struct{}

func (s *K8sIngressRouter) Initk8sIngressRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sIngressRouterWithoutRecord := Router.Group("kubernetes")
	var k8sIngressApi = v1.ApiGroupApp.Ingress.K8sIngressApi
	{
		k8sNodeRouter.POST("ingress", k8sIngressApi.CreateIngress)
		k8sNodeRouter.DELETE("ingress", k8sIngressApi.DeleteIngress)
		k8sNodeRouter.PUT("ingress", k8sIngressApi.UpdateIngress)
	}
	{
		k8sIngressRouterWithoutRecord.GET("ingress", k8sIngressApi.GetIngressList)
		k8sIngressRouterWithoutRecord.GET("ingressDetails", k8sIngressApi.DescribeIngressInfo)
	}
}
