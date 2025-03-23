package system

import (
	"KubeGale/common/casbin"
	"KubeGale/global"
	"KubeGale/model/system"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type AuthorityService struct{}

var AuthorityServiceApp = new(AuthorityService)

// AssignRole 为角色分配权限
func (p *AuthorityService) AssignRole(roleId int, apiIds []int) error {
	// 参数校验
	if roleId <= 0 {
		global.KUBEGALE_LOG.Warn("角色ID无效", zap.Int("roleId", roleId))
		return errors.New("角色ID无效")
	}

	// 检查角色是否存在
	var role system.Role
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", roleId, 0).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Warn("角色不存在", zap.Int("roleId", roleId))
			return errors.New("角色不存在")
		}
		global.KUBEGALE_LOG.Error("查询角色失败", zap.Error(err), zap.Int("roleId", roleId))
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 使用事务保证数据一致性
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 获取当前角色已有的API权限
		var existingApiIds []int
		if err := tx.Table("role_apis").Where("role_id = ?", roleId).Pluck("api_id", &existingApiIds).Error; err != nil {
			global.KUBEGALE_LOG.Error("获取角色现有权限失败", zap.Error(err), zap.Int("roleId", roleId))
			return fmt.Errorf("获取角色现有权限失败: %w", err)
		}

		// 2. 计算需要添加和删除的API权限
		toAdd, toRemove := calculateDiff(existingApiIds, apiIds)

		// 3. 记录操作日志
		global.KUBEGALE_LOG.Info("角色权限变更",
			zap.Int("roleId", roleId),
			zap.Ints("现有权限", existingApiIds),
			zap.Ints("新权限", apiIds),
			zap.Ints("添加权限", toAdd),
			zap.Ints("移除权限", toRemove))

		// 4. 删除需要移除的权限
		if len(toRemove) > 0 {
			if err := tx.Exec("DELETE FROM role_apis WHERE role_id = ? AND api_id IN ?", roleId, toRemove).Error; err != nil {
				global.KUBEGALE_LOG.Error("移除角色API权限失败", zap.Error(err), zap.Int("roleId", roleId))
				return fmt.Errorf("移除角色API权限失败: %w", err)
			}
		}

		// 5. 添加新的权限
		if len(toAdd) > 0 {
			// 检查API是否存在
			var count int64
			if err := tx.Model(&system.Api{}).Where("id IN ? AND is_deleted = ?", toAdd, 0).Count(&count).Error; err != nil {
				global.KUBEGALE_LOG.Error("检查API是否存在失败", zap.Error(err))
				return fmt.Errorf("检查API是否存在失败: %w", err)
			}
			if int(count) != len(toAdd) {
				global.KUBEGALE_LOG.Warn("部分API不存在", zap.Ints("apiIds", toAdd), zap.Int64("existCount", count))
				return errors.New("部分API不存在")
			}

			// 批量插入新权限
			values := make([]string, 0, len(toAdd))
			args := make([]interface{}, 0, len(toAdd)*2)
			for _, apiId := range toAdd {
				values = append(values, "(?, ?)")
				args = append(args, roleId, apiId)
			}

			sql := fmt.Sprintf("INSERT INTO role_apis (role_id, api_id) VALUES %s", strings.Join(values, ", "))
			if err := tx.Exec(sql, args...).Error; err != nil {
				global.KUBEGALE_LOG.Error("添加角色API权限失败", zap.Error(err), zap.Int("roleId", roleId))
				return fmt.Errorf("添加角色API权限失败: %w", err)
			}
		}

		// 6. 更新Casbin策略
		if err := p.updateCasbinPolicy(tx, roleId, apiIds, toAdd, toRemove); err != nil {
			global.KUBEGALE_LOG.Error("更新Casbin策略失败", zap.Error(err), zap.Int("roleId", roleId))
			return fmt.Errorf("更新Casbin策略失败: %w", err)
		}

		// 7. 更新角色的更新时间
		if err := tx.Model(&system.Role{}).Where("id = ?", roleId).Update("update_time", time.Now().Unix()).Error; err != nil {
			global.KUBEGALE_LOG.Warn("更新角色时间戳失败", zap.Error(err), zap.Int("roleId", roleId))
			// 不返回错误，因为这不是关键操作
		}

		global.KUBEGALE_LOG.Info("角色权限分配成功", zap.Int("roleId", roleId), zap.Ints("apiIds", apiIds))
		return nil
	})
}

