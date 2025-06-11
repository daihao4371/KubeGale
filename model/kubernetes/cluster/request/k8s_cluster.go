package request

import (
	"KubeGale/model/common/request"
	"KubeGale/model/kubernetes/cluster"
	"time"
)

// K8sClusterSearch 集群搜索结构体
// 用于定义集群搜索的查询参数，包括时间范围、别名、地理位置等信息
type K8sClusterSearch struct {
	StartCreatedAt   *time.Time `json:"startCreatedAt" form:"startCreatedAt"` // 创建时间范围开始
	EndCreatedAt     *time.Time `json:"endCreatedAt" form:"endCreatedAt"`     // 创建时间范围结束
	Alias            string     `json:"alias" form:"alias"`                   // 集群别名
	City             string     `json:"city" form:"city"`                     // 所在城市
	District         string     `json:"district" form:"district"`             // 所在区域
	Name             string     `json:"name" form:"name"`                     // 集群名称
	request.PageInfo            // 分页信息
}

// ClusterRoleType 集群角色类型结构体
// 用于定义集群角色的类型和关联的集群ID
type ClusterRoleType struct {
	RoleType  string `json:"role_type"`  // 角色类型
	ClusterId uint   `json:"cluster_id"` // 集群ID
}

// ClusterApiGroups 集群API组结构体
// 用于定义集群API的类型和关联的集群ID
type ClusterApiGroups struct {
	ApiType   string `json:"api_type"`   // API类型
	ClusterId uint   `json:"cluster_id"` // 集群ID
}

// NamespaceRoles 命名空间角色结构体
// 用于定义命名空间及其对应的角色列表
type NamespaceRoles struct {
	Namespace string   `json:"namespace"` // 命名空间名称
	Roles     []string `json:"roles"`     // 角色列表
}

// CreateClusterRole 创建集群角色结构体
// 用于定义创建集群角色时的请求参数，包括集群角色、命名空间角色和用户UUID等信息
type CreateClusterRole struct {
	ClusterRoles   []string         `json:"cluster_roles"`   // 集群角色列表
	NamespaceRoles []NamespaceRoles `json:"namespace_roles"` // 命名空间角色列表
	UserUuids      []string         `json:"user_uuids"`      // 用户UUID列表
	cluster.User                    // 用户信息
	ClusterId      int              `json:"cluster_id"` // 集群ID
}

// DeleteClusterRole 删除集群角色结构体
// 用于定义删除集群角色时的请求参数，包括用户UUID和集群ID
type DeleteClusterRole struct {
	UserUuids []string `json:"user_uuids"` // 用户UUID列表
	ClusterId int      `json:"cluster_id"` // 集群ID
}
