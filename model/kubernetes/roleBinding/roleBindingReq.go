package roleBinding

import "KubeGale/model/common/request"

// GetRoleBindingListReq 获取角色绑定列表请求结构体
type GetRoleBindingListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetRoleBindingListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeRoleBindingReq 获取角色绑定详情请求结构体
type DescribeRoleBindingReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string `json:"namespace" form:"namespace"`             // 命名空间
	RoleBindingName string `json:"roleBindingName" form:"roleBindingName"` // 角色绑定名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteRoleBindingReq 删除角色绑定请求结构体
type DeleteRoleBindingReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string `json:"namespace" form:"namespace"`             // 命名空间
	RoleBindingName string `json:"roleBindingName" form:"roleBindingName"` // 角色绑定名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateRoleBindingReq 更新角色绑定请求结构体
type UpdateRoleBindingReq struct {
	ClusterId       int         `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string      `json:"namespace" form:"namespace"`             // 命名空间
	RoleBindingName string      `json:"roleBindingName" form:"roleBindingName"` // 角色绑定名称
	Content         interface{} `json:"content" form:"content"`                 // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}

// CreateRoleBindingReq 创建角色绑定请求结构体
type CreateRoleBindingReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateRoleBindingReq) GetClusterID() int {
	return r.ClusterId
}
