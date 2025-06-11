package request

import (
	"KubeGale/model/common/request"
)

// K8sVeleroRestoresSearchReq Velero恢复搜索请求结构体
type K8sVeleroRestoresSearchReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *K8sVeleroRestoresSearchReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeVeleroRestoreReq 获取Velero恢复详情请求结构体
type DescribeVeleroRestoreReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string `json:"namespace" form:"namespace"`                 // 命名空间
	VeleroRestoreName string `json:"VeleroRestoreName" form:"VeleroRestoreName"` // Velero恢复名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteVeleroRestoreReq 删除Velero恢复请求结构体
type DeleteVeleroRestoreReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string `json:"namespace" form:"namespace"`                 // 命名空间
	VeleroRestoreName string `json:"VeleroRestoreName" form:"VeleroRestoreName"` // Velero恢复名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateVeleroRestoreReq 更新Velero恢复请求结构体
type UpdateVeleroRestoreReq struct {
	ClusterId         int         `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string      `json:"namespace" form:"namespace"`                 // 命名空间
	VeleroRestoreName string      `json:"VeleroRestoreName" form:"VeleroRestoreName"` // Velero恢复名称
	Content           interface{} `json:"content" form:"content"`                     // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}

// CreateVeleroRestoreReq 创建Velero恢复请求结构体
type CreateVeleroRestoreReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateVeleroRestoreReq) GetClusterID() int {
	return r.ClusterId
}
