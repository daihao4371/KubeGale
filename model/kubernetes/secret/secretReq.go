package secret

import "KubeGale/model/common/request"

// GetSecretList 获取密钥列表请求结构体
type GetSecretList struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`   // 命名空间
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetSecretList) GetClusterID() int {
	return r.ClusterId
}

// DescribeSecretReq 获取密钥详情请求结构体
type DescribeSecretReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace  string `json:"namespace" form:"namespace"`   // 命名空间
	SecretName string `json:"secretName" form:"secretName"` // 密钥名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeSecretReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteSecretReq 删除密钥请求结构体
type DeleteSecretReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace  string `json:"namespace" form:"namespace"`   // 命名空间
	SecretName string `json:"secretName" form:"secretName"` // 密钥名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteSecretReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateSecretReq 更新密钥请求结构体
type UpdateSecretReq struct {
	ClusterId  int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace  string      `json:"namespace" form:"namespace"`   // 命名空间
	SecretName string      `json:"secretName" form:"secretName"` // 密钥名称
	Content    interface{} `json:"content" form:"content"`       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateSecretReq) GetClusterID() int {
	return r.ClusterId
}

// CreateSecretReq 创建密钥请求结构体
type CreateSecretReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateSecretReq) GetClusterID() int {
	return r.ClusterId
}
