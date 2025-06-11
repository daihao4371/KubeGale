package serviceAccount

import "KubeGale/model/common/request"

// GetServiceAccountReq 获取服务账号列表请求结构体
type GetServiceAccountReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeServiceAccountReq 获取服务账号详情请求结构体
type DescribeServiceAccountReq struct {
	ClusterId          int    `json:"cluster_id" form:"cluster_id"`                 // 集群ID
	Namespace          string `json:"namespace" form:"namespace"`                   // 命名空间
	ServiceAccountName string `json:"serviceAccountName" form:"serviceAccountName"` // 服务账号名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteServiceAccountReq 删除服务账号请求结构体
type DeleteServiceAccountReq struct {
	ClusterId          int    `json:"cluster_id" form:"cluster_id"`                 // 集群ID
	Namespace          string `json:"namespace" form:"namespace"`                   // 命名空间
	ServiceAccountName string `json:"serviceAccountName" form:"serviceAccountName"` // 服务账号名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateServiceAccountReq 更新服务账号请求结构体
type UpdateServiceAccountReq struct {
	ClusterId          int         `json:"cluster_id" form:"cluster_id"`                 // 集群ID
	Namespace          string      `json:"namespace" form:"namespace"`                   // 命名空间
	ServiceAccountName string      `json:"serviceAccountName" form:"serviceAccountName"` // 服务账号名称
	Content            interface{} `json:"content" form:"content"`                       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}

// CreateServiceAccountReq 创建服务账号请求结构体
type CreateServiceAccountReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateServiceAccountReq) GetClusterID() int {
	return r.ClusterId
}
