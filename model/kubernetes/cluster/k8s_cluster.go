// 自动生成模板K8sCluster
package cluster

import (
	"KubeGale/global"
	"github.com/gofrs/uuid/v5"
	rbacV1 "k8s.io/api/rbac/v1"
)

// K8sCluster Kubernetes集群信息结构体
type K8sCluster struct {
	global.KUBEGALE_MODEL
	ID                 uint      `json:"id" gorm:"not null;unique;primary_key"`                                                  // 主键ID
	UUID               uuid.UUID `json:"uuid" gorm:"comment:集群UUID"`                                                             // 集群唯一标识
	Name               string    `json:"name" form:"name" gorm:"comment:集群名称"`                                                   // 集群名称
	KubeType           uint      `json:"kube_type" form:"kube_type" gorm:"comment:凭据类型1:KubeConfig,2:Token"`                     // 认证类型
	KubeConfig         string    `gorm:"type:longText" json:"kube_config" form:"kube_config" gorm:"comment:kube_config"`         // KubeConfig配置内容
	KubeCaCrt          string    `gorm:"type:longText; comment:ca.crt" json:"kube_ca_crt" form:"kube_ca_crt"`                    // CA证书内容
	ApiAddress         string    `gorm:"type:longText" json:"api_address" form:"api_address" gorm:"comment:api_address"`         // API服务器地址
	PrometheusUrl      string    `gorm:"type:longText" json:"prometheus_url" form:"prometheus_url" gorm:"comment:prometheus 地址"` // Prometheus监控地址
	PrometheusAuthType uint      `json:"prometheus_auth_type" form:"prometheus_auth_type" gorm:"comment: 认证类型"`                  // Prometheus认证类型
	PrometheusUser     string    `gorm:"type:longText" json:"prometheus_user" form:"prometheus_user" gorm:"comment:用户名"`         // Prometheus用户名
	PrometheusPwd      string    `gorm:"type:longText" json:"prometheus_pwd" form:"prometheus_pwd" gorm:"comment:密码"`            // Prometheus密码
	Users              []User    `json:"users" gorm:"foreignKey:ClusterId;"`                                                     // 关联的用户列表
	CreatedBy          uint      `gorm:"column:created_by;comment:创建者"`                                                          // 创建者ID
	UpdatedBy          uint      `gorm:"column:updated_by;comment:更新者"`                                                          // 更新者ID
	DeletedBy          uint      `gorm:"column:deleted_by;comment:删除者"`                                                          // 删除者ID
}

// RoleData 集群角色数据结构
type RoleData struct {
	ClusterId          uint `json:"cluster_id"` // 集群ID
	rbacV1.ClusterRole      // Kubernetes集群角色
}

// GetClusterById 通过ID获取集群信息的请求结构
type GetClusterById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

// TableName 指定K8sCluster结构体对应的数据库表名
func (K8sCluster) TableName() string {
	return "k8s_clusters"
}
