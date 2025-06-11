package cronjob

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/batch/v1"
)

// CronJobListResponse CronJob列表响应结构体
// 用于返回CronJob的列表信息，包含分页信息
type CronJobListResponse struct {
	Items            *[]v1.CronJob `json:"items" form:"items"` // CronJob列表
	Total            int           `json:"total" form:"total"` // 总数
	request.PageInfo               // 分页信息
}

// DescribeCronJobResponse CronJob详情响应结构体
// 用于返回单个CronJob的详细信息
type DescribeCronJobResponse struct {
	Items *v1.CronJob `json:"items" form:"items"` // CronJob详情
}
