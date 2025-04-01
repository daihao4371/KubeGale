package system

import (
	"KubeGale/common/casbin"
	"KubeGale/global"
	"KubeGale/model/system"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleService struct{}

var RoleServiceApp = new(RoleService)

// CreateRole 创建角色
func (r *RoleService) CreateRole(role *system.Role, apiIds []int) error {
	// 参数验证
	if role == nil {
		global.KUBEGALE_LOG.Warn("角色对象不能为空")
		return errors.New("角色对象不能为空")
	}

	if role.Name == "" {
		global.KUBEGALE_LOG.Warn("角色名称不能为空")
		return errors.New("角色名称不能为空")
	}

	// 验证角色类型
	if role.RoleType != 1 && role.RoleType != 2 {
		global.KUBEGALE_LOG.Warn("无效的角色类型", zap.Int("roleType", role.RoleType))
		return errors.New("无效的角色类型，必须为1(系统角色)或2(自定义角色)")
	}

	var roleId int
	err := global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		var count int64

		// 检查角色名是否已存在
		if err := tx.Model(&system.Role{}).Where("name = ? AND is_deleted = ?", role.Name, 0).Count(&count).Error; err != nil {
			global.KUBEGALE_LOG.Error("检查角色名称失败", zap.Error(err), zap.String("name", role.Name))
			return fmt.Errorf("检查角色名称失败: %w", err)
		}
		if count > 0 {
			global.KUBEGALE_LOG.Warn("角色名称已存在", zap.String("name", role.Name))
			return errors.New("角色名称已存在")
		}

		// 如果是默认角色，检查是否已有其他默认角色
		if role.IsDefault == 1 {
			var defaultCount int64
			if err := tx.Model(&system.Role{}).Where("is_default = ? AND is_deleted = ?", 1, 0).Count(&defaultCount).Error; err != nil {
				global.KUBEGALE_LOG.Error("检查默认角色失败", zap.Error(err))
				return fmt.Errorf("检查默认角色失败: %w", err)
			}
			if defaultCount > 0 {
				global.KUBEGALE_LOG.Warn("已存在默认角色，一个系统中只能有一个默认角色")
				return errors.New("已存在默认角色，一个系统中只能有一个默认角色")
			}
		}

		// 设置创建时间和更新时间
		now := time.Now().Unix()
		role.CreateTime = now
		role.UpdateTime = now
		role.IsDeleted = 0

		// 创建角色并返回ID
		result := tx.Create(role)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("创建角色失败", zap.Error(result.Error), zap.String("name", role.Name))
			return fmt.Errorf("创建角色失败: %w", result.Error)
		}

		roleId = role.ID
		global.KUBEGALE_LOG.Info("角色创建成功", zap.Int("id", roleId), zap.String("name", role.Name))

		return nil
	})

	if err != nil {
		return err
	}

	// 分配权限
	if len(apiIds) > 0 {
		global.KUBEGALE_LOG.Info("开始为角色分配API权限", zap.Int("roleId", roleId), zap.Ints("apiIds", apiIds))

		// 使用权限服务分配API权限
		if err := AuthorityServiceApp.AssignRole(roleId, apiIds); err != nil {
			global.KUBEGALE_LOG.Error("分配权限失败", zap.Error(err), zap.Int("roleId", roleId))
			return fmt.Errorf("分配权限失败: %w", err)
		}

		global.KUBEGALE_LOG.Info("角色权限分配成功", zap.Int("roleId", roleId))
	}

	return nil
}

// GetRoleById 根据ID获取角色
func (r *RoleService) GetRoleById(id int) (*system.Role, error) {
	if id <= 0 {
		global.KUBEGALE_LOG.Warn("无效的角色ID", zap.Int("id", id))
		return nil, errors.New("无效的角色ID")
	}

	var role system.Role
	// 预加载角色关联的API权限
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", id, 0).
		Preload("Apis").First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Warn("角色不存在", zap.Int("id", id))
			return nil, errors.New("角色不存在")
		}
		global.KUBEGALE_LOG.Error("查询角色失败", zap.Error(err), zap.Int("id", id))
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	global.KUBEGALE_LOG.Info("获取角色成功", zap.Int("id", id), zap.String("name", role.Name))
	return &role, nil
}

