package job

import "KubeGale/model/common/request"

// GetJobListReq 获取Job列表请求结构体
type GetJobListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	Keyword          string `json:"keyword" form:"keyword"`             // 关键字搜索
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetJobListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeJobReq 获取Job详情请求结构体
type DescribeJobReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	JobName   string `json:"jobName" form:"jobName"`       // Job名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeJobReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteJobReq 删除Job请求结构体
type DeleteJobReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	JobName   string `json:"jobName" form:"jobName"`       // Job名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteJobReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateJobReq 更新Job请求结构体
type UpdateJobReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	JobName   string      `json:"jobName" form:"jobName"`       // Job名称
	Content   interface{} `json:"content" form:"content"`       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateJobReq) GetClusterID() int {
	return r.ClusterId
}

// CreateJobReq 创建Job请求结构体
type CreateJobReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateJobReq) GetClusterID() int {
	return r.ClusterId
}
