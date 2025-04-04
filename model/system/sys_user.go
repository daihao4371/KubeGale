package system

import (
	"KubeGale/global"
)

type User struct {
	global.Model         // 软删除字段，自动管理
	Username     string  `json:"username" gorm:"type:varchar(100);not null;comment:用户登录名"`                            // 用户登录名，唯一且非空
	Password     string  `json:"password" gorm:"type:varchar(255);not null;comment:用户登录密码"`                           // 用户登录密码，非空
	RealName     string  `json:"realName" gorm:"type:varchar(100);comment:用户真实姓名"`                                    // 用户真实姓名
	Desc         string  `json:"desc" gorm:"type:text;comment:用户描述"`                                                  // 用户描述，支持较长文本
	Mobile       string  `json:"mobile" gorm:"type:varchar(20);comment:手机号"`                                          // 手机号，唯一索引
	FeiShuUserId string  `json:"feiShuUserId" gorm:"type:varchar(50);comment:飞书用户ID"`                                 // 飞书用户ID
	AccountType  int     `json:"accountType" gorm:"default:1;comment:账号类型 1普通用户 2服务账号" binding:"omitempty,oneof=1 2"` // 账号类型，默认为普通用户
	HomePath     string  `json:"homePath" gorm:"type:varchar(255);comment:登录后的默认首页"`                                  // 登录后的默认首页
	Enable       int     `json:"enable" gorm:"default:1;comment:用户状态 1正常 2冻结" binding:"omitempty,oneof=1 2"`          // 用户状态，默认为正常
	Roles        []*Role `json:"roles" gorm:"many2many:user_roles;comment:关联角色"`                                      // 多对多关联角色
	Menus        []*Menu `json:"menus" gorm:"many2many:user_menus;comment:关联菜单"`                                      // 多对多关联菜单
	Apis         []*Api  `json:"apis" gorm:"many2many:user_apis;comment:关联接口"`                                        // 多对多关联接口
}

type TokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ChangePasswordRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type WriteOffRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	UserId       int    `json:"user_id"`
	RealName     string `json:"real_name"`
	Desc         string `json:"desc"`
	Mobile       string `json:"mobile"`
	FeiShuUserId string `json:"fei_shu_user_id"`
	AccountType  int    `json:"account_type"`
	HomePath     string `json:"home_path"`
	Enable       int    `json:"enable"`
}

type DeleteUserRequest struct {
	UserId int `json:"user_id"`
}

func (User) TableName() string {
	return "sys_users"
}
