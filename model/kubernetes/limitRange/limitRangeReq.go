package limitRange

import "KubeGale/model/common/request"

// GetLimitRangeListReq 获取LimitRange列表请求结构体
type GetLimitRangeListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetLimitRangeListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeLimitRangeReq 获取LimitRange详情请求结构体
type DescribeLimitRangeReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`         // 集群ID
	Namespace      string `json:"namespace" form:"namespace"`           // 命名空间
	LimitRangeName string `json:"limitRangeName" form:"limitRangeName"` // LimitRange名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteLimitRangeReq 删除LimitRange请求结构体
type DeleteLimitRangeReq struct {
	ClusterId      int    `json:"cluster_id" form:"cluster_id"`         // 集群ID
	Namespace      string `json:"namespace" form:"namespace"`           // 命名空间
	LimitRangeName string `json:"limitRangeName" form:"limitRangeName"` // LimitRange名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateLimitRangeReq 更新LimitRange请求结构体
type UpdateLimitRangeReq struct {
	ClusterId      int         `json:"cluster_id" form:"cluster_id"`         // 集群ID
	Namespace      string      `json:"namespace" form:"namespace"`           // 命名空间
	LimitRangeName string      `json:"limitRangeName" form:"limitRangeName"` // LimitRange名称
	Content        interface{} `json:"content" form:"content"`               // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}

// CreateLimitRangeReq 创建LimitRange请求结构体
type CreateLimitRangeReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateLimitRangeReq) GetClusterID() int {
	return r.ClusterId
}
