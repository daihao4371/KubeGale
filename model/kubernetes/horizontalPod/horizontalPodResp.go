package horizontalPod

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/autoscaling/v1"
)

// HorizontalPodResponse 水平Pod自动伸缩器列表响应结构体
// 用于返回水平Pod自动伸缩器的列表信息，包含分页信息
type HorizontalPodResponse struct {
	Items            *[]v1.HorizontalPodAutoscaler `json:"items" form:"items"` // 水平Pod自动伸缩器列表
	Total            int                           `json:"total" form:"total"` // 总数
	request.PageInfo                               // 分页信息
}

// DescribeHorizontalPodResponse 水平Pod自动伸缩器详情响应结构体
// 用于返回单个水平Pod自动伸缩器的详细信息
type DescribeHorizontalPodResponse struct {
	Items *v1.HorizontalPodAutoscaler `json:"items" form:"items"` // 水平Pod自动伸缩器详情
}
