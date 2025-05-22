package cmdb

import (
	"KubeGale/global"
	"time"
)

// CmdbProjects 结构体
type CmdbProjects struct {
	global.KUBEGALE_MODEL
	Name        string     `json:"name" form:"name" gorm:"column:name;comment:项目名称;"`
	Description string     `json:"description" form:"description" gorm:"column:description;comment:项目描述;"`
	Manager     string     `json:"manager" form:"manager" gorm:"column:manager;comment:项目负责人;"`
	CreatedBy   uint       `json:"created_by" form:"created_by" gorm:"column:created_by;comment:创建人;"`
	UpdatedBy   uint       `json:"updated_by" form:"updated_by" gorm:"column:updated_by;comment:更新人;"`
	DeletedBy   uint       `json:"deleted_by" form:"deleted_by" gorm:"column:deleted_by;comment:删除人;"`
	DeletedAt   *time.Time `json:"deleted_at" form:"deleted_at" gorm:"column:deleted_at;comment:删除时间;"`
}

// TableName 表名
func (CmdbProjects) TableName() string {
	return "cmdb_projects"
}
