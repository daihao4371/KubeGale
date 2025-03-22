package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// 使用全局常量
var (
	NewError                 = global.NewError
	NewErrorWithMsg          = global.NewErrorWithMsg
	ERROR_USER_ID_INVALID    = global.ERROR_USER_ID_INVALID
	ERROR_PARAM_INVALID      = global.ERROR_PARAM_INVALID
	ERROR_SAME_PASSWORD      = global.ERROR_SAME_PASSWORD
	ERROR_USER_NOT_EXIST     = global.ERROR_USER_NOT_EXIST
	ERROR_OLD_PASSWORD_WRONG = global.ERROR_OLD_PASSWORD_WRONG
	ERROR_USER_ALREADY_EXIST = global.ERROR_USER_ALREADY_EXIST
	ERROR_MOBILE_INVALID     = global.ERROR_MOBILE_INVALID
	ERROR_MOBILE_USED        = global.ERROR_MOBILE_USED
	ERROR_USER_DISABLED      = global.ERROR_USER_DISABLED
	ERROR_PASSWORD_WRONG     = global.ERROR_PASSWORD_WRONG
)

type UserService struct{}

var UserServiceApp = new(UserService)

// SignUp 用户注册
func (us *UserService) SignUp(user *system.User) error {
	// 检查用户名是否已存在
	var count int64
	if err := global.KUBEGALE_DB.Model(&system.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return fmt.Errorf("检查用户名失败: %w", err)
	}
	if count > 0 {
		return NewError(ERROR_USER_ALREADY_EXIST)
	}

	// 验证手机号格式（如果提供）
	if user.Mobile != "" {
		// 简单的手机号验证，可以根据需求调整
		if len(user.Mobile) != 11 {
			return NewError(ERROR_MOBILE_INVALID)
		}

		// 检查手机号是否已被使用
		var mobileCount int64
		if err := global.KUBEGALE_DB.Model(&system.User{}).Where("mobile = ?", user.Mobile).Count(&mobileCount).Error; err != nil {
			return fmt.Errorf("检查手机号失败: %w", err)
		}
		if mobileCount > 0 {
			return NewError(ERROR_MOBILE_USED)
		}
	}

	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	user.Password = string(hash)

	// 设置默认值
	if user.AccountType == 0 {
		user.AccountType = 1 // 默认为普通用户
	}
	if user.Enable == 0 {
		user.Enable = 1 // 默认启用
	}

	// 使用事务创建用户
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}

		// 分配默认角色
		var defaultRole system.Role
		if err := tx.Where("is_default = ?", 1).First(&defaultRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果没有默认角色，可以跳过或创建一个基本角色
				return nil
			}
			return fmt.Errorf("查询默认角色失败: %w", err)
		}

		// 创建用户-角色关联
		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", user.ID, defaultRole.ID).Error; err != nil {
			return fmt.Errorf("分配默认角色失败: %w", err)
		}

		return nil
	})
}

// Login 用户登录
func (us *UserService) Login(user *system.User) (*system.User, error) {
	// 参数验证
	if user.Username == "" {
		return nil, NewErrorWithMsg(ERROR_PARAM_INVALID, "用户名不能为空")
	}
	if user.Password == "" {
		return nil, NewErrorWithMsg(ERROR_PARAM_INVALID, "密码不能为空")
	}

	// 查询用户
	var u system.User
	if err := global.KUBEGALE_DB.Where("username = ? AND is_deleted = 0", user.Username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewError(ERROR_USER_NOT_EXIST)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户状态
	if u.Enable != 1 {
		return nil, NewError(ERROR_USER_DISABLED)
	}

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return nil, NewError(ERROR_PASSWORD_WRONG)
	}

	// 更新最后登录时间
	if err := global.KUBEGALE_DB.Model(&u).Update("last_login_time", time.Now().Unix()).Error; err != nil {
		global.KUBEGALE_LOG.Warn("更新用户最后登录时间失败", zap.Error(err))
	}

	// 预加载用户关联的角色、菜单和API权限
	if err := global.KUBEGALE_DB.Preload("Roles").Preload("Menus").Preload("Apis").Where("id = ?", u.ID).First(&u).Error; err != nil {
		global.KUBEGALE_LOG.Warn("加载用户权限信息失败", zap.Error(err))
	}

	return &u, nil
}

