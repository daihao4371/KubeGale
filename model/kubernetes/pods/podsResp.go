package pods

import (
	"KubeGale/model/common/request"
	"KubeGale/utils/kubernetes/podtool"

	corev1 "k8s.io/api/core/v1"
)

// PodListResponse Pod列表响应结构体
// 用于返回Pod的列表信息，包含分页信息
type PodListResponse struct {
	Items *[]corev1.Pod `json:"items" form:"items"` // Pod列表
	Total int           `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribePodInfoResponse Pod详情响应结构体
// 用于返回单个Pod的详细信息
type DescribePodInfoResponse struct {
	Items *corev1.Pod `json:"items" form:"items"` // Pod详情
}

// EventInfoResponse Pod事件响应结构体
// 用于返回Pod相关的事件信息
type EventInfoResponse struct {
	Items *[]corev1.Event `json:"items" form:"items"` // 事件列表
	Total int             `json:"total" form:"total"` // 总数
}

// PodFilesResponse Pod文件响应结构体
// 用于返回Pod内文件的信息
type PodFilesResponse struct {
	Files []podtool.File `json:"files" form:"files"` // 文件列表
}