// calculateDiff 计算两个切片的差异，返回需要添加和删除的元素
func calculateDiff(existing, new []int) (toAdd, toRemove []int) {
	existingMap := make(map[int]bool)
	for _, id := range existing {
		existingMap[id] = true
	}

	newMap := make(map[int]bool)
	for _, id := range new {
		newMap[id] = true
		if !existingMap[id] {
			toAdd = append(toAdd, id)
		}
	}

	for _, id := range existing {
		if !newMap[id] {
			toRemove = append(toRemove, id)
		}
	}

	return toAdd, toRemove
}

// updateCasbinPolicy 更新Casbin策略，只更新变化的部分
func (p *AuthorityService) updateCasbinPolicy(tx *gorm.DB, roleId int, allApiIds, addedApiIds, removedApiIds []int) error {
	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer() // 使用新包中的函数
	if err != nil {
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	// 获取角色关联的所有用户
	var userIds []int
	if err := tx.Table("user_roles").Where("role_id = ?", roleId).Pluck("user_id", &userIds).Error; err != nil {
		return fmt.Errorf("获取角色关联用户失败: %w", err)
	}

	if len(userIds) == 0 {
		global.KUBEGALE_LOG.Info("角色没有关联用户，跳过Casbin策略更新", zap.Int("roleId", roleId))
		return nil
	}

	// 优化1: 批量获取API信息，减少数据库查询次数
	var allApis []system.Api
	var apiMap = make(map[int]system.Api)

	if len(removedApiIds) > 0 || len(addedApiIds) > 0 {
		// 合并需要查询的API ID
		apiIdsToQuery := make([]int, 0, len(removedApiIds)+len(addedApiIds))
		apiIdsToQuery = append(apiIdsToQuery, removedApiIds...)
		apiIdsToQuery = append(apiIdsToQuery, addedApiIds...)

		// 一次性查询所有需要的API
		if err := tx.Where("id IN ? AND is_deleted = ?", apiIdsToQuery, 0).Find(&allApis).Error; err != nil {
			return fmt.Errorf("获取API详情失败: %w", err)
		}

		// 构建API ID到API对象的映射，方便后续使用
		for _, api := range allApis {
			apiMap[api.ID] = api
		}
	}

	// 优化2: 批量加载Casbin策略，避免频繁加载
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载Casbin策略失败: %w", err)
	}

	// 优化3: 使用批量操作处理策略
	var policiesToRemove [][]string
	var policiesToAdd [][]string

	// 处理需要移除的API权限
	if len(removedApiIds) > 0 {
		// 为每个用户检查并准备移除策略
		for _, userId := range userIds {
			userIdStr := fmt.Sprintf("%d", userId)

			for _, apiId := range removedApiIds {
				api, exists := apiMap[apiId]
				if !exists {
					continue // 跳过不存在的API
				}

				methodName := casbin.GetMethodName(api.Method) // 使用新包中的函数

				// 检查是否有其他角色或直接授权提供相同的权限
				var otherCount int64
				if err := tx.Raw(`
					SELECT COUNT(*) FROM (
						SELECT 1 FROM role_apis ra 
						JOIN user_roles ur ON ra.role_id = ur.role_id 
						JOIN sys_apis a ON ra.api_id = a.id 
						WHERE ur.user_id = ? AND a.path = ? AND a.method = ? AND ra.role_id != ?
						UNION ALL
						SELECT 1 FROM user_apis ua 
						JOIN sys_apis a ON ua.api_id = a.id 
						WHERE ua.user_id = ? AND a.path = ? AND a.method = ?
					) AS combined_perms
				`, userId, api.Path, api.Method, roleId, userId, api.Path, api.Method).Count(&otherCount).Error; err != nil {
					return fmt.Errorf("检查其他权限来源失败: %w", err)
				}

				// 只有当没有其他权限来源时才移除策略
				if otherCount == 0 {
					// 直接添加到待移除策略列表，不检查是否存在
					policiesToRemove = append(policiesToRemove, []string{userIdStr, api.Path, methodName})
				}
			}
		}
	}

	// 处理需要添加的API权限
	if len(addedApiIds) > 0 {
		// 为每个用户添加策略
		for _, userId := range userIds {
			userIdStr := fmt.Sprintf("%d", userId)

			for _, apiId := range addedApiIds {
				api, exists := apiMap[apiId]
				if !exists {
					continue // 跳过不存在的API
				}

				methodName := casbin.GetMethodName(api.Method) // 使用新包中的函数

				// 直接添加到待添加策略列表，让Casbin处理重复问题
				policiesToAdd = append(policiesToAdd, []string{userIdStr, api.Path, methodName})
			}
		}
	}

	// 优化4: 批量移除和添加策略
	if len(policiesToRemove) > 0 {
		_, err := enforcer.RemovePolicies(policiesToRemove)
		if err != nil {
			return fmt.Errorf("批量移除Casbin策略失败: %w", err)
		}
	}

	if len(policiesToAdd) > 0 {
		_, err := enforcer.AddPolicies(policiesToAdd)
		if err != nil {
			return fmt.Errorf("批量添加Casbin策略失败: %w", err)
		}
	}

	// 保存策略到适配器
	if len(policiesToRemove) > 0 || len(policiesToAdd) > 0 {
		if err := enforcer.SavePolicy(); err != nil {
			return fmt.Errorf("保存Casbin策略失败: %w", err)
		}
	}

	global.KUBEGALE_LOG.Info("Casbin策略更新成功",
		zap.Int("roleId", roleId),
		zap.Int("移除策略数", len(policiesToRemove)),
		zap.Int("添加策略数", len(policiesToAdd)))
	return nil
}

