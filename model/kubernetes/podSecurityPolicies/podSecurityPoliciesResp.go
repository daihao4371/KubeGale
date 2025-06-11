package podSecurityPolicies

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// PodSecurityPoliciesListResponse Pod安全策略列表响应结构体
// 用于返回Pod安全策略的列表信息，包含分页信息
type PodSecurityPoliciesListResponse struct {
	Items *[]v1.PodSecurityContext `json:"items" form:"items"` // Pod安全策略列表
	Total int                      `json:"total" form:"total"` // 总数
	request.PageInfo
}

//type DescribeStorageClassResponse struct {
//	Items *v1.StorageClass `json:"items" form:"items"`
//}
