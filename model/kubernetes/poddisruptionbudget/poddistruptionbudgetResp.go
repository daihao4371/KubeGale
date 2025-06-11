package poddisruptionbudget

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/policy/v1"
)

// PoddisruptionbudgetListResponse Pod中断预算列表响应结构体
// 用于返回Pod中断预算的列表信息，包含分页信息
type PoddisruptionbudgetListResponse struct {
	Items *[]v1.PodDisruptionBudget `json:"items" form:"items"` // Pod中断预算列表
	Total int                       `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribePoddisruptionbudgetResponse Pod中断预算详情响应结构体
// 用于返回单个Pod中断预算的详细信息
type DescribePoddisruptionbudgetResponse struct {
	Items *v1.PodDisruptionBudget `json:"items" form:"items"` // Pod中断预算详情
}
