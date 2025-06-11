package kubernetes

import (
	"KubeGale/utils"

	v1 "k8s.io/api/authorization/v1"
)

// 验证规则变量定义
var (
	// ProxyVerify 代理验证规则
	ProxyVerify = utils.Rules{"Path": {utils.NotEmpty()}}
	// TerminalVerify 终端验证规则
	TerminalVerify = utils.Rules{"Name": {utils.NotEmpty()}, "PodName": {utils.NotEmpty()}, "Namespace": {utils.NotEmpty()}}
	// RoleTypeVerify 角色类型验证规则
	RoleTypeVerify = utils.Rules{"RoleType": {utils.NotEmpty()}}
	// ApiGroupsVerify API组验证规则
	ApiGroupsVerify = utils.Rules{"ApiType": {utils.NotEmpty()}}
	// RoleVerify 角色验证规则
	RoleVerify = utils.Rules{"Rules": {utils.NotEmpty()}, "Metadata": {utils.NotEmpty()}}
	// PodVerify Pod验证规则
	PodVerify = utils.Rules{"ClusterId": {utils.NotEmpty()}}
)

// Kubeconfig Kubernetes配置文件结构体
// 用于解析和生成kubeconfig文件
type Kubeconfig struct {
	APIVersion     string         `yaml:"apiVersion" json:"apiVersion"`           // API版本
	Kind           string         `yaml:"kind" json:"kind"`                       // 资源类型
	Clusters       []ClusterEntry `yaml:"clusters" json:"clusters"`               // 集群配置列表
	Contexts       []ContextEntry `yaml:"contexts" json:"contexts"`               // 上下文配置列表
	CurrentContext string         `yaml:"current-context" json:"current-context"` // 当前上下文
	Preferences    struct{}       `yaml:"preferences" json:"preferences"`         // 偏好设置
	Users          []UserEntry    `yaml:"users" json:"users"`                     // 用户配置列表
}

// KubeCluster 集群配置结构体
type KubeCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data" json:"certificate-authority-data"` // 证书颁发机构数据
	Server                   string `yaml:"server" json:"server"`                                         // 服务器地址
}

// ClusterEntry 集群条目结构体
type ClusterEntry struct {
	Name    string      `yaml:"name" json:"name"`       // 集群名称
	Cluster KubeCluster `yaml:"cluster" json:"cluster"` // 集群配置
}

// Context 上下文结构体
type Context struct {
	Cluster string `yaml:"cluster" json:"cluster"` // 集群名称
	User    string `yaml:"user" json:"user"`       // 用户名称
}

// ContextEntry 上下文条目结构体
type ContextEntry struct {
	Name    string  `yaml:"name" json:"name"`       // 上下文名称
	Context Context `yaml:"context" json:"context"` // 上下文配置
}

// KubeUser 用户配置结构体
type KubeUser struct {
	ClientCertificateData string `yaml:"client-certificate-data" json:"client-certificate-data"` // 客户端证书数据
	ClientKeyData         string `yaml:"client-key-data" json:"client-key-data"`                 // 客户端密钥数据
}

// UserEntry 用户条目结构体
type UserEntry struct {
	Name string   `yaml:"name" json:"name"` // 用户名称
	User KubeUser `yaml:"user" json:"user"` // 用户配置
}

// PermissionCheckResult 权限检查结果结构体
type PermissionCheckResult struct {
	Resource v1.ResourceAttributes // 资源属性
	Allowed  bool                  // 是否允许访问
}

// ApiGroupOption API组选项结构体
type ApiGroupOption struct {
	Group     string              `json:"group"`     // API组名称
	Resources []ApiResourceOption `json:"resources"` // 资源选项列表
}

// ApiResourceOption API资源选项结构体
type ApiResourceOption struct {
	Resource string   `json:"resource"` // 资源名称
	Verbs    []string `json:"verbs"`    // 允许的操作列表
}
