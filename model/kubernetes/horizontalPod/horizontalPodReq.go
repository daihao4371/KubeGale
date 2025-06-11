package horizontalPod

import "KubeGale/model/common/request"

// GetHorizontalPodListReq 获取水平Pod自动伸缩器列表请求结构体
type GetHorizontalPodListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetHorizontalPodListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeHorizontalPodReq 获取水平Pod自动伸缩器详情请求结构体
type DescribeHorizontalPodReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string `json:"namespace" form:"namespace"`                 // 命名空间
	HorizontalPodName string `json:"HorizontalPodName" form:"HorizontalPodName"` // 水平Pod自动伸缩器名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteHorizontalPodReq 删除水平Pod自动伸缩器请求结构体
type DeleteHorizontalPodReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string `json:"namespace" form:"namespace"`                 // 命名空间
	HorizontalPodName string `json:"HorizontalPodName" form:"HorizontalPodName"` // 水平Pod自动伸缩器名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateHorizontalPodReq 更新水平Pod自动伸缩器请求结构体
type UpdateHorizontalPodReq struct {
	ClusterId         int         `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string      `json:"namespace" form:"namespace"`                 // 命名空间
	HorizontalPodName string      `json:"HorizontalPodName" form:"HorizontalPodName"` // 水平Pod自动伸缩器名称
	Content           interface{} `json:"content" form:"content"`                     // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}

// CreateHorizontalPodReq 创建水平Pod自动伸缩器请求结构体
type CreateHorizontalPodReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateHorizontalPodReq) GetClusterID() int {
	return r.ClusterId
}
