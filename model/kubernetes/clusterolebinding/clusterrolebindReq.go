package clusterolebinding

import "KubeGale/model/common/request"

// GetClusterRoleBindingListReq 获取集群角色绑定列表请求结构体
type GetClusterRoleBindingListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetClusterRoleBindingListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeClusterRoleBindingReq 获取集群角色绑定详情请求结构体
type DescribeClusterRoleBindingReq struct {
	ClusterId              int    `json:"cluster_id" form:"cluster_id"`                         // 集群ID
	Namespace              string `json:"namespace" form:"namespace"`                           // 命名空间
	ClusterRoleBindingName string `json:"clusterRoleBindingName" form:"clusterRoleBindingName"` // 集群角色绑定名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteClusterRoleBindingReq 删除集群角色绑定请求结构体
type DeleteClusterRoleBindingReq struct {
	ClusterId              int    `json:"cluster_id" form:"cluster_id"`                         // 集群ID
	Namespace              string `json:"namespace" form:"namespace"`                           // 命名空间
	ClusterRoleBindingName string `json:"clusterRoleBindingName" form:"clusterRoleBindingName"` // 集群角色绑定名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateClusterRoleBindingReq 更新集群角色绑定请求结构体
type UpdateClusterRoleBindingReq struct {
	ClusterId              int         `json:"cluster_id" form:"cluster_id"`                         // 集群ID
	Namespace              string      `json:"namespace" form:"namespace"`                           // 命名空间
	ClusterRoleBindingName string      `json:"clusterRoleBindingName" form:"clusterRoleBindingName"` // 集群角色绑定名称
	Content                interface{} `json:"content" form:"content"`                               // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

// CreateClusterRoleBindingReq 创建集群角色绑定请求结构体
type CreateClusterRoleBindingReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateClusterRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}
