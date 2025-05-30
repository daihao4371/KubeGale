package system

import (
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"time"

	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/system"
	"KubeGale/utils"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

// @function: Register
// @description: 用户注册
type UserService struct{}

var UserServiceApp = new(UserService)

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.KUBEGALE_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	err = global.KUBEGALE_DB.Create(&u).Error
	return u, err
}

//@function: Login
//@description: 用户登录

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.KUBEGALE_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err = global.KUBEGALE_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
		MenuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}

//@function: ChangePassword
//@description: 修改用户密码

func (userService *UserService) ChangePassword(u *system.SysUser, newPassword string) (userInter *system.SysUser, err error) {
	var user system.SysUser
	if err = global.KUBEGALE_DB.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.KUBEGALE_DB.Save(&user).Error
	return &user, err

}

//@function: GetUserInfoList
//@description: 分页获取数据

func (userService *UserService) GetUserInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

//@function: SetUserAuthority
//@description: 设置一个用户的权限

func (userService *UserService) SetUserAuthority(id uint, authorityId uint) (err error) {

	assignErr := global.KUBEGALE_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}

	var authority system.SysAuthority
	err = global.KUBEGALE_DB.Where("authority_id = ?", authorityId).First(&authority).Error
	if err != nil {
		return err
	}
	var authorityMenu []system.SysAuthorityMenu
	var authorityMenuIDs []string
	err = global.KUBEGALE_DB.Where("sys_authority_authority_id = ?", authorityId).Find(&authorityMenu).Error
	if err != nil {
		return err
	}

	for i := range authorityMenu {
		authorityMenuIDs = append(authorityMenuIDs, authorityMenu[i].MenuId)
	}

	var authorityMenus []system.SysBaseMenu
	err = global.KUBEGALE_DB.Preload("Parameters").Where("id in (?)", authorityMenuIDs).Find(&authorityMenus).Error
	if err != nil {
		return err
	}
	hasMenu := false
	for i := range authorityMenus {
		if authorityMenus[i].Name == authority.DefaultRouter {
			hasMenu = true
			break
		}
	}
	if !hasMenu {
		return errors.New("找不到默认路由,无法切换本角色")
	}

	err = global.KUBEGALE_DB.Model(&system.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
	return err
}

//@function: SetUserAuthorities
//@description: 设置一个用户的权限

func (userService *UserService) SetUserAuthorities(adminAuthorityID, id uint, authorityIds []uint) (err error) {
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		var user system.SysUser
		TxErr := tx.Where("id = ?", id).First(&user).Error
		if TxErr != nil {
			global.KUBEGALE_LOG.Debug(TxErr.Error())
			return errors.New("查询用户数据失败")
		}
		TxErr = tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			e := AuthorityServiceApp.CheckAuthorityIDAuth(adminAuthorityID, v)
			if e != nil {
				return e
			}
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Model(&user).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		// 返回 nil 提交事务
		return nil
	})
}

//@function: DeleteUser
//@description: 删除用户

func (userService *UserService) DeleteUser(id int) (err error) {
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&system.SysUser{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser

func (userService *UserService) SetUserInfo(req system.SysUser) error {
	return global.KUBEGALE_DB.Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"enable":     req.Enable,
		}).Error
}

//@function: SetSelfInfo
//@description: 设置用户信息

func (userService *UserService) SetSelfInfo(req system.SysUser) error {
	return global.KUBEGALE_DB.Model(&system.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error
}

//@function: SetSelfSetting
//@description: 设置用户配置

func (userService *UserService) SetSelfSetting(req *datatypes.JSON, uid uint) error {
	return global.KUBEGALE_DB.Model(&system.SysUser{}).Where("id = ?", uid).Update("origin_setting", req).Error
}

//@function: GetUserInfo
//@description: 获取用户信息

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (user system.SysUser, err error) {
	var reqUser system.SysUser
	err = global.KUBEGALE_DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	MenuServiceApp.UserAuthorityDefaultRouter(&reqUser)
	return reqUser, err
}

//@function: FindUserById
//@description: 通过id获取用户信息

func (userService *UserService) FindUserById(id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = global.KUBEGALE_DB.Where("id = ?", id).First(&u).Error
	return &u, err
}

//@function: FindUserByUuid
//@description: 通过uuid获取用户信息

func (userService *UserService) FindUserByUuid(uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = global.KUBEGALE_DB.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

// @function: ResetPassword
// @description: 修改用户密码
func (userService *UserService) ResetPassword(ID uint) (err error) {
	err = global.KUBEGALE_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error
	return err
}
