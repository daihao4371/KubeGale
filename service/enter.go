package service

import (
	"KubeGale/service/cloudCmdb" // Added this import
	"KubeGale/service/cmdb"
	"KubeGale/service/im"
	"KubeGale/service/kubernetes/cloudtty"
	"KubeGale/service/kubernetes/configManager/configmap"
	"KubeGale/service/kubernetes/configManager/horizontalPod"
	"KubeGale/service/kubernetes/configManager/limitRange"
	"KubeGale/service/kubernetes/configManager/poddistruptionbudget"
	"KubeGale/service/kubernetes/configManager/resourceQuota"
	"KubeGale/service/kubernetes/configManager/secret"
	"KubeGale/service/kubernetes/metrics"
	"KubeGale/service/kubernetes/namespace"
	"KubeGale/service/kubernetes/network/Ingress"
	"KubeGale/service/kubernetes/network/endpoint"
	"KubeGale/service/kubernetes/network/service"
	"KubeGale/service/kubernetes/node"
	"KubeGale/service/kubernetes/rolesManager/clusterRole"
	"KubeGale/service/kubernetes/rolesManager/clusterRoleBinding"
	"KubeGale/service/kubernetes/rolesManager/role"
	"KubeGale/service/kubernetes/rolesManager/roleBinding"
	"KubeGale/service/kubernetes/serviceAccount"
	"KubeGale/service/kubernetes/storageManager/pv"
	"KubeGale/service/kubernetes/storageManager/pvc"
	"KubeGale/service/kubernetes/storageManager/storageClass"
	"KubeGale/service/kubernetes/velero"
	"KubeGale/service/kubernetes/workload/cronjob"
	"KubeGale/service/kubernetes/workload/daemonset"
	"KubeGale/service/kubernetes/workload/deployment"
	"KubeGale/service/kubernetes/workload/job"
	"KubeGale/service/kubernetes/workload/pod"
	"KubeGale/service/kubernetes/workload/replicaSet"
	"KubeGale/service/kubernetes/workload/statefulSet"
	"KubeGale/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup               system.ServiceGroup
	ImServiceGroup                   im.ServiceGroup
	CmdbServiceGroup                 cmdb.ServiceGroup
	CloudCmdbServiceGroup            cloudCmdb.ServiceGroup
	NodeServiceGroup                 node.ServiceGroup
	PodServiceGroup                  pod.ServiceGroup
	DeploymentServiceGroup           deployment.ServiceGroup
	PvcServiceGroup                  pvc.ServiceGroup
	SecretServiceGroup               secret.ServiceGroup
	ConfigMapServiceGroup            configmap.ServiceGroup
	ServiceAccountServiceGroup       serviceAccount.ServiceGroup
	IngressServiceGroup              Ingress.ServiceGroup
	SvcServiceGroup                  service.ServiceGroup
	ReplicaSetServiceGroup           replicaSet.ServiceGroup
	DaemonSetServiceGroup            daemonset.ServiceGroup
	StatefulSetServiceGroup          statefulSet.ServiceGroup
	JobServiceGroup                  job.ServiceGroup
	CronJobServiceGroup              cronjob.ServiceGroup
	CloudTTYServiceGroup             cloudtty.ServiceGroup
	NamespaceServiceGroup            namespace.ServiceGroup
	EndPointServiceGroup             endpoint.ServiceGroup
	ResourceQuotaServiceGroup        resourceQuota.ServiceGroup
	LimitRangeServiceGroup           limitRange.ServiceGroup
	HorizontalPodServiceGroup        horizontalPod.ServiceGroup
	PoddistruptionbudgetServiceGroup poddistruptionbudget.ServiceGroup
	PvServiceGroup                   pv.ServiceGroup
	StorageClassServiceGroup         storageClass.ServiceGroup
	ClusterRoleServiceGroup          clusterRole.ServiceGroup
	ClusterRoleBindingServiceGroup   clusterRoleBinding.ServiceGroup
	RoleServiceGroup                 role.ServiceGroup
	RoleBindingServiceGroup          roleBinding.ServiceGroup
	VeleroServiceGroup               velero.ServiceGroup
	MetricsServiceGroup              metrics.ServiceGroup // Added this line
}

var ServiceGroupApp = new(ServiceGroup)
