package service

import "KubeGale/model/common/request"

// GetServiceListReq 获取服务列表请求结构体
type GetServiceListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetServiceListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeServiceReq 获取服务详情请求结构体
type DescribeServiceReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string `json:"namespace" form:"namespace"`     // 命名空间
	ServiceName string `json:"serviceName" form:"serviceName"` // 服务名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeServiceReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteServiceReq 删除服务请求结构体
type DeleteServiceReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`   // 集群ID
	ServiceName string `json:"serviceName" form:"serviceName"` // 服务名称
	Namespace   string `json:"namespace" form:"namespace"`     // 命名空间
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteServiceReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateServiceReq 更新服务请求结构体
type UpdateServiceReq struct {
	ClusterId   int         `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string      `json:"namespace" form:"namespace"`     // 命名空间
	ServiceName string      `json:"serviceName" form:"serviceName"` // 服务名称
	Content     interface{} `json:"content" form:"content"`         // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateServiceReq) GetClusterID() int {
	return r.ClusterId
}

// CreateServiceReq 创建服务请求结构体
type CreateServiceReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateServiceReq) GetClusterID() int {
	return r.ClusterId
}