// 获取用户信息
func (us *UserService) GetProfile(uid int) (*system.User, error) {
	if uid <= 0 {
		return nil, NewError(ERROR_USER_ID_INVALID)
	}

	var user system.User
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = 0", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewError(ERROR_USER_NOT_EXIST)
		}
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 预加载用户关联的角色、菜单和API权限
	if err := global.KUBEGALE_DB.Preload("Roles").Preload("Menus").Preload("Apis").Where("id = ?", uid).First(&user).Error; err != nil {
		global.KUBEGALE_LOG.Warn("加载用户权限信息失败", zap.Error(err))
	}

	// 出于安全考虑，不返回密码
	user.Password = ""

	return &user, nil
}

// GetPermCode 获取用户权限
func (us *UserService) GetPermCode(uid int) ([]string, error) {
	if uid <= 0 {
		return nil, NewError(ERROR_USER_ID_INVALID)
	}

	// 查询用户是否存在
	var count int64
	if err := global.KUBEGALE_DB.Model(&system.User{}).Where("id = ? AND is_deleted = 0", uid).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if count == 0 {
		return nil, NewError(ERROR_USER_NOT_EXIST)
	}

	var permCodes []string

	// 1. 获取用户直接关联的API权限
	var userApis []system.Api
	if err := global.KUBEGALE_DB.Table("sys_apis").
		Joins("JOIN user_apis ON user_apis.api_id = sys_apis.id").
		Where("user_apis.user_id = ? AND sys_apis.is_deleted = 0", uid).
		Find(&userApis).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取用户API权限失败", zap.Error(err))
		return nil, fmt.Errorf("获取用户API权限失败: %w", err)
	}

	// 2. 获取用户通过角色关联的API权限
	var roleApis []system.Api
	if err := global.KUBEGALE_DB.Table("sys_apis").
		Joins("JOIN role_apis ON role_apis.api_id = sys_apis.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_apis.role_id").
		Where("user_roles.user_id = ? AND sys_apis.is_deleted = 0", uid).
		Find(&roleApis).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取角色API权限失败", zap.Error(err))
		return nil, fmt.Errorf("获取角色API权限失败: %w", err)
	}

	// 合并权限并去重
	apiMap := make(map[string]bool)

	// 添加用户直接关联的API权限
	for _, api := range userApis {
		permCode := fmt.Sprintf("%s:%s", api.Path, getMethodName(api.Method))
		apiMap[permCode] = true
	}

	// 添加用户通过角色关联的API权限
	for _, api := range roleApis {
		permCode := fmt.Sprintf("%s:%s", api.Path, getMethodName(api.Method))
		apiMap[permCode] = true
	}

	// 将map转换为切片
	for code := range apiMap {
		permCodes = append(permCodes, code)
	}

	return permCodes, nil
}

// getMethodName 根据方法ID获取方法名称
func getMethodName(methodID int) string {
	switch methodID {
	case 1:
		return "GET"
	case 2:
		return "POST"
	case 3:
		return "PUT"
	case 4:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

// GetUserList 获取用户列表
func (us *UserService) GetUserList(page, pageSize int, keyword string) (users []*system.User, total int64, err error) {
	// 创建查询构建器
	query := global.KUBEGALE_DB.Model(&system.User{}).Where("is_deleted = ?", 0)

	// 如果有关键字，添加搜索条件
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR mobile LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总记录数
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %w", err)
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// 排序
	query = query.Order("id DESC")

	// 预加载关联数据
	query = query.Preload("Roles")

	// 执行查询
	var userList []*system.User
	if err = query.Find(&userList).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 出于安全考虑，清除密码字段
	for _, user := range userList {
		user.Password = ""
	}

	return userList, total, nil
}

// ChangePassword 修改密码
func (us *UserService) ChangePassword(uid int, oldPassword string, newPassword string) error {
	// 参数验证
	if uid <= 0 {
		return NewError(ERROR_USER_ID_INVALID)
	}
	if oldPassword == "" {
		return NewErrorWithMsg(ERROR_PARAM_INVALID, "旧密码不能为空")
	}
	if newPassword == "" {
		return NewErrorWithMsg(ERROR_PARAM_INVALID, "新密码不能为空")
	}
	if oldPassword == newPassword {
		return NewError(ERROR_SAME_PASSWORD)
	}

	// 查询用户
	var user system.User
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = 0", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewError(ERROR_USER_NOT_EXIST)
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证旧密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return NewError(ERROR_OLD_PASSWORD_WRONG)
	}

	// 生成新密码的哈希值
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := global.KUBEGALE_DB.Model(&user).Update("password", string(hash)).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 记录密码修改日志
	global.KUBEGALE_LOG.Info("用户密码已修改", zap.Int("user_id", uid))

	return nil
}

// UpdateProfile 修改用户信息
func (us *UserService) UpdateProfile(uid int, req *system.User) error {
	// 参数验证
	if uid <= 0 {
		return NewError(ERROR_USER_ID_INVALID)
	}

	// 验证手机号格式（如果提供）
	if req.Mobile != "" {
		// 简单的手机号验证
		if len(req.Mobile) != 11 {
			return NewError(ERROR_MOBILE_INVALID)
		}

		// 检查手机号是否已被其他用户使用
		var mobileCount int64
		if err := global.KUBEGALE_DB.Model(&system.User{}).Where("mobile = ? AND id != ? AND is_deleted = 0", req.Mobile, uid).Count(&mobileCount).Error; err != nil {
			return fmt.Errorf("检查手机号失败: %w", err)
		}
		if mobileCount > 0 {
			return NewError(ERROR_MOBILE_USED)
		}
	}

	// 查询用户是否存在
	var user system.User
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = 0", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewError(ERROR_USER_NOT_EXIST)
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 准备更新的字段
	updates := map[string]interface{}{
		"real_name":       req.RealName,
		"desc":            req.Desc,
		"mobile":          req.Mobile,
		"fei_shu_user_id": req.FeiShuUserId,
		"account_type":    req.AccountType,
		"home_path":       req.HomePath,
		"enable":          req.Enable,
		"updated_at":      time.Now(),
	}

	// 更新用户信息
	if err := global.KUBEGALE_DB.Model(&user).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新用户信息失败: %w", err)
	}

	// 记录日志
	global.KUBEGALE_LOG.Info("用户信息已更新", zap.Int("user_id", uid))

	return nil
}

