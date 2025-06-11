package clusterolebinding

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/rbac/v1"
)

// ClusterRoleBindingListResponse 集群角色绑定列表响应结构体
// 用于返回集群角色绑定的列表信息，包含分页信息
type ClusterRoleBindingListResponse struct {
	Items            *[]v1.ClusterRoleBinding `json:"items" form:"items"` // 集群角色绑定列表
	Total            int                      `json:"total" form:"total"` // 总数
	request.PageInfo                          // 分页信息
}

// DescribeClusterRoleBindingResponse 集群角色绑定详情响应结构体
// 用于返回单个集群角色绑定的详细信息
type DescribeClusterRoleBindingResponse struct {
	Items *v1.ClusterRoleBinding `json:"items" form:"items"` // 集群角色绑定详情
}
