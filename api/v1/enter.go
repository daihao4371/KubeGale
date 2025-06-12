package v1

import (
	"KubeGale/api/v1/cmdb"
	"KubeGale/api/v1/im"
	"KubeGale/api/v1/kubernetes/cloudtty"
	"KubeGale/api/v1/kubernetes/clusterManager/cluster"
	"KubeGale/api/v1/kubernetes/configManager/configmap"
	"KubeGale/api/v1/kubernetes/configManager/horizontalpodautoscalers"
	"KubeGale/api/v1/kubernetes/configManager/limitRange"
	"KubeGale/api/v1/kubernetes/configManager/poddisruptionbudget"
	"KubeGale/api/v1/kubernetes/configManager/resourceQuota"
	"KubeGale/api/v1/kubernetes/configManager/secret"
	"KubeGale/api/v1/kubernetes/metrics"
	"KubeGale/api/v1/kubernetes/namespaceManager/namespaces"
	"KubeGale/api/v1/kubernetes/networks/endpoint"
	"KubeGale/api/v1/kubernetes/networks/ingress"
	"KubeGale/api/v1/kubernetes/networks/service"
	"KubeGale/api/v1/kubernetes/nodeManager/node"
	"KubeGale/api/v1/kubernetes/record"
	"KubeGale/api/v1/kubernetes/rolesMangager/clusterrole"
	"KubeGale/api/v1/kubernetes/rolesMangager/clusterrolebinding"
	"KubeGale/api/v1/kubernetes/rolesMangager/rolebindings"
	"KubeGale/api/v1/kubernetes/rolesMangager/roles"
	"KubeGale/api/v1/kubernetes/rolesMangager/serviceaccount"
	"KubeGale/api/v1/kubernetes/storageManager/pv"
	"KubeGale/api/v1/kubernetes/storageManager/pvc"
	"KubeGale/api/v1/kubernetes/storageManager/storageClass"
	"KubeGale/api/v1/kubernetes/workload/cronjob"
	"KubeGale/api/v1/kubernetes/workload/daemonSet"
	"KubeGale/api/v1/kubernetes/workload/deployment"
	"KubeGale/api/v1/kubernetes/workload/job"
	"KubeGale/api/v1/kubernetes/workload/pod"
	"KubeGale/api/v1/kubernetes/workload/replicaset"
	"KubeGale/api/v1/kubernetes/workload/statefulSet"
	"KubeGale/api/v1/system"
	"KubeGale/api/v1/ws"
	"KubeGale/service/kubernetes/configManager/poddistruptionbudget"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup                   system.ApiGroup
	ImApiGroup                       im.ApiGroup
	CmdbApiGroup                     cmdb.ApiGroup
	ClusterApiGroup                  cluster.ApiGroup
	NodeApiGroup                     node.ApiGroup
	PodApiGroup                      pod.ApiGroup
	DeploymentGroup                  deployment.ApiGroup
	ConfigMapGroup                   configmap.ApiGroup
	PvcGroup                         pvc.ApiGroup
	ServiceAccount                   serviceaccount.ApiGroup
	Secret                           secret.ApiGroup
	Ingress                          ingress.ApiGroup
	Service                          service.ApiGroup
	Replicaset                       replicaset.ApiGroup
	WsApi                            ws.ApiGroup
	DaemonSet                        daemonSet.ApiGroup
	StatefulSet                      statefulSet.ApiGroup
	Job                              job.ApiGroup
	CronJob                          cronjob.ApiGroup
	CloudTTY                         cloudtty.ApiGroup
	Namespace                        namespaces.ApiGroup
	Endpoint                         endpoint.ApiGroup
	ResourceQuota                    resourceQuota.ApiGroup
	LimitRange                       limitRange.ApiGroup
	HorizontalPod                    horizontalpodautoscalers.ApiGroup
	PoddistruptionbudgetServiceGroup poddistruptionbudget.ServiceGroup
	Poddisruptionbudget              poddisruptionbudget.ApiGroup
	Pv                               pv.ApiGroup
	StorageClass                     storageClass.ApiGroup
	ClusterRole                      clusterrole.ApiGroup
	ClusterRoleBinding               clusterrolebinding.ApiGroup
	Role                             roles.ApiGroup
	RoleBinding                      rolebindings.ApiGroup
	Metrics                          metrics.ApiGroup
	Record                           record.K8sRecordApi
}
