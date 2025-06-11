package resourceQuota

import "KubeGale/model/common/request"

// GetResourceQuotaListReq 获取资源配额列表请求结构体
type GetResourceQuotaListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetResourceQuotaListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeResourceQuotaReq 获取资源配额详情请求结构体
type DescribeResourceQuotaReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string `json:"namespace" form:"namespace"`                 // 命名空间
	ResourceQuotaName string `json:"resourcequotaName" form:"resourcequotaName"` // 资源配额名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteResourceQuotaReq 删除资源配额请求结构体
type DeleteResourceQuotaReq struct {
	ClusterId         int    `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string `json:"namespace" form:"namespace"`                 // 命名空间
	ResourceQuotaName string `json:"resourcequotaName" form:"resourcequotaName"` // 资源配额名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateResourceQuotaReq 更新资源配额请求结构体
type UpdateResourceQuotaReq struct {
	ClusterId         int         `json:"cluster_id" form:"cluster_id"`               // 集群ID
	Namespace         string      `json:"namespace" form:"namespace"`                 // 命名空间
	ResourceQuotaName string      `json:"resourcequotaName" form:"resourcequotaName"` // 资源配额名称
	Content           interface{} `json:"content" form:"content"`                     // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}

// CreateResourceQuotaReq 创建资源配额请求结构体
type CreateResourceQuotaReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateResourceQuotaReq) GetClusterID() int {
	return r.ClusterId
}
