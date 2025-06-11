package namespaces

import "KubeGale/model/common/request"

// GetNamespaceListReq 获取命名空间列表请求结构体
type GetNamespaceListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	Keyword          string `json:"keyword" form:"keyword"`             // 关键字搜索
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetNamespaceListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeNamespaceReq 获取命名空间详情请求结构体
type DescribeNamespaceReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string `json:"namespace" form:"namespace"`         // 命名空间
	NamespaceName string `json:"NamespaceName" form:"NamespaceName"` // 命名空间名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeNamespaceReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteNamespaceReq 删除命名空间请求结构体
type DeleteNamespaceReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	NamespaceName string `json:"NamespaceName" form:"NamespaceName"` // 命名空间名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteNamespaceReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateNamespaceReq 更新命名空间请求结构体
type UpdateNamespaceReq struct {
	ClusterId     int         `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string      `json:"namespace" form:"namespace"`         // 命名空间
	NamespaceName string      `json:"NamespaceName" form:"NamespaceName"` // 命名空间名称
	Content       interface{} `json:"content" form:"content"`             // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateNamespaceReq) GetClusterID() int {
	return r.ClusterId
}

// CreateNamespaceReq 创建命名空间请求结构体
type CreateNamespaceReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateNamespaceReq) GetClusterID() int {
	return r.ClusterId
}
