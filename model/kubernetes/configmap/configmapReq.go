package configmap

import "KubeGale/model/common/request"

// GetConfigMapListReq 获取ConfigMap列表请求结构体
type GetConfigMapListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetConfigMapListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeConfigMapReq 获取ConfigMap详情请求结构体
type DescribeConfigMapReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string `json:"namespace" form:"namespace"`         // 命名空间
	ConfigMapName string `json:"configmapName" form:"configmapName"` // ConfigMap名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeConfigMapReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteConfigMapReq 删除ConfigMap请求结构体
type DeleteConfigMapReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string `json:"namespace" form:"namespace"`         // 命名空间
	ConfigMapName string `json:"configmapName" form:"configmapName"` // ConfigMap名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteConfigMapReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateConfigMapReq 更新ConfigMap请求结构体
type UpdateConfigMapReq struct {
	ClusterId     int         `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string      `json:"namespace" form:"namespace"`         // 命名空间
	ConfigMapName string      `json:"configmapName" form:"configmapName"` // ConfigMap名称
	Content       interface{} `json:"content" form:"content"`             // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateConfigMapReq) GetClusterID() int {
	return r.ClusterId
}

// CreateConfigMapReq 创建ConfigMap请求结构体
type CreateConfigMapReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateConfigMapReq) GetClusterID() int {
	return r.ClusterId
}
