package request

import "KubeGale/model/common/request"

// GetVeleroListReq 获取Velero列表请求结构体
type GetVeleroListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	Keyword          string `json:"keyword" form:"keyword"`             // 关键字
	request.PageInfo        // 分页信息
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *GetVeleroListReq) GetClusterID() int {
	return r.ClusterId
}

// DescribeVeleroReq 获取Velero详情请求结构体
type DescribeVeleroReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace  string `json:"namespace" form:"namespace"`   // 命名空间
	VeleroName string `json:"VeleroName" form:"VeleroName"` // Velero名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DescribeVeleroReq) GetClusterID() int {
	return r.ClusterId
}

// DeleteVeleroReq 删除Velero请求结构体
type DeleteVeleroReq struct {
	ClusterId  int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace  string `json:"namespace" form:"namespace"`   // 命名空间
	VeleroName string `json:"VeleroName" form:"VeleroName"` // Velero名称
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *DeleteVeleroReq) GetClusterID() int {
	return r.ClusterId
}

// UpdateVeleroReq 更新Velero请求结构体
type UpdateVeleroReq struct {
	ClusterId  int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace  string      `json:"namespace" form:"namespace"`   // 命名空间
	VeleroName string      `json:"VeleroName" form:"VeleroName"` // Velero名称
	Content    interface{} `json:"content" form:"content"`       // 更新内容
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *UpdateVeleroReq) GetClusterID() int {
	return r.ClusterId
}

// VeleroModel Velero模型结构体
// 用于存储Velero的配置信息
type VeleroModel struct {
	ClusterId   int    `json:"cluster_id" form:"cluster_id"`   // 集群ID
	S3Region    string `json:"s3Region" form:"s3Region"`       // S3区域
	S3Address   string `json:"s3Address" form:"s3Address"`     // S3地址
	S3Key       string `json:"s3Key" form:"s3Key"`             // S3密钥
	S3Secret    string `json:"s3Secret" form:"s3Secret"`       // S3密钥
	S3Bucket    string `json:"s3Bucket" form:"s3Bucket"`       // S3存储桶
	Provider    string `json:"provider"`                       // 提供商
	VeleroImage string `json:"veleroImage" form:"veleroImage"` // Velero镜像
	PluginImage string `json:"pluginImage" form:"pulginImage"` // 插件镜像
}

// TableName 获取表名
// 返回数据库表名
func (v *VeleroModel) TableName() string {
	return "k8s_velero"
}

// GetClusterID 获取集群ID
// 实现获取集群ID的接口方法
func (r *VeleroModel) GetClusterID() int {
	return r.ClusterId
}
