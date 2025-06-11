package ingress

import "KubeGale/model/common/request"

// GetIngressListReq 获取Ingress列表请求结构体
type GetIngressListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetIngressListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeIngressReq 获取Ingress详情请求结构体
type DescribeIngressReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string `json:"namespace" form:"namespace"`     // 命名空间
	IngressName string `json:"ingressName" form:"ingressName"` // Ingress名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeIngressReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteIngressReq 删除Ingress请求结构体
type DeleteIngressReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string `json:"namespace" form:"namespace"`     // 命名空间
	IngressName string `json:"ingressName" form:"ingressName"` // Ingress名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteIngressReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateIngressReq 更新Ingress请求结构体
type UpdateIngressReq struct {
	ClusterId   int         `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string      `json:"namespace" form:"namespace"`     // 命名空间
	IngressName string      `json:"ingressName" form:"ingressName"` // Ingress名称
	Content     interface{} `json:"content" form:"content"`         // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateIngressReq) GetClusterID() int {
	return r.ClusterId
}

// CreateIngressReq 创建Ingress请求结构体
type CreateIngressReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateIngressReq) GetClusterID() int {
	return r.ClusterId
}
