package clusterrole

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/rbac/v1"
)

// ClusterRoleListResponse 集群角色列表响应结构体
// 用于返回集群角色的列表信息，包含分页信息
type ClusterRoleListResponse struct {
	Items            *[]v1.ClusterRole `json:"items" form:"items"` // 集群角色列表
	Total            int               `json:"total" form:"total"` // 总数
	request.PageInfo                   // 分页信息
}

// DescribeClusterRoleResponse 集群角色详情响应结构体
// 用于返回单个集群角色的详细信息
type DescribeClusterRoleResponse struct {
	Items *v1.ClusterRole `json:"items" form:"items"` // 集群角色详情
}
