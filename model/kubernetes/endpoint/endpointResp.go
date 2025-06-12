package endpoint

import (
	"KubeGale/model/common/request"

	discoveryv1 "k8s.io/api/discovery/v1"
)

type EndPointListResponse struct {
	Items *[]discoveryv1.EndpointSlice `json:"items" form:"items"`
	Total int                          `json:"total" form:"total"`
	request.PageInfo
}

type DescribeEndPointResponse struct {
	Items *discoveryv1.EndpointSlice `json:"items" form:"items"`
}
