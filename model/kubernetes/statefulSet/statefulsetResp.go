package statefulSet

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/apps/v1"
)

// StatefulSetListResponse 有状态应用列表响应结构体
// 用于返回有状态应用的列表信息，包含分页信息
type StatefulSetListResponse struct {
	Items *[]v1.StatefulSet `json:"items" form:"items"` // 有状态应用列表
	Total int               `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeStatefulSetResponse 有状态应用详情响应结构体
// 用于返回单个有状态应用的详细信息
type DescribeStatefulSetResponse struct {
	Items *v1.StatefulSet `json:"items" form:"items"` // 有状态应用详情
}