// UpdateRole 更新角色信息
func (r *RoleService) UpdateRole(role *system.Role, apiIds []int) error {
	if role == nil {
		global.KUBEGALE_LOG.Warn("角色对象不能为空")
		return errors.New("角色对象不能为空")
	}
	if role.ID <= 0 {
		global.KUBEGALE_LOG.Warn("无效的角色ID", zap.Int("id", role.ID))
		return errors.New("无效的角色ID")
	}
	if role.Name == "" {
		global.KUBEGALE_LOG.Warn("角色名称不能为空")
		return errors.New("角色名称不能为空")
	}

	// 获取原角色信息
	var oldRole system.Role
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", role.ID, 0).First(&oldRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Warn("角色不存在", zap.Int("id", role.ID))
			return errors.New("角色不存在")
		}
		global.KUBEGALE_LOG.Error("获取原角色信息失败", zap.Error(err), zap.Int("id", role.ID))
		return fmt.Errorf("获取原角色信息失败: %w", err)
	}

	// 检查角色名是否已被其他角色使用
	var count int64
	if err := global.KUBEGALE_DB.Model(&system.Role{}).
		Where("name = ? AND id != ?", role.Name, role.ID).
		Count(&count).Error; err != nil {
		global.KUBEGALE_LOG.Error("检查角色名称失败", zap.Error(err), zap.String("name", role.Name))
		return fmt.Errorf("检查角色名称失败: %w", err)
	}
	if count > 0 {
		global.KUBEGALE_LOG.Warn("角色名称已被使用", zap.String("name", role.Name))
		return errors.New("角色名称已被使用")
	}

	// 如果是默认角色，检查是否已有其他默认角色
	if role.IsDefault == 1 && oldRole.IsDefault != 1 {
		var defaultCount int64
		if err := global.KUBEGALE_DB.Model(&system.Role{}).
			Where("is_default = ? AND id != ? ", 1, role.ID).
			Count(&defaultCount).Error; err != nil {
			global.KUBEGALE_LOG.Error("检查默认角色失败", zap.Error(err))
			return fmt.Errorf("检查默认角色失败: %w", err)
		}
		if defaultCount > 0 {
			global.KUBEGALE_LOG.Warn("已存在默认角色，一个系统中只能有一个默认角色")
			return errors.New("已存在默认角色，一个系统中只能有一个默认角色")
		}
	}

	updates := map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
		"role_type":   role.RoleType,
		"is_default":  role.IsDefault,
		"update_time": time.Now().Unix(),
	}

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	err = global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 更新角色基本信息
		result := tx.Model(&system.Role{}).
			Where("id = ? AND is_deleted = ?", role.ID, 0).
			Updates(updates)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("更新角色失败", zap.Error(result.Error), zap.Int("id", role.ID))
			return fmt.Errorf("更新角色失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			global.KUBEGALE_LOG.Warn("角色不存在或已被删除", zap.Int("id", role.ID))
			return errors.New("角色不存在或已被删除")
		}

		// 如果角色名发生变化，需要更新Casbin中的权限
		if oldRole.Name != role.Name {
			global.KUBEGALE_LOG.Info("角色名称已变更，更新Casbin策略",
				zap.String("oldName", oldRole.Name),
				zap.String("newName", role.Name))

			policies, err := enforcer.GetFilteredPolicy(0, oldRole.Name)
			if err != nil {
				global.KUBEGALE_LOG.Error("获取原角色权限失败", zap.Error(err), zap.String("role", oldRole.Name))
				return fmt.Errorf("获取原角色权限失败: %w", err)
			}

			if _, err := enforcer.DeleteRole(oldRole.Name); err != nil {
				global.KUBEGALE_LOG.Error("删除原角色权限失败", zap.Error(err), zap.String("role", oldRole.Name))
				return fmt.Errorf("删除原角色权限失败: %w", err)
			}

			for _, policy := range policies {
				policy[0] = role.Name
				if _, err := enforcer.AddPolicy(policy); err != nil {
					global.KUBEGALE_LOG.Error("添加新角色权限失败", zap.Error(err), zap.Strings("policy", policy))
					return fmt.Errorf("添加新角色权限失败: %w", err)
				}
			}

			if err := enforcer.SavePolicy(); err != nil {
				global.KUBEGALE_LOG.Error("保存权限策略失败", zap.Error(err))
				return fmt.Errorf("保存权限策略失败: %w", err)
			}
		}

		global.KUBEGALE_LOG.Info("角色更新成功", zap.Int("id", role.ID), zap.String("name", role.Name))
		return nil
	})

	return err
}

// DeleteRole 删除角色
func (r *RoleService) DeleteRole(id int) error {
	if id <= 0 {
		global.KUBEGALE_LOG.Warn("无效的角色ID", zap.Int("id", id))
		return errors.New("无效的角色ID")
	}

	// 检查是否为默认角色
	var role system.Role
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", id, 0).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Warn("角色不存在", zap.Int("id", id))
			return errors.New("角色不存在")
		}
		global.KUBEGALE_LOG.Error("查询角色失败", zap.Error(err), zap.Int("id", id))
		return fmt.Errorf("查询角色失败: %w", err)
	}

	if role.IsDefault == 1 {
		global.KUBEGALE_LOG.Warn("默认角色不能删除", zap.Int("id", id), zap.String("name", role.Name))
		return errors.New("默认角色不能删除")
	}

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	// 使用事务保证数据一致性
	err = global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 软删除角色
		updates := map[string]interface{}{
			"is_deleted":  1,
			"update_time": time.Now().Unix(),
		}

		result := tx.Model(&system.Role{}).Where("id = ? AND is_deleted = ?", id, 0).Updates(updates)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("删除角色失败", zap.Error(result.Error), zap.Int("id", id))
			return fmt.Errorf("删除角色失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			global.KUBEGALE_LOG.Warn("角色不存在或已被删除", zap.Int("id", id))
			return errors.New("角色不存在或已被删除")
		}

		// 获取所有策略
		allPolicies, _ := enforcer.GetPolicy()

		// 删除与角色关联的API权限策略
		for _, policy := range allPolicies {
			if len(policy) >= 3 && strings.HasPrefix(policy[0], fmt.Sprintf("role_%d_", id)) {
				_, err := enforcer.RemovePolicy(policy)
				if err != nil {
					global.KUBEGALE_LOG.Error("删除Casbin策略失败",
						zap.Error(err),
						zap.Strings("policy", policy))
					return fmt.Errorf("删除Casbin策略失败: %w", err)
				}
			}
		}

		// 保存Casbin策略更改
		if err := enforcer.SavePolicy(); err != nil {
			global.KUBEGALE_LOG.Error("保存Casbin策略失败", zap.Error(err))
			return fmt.Errorf("保存Casbin策略失败: %w", err)
		}

		global.KUBEGALE_LOG.Info("角色删除成功", zap.Int("id", id), zap.String("name", role.Name))
		return nil
	})

	return err
}

