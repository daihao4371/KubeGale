package roleBinding

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/rbac/v1"
)

// RoleBindingListResponse 角色绑定列表响应结构体
// 用于返回角色绑定的列表信息，包含分页信息
type RoleBindingListResponse struct {
	Items *[]v1.RoleBinding `json:"items" form:"items"` // 角色绑定列表
	Total int               `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeRoleBindingResponse 角色绑定详情响应结构体
// 用于返回单个角色绑定的详细信息
type DescribeRoleBindingResponse struct {
	Items *v1.RoleBinding `json:"items" form:"items"` // 角色绑定详情
}
