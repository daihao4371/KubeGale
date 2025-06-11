package statefulSet

import "KubeGale/model/common/request"

// GetStatefulSetListReq 获取有状态应用列表请求结构体
type GetStatefulSetListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	Keyword          string `json:"keyword" form:"keyword"`             // 关键字
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetStatefulSetListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeStatefulSetReq 获取有状态应用详情请求结构体
type DescribeStatefulSetReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string `json:"namespace" form:"namespace"`             // 命名空间
	StatefulSetName string `json:"statefulsetName" form:"statefulsetName"` // 有状态应用名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteStatefulSetReq 删除有状态应用请求结构体
type DeleteStatefulSetReq struct {
	ClusterId       int    `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string `json:"namespace" form:"namespace"`             // 命名空间
	StatefulSetName string `json:"statefulsetName" form:"statefulsetName"` // 有状态应用名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateStatefulSetReq 更新有状态应用请求结构体
type UpdateStatefulSetReq struct {
	ClusterId       int         `json:"cluster_id" form:"cluster_id"`           // 集群ID
	Namespace       string      `json:"namespace" form:"namespace"`             // 命名空间
	StatefulSetName string      `json:"statefulsetName" form:"statefulsetName"` // 有状态应用名称
	Content         interface{} `json:"content" form:"content"`                 // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}

// CreateStatefulSetReq 创建有状态应用请求结构体
type CreateStatefulSetReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateStatefulSetReq) GetClusterID() int {
	return r.ClusterId
}
