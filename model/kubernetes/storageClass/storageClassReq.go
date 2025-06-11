package storageClass

import "KubeGale/model/common/request"

// GetStorageClassListReq 获取存储类列表请求结构体
type GetStorageClassListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetStorageClassListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeStorageClassReq 获取存储类详情请求结构体
type DescribeStorageClassReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`             // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`               // 命名空间
	StorageClassName string `json:"storageClassName" form:"storageClassName"` // 存储类名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeStorageClassReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteStorageClassReq 删除存储类请求结构体
type DeleteStorageClassReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`             // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`               // 命名空间
	StorageClassName string `json:"storageClassName" form:"storageClassName"` // 存储类名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteStorageClassReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateStorageClassReq 更新存储类请求结构体
type UpdateStorageClassReq struct {
	ClusterId        int         `json:"cluster_id" form:"cluster_id"`             // 集群ID
	Namespace        string      `json:"namespace" form:"namespace"`               // 命名空间
	StorageClassName string      `json:"storageClassName" form:"storageClassName"` // 存储类名称
	Content          interface{} `json:"content" form:"content"`                   // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateStorageClassReq) GetClusterID() int {
	return r.ClusterId
}

// CreateStorageClassReq 创建存储类请求结构体
type CreateStorageClassReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateStorageClassReq) GetClusterID() int {
	return r.ClusterId
}
