package job

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/batch/v1"
)

// JobListResponse Job列表响应结构体
// 用于返回Job的列表信息，包含分页信息
type JobListResponse struct {
	Items            *[]v1.Job `json:"items" form:"items"` // Job列表
	Total            int       `json:"total" form:"total"` // 总数
	request.PageInfo           // 分页信息
}

// DescribeJobResponse Job详情响应结构体
// 用于返回单个Job的详细信息
type DescribeJobResponse struct {
	Items *v1.Job `json:"items" form:"items"` // Job详情
}
