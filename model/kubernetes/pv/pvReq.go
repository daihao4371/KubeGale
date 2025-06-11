package pv

import "KubeGale/model/common/request"

// GetPVListReq 获取持久卷列表请求结构体
type GetPVListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetPVListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribePVReq 获取持久卷详情请求结构体
type DescribePVReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	PVName    string `json:"pvName" form:"pvName"`         // 持久卷名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribePVReq) GetClusterID() int {
	return r.ClusterId
}

// DeletePVReq 删除持久卷请求结构体
type DeletePVReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	PVName    string `json:"pvName" form:"pvName"`         // 持久卷名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeletePVReq) GetClusterID() int {
	return r.ClusterId
}

// UpdatePVReq 更新持久卷请求结构体
type UpdatePVReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	PVName    string      `json:"pvName" form:"pvName"`         // 持久卷名称
	Content   interface{} `json:"content" form:"content"`       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdatePVReq) GetClusterID() int {
	return r.ClusterId
}

// CreatePVReq 创建持久卷请求结构体
type CreatePVReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreatePVReq) GetClusterID() int {
	return r.ClusterId
}
