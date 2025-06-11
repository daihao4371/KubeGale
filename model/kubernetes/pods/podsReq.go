package pods

import (
	"KubeGale/model/common/request"
	"io"
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

// PodRecord Pod操作记录结构体
// 用于记录Pod相关的操作历史
type PodRecord struct {
	ID            uint           `json:"id" gorm:"not null;unique;primary_key"`                       // 记录ID
	Cluster       string         `json:"cluster" gorm:"comment:集群名称"`                                 // 集群名称
	Namespace     string         `json:"namespace" gorm:"comment:命名空间"`                               // 命名空间
	PodName       string         `json:"pod_name" gorm:"comment:Pod名称"`                               // Pod名称
	ContainerName string         `json:"container_name" gorm:"comment:容器名称"`                          // 容器名称
	UUID          uuid.UUID      `json:"uuid" gorm:"comment:UUID"`                                    // UUID
	Username      string         `json:"userName" gorm:"comment:操作用户"`                                // 操作用户
	NickName      string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                   // 用户昵称
	Records       []byte         `json:"records" gorm:"type:longblob;comment:'操作记录(二进制存储)';size:128"` // 操作记录
	CreatedAt     time.Time      `json:"created_at"`                                                  // 创建时间
	UpdatedAt     time.Time      `json:"updated_at"`                                                  // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                                              // 删除时间
}

// PodsFilesRequest Pod文件操作请求结构体
// 用于处理Pod内文件相关的操作请求
type PodsFilesRequest struct {
	ClusterId     int       `json:"cluster_id" form:"cluster_id" validate:"required"` // 集群ID
	Folder        string    `json:"folder" form:"folder"`                             // 文件夹路径
	PodName       string    `json:"podName" form:"podName" validate:"required"`       // Pod名称
	ContainerName string    `json:"containerName" form:"containerName"`               // 容器名称
	Namespace     string    `json:"namespace" form:"namespace" validate:"required"`   // 命名空间
	Path          string    `json:"path" form:"path"`                                 // 文件路径
	OldPath       string    `json:"oldPath" form:"oldPath"`                           // 原文件路径
	Commands      []string  `json:"-" form:"Commands"`                                // 命令列表
	Stdin         io.Reader `json:"-" form:"Stdin"`                                   // 标准输入
	Content       string    `json:"content" form:"content"`                           // 文件内容
	FilePath      string    `json:"filePath" form:"filePath"`                         // 文件路径
	XToken        string    `json:"x-token" form:"x-token"`                           // 认证令牌
}

// TableName 指定表名
// 返回Pod记录的表名
func (PodRecord) TableName() string {
	return "k8s_pod_records"
}

// PodListReq 获取Pod列表请求结构体
type PodListReq struct {
	ClusterId        int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace        string `json:"namespace" form:"namespace"`         // 命名空间
	LabelSelector    string `json:"labelSelector" form:"labelSelector"` // 标签选择器
	FieldSelector    string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
	request.PageInfo        // 分页信息
}

// PodMetricsReq 获取Pod指标请求结构体
type PodMetricsReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
}

// DescribePodInfo 获取Pod详情请求结构体
type DescribePodInfo struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	PodName   string `json:"podName" form:"podName"`       // Pod名称
}

// CreatePodReq 创建Pod请求结构体
type CreatePodReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	Content   interface{} `json:"content" form:"content"`       // 创建内容
}

// DeletePodReq 删除Pod请求结构体
type DeletePodReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string `json:"namespace" form:"namespace"`   // 命名空间
	PodName   string `json:"podName" form:"podName"`       // Pod名称
}

// UpdatePodReq 更新Pod请求结构体
type UpdatePodReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"` // 集群ID
	Namespace string      `json:"namespace" form:"namespace"`   // 命名空间
	PodName   string      `json:"podName" form:"podName"`       // Pod名称
	Content   interface{} `json:"content" form:"content"`       // 更新内容
}

// PodEventsReq 获取Pod事件请求结构体
type PodEventsReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string `json:"namespace" form:"namespace"`         // 命名空间
	PodName       string `json:"podName" form:"podName"`             // Pod名称
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"` // 字段选择器
}

// ListPodFiles 获取Pod文件列表请求结构体
type ListPodFiles struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`       // 集群ID
	Namespace     string `json:"namespace" form:"namespace"`         // 命名空间
	PodName       string `json:"podName" form:"podName"`             // Pod名称
	ContainerName string `json:"containerName" form:"containerName"` // 容器名称
	Path          string `json:"path" form:"path"`                   // 文件路径
}
