package roles

import "KubeGale/model/common/request"

// GetRolesListReq 获取角色列表请求结构体
type GetRolesListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetRolesListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeRolesReq 获取角色详情请求结构体
type DescribeRolesReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	RolesName string `json:"roleName" form:"roleName"`     // 角色名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeRolesReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteRolesReq 删除角色请求结构体
type DeleteRolesReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	RolesName string `json:"roleName" form:"roleName"`     // 角色名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteRolesReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateRolesReq 更新角色请求结构体
type UpdateRolesReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	RolesName string      `json:"roleName" form:"roleName"`     // 角色名称
	Content   interface{} `json:"content" form:"content"`       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateRolesReq) GetClusterID() int {
	return r.ClusterId
}

// CreateRolesReq 创建角色请求结构体
type CreateRolesReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateRolesReq) GetClusterID() int {
	return r.ClusterId
}
