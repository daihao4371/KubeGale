package router

import (
	"KubeGale/router/cmdb"
	"KubeGale/router/im"
	"KubeGale/router/kubernetes/cloudtty"
	"KubeGale/router/kubernetes/clusterManager/cluster"
	"KubeGale/router/kubernetes/configManager/configmap"
	"KubeGale/router/kubernetes/configManager/horizontalPod"
	"KubeGale/router/kubernetes/configManager/limitRange"
	"KubeGale/router/kubernetes/configManager/poddistruptionbudget"
	"KubeGale/router/kubernetes/configManager/resourceQuota"
	"KubeGale/router/kubernetes/configManager/secret"
	"KubeGale/router/kubernetes/example"
	"KubeGale/router/kubernetes/metrics"
	"KubeGale/router/kubernetes/namespaceManager/namespaces"
	"KubeGale/router/kubernetes/network/Ingress"
	"KubeGale/router/kubernetes/network/endpoint"
	"KubeGale/router/kubernetes/network/service"
	"KubeGale/router/kubernetes/nodeManager/node"
	"KubeGale/router/kubernetes/rolesManager/clusterRole"
	"KubeGale/router/kubernetes/rolesManager/clusterRoleBinding"
	"KubeGale/router/kubernetes/rolesManager/rolebinding"
	"KubeGale/router/kubernetes/rolesManager/roles"
	"KubeGale/router/kubernetes/rolesManager/serviceAccount"
	"KubeGale/router/kubernetes/storageManager/pv"
	"KubeGale/router/kubernetes/storageManager/pvc"
	"KubeGale/router/kubernetes/storageManager/storageClass"
	"KubeGale/router/kubernetes/workload/cronjob"
	"KubeGale/router/kubernetes/workload/daemonset"
	"KubeGale/router/kubernetes/workload/deployment"
	"KubeGale/router/kubernetes/workload/job"
	"KubeGale/router/kubernetes/workload/pod"
	"KubeGale/router/kubernetes/workload/replicaSet"
	"KubeGale/router/kubernetes/workload/statefulSet"
	"KubeGale/router/system"
	"KubeGale/router/velero"
	"KubeGale/router/ws"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System system.RouterGroup
	Im     im.RouterGroup
	Cmdb   cmdb.RouterGroup

	Example              example.RouterGroup
	Cluster              cluster.RouterGroup
	Node                 node.RouterGroup
	Pod                  pod.RouterGroup
	Deployment           deployment.RouterGroup
	Secret               secret.RouterGroup
	Configmap            configmap.RouterGroup
	ServiceAccount       serviceAccount.RouterGroup
	Pvc                  pvc.RouterGroup
	Service              service.RouterGroup
	Ingress              Ingress.RouterGroup
	ReplicaSet           replicaSet.RouterGroup
	Ws                   ws.RouterGroup
	DaemonSet            daemonset.RouterGroup
	StatefulSet          statefulSet.RouterGroup
	Job                  job.RouterGroup
	CronJob              cronjob.RouterGroup
	CloudTTY             cloudtty.RouterGroup
	Namespace            namespaces.RouterGroup
	EndPoint             endpoint.RouterGroup
	ResourceQuota        resourceQuota.RouterGroup
	LimitRange           limitRange.RouterGroup
	HorizontalPod        horizontalPod.RouterGroup
	Poddistruptionbudget poddistruptionbudget.RouterGroup
	Pv                   pv.RouterGroup
	StorageClass         storageClass.RouterGroup
	ClusterRole          clusterRole.RouterGroup
	ClusterRoleBinding   clusterRoleBinding.RouterGroup
	Role                 roles.RouterGroup
	RoleBinding          rolebinding.RouterGroup
	Velero               velero.K8sVeleroTasksRouter
	Metrics              metrics.RouterGroup
}
