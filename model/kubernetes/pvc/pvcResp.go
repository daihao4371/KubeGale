package pvc

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// PvcListResponse 持久卷声明列表响应结构体
// 用于返回持久卷声明的列表信息，包含分页信息
type PvcListResponse struct {
	Items *[]v1.PersistentVolumeClaim `json:"items" form:"items"` // 持久卷声明列表
	Total int                         `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribePVCResponse 持久卷声明详情响应结构体
// 用于返回单个持久卷声明的详细信息
type DescribePVCResponse struct {
	Items *v1.PersistentVolumeClaim `json:"items" form:"items"` // 持久卷声明详情
}
