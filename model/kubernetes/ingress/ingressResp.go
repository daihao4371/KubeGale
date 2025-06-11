package ingress

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/networking/v1"
)

// IngressListResponse Ingress列表响应结构体
// 用于返回Ingress的列表信息，包含分页信息
type IngressListResponse struct {
	Items            *[]v1.Ingress `json:"items" form:"items"` // Ingress列表
	Total            int           `json:"total" form:"total"` // 总数
	request.PageInfo               // 分页信息
}

// DescribeIngressResponse Ingress详情响应结构体
// 用于返回单个Ingress的详细信息
type DescribeIngressResponse struct {
	Items *v1.Ingress `json:"items" form:"items"` // Ingress详情
}