// AssignRoleToUser 为用户分配角色和权限
func (p *AuthorityService) AssignRoleToUser(userId int, roleIds []int, apiIds []int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 获取角色信息
		var roles []*system.Role
		if len(roleIds) > 0 {
			if err := tx.Where("id IN ? AND is_deleted = ?", roleIds, 0).Find(&roles).Error; err != nil {
				return fmt.Errorf("获取角色信息失败: %v", err)
			}
		}

		// 获取Casbin强制器
		enforcer, err := casbin.GetEnforcer()
		if err != nil {
			global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
			return fmt.Errorf("获取Casbin强制器失败: %w", err)
		}

		// 先移除用户现有的角色关联和权限
		userStr := fmt.Sprintf("%d", userId)
		_, err = enforcer.RemoveFilteredGroupingPolicy(0, userStr)
		if err != nil {
			global.KUBEGALE_LOG.Error("移除用户现有角色关联失败", zap.Error(err))
			return fmt.Errorf("移除用户现有角色关联失败: %v", err)
		}

		_, err = enforcer.RemoveFilteredPolicy(0, userStr)
		if err != nil {
			global.KUBEGALE_LOG.Error("移除用户现有权限失败", zap.Error(err))
			return fmt.Errorf("移除用户现有权限失败: %v", err)
		}

		// 获取用户信息
		var user system.User
		if err := tx.First(&user, userId).Error; err != nil {
			return fmt.Errorf("获取用户信息失败: %v", err)
		}

		// 更新用户的角色关联
		if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
			return fmt.Errorf("更新用户角色关联失败: %v", err)
		}

		// 更新用户的API关联
		if len(apiIds) > 0 {
			var apis []*system.Api
			if err := tx.Where("id IN ? AND is_deleted = ?", apiIds, 0).Find(&apis).Error; err != nil {
				return fmt.Errorf("获取API信息失败: %v", err)
			}
			if err := tx.Model(&user).Association("Apis").Replace(apis); err != nil {
				return fmt.Errorf("更新用户API关联失败: %v", err)
			}
		}

		// 添加角色关联策略
		if len(roles) > 0 {
			rolePolicies := make([][]string, 0, len(roles))
			for _, role := range roles {
				rolePolicies = append(rolePolicies, []string{userStr, fmt.Sprintf("role_%d", role.ID)})
			}

			// 批量添加角色关联策略
			_, err = enforcer.AddGroupingPolicies(rolePolicies)
			if err != nil {
				global.KUBEGALE_LOG.Error("添加用户角色关联失败", zap.Error(err))
				return fmt.Errorf("添加用户角色关联失败: %v", err)
			}
		}

		// 添加API权限
		if len(apiIds) > 0 {
			// 批量获取API信息
			var apis []system.Api
			if err := tx.Where("id IN ? AND is_deleted = ?", apiIds, 0).Find(&apis).Error; err != nil {
				global.KUBEGALE_LOG.Error("获取API信息失败", zap.Error(err))
				return fmt.Errorf("获取API信息失败: %v", err)
			}

			// 构建API权限策略
			apiPolicies := make([][]string, 0, len(apis))
			for _, api := range apis {
				methodName := casbin.GetMethodName(api.Method)
				apiPolicies = append(apiPolicies, []string{userStr, api.Path, methodName})
			}

			// 批量添加API权限策略
			if len(apiPolicies) > 0 {
				_, err = enforcer.AddPolicies(apiPolicies)
				if err != nil {
					global.KUBEGALE_LOG.Error("添加API权限策略失败", zap.Error(err))
					return fmt.Errorf("添加API权限策略失败: %v", err)
				}
			}
		}

		// 保存策略
		if err := enforcer.SavePolicy(); err != nil {
			global.KUBEGALE_LOG.Error("保存Casbin策略失败", zap.Error(err))
			return fmt.Errorf("保存策略失败: %v", err)
		}

		global.KUBEGALE_LOG.Info("用户角色和权限分配成功",
			zap.Int("userId", userId),
			zap.Ints("roleIds", roleIds),
			zap.Ints("apiIds", apiIds))
		return nil
	})
}

