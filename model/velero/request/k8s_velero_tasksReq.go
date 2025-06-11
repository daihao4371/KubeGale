package request

import (
	"KubeGale/model/common/request"
)

// K8sVeleroTasksSearchReq Velero任务搜索请求结构体
type K8sVeleroTasksSearchReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *K8sVeleroTasksSearchReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeVeleroTaskReq 获取Velero任务详情请求结构体
type DescribeVeleroTaskReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`         // 集群ID
	Namespace      string `json:"namespace" form:"namespace"`           // 命名空间
	VeleroTaskName string `json:"VeleroTaskName" form:"VeleroTaskName"` // Velero任务名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteVeleroTaskReq 删除Velero任务请求结构体
type DeleteVeleroTaskReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`         // 集群ID
	Namespace      string `json:"namespace" form:"namespace"`           // 命名空间
	VeleroTaskName string `json:"VeleroTaskName" form:"VeleroTaskName"` // Velero任务名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateVeleroTaskReq 更新Velero任务请求结构体
type UpdateVeleroTaskReq struct {
	ClusterId      int         `json:"cluster_id" form:"cluster_id"`         // 集群ID
	Namespace      string      `json:"namespace" form:"namespace"`           // 命名空间
	VeleroTaskName string      `json:"VeleroTaskName" form:"VeleroTaskName"` // Velero任务名称
	Content        interface{} `json:"content" form:"content"`               // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}

// CreateVeleroTaskReq 创建Velero任务请求结构体
type CreateVeleroTaskReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateVeleroTaskReq) GetClusterID() int {
	return r.ClusterId
}
