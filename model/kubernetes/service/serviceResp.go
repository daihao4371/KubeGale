package service

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// ServiceListResponse 服务列表响应结构体
// 用于返回服务的列表信息，包含分页信息
type ServiceListResponse struct {
	Items *[]v1.Service `json:"items" form:"items"` // 服务列表
	Total int           `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeServiceResponse 服务详情响应结构体
// 用于返回单个服务的详细信息
type DescribeServiceResponse struct {
	Items *v1.Service `json:"items" form:"items"` // 服务详情
}
