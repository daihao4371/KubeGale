package response

import (
	"KubeGale/model/kubernetes"
	cluster2 "KubeGale/model/kubernetes/cluster"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
)

// ClusterResponse 表示单个Kubernetes集群的响应结构
type ClusterResponse struct {
	Cluster cluster2.K8sCluster `json:"cluster"` // 集群详细信息
}

// ClusterUserResponse 表示集群用户列表的响应结构
type ClusterUserResponse struct {
	User []cluster2.User `json:"user"` // 用户列表
}

// RolesResponse 表示集群角色列表的响应结构
type RolesResponse struct {
	Roles []v1.ClusterRole `json:"roles"` // 集群角色列表
}

// ApiGroupResponse 表示API分组选项的响应结构
type ApiGroupResponse struct {
	Groups []kubernetes.ApiGroupOption `json:"groups"` // API分组选项列表
}

// ClusterUserNamespace 表示集群用户可访问的命名空间列表响应结构
type ClusterUserNamespace struct {
	Namespaces []string `json:"namespaces"` // 命名空间名称列表
}

// ClusterListNamespace 表示集群所有命名空间详细信息的响应结构
type ClusterListNamespace struct {
	Namespaces []corev1.Namespace `json:"namespaces"` // 命名空间详细信息列表
}
