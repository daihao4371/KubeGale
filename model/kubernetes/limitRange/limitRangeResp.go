package limitRange

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// LimitRangeListResponse LimitRange列表响应结构体
// 用于返回LimitRange的列表信息，包含分页信息
type LimitRangeListResponse struct {
	Items            *[]v1.LimitRange `json:"items" form:"items"` // LimitRange列表
	Total            int              `json:"total" form:"total"` // 总数
	request.PageInfo                  // 分页信息
}

// DescribeLimitRangeResponse LimitRange详情响应结构体
// 用于返回单个LimitRange的详细信息
type DescribeLimitRangeResponse struct {
	Items *v1.LimitRange `json:"items" form:"items"` // LimitRange详情
}
