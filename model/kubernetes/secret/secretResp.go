package secret

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// SecretListResponse 密钥列表响应结构体
// 用于返回密钥的列表信息，包含分页信息
type SecretListResponse struct {
	Items *[]v1.Secret `json:"items" form:"items"` // 密钥列表
	Total int          `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeSecretResponse 密钥详情响应结构体
// 用于返回单个密钥的详细信息
type DescribeSecretResponse struct {
	Items *v1.Secret `json:"items"` // 密钥详情
}
