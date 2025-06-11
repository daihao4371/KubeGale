package clusterrole

import "KubeGale/model/common/request"

// GetClusterRoleListReq 获取集群角色列表请求结构体
type GetClusterRoleListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetClusterRoleListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeClusterRoleReq 获取集群角色详情请求结构体
type DescribeClusterRoleReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string `json:"namespace" form:"namespace"`             // 命名空间
	ClusterRoleName string `json:"clusterRoleName" form:"clusterRoleName"` // 集群角色名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteClusterRoleReq 删除集群角色请求结构体
type DeleteClusterRoleReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string `json:"namespace" form:"namespace"`             // 命名空间
	ClusterRoleName string `json:"clusterRoleName" form:"clusterRoleName"` // 集群角色名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateClusterRoleReq 更新集群角色请求结构体
type UpdateClusterRoleReq struct {
	ClusterId       int         `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string      `json:"namespace" form:"namespace"`             // 命名空间
	ClusterRoleName string      `json:"clusterRoleName" form:"clusterRoleName"` // 集群角色名称
	Content         interface{} `json:"content" form:"content"`                 // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}

// CreateClusterRoleReq 创建集群角色请求结构体
type CreateClusterRoleReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateClusterRoleReq) GetClusterID() int {
	return r.ClusterId
}
