package serviceAccount

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// ServiceAccountListResponse 服务账号列表响应结构体
// 用于返回服务账号的列表信息，包含分页信息
type ServiceAccountListResponse struct {
	Items *[]v1.ServiceAccount `json:"items" form:"items"` // 服务账号列表
	Total int                  `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeServiceAccountResponse 服务账号详情响应结构体
// 用于返回单个服务账号的详细信息
type DescribeServiceAccountResponse struct {
	Items *v1.ServiceAccount `json:"items" form:"items"` // 服务账号详情
}
