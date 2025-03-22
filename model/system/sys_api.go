package system

type Api struct {
	ID          int    `json:"id" gorm:"primaryKey;column:id;comment:主键ID"`
	Name        string `json:"name" gorm:"column:name;type:varchar(50);not null;comment:API名称"`
	Path        string `json:"path" gorm:"column:path;type:varchar(255);not null;comment:API路径"`
	Method      int    `json:"method" gorm:"column:method;type:tinyint(1);not null;comment:HTTP请求方法(1:GET,2:POST,3:PUT,4:DELETE)"`
	Description string `json:"description" gorm:"column:description;type:varchar(500);comment:API描述"`
	Version     string `json:"version" gorm:"column:version;type:varchar(20);default:v1;comment:API版本"`
	Category    int    `json:"category" gorm:"column:category;type:tinyint(1);not null;comment:API分类(1:系统,2:业务)"`
	IsPublic    int    `json:"is_public" gorm:"column:is_public;type:tinyint(1);default:0;comment:是否公开(0:否,1:是)"`
	CreateTime  int64  `json:"create_time" gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime  int64  `json:"update_time" gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	IsDeleted   int    `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;comment:是否删除(0:否,1:是)"`
}

type CreateApiRequest struct {
	Name        string `json:"name" binding:"required"`       // API名称
	Path        string `json:"path" binding:"required"`       // API路径
	Method      int    `json:"method" binding:"required"`     // 请求方法
	Description string `json:"description"`                   // API描述
	Version     string `json:"version"`                       // API版本
	Category    int    `json:"category"`                      // API分类
	IsPublic    int    `json:"is_public" binding:"oneof=0 1"` // 是否公开
}

type GetApiRequest struct {
	ID int `json:"id" binding:"required,gt=0"` // API ID
}

type UpdateApiRequest struct {
	ID          int    `json:"id" binding:"required,gt=0"`    // API ID
	Name        string `json:"name" binding:"required"`       // API名称
	Path        string `json:"path" binding:"required"`       // API路径
	Method      int    `json:"method" binding:"required"`     // 请求方法
	Description string `json:"description"`                   // API描述
	Version     string `json:"version"`                       // API版本
	Category    int    `json:"category"`                      // API分类
	IsPublic    int    `json:"is_public" binding:"oneof=0 1"` // 是否公开
}

type ListApisRequest struct {
	PageNumber int `json:"page_number" binding:"required,gt=0"` // 页码
	PageSize   int `json:"page_size" binding:"required,gt=0"`   // 每页数量
}

func (Api) TableName() string {
	return "sys_apis"
}
