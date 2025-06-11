package record

import (
	"KubeGale/model/common/request"
)

// GetRecordListReq 获取记录列表请求结构体
type GetRecordListReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"` // 集群ID
	request.PageInfo
}

// DescribeRecordReq 获取记录详情请求结构体
type DescribeRecordReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"` // 集群ID
}

// UpdateRecordReq 更新记录请求结构体
type UpdateRecordReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"` // 集群ID
}

// CreateRecordReq 创建记录请求结构体
type CreateRecordReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"` // 集群ID
}

// DeleteRecordReq 删除记录请求结构体
type DeleteRecordReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"` // 集群ID
}
