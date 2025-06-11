package resourceQuota

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// ResourceQuotaListResponse 资源配额列表响应结构体
// 用于返回资源配额的列表信息，包含分页信息
type ResourceQuotaListResponse struct {
	Items *[]v1.ResourceQuota `json:"items" form:"items"` // 资源配额列表
	Total int                 `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeResourceQuotaResponse 资源配额详情响应结构体
// 用于返回单个资源配额的详细信息
type DescribeResourceQuotaResponse struct {
	Items *v1.ResourceQuota `json:"items" form:"items"` // 资源配额详情
}
