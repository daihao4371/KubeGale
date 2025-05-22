package request

import (
	"KubeGale/model/common/request"
	"time"
)

type CmdbHostsSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Project        int        `json:"project" form:"project"`
	request.PageInfo
}

// DeleteCmdbHostsRequest 删除请求结构体
type DeleteCmdbHostsRequest struct {
	ID uint `json:"id" binding:"required"`
}

// DeleteCmdbHostsIdsRequest 批量删除请求结构体
type DeleteCmdbHostsIdsRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}
