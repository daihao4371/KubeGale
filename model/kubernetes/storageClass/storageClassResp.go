package storageClass

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/storage/v1"
)

// StorageClassListResponse 存储类列表响应结构体
// 用于返回存储类的列表信息，包含分页信息
type StorageClassListResponse struct {
	Items *[]v1.StorageClass `json:"items" form:"items"` // 存储类列表
	Total int                `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeStorageClassResponse 存储类详情响应结构体
// 用于返回单个存储类的详细信息
type DescribeStorageClassResponse struct {
	Items *v1.StorageClass `json:"items" form:"items"` // 存储类详情
}
