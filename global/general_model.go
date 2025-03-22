package global

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Model struct {
	ID        int                   `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"` // 主键ID，自增
	CreatedAt time.Time             `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`   // 创建时间，自动记录
	UpdatedAt time.Time             `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`   // 更新时间，自动记录
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"index;comment:删除时间"`            // 软删除时间，使用普通索引
}
