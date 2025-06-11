package daemonset

import "KubeGale/model/common/request"

// GetDaemonSetListReq 获取DaemonSet列表请求结构体
type GetDaemonSetListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	Keyword          string `json:"keyword" form:"keyword"`             // 关键字搜索
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetDaemonSetListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeDaemonSetReq 获取DaemonSet详情请求结构体
type DescribeDaemonSetReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string `json:"namespace" form:"namespace"`         // 命名空间
	DaemonsetName string `json:"daemonsetName" form:"daemonsetName"` // DaemonSet名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteDaemonSetReq 删除DaemonSet请求结构体
type DeleteDaemonSetReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string `json:"namespace" form:"namespace"`         // 命名空间
	DaemonsetName string `json:"daemonsetName" form:"daemonsetName"` // DaemonSet名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateDaemonSetReq 更新DaemonSet请求结构体
type UpdateDaemonSetReq struct {
	ClusterId     int         `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string      `json:"namespace" form:"namespace"`         // 命名空间
	DaemonsetName string      `json:"daemonsetName" form:"daemonsetName"` // DaemonSet名称
	Content       interface{} `json:"content" form:"content"`             // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}

// CreateDaemonSetReq 创建DaemonSet请求结构体
type CreateDaemonSetReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateDaemonSetReq) GetClusterID() int {
	return r.ClusterId
}
