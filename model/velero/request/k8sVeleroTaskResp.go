package request

import (
	"KubeGale/model/common/request"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

// K8sVeleroTaskListResponse Velero任务列表响应结构体
// 用于返回Velero任务的列表信息，包含分页信息
type K8sVeleroTaskListResponse struct {
	Items *[]v1.Schedule `json:"items" form:"items"` // Velero任务列表
	Total int            `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeK8sVeleroTaskResponse Velero任务详情响应结构体
// 用于返回单个Velero任务的详细信息
type DescribeK8sVeleroTaskResponse struct {
	Items *v1.Schedule `json:"items" form:"items"` // Velero任务详情
}
