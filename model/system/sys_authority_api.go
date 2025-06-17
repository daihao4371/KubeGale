package system

// SysAuthorityApi 角色-API权限关联表
// 用于描述角色（Authority）与API权限（Api）之间的多对多关系

type SysAuthorityApi struct {
	AuthorityId uint `gorm:"column:authority_id;not null;index" json:"authorityId"` // 角色ID
	ApiId       uint `gorm:"column:api_id;not null;index" json:"apiId"`             // API权限ID
}

func (SysAuthorityApi) TableName() string {
	return "sys_authority_apis"
}
