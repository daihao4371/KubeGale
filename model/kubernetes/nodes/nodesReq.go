package nodes

import (
	"KubeGale/model/common/request"
)

// NodeListReq 获取节点列表请求结构体
type NodeListReq struct {
	ClusterId        int `json:"cluster_id" form:"cluster_id"` // 集群ID
	request.PageInfo     // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *NodeListReq) GetClusterID() int {
	return r.ClusterId
}

// NodeMetricsReq 获取节点指标请求结构体
type NodeMetricsReq struct {
	ClusterId int `json:"cluster_id" form:"cluster_id"` // 集群ID
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *NodeMetricsReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeNodeReq 获取节点详情请求结构体
type DescribeNodeReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	NodeName  string `json:"nodeName" form:"nodeName"`     // 节点名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeNodeReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateNodeReq 更新节点请求结构体
type UpdateNodeReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	NodeName  string      `json:"nodeName" form:"nodeName"`     // 节点名称
	Content   interface{} `json:"content" form:"content"`       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateNodeReq) GetClusterID() int {
	return r.ClusterId
}

// EvictAllNodePodReq 驱逐节点上所有Pod的请求结构体
type EvictAllNodePodReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	NodeName      string `json:"nodeName" form:"nodeName"`           // 节点名称
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *EvictAllNodePodReq) GetClusterID() int {
	return r.ClusterId
}
