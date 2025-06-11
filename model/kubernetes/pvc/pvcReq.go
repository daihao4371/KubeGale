package pvc

import "KubeGale/model/common/request"

// GetPvcListReq 获取持久卷声明列表请求结构体
type GetPvcListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetPvcListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribePVCReq 获取持久卷声明详情请求结构体
type DescribePVCReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	PVCName   string `json:"pvcName" form:"pvcName"`       // 持久卷声明名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribePVCReq) GetClusterID() int {
	return r.ClusterId
}

// DeletePVCReq 删除持久卷声明请求结构体
type DeletePVCReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	PVCName   string `json:"pvcName" form:"pvcName"`       // 持久卷声明名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeletePVCReq) GetClusterID() int {
	return r.ClusterId
}

// UpdatePVCReq 更新持久卷声明请求结构体
type UpdatePVCReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	PVCName   string      `json:"pvcName" form:"pvcName"`       // 持久卷声明名称
	Content   interface{} `json:"content" form:"content"`       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdatePVCReq) GetClusterID() int {
	return r.ClusterId
}

// CreatePVCReq 创建持久卷声明请求结构体
type CreatePVCReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreatePVCReq) GetClusterID() int {
	return r.ClusterId
}