// ListRoles 获取角色列表
func (r *RoleService) ListRoles(page, pageSize int) ([]*system.Role, int, error) {
	// 添加默认分页参数
	if page <= 0 {
		page = 1
		global.KUBEGALE_LOG.Warn("使用默认页码", zap.Int("page", page))
	}
	if pageSize <= 0 {
		pageSize = 10
		global.KUBEGALE_LOG.Warn("使用默认每页数量", zap.Int("pageSize", pageSize))
	}

	var roles []*system.Role
	var total int64

	// 构建查询，只查询未删除的角色
	db := global.KUBEGALE_DB.Model(&system.Role{}).Where("is_deleted = ?", 0)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取角色总数失败", zap.Error(err))
		return nil, 0, fmt.Errorf("获取角色总数失败: %w", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&roles).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取角色列表失败", zap.Error(err))
		return nil, 0, fmt.Errorf("获取角色列表失败: %w", err)
	}

	global.KUBEGALE_LOG.Info("获取角色列表成功",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("count", len(roles)),
		zap.Int64("total", total))

	return roles, int(total), nil
}

// GetRole 获取角色详细信息(包含权限)
func (r *RoleService) GetRole(roleId int) (*system.Role, error) {
	if roleId <= 0 {
		return nil, errors.New("无效的角色ID")
	}

	var role system.Role
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", roleId, 0).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}

	// 使用Preload直接加载角色关联的API权限
	if err := global.KUBEGALE_DB.Where("id = ?", roleId).Preload("Apis", "is_deleted = 0").First(&role).Error; err != nil {
		global.KUBEGALE_LOG.Error("加载角色API权限失败", zap.Error(err), zap.Int("roleId", roleId))
		return nil, fmt.Errorf("加载角色API权限失败: %v", err)
	}

	global.KUBEGALE_LOG.Info("获取角色详情成功",
		zap.Int("roleId", roleId),
		zap.String("name", role.Name),
		zap.Int("apiCount", len(role.Apis)))

	return &role, nil
}

// GetUserRole 获取用户的角色信息
func (r *RoleService) GetUserRole(userId int) (*system.Role, error) {
	if userId <= 0 {
		return nil, errors.New("无效的用户ID")
	}

	// 先从数据库中获取用户的角色
	var user system.User
	if err := global.KUBEGALE_DB.Preload("Roles", "is_deleted = ?", 0).Where("id = ? AND is_deleted = ?", userId, 0).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	// 如果用户没有角色,返回nil
	if user.Roles == nil || len(user.Roles) == 0 {
		return nil, nil
	}

	// 获取第一个角色
	role := user.Roles[0]

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return nil, fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	// 获取用户的所有权限(包括直接权限和角色权限)
	allPolicies := make([][]string, 0)

	// 获取用户直接权限
	userStr := fmt.Sprintf("%d", userId)
	userPolicies, _ := enforcer.GetFilteredPolicy(0, userStr)
	if len(userPolicies) > 0 {
		allPolicies = append(allPolicies, userPolicies...)
	}

	// 获取角色权限
	if role.ID > 0 {
		roleStr := fmt.Sprintf("role_%d", role.ID)
		rolePolicies, _ := enforcer.GetFilteredPolicy(0, roleStr)
		if len(rolePolicies) > 0 {
			allPolicies = append(allPolicies, rolePolicies...)
		}
	}

	// 解析权限策略获取API的ID
	apiIdsMap := make(map[int]struct{})

	for _, policy := range allPolicies {
		if len(policy) < 2 {
			continue
		}
		if id, err := strconv.Atoi(policy[1]); err == nil {
			apiIdsMap[id] = struct{}{}
		}
	}

	// 转换为切片
	apiIds := make([]int, 0, len(apiIdsMap))
	for id := range apiIdsMap {
		apiIds = append(apiIds, id)
	}

	// 查询API详细信息
	var apis []*system.Api

	if len(apiIds) > 0 {
		if err := global.KUBEGALE_DB.Where("id IN ? AND is_deleted = ?", apiIds, 0).Find(&apis).Error; err != nil {
			return nil, fmt.Errorf("查询API失败: %v", err)
		}
		role.Apis = apis
	}

	return role, nil
}
