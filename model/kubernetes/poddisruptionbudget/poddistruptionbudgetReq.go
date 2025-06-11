package poddisruptionbudget

import "KubeGale/model/common/request"

// GetPoddisruptionbudgetListReq 获取Pod中断预算列表请求结构体
type GetPoddisruptionbudgetListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetPoddisruptionbudgetListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribePoddisruptionbudgetReq 获取Pod中断预算详情请求结构体
type DescribePoddisruptionbudgetReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`                           // 集群ID
	Namespace               string `json:"namespace" form:"namespace"`                             // 命名空间
	PoddisruptionbudgetName string `json:"poddisruptionbudgetName" form:"poddisruptionbudgetName"` // Pod中断预算名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}

// DeletePoddisruptionbudgetReq 删除Pod中断预算请求结构体
type DeletePoddisruptionbudgetReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`                           // 集群ID
	Namespace               string `json:"namespace" form:"namespace"`                             // 命名空间
	PoddisruptionbudgetName string `json:"poddisruptionbudgetName" form:"poddisruptionbudgetName"` // Pod中断预算名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeletePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}

// UpdatePoddisruptionbudgetReq 更新Pod中断预算请求结构体
type UpdatePoddisruptionbudgetReq struct {
	ClusterId               int         `json:"cluster_id" form:"cluster_id"`                           // 集群ID
	Namespace               string      `json:"namespace" form:"namespace"`                             // 命名空间
	PoddisruptionbudgetName string      `json:"poddisruptionbudgetName" form:"poddisruptionbudgetName"` // Pod中断预算名称
	Content                 interface{} `json:"content" form:"content"`                                 // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdatePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}

// CreatePoddisruptionbudgetReq 创建Pod中断预算请求结构体
type CreatePoddisruptionbudgetReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreatePoddisruptionbudgetReq) GetClusterID() int {
	return r.ClusterId
}
