package roles

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/rbac/v1"
)

// RoleListResponse 角色列表响应结构体
// 用于返回角色的列表信息，包含分页信息
type RoleListResponse struct {
	Items *[]v1.Role `json:"items" form:"items"` // 角色列表
	Total int        `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeRoleResponse 角色详情响应结构体
// 用于返回单个角色的详细信息
type DescribeRoleResponse struct {
	Items *v1.Role `json:"items" form:"items"` // 角色详情
}
