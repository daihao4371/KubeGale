package request

import (
	"KubeGale/model/common/request"
	"time"
)

type CmdbProjectsSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

// DeleteCmdbProjectsRequest 删除请求结构体
type DeleteCmdbProjectsRequest struct {
	ID uint `json:"id" binding:"required"`
}

// DeleteCmdbProjectsIdsRequest 批量删除请求结构体
type DeleteCmdbProjectsIdsRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}