// RemoveUserPermissions 移除用户权限
func (p *AuthorityService) RemoveUserPermissions(userId int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}

	// 不允许删除userId为1的权限
	if userId == 1 {
		return errors.New("不允许删除超级管理员权限")
	}

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	userStr := fmt.Sprintf("%d", userId)

	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 获取用户信息
		var user system.User
		if err := tx.First(&user, userId).Error; err != nil {
			return fmt.Errorf("获取用户信息失败: %v", err)
		}

		// 获取用户当前的API列表
		var apis []*system.Api
		if err := tx.Model(&user).Association("Apis").Find(&apis); err != nil {
			return fmt.Errorf("获取用户API关联失败: %v", err)
		}

		global.KUBEGALE_LOG.Info("准备移除用户API权限",
			zap.Int("userId", userId),
			zap.Int("apiCount", len(apis)))

		// 批量移除API权限策略
		if len(apis) > 0 {
			// 构建需要移除的策略列表
			policiesToRemove := make([][]string, 0, len(apis))

			// HTTP方法映射表
			methodMap := map[int]string{
				1: "GET",
				2: "POST",
				3: "PUT",
				4: "DELETE",
				5: "PATCH",
				6: "OPTIONS",
				7: "HEAD",
			}

			for _, api := range apis {
				method := methodMap[api.Method]
				policiesToRemove = append(policiesToRemove, []string{userStr, api.Path, method})
			}

			// 批量移除策略
			_, err := enforcer.RemovePolicies(policiesToRemove)
			if err != nil {
				global.KUBEGALE_LOG.Error("批量移除用户API权限失败", zap.Error(err))
				return fmt.Errorf("移除用户API权限失败: %v", err)
			}
		}

		// 移除用户的角色关联
		_, err := enforcer.RemoveFilteredGroupingPolicy(0, userStr)
		if err != nil {
			global.KUBEGALE_LOG.Error("移除用户角色关联失败", zap.Error(err))
			return fmt.Errorf("移除用户角色关联失败: %v", err)
		}

		// 清空用户的关联
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return fmt.Errorf("清空用户角色关联失败: %v", err)
		}

		if err := tx.Model(&user).Association("Apis").Clear(); err != nil {
			return fmt.Errorf("清空用户API关联失败: %v", err)
		}

		// 保存策略变更
		if err := enforcer.SavePolicy(); err != nil {
			global.KUBEGALE_LOG.Error("保存Casbin策略失败", zap.Error(err))
			return fmt.Errorf("保存策略失败: %v", err)
		}

		global.KUBEGALE_LOG.Info("用户权限移除成功", zap.Int("userId", userId))
		return nil
	})
}

