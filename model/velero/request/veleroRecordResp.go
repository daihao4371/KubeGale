package request

import (
	"KubeGale/model/common/request"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

// K8sVeleroRecordListResponse Velero记录列表响应结构体
// 用于返回Velero记录的列表信息，包含分页信息
type K8sVeleroRecordListResponse struct {
	Items *[]v1.Backup `json:"items" form:"items"` // Velero记录列表
	Total int          `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeK8sVeleroRecordResponse Velero记录详情响应结构体
// 用于返回单个Velero记录的详细信息
type DescribeK8sVeleroRecordResponse struct {
	Items *v1.Backup `json:"items" form:"items"` // Velero记录详情
}
