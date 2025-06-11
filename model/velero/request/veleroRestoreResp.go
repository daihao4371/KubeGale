package request

import (
	"KubeGale/model/common/request"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

// K8sVeleroRestoreListResponse Velero恢复列表响应结构体
// 用于返回Velero恢复的列表信息，包含分页信息
type K8sVeleroRestoreListResponse struct {
	Items *[]v1.Restore `json:"items" form:"items"` // Velero恢复列表
	Total int           `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeK8sVeleroRestoreResponse Velero恢复详情响应结构体
// 用于返回单个Velero恢复的详细信息
type DescribeK8sVeleroRestoreResponse struct {
	Items *v1.Restore `json:"items" form:"items"` // Velero恢复详情
}
