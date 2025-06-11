package configmap

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// ConfigMapListResponse ConfigMap列表响应结构体
// 用于返回ConfigMap的列表信息，包含分页信息
type ConfigMapListResponse struct {
	Items            *[]v1.ConfigMap `json:"items" form:"items"` // ConfigMap列表
	Total            int             `json:"total" form:"total"` // 总数
	request.PageInfo                 // 分页信息
}

// DescribeConfigMapResponse ConfigMap详情响应结构体
// 用于返回单个ConfigMap的详细信息
type DescribeConfigMapResponse struct {
	Items *v1.ConfigMap `json:"items"` // ConfigMap详情
}
