package request

import (
	"KubeGale/model/common/request"
)

// K8sVeleroRecordsSearchReq Velero记录搜索请求结构体
type K8sVeleroRecordsSearchReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *K8sVeleroRecordsSearchReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeVeleroRecordReq 获取Velero记录详情请求结构体
type DescribeVeleroRecordReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`             // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`               // 命名空间
	VeleroRecordName string `json:"VeleroRecordName" form:"VeleroRecordName"` // Velero记录名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteVeleroRecordReq 删除Velero记录请求结构体
type DeleteVeleroRecordReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`             // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`               // 命名空间
	VeleroRecordName string `json:"VeleroRecordName" form:"VeleroRecordName"` // Velero记录名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateVeleroRecordReq 更新Velero记录请求结构体
type UpdateVeleroRecordReq struct {
	ClusterId        int         `json:"cluster_id" form:"cluster_id"`             // 集群ID
	Namespace        string      `json:"namespace" form:"namespace"`               // 命名空间
	VeleroRecordName string      `json:"VeleroRecordName" form:"VeleroRecordName"` // Velero记录名称
	Content          interface{} `json:"content" form:"content"`                   // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}

// CreateVeleroRecordReq 创建Velero记录请求结构体
type CreateVeleroRecordReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateVeleroRecordReq) GetClusterID() int {
	return r.ClusterId
}