// RemoveRolePermissions 批量移除角色对应api权限
func (p *AuthorityService) RemoveRolePermissions(roleId int) error {
	if roleId <= 0 {
		return nil
	}

	// 查询角色名称
	var role system.Role
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", roleId, 0).First(&role).Error; err != nil {
		return fmt.Errorf("查询角色失败: %v", err)
	}

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	// 直接移除角色的所有权限策略
	if _, err := enforcer.RemoveFilteredPolicy(0, role.Name); err != nil {
		return fmt.Errorf("移除角色权限失败: %v", err)
	}

	// 重新加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载策略失败: %v", err)
	}

	return nil
}

// assignAPIPermissions 分配API权限
func (p *AuthorityService) assignAPIPermissions(roleName string, apiIds []int, batchSize int) error {
	if roleName == "" {
		return errors.New("角色名称不能为空")
	}

	// 如果API ID列表为空,直接返回
	if len(apiIds) == 0 {
		return nil
	}

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	// 批量获取API信息
	var apis []system.Api
	if err := global.KUBEGALE_DB.Where("id IN ? AND is_deleted = ?", apiIds, 0).Find(&apis).Error; err != nil {
		global.KUBEGALE_LOG.Error("批量获取API信息失败", zap.Error(err), zap.Ints("apiIds", apiIds))
		return fmt.Errorf("批量获取API信息失败: %v", err)
	}

	// 检查是否所有API都存在
	if len(apis) != len(apiIds) {
		global.KUBEGALE_LOG.Warn("部分API不存在",
			zap.Ints("请求的API", apiIds),
			zap.Int("找到的API数量", len(apis)))
	}

	// HTTP方法映射表
	methodMap := map[int]string{
		1: "GET",
		2: "POST",
		3: "PUT",
		4: "DELETE",
		5: "PATCH",
		6: "OPTIONS",
		7: "HEAD",
	}

	// 构建casbin策略规则
	policies := make([][]string, 0, len(apis))
	for _, api := range apis {
		// 获取HTTP方法
		method, ok := methodMap[api.Method]
		if !ok {
			global.KUBEGALE_LOG.Warn("无效的HTTP方法", zap.Int("method", api.Method), zap.Int("apiId", api.ID))
			continue // 跳过无效的方法，而不是直接返回错误
		}

		policies = append(policies, []string{roleName, api.Path, method})
	}

	// 批量添加策略
	if len(policies) > 0 {
		// 分批处理策略
		for i := 0; i < len(policies); i += batchSize {
			end := i + batchSize
			if end > len(policies) {
				end = len(policies)
			}

			batch := policies[i:end]
			_, err := enforcer.AddPolicies(batch)
			if err != nil {
				global.KUBEGALE_LOG.Error("批量添加策略失败",
					zap.Error(err),
					zap.String("roleName", roleName),
					zap.Int("batchSize", len(batch)))
				return fmt.Errorf("批量添加策略失败: %v", err)
			}
		}

		// 保存策略
		if err := enforcer.SavePolicy(); err != nil {
			global.KUBEGALE_LOG.Error("保存策略失败", zap.Error(err))
			return fmt.Errorf("保存策略失败: %v", err)
		}

		global.KUBEGALE_LOG.Info("API权限分配成功",
			zap.String("roleName", roleName),
			zap.Int("apiCount", len(policies)))
	}

	return nil
}

// batchAddPolicies 批量添加策略
func (p *AuthorityService) batchAddPolicies(policies [][]string, batchSize int) error {
	if len(policies) == 0 {
		return nil
	}

	if batchSize <= 0 {
		return errors.New("无效的批次大小")
	}

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	// 按批次处理策略规则
	for i := 0; i < len(policies); i += batchSize {
		end := i + batchSize
		if end > len(policies) {
			end = len(policies)
		}

		// 添加一批策略规则
		if _, err := enforcer.AddPolicies(policies[i:end]); err != nil {
			return fmt.Errorf("添加权限策略失败: %v", err)
		}
	}

	return nil
}

