package podSecurityPolicies

import "KubeGale/model/common/request"

// GetPodSecurityPoliciesListReq 获取Pod安全策略列表请求结构体
type GetPodSecurityPoliciesListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetPodSecurityPoliciesListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribePodSecurityPoliciesReq 获取Pod安全策略详情请求结构体
type DescribePodSecurityPoliciesReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`                           // 集群ID
	Namespace               string `json:"namespace" form:"namespace"`                             // 命名空间
	PodSecurityPoliciesName string `json:"podSecurityPoliciesName" form:"podSecurityPoliciesName"` // Pod安全策略名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}

// DeletePodSecurityPoliciesReq 删除Pod安全策略请求结构体
type DeletePodSecurityPoliciesReq struct {
	ClusterId               int    `json:"cluster_id" form:"cluster_id"`                           // 集群ID
	Namespace               string `json:"namespace" form:"namespace"`                             // 命名空间
	PodSecurityPoliciesName string `json:"podSecurityPoliciesName" form:"podSecurityPoliciesName"` // Pod安全策略名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeletePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}

// UpdatePodSecurityPoliciesReq 更新Pod安全策略请求结构体
type UpdatePodSecurityPoliciesReq struct {
	ClusterId               int         `json:"cluster_id" form:"cluster_id"`                           // 集群ID
	Namespace               string      `json:"namespace" form:"namespace"`                             // 命名空间
	PodSecurityPoliciesName string      `json:"podSecurityPoliciesName" form:"podSecurityPoliciesName"` // Pod安全策略名称
	Content                 interface{} `json:"content" form:"content"`                                 // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdatePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}

// CreatePodSecurityPoliciesReq 创建Pod安全策略请求结构体
type CreatePodSecurityPoliciesReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreatePodSecurityPoliciesReq) GetClusterID() int {
	return r.ClusterId
}
