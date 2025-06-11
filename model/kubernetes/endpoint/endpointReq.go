package endpoint

import "KubeGale/model/common/request"

// GetEndPointListReq 获取EndpointSlice列表请求结构体
type GetEndPointListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	Keyword          string `json:"keyword" form:"keyword"`             // 关键字搜索
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetEndPointListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeEndPointReq 获取EndpointSlice详情请求结构体
type DescribeEndPointReq struct {
	ClusterId    int    `json:"cluster_id" form:"cluster_id"`     // 集群ID
	Namespace    string `json:"namespace" form:"namespace"`       // 命名空间
	EndPointName string `json:"endpointName" form:"endpointName"` // EndpointSlice名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeEndPointReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteEndPointReq 删除EndpointSlice请求结构体
type DeleteEndPointReq struct {
	ClusterId    int    `json:"cluster_id" form:"cluster_id"`     // 集群ID
	Namespace    string `json:"namespace" form:"namespace"`       // 命名空间
	EndPointName string `json:"endpointName" form:"endpointName"` // EndpointSlice名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteEndPointReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateEndPointReq 更新EndpointSlice请求结构体
type UpdateEndPointReq struct {
	ClusterId    int         `json:"cluster_id" form:"cluster_id"`     // 集群ID
	Namespace    string      `json:"namespace" form:"namespace"`       // 命名空间
	EndPointName string      `json:"endpointName" form:"endpointName"` // EndpointSlice名称
	Content      interface{} `json:"content" form:"content"`           // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateEndPointReq) GetClusterID() int {
	return r.ClusterId
}

// CreateEndPointReq 创建EndpointSlice请求结构体
type CreateEndPointReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateEndPointReq) GetClusterID() int {
	return r.ClusterId
}
