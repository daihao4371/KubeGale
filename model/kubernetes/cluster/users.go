package cluster

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

// User Kubernetes集群用户信息结构体
// 用于存储和管理Kubernetes集群中的用户信息，包括用户认证、权限等配置
type User struct {
	ID             uint           `json:"id" gorm:"not null;unique;primary_key"`                                          // 主键ID，唯一标识用户
	UUID           uuid.UUID      `json:"uuid" gorm:"comment:用户UUID"`                                                     // 用户UUID，用于全局唯一标识
	Username       string         `json:"userName" gorm:"comment:用户登录名"`                                                  // 用户登录名，用于系统登录
	NickName       string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                      // 用户昵称，显示名称
	KubeConfig     string         `gorm:"type:longText" json:"kube_config" form:"kube_config" gorm:"comment:kube_config"` // 用户的Kubernetes配置文件内容
	ClusterRoles   string         `gorm:"type:longText" json:"cluster_roles"`                                             // 用户在集群级别的角色配置
	NamespaceRoles string         `json:"namespace_roles" gorm:"comment:命名空间权限"`                                          // 用户在命名空间级别的权限配置
	ClusterId      uint           `json:"cluster_id" gorm:"comment:集群ID"`                                                 // 用户所属的集群ID
	CreatedAt      time.Time      `json:"created_at"`                                                                     // 记录创建时间
	UpdatedAt      time.Time      `json:"updated_at"`                                                                     // 记录最后更新时间
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                                                                 // 软删除标记，用于记录删除时间
}

// TableName 指定User结构体对应的数据库表名
// 返回值为"k8s_users"，表示该结构体映射到数据库中的k8s_users表
func (User) TableName() string {
	return "k8s_users"
}
