package initialize

import (
	"KubeGale/global"
	cloudCmdb "KubeGale/model/cloudCmdb"
	cmdb "KubeGale/model/cmdb"
	cmdbReq "KubeGale/model/cmdb/request"
	example "KubeGale/model/example"
	"KubeGale/model/im"
	"KubeGale/model/kubernetes/cluster"
	"KubeGale/model/velero"
	"go.uber.org/zap"
)

// RegisterIMTables 注册IM相关数据库表
func bizModel() error {
	db := global.KUBEGALE_DB
	err := db.AutoMigrate(
		im.NotificationConfig{},
		im.FeiShuConfig{},
		im.CardContentConfig{},
		im.DingTalkConfig{},

		// 资产管理自建机房资源
		cmdb.CmdbHosts{},
		cmdb.CmdbProjects{},
		cmdbReq.CommandExecutionLog{},

		// 云资源
		cloudCmdb.LoadBalancer{},
		cloudCmdb.CloudPlatform{},
		cloudCmdb.RDS{},
		cloudCmdb.CloudRegions{},
		cloudCmdb.VirtualMachine{},

		// 文件上传
		example.ExaFileUploadAndDownload{},

		// 用于存储Velero备份任务的配置信息，包括备份范围、策略、保留时间等
		velero.K8sVeleroTasks{},

		// User Kubernetes集群用户信息结构体
		// 用于存储和管理Kubernetes集群中的用户信息，包括用户认证、权限等配置
		cluster.User{},
		cluster.K8sCluster{}, // Kubernetes集群信息结构体

	)
	if err != nil {
		global.KUBEGALE_LOG.Error("register  tables failed", zap.Error(err))
		return err
	}
	global.KUBEGALE_LOG.Info("register  tables success")
	return nil
}