// AssignRoleToUsers 批量为用户分配角色
func (p *AuthorityService) AssignRoleToUsers(userIds []int, roleIds []int) error {
	const batchSize = 1000

	if len(userIds) == 0 {
		return errors.New("用户ID列表不能为空")
	}

	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 获取角色信息
		var roles []*system.Role
		if len(roleIds) > 0 {
			if err := tx.Where("id IN ? AND is_deleted = ?", roleIds, 0).Find(&roles).Error; err != nil {
				return fmt.Errorf("获取角色信息失败: %v", err)
			}
		}

		// 获取Casbin强制器
		enforcer, err := casbin.GetEnforcer()
		if err != nil {
			global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
			return fmt.Errorf("获取Casbin强制器失败: %w", err)
		}

		// 为每个用户添加角色
		for _, userId := range userIds {
			userStr := fmt.Sprintf("%d", userId)

			// 先移除用户现有的角色关联
			if _, err := enforcer.RemoveFilteredGroupingPolicy(0, userStr); err != nil {
				return fmt.Errorf("移除用户现有角色关联失败: %v", err)
			}

			// 更新用户的角色关联
			var user system.User
			if err := tx.First(&user, userId).Error; err != nil {
				return fmt.Errorf("获取用户信息失败: %v", err)
			}

			if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
				return fmt.Errorf("更新用户角色关联失败: %v", err)
			}

			// 添加角色关联策略
			if len(roles) > 0 {
				rolePolicies := make([][]string, 0, len(roles))
				for _, role := range roles {
					rolePolicies = append(rolePolicies, []string{userStr, role.Name})
				}
				if _, err := enforcer.AddGroupingPolicies(rolePolicies); err != nil {
					return fmt.Errorf("添加用户角色关联失败: %v", err)
				}
			}
		}

		// 加载最新的策略
		if err := enforcer.LoadPolicy(); err != nil {
			return fmt.Errorf("加载策略失败: %v", err)
		}

		return nil
	})
}

// RemoveUsersPermissions 批量移除用户权限
func (p *AuthorityService) RemoveUsersPermissions(userIds []int) error {
	if len(userIds) == 0 {
		return errors.New("用户ID列表不能为空")
	}

	// 获取Casbin强制器
	enforcer, err := casbin.GetEnforcer()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取Casbin强制器失败", zap.Error(err))
		return fmt.Errorf("获取Casbin强制器失败: %w", err)
	}

	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		for _, userId := range userIds {
			// 不允许删除userId为1的权限
			if userId == 1 {
				continue
			}

			userStr := fmt.Sprintf("%d", userId)

			// 获取用户信息
			var user system.User
			if err := tx.First(&user, userId).Error; err != nil {
				return fmt.Errorf("获取用户信息失败: %v", err)
			}

			// 获取用户当前的API列表
			var apis []*system.Api
			if err := tx.Model(&user).Association("Apis").Find(&apis); err != nil {
				return fmt.Errorf("获取用户API关联失败: %v", err)
			}

			// 移除每个API的权限策略
			for _, api := range apis {
				// HTTP方法映射表
				methodMap := map[int]string{
					1: "GET",
					2: "POST",
					3: "PUT",
					4: "DELETE",
					5: "PATCH",
					6: "OPTIONS",
					7: "HEAD",
				}
				method := methodMap[api.Method]

				if _, err := enforcer.RemovePolicy(userStr, api.Path, method); err != nil {
					return fmt.Errorf("移除用户API权限失败: %v", err)
				}
			}

			// 移除用户的角色关联
			if _, err := enforcer.RemoveFilteredGroupingPolicy(0, userStr); err != nil {
				return fmt.Errorf("移除用户角色关联失败: %v", err)
			}

			// 清空用户的关联
			if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
				return fmt.Errorf("清空用户角色关联失败: %v", err)
			}

			if err := tx.Model(&user).Association("Apis").Clear(); err != nil {
				return fmt.Errorf("清空用户API关联失败: %v", err)
			}
		}

		// 加载最新的策略
		if err := enforcer.LoadPolicy(); err != nil {
			return fmt.Errorf("加载策略失败: %v", err)
		}

		return nil
	})
}
