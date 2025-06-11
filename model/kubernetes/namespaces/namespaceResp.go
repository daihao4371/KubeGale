package namespaces

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// NamespaceListResponse 命名空间列表响应结构体
// 用于返回命名空间的列表信息，包含分页信息
type NamespaceListResponse struct {
	Items *[]v1.Namespace `json:"items" form:"items"` // 命名空间列表
	Total int             `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeNamespaceResponse 命名空间详情响应结构体
// 用于返回单个命名空间的详细信息
type DescribeNamespaceResponse struct {
	Items *v1.Namespace `json:"items" form:"items"` // 命名空间详情
}