// WriteOff 注销账号
func (us *UserService) WriteOff(username string, password string) error {
	// 参数验证
	if username == "" {
		return NewErrorWithMsg(ERROR_PARAM_INVALID, "用户名不能为空")
	}
	if password == "" {
		return NewErrorWithMsg(ERROR_PARAM_INVALID, "密码不能为空")
	}

	// 验证用户是否存在
	var user system.User
	if err := global.KUBEGALE_DB.Where("username = ? AND is_deleted = 0", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewError(ERROR_USER_NOT_EXIST)
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return NewError(ERROR_PASSWORD_WRONG)
	}

	// 使用事务进行注销操作
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 软删除用户记录
		if err := tx.Model(&user).Update("is_deleted", 1).Error; err != nil {
			return fmt.Errorf("注销用户失败: %w", err)
		}

		// 记录注销时间
		if err := tx.Model(&user).Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("更新注销时间失败: %w", err)
		}

		// 可选：清除用户关联的角色、菜单和API权限
		// 如果需要保留历史记录，可以不执行这些操作
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", user.ID).Error; err != nil {
			return fmt.Errorf("清除用户角色关联失败: %w", err)
		}

		if err := tx.Exec("DELETE FROM user_menus WHERE user_id = ?", user.ID).Error; err != nil {
			return fmt.Errorf("清除用户菜单关联失败: %w", err)
		}

		if err := tx.Exec("DELETE FROM user_apis WHERE user_id = ?", user.ID).Error; err != nil {
			return fmt.Errorf("清除用户API关联失败: %w", err)
		}

		// 记录日志
		global.KUBEGALE_LOG.Info("用户已注销", zap.String("username", username), zap.Int("user_id", user.ID))

		return nil
	})
}

// DeleteUser 删除用户
func (us *UserService) DeleteUser(uid int) error {
	// 参数验证
	if uid <= 0 {
		return NewError(ERROR_USER_ID_INVALID)
	}

	// 查询用户是否存在
	var user system.User
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = 0", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewError(ERROR_USER_NOT_EXIST)
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 使用事务进行删除操作
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 软删除用户记录
		if err := tx.Model(&user).Update("is_deleted", 1).Error; err != nil {
			return fmt.Errorf("删除用户失败: %w", err)
		}

		// 记录删除时间
		if err := tx.Model(&user).Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("更新删除时间失败: %w", err)
		}

		// 清除用户关联的角色、菜单和API权限
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", uid).Error; err != nil {
			return fmt.Errorf("清除用户角色关联失败: %w", err)
		}

		if err := tx.Exec("DELETE FROM user_menus WHERE user_id = ?", uid).Error; err != nil {
			return fmt.Errorf("清除用户菜单关联失败: %w", err)
		}

		if err := tx.Exec("DELETE FROM user_apis WHERE user_id = ?", uid).Error; err != nil {
			return fmt.Errorf("清除用户API关联失败: %w", err)
		}

		// 记录日志
		global.KUBEGALE_LOG.Info("用户已删除", zap.Int("user_id", uid))

		return nil
	})
}
