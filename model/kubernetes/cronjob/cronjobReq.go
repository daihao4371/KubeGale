package cronjob

import "KubeGale/model/common/request"

// GetCronJobListReq 获取CronJob列表请求结构体
type GetCronJobListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	Keyword          string `json:"keyword" form:"keyword"`             // 关键字搜索
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetCronJobListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeCronJobReq 获取CronJob详情请求结构体
type DescribeCronJobReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string `json:"namespace" form:"namespace"`     // 命名空间
	CronJobName string `json:"cronjobName" form:"cronjobName"` // CronJob名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeCronJobReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteCronJobReq 删除CronJob请求结构体
type DeleteCronJobReq struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string `json:"namespace" form:"namespace"`     // 命名空间
	CronJobName string `json:"cronjobName" form:"cronjobName"` // CronJob名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteCronJobReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateCronJobReq 更新CronJob请求结构体
type UpdateCronJobReq struct {
	ClusterId   int         `json:"cluster_id" form:"cluster_id"`   // 集群ID
	Namespace   string      `json:"namespace" form:"namespace"`     // 命名空间
	CronJobName string      `json:"cronjobName" form:"cronjobName"` // CronJob名称
	Content     interface{} `json:"content" form:"content"`         // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateCronJobReq) GetClusterID() int {
	return r.ClusterId
}

// CreateCronJobReq 创建CronJob请求结构体
type CreateCronJobReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *CreateCronJobReq) GetClusterID() int {
	return r.ClusterId
}
