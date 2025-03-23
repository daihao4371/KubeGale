package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"time"
)

type MenuService struct{}

var MenuServiceApp = new(MenuService)

// GetMenus 获取菜单列表,支持分页和树形结构
func (m *MenuService) GetMenus(pageNum, pageSize int) ([]*system.Menu, int64, error) {
	// 参数验证
	if pageNum < 1 || pageSize < 1 {
		global.KUBEGALE_LOG.Warn("分页参数无效", zap.Int("页码", pageNum), zap.Int("每页数量", pageSize))
		return nil, 0, errors.New("分页参数无效")
	}

	// 创建查询构建器
	query := global.KUBEGALE_DB.Model(&system.Menu{}).Where("is_deleted = ?", 0)

	// 获取总记录数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取菜单总数失败", zap.Error(err))
		return nil, 0, fmt.Errorf("获取菜单总数失败: %w", err)
	}

	// 查询所有菜单
	var allMenus []*system.Menu
	if err := query.Find(&allMenus).Error; err != nil {
		global.KUBEGALE_LOG.Error("查询菜单列表失败", zap.Error(err))
		return nil, 0, fmt.Errorf("查询菜单列表失败: %w", err)
	}

	// 构建菜单树
	menuTree := buildMenuTree(allMenus, 0)

	// 如果需要分页，对树形结构进行分页处理
	if pageNum > 0 && pageSize > 0 {
		start := (pageNum - 1) * pageSize
		end := start + pageSize

		// 确保不越界
		if start >= len(menuTree) {
			return []*system.Menu{}, total, nil
		}
		if end > len(menuTree) {
			end = len(menuTree)
		}

		menuTree = menuTree[start:end]
	}

	return menuTree, total, nil
}

// buildMenuTree 构建菜单树形结构
func buildMenuTree(menus []*system.Menu, parentID int) []*system.Menu {
	var tree []*system.Menu

	for _, menu := range menus {
		if menu.ParentID == parentID {
			// 递归查找子菜单
			menu.Children = buildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}

	// 按照Meta.Order排序
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].Meta.Order < tree[j].Meta.Order
	})

	return tree
}

// CreateMenu 创建新菜单
func (m *MenuService) CreateMenu(menu *system.Menu) error {
	// 参数验证
	if menu == nil {
		global.KUBEGALE_LOG.Warn("菜单不能为空")
		return errors.New("菜单不能为空")
	}

	// 检查必填字段
	if menu.Name == "" {
		return errors.New("菜单名称不能为空")
	}

	if menu.Path == "" {
		return errors.New("菜单路径不能为空")
	}

	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 检查父菜单是否存在
		if menu.ParentID != 0 {
			var count int64
			if err := tx.Model(&system.Menu{}).Where("id = ? AND is_deleted = ?", menu.ParentID, 0).Count(&count).Error; err != nil {
				global.KUBEGALE_LOG.Error("检查父菜单失败", zap.Error(err))
				return fmt.Errorf("检查父菜单失败: %w", err)
			}
			if count == 0 {
				return errors.New("父菜单不存在")
			}
		}

		// 检查同级菜单名称是否重复
		var count int64
		if err := tx.Model(&system.Menu{}).Where("name = ? AND parent_id = ? AND is_deleted = ?", menu.Name, menu.ParentID, 0).Count(&count).Error; err != nil {
			global.KUBEGALE_LOG.Error("检查菜单名称失败", zap.Error(err))
			return fmt.Errorf("检查菜单名称失败: %w", err)
		}
		if count > 0 {
			return errors.New("同级菜单名称已存在")
		}

		// 创建菜单
		if err := tx.Create(menu).Error; err != nil {
			global.KUBEGALE_LOG.Error("创建菜单失败", zap.Error(err))
			return fmt.Errorf("创建菜单失败: %w", err)
		}

		// 如果有子菜单，递归创建子菜单
		if len(menu.Children) > 0 {
			for _, child := range menu.Children {
				child.ParentID = menu.ID // 设置父菜单ID
				if err := tx.Create(child).Error; err != nil {
					global.KUBEGALE_LOG.Error("创建子菜单失败", zap.Error(err))
					return fmt.Errorf("创建子菜单失败: %w", err)
				}
			}
		}

		global.KUBEGALE_LOG.Info("菜单创建成功", zap.Int("id", menu.ID), zap.String("name", menu.Name))
		return nil
	})
}

// GetMenuById 根据ID获取菜单
func (m *MenuService) GetMenuById(id int) (*system.Menu, error) {
	// 参数验证
	if id <= 0 {
		global.KUBEGALE_LOG.Warn("菜单ID无效", zap.Int("ID", id))
		return nil, errors.New("菜单ID无效")
	}

	var menu system.Menu
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", id, 0).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Warn("菜单不存在", zap.Int("ID", id))
			return nil, errors.New("菜单不存在")
		}
		global.KUBEGALE_LOG.Error("查询菜单失败", zap.Error(err), zap.Int("ID", id))
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}

	// 查询子菜单
	var children []*system.Menu
	if err := global.KUBEGALE_DB.Where("parent_id = ? AND is_deleted = ?", id, 0).Find(&children).Error; err != nil {
		global.KUBEGALE_LOG.Error("查询子菜单失败", zap.Error(err), zap.Int("parentID", id))
		return nil, fmt.Errorf("查询子菜单失败: %w", err)
	}

	// 为子菜单递归查询它们的子菜单
	for i := range children {
		childMenu, err := m.GetMenuById(children[i].ID)
		if err != nil {
			global.KUBEGALE_LOG.Warn("获取子菜单详情失败", zap.Error(err), zap.Int("childID", children[i].ID))
			continue
		}
		children[i] = childMenu
	}

	// 设置子菜单
	menu.Children = children

	return &menu, nil
}

// UpdateMenu 更新菜单
func (m *MenuService) UpdateMenu(menu *system.Menu) error {
	// 参数验证
	if menu == nil {
		global.KUBEGALE_LOG.Warn("菜单对象不能为空")
		return errors.New("菜单对象不能为空")
	}
	if menu.ID <= 0 {
		global.KUBEGALE_LOG.Warn("无效的菜单ID", zap.Int("ID", menu.ID))
		return errors.New("无效的菜单ID")
	}
	if menu.Name == "" {
		return errors.New("菜单名称不能为空")
	}
	if menu.Path == "" {
		return errors.New("菜单路径不能为空")
	}

	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 检查菜单是否存在
		var existingMenu system.Menu
		if err := tx.Where("id = ? AND is_deleted = ?", menu.ID, 0).First(&existingMenu).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				global.KUBEGALE_LOG.Warn("菜单不存在", zap.Int("ID", menu.ID))
				return errors.New("菜单不存在")
			}
			global.KUBEGALE_LOG.Error("检查菜单是否存在失败", zap.Error(err), zap.Int("ID", menu.ID))
			return fmt.Errorf("检查菜单是否存在失败: %w", err)
		}

		// 检查是否有子菜单
		var childCount int64
		if err := tx.Model(&system.Menu{}).Where("parent_id = ? AND is_deleted = ?", menu.ID, 0).Count(&childCount).Error; err != nil {
			global.KUBEGALE_LOG.Error("检查子菜单失败", zap.Error(err), zap.Int("ID", menu.ID))
			return fmt.Errorf("检查子菜单失败: %w", err)
		}

		// 如果有子菜单且尝试修改父级菜单ID,则不允许修改
		if childCount > 0 && existingMenu.ParentID != menu.ParentID {
			global.KUBEGALE_LOG.Warn("当前菜单存在子菜单,不能修改父级菜单",
				zap.Int("ID", menu.ID),
				zap.Int("oldParentID", existingMenu.ParentID),
				zap.Int("newParentID", menu.ParentID))
			return errors.New("当前菜单存在子菜单,不能修改父级菜单")
		}

		// 检查父菜单是否存在且不能将菜单设置为自己的子菜单
		if menu.ParentID != 0 {
			if menu.ParentID == menu.ID {
				global.KUBEGALE_LOG.Warn("不能将菜单设置为自己的子菜单", zap.Int("ID", menu.ID))
				return errors.New("不能将菜单设置为自己的子菜单")
			}
			var count int64
			if err := tx.Model(&system.Menu{}).Where("id = ? AND is_deleted = ?", menu.ParentID, 0).Count(&count).Error; err != nil {
				global.KUBEGALE_LOG.Error("检查父菜单失败", zap.Error(err), zap.Int("parentID", menu.ParentID))
				return fmt.Errorf("检查父菜单失败: %w", err)
			}
			if count == 0 {
				global.KUBEGALE_LOG.Warn("父菜单不存在", zap.Int("parentID", menu.ParentID))
				return errors.New("父菜单不存在")
			}
		}

		// 检查同级菜单名称是否重复
		var count int64
		if err := tx.Model(&system.Menu{}).Where("name = ? AND parent_id = ? AND id != ? AND is_deleted = ?",
			menu.Name, menu.ParentID, menu.ID, 0).Count(&count).Error; err != nil {
			global.KUBEGALE_LOG.Error("检查菜单名称失败", zap.Error(err), zap.String("name", menu.Name))
			return fmt.Errorf("检查菜单名称失败: %w", err)
		}
		if count > 0 {
			global.KUBEGALE_LOG.Warn("同级菜单名称已存在", zap.String("name", menu.Name), zap.Int("parentID", menu.ParentID))
			return errors.New("同级菜单名称已存在")
		}

		// 保留原始的创建时间
		menu.CreateTime = existingMenu.CreateTime
		// 更新修改时间
		menu.UpdateTime = time.Now().Unix()

		// 更新菜单信息
		result := tx.Model(&system.Menu{}).Where("id = ? AND is_deleted = ?", menu.ID, 0).Updates(menu)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("更新菜单失败", zap.Error(result.Error), zap.Int("ID", menu.ID))
			return fmt.Errorf("更新菜单失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			global.KUBEGALE_LOG.Warn("菜单不存在或已被删除", zap.Int("ID", menu.ID))
			return errors.New("菜单不存在或已被删除")
		}

		global.KUBEGALE_LOG.Info("菜单更新成功", zap.Int("id", menu.ID), zap.String("name", menu.Name))
		return nil
	})
}

// DeleteMenu 删除菜单
func (m *MenuService) DeleteMenu(id int) error {
	// 参数验证
	if id <= 0 {
		global.KUBEGALE_LOG.Warn("无效的菜单ID", zap.Int("ID", id))
		return errors.New("无效的菜单ID")
	}

	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 检查菜单是否存在
		var menu system.Menu
		if err := tx.Where("id = ? AND is_deleted = ?", id, 0).First(&menu).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				global.KUBEGALE_LOG.Warn("菜单不存在", zap.Int("ID", id))
				return errors.New("菜单不存在")
			}
			global.KUBEGALE_LOG.Error("检查菜单是否存在失败", zap.Error(err), zap.Int("ID", id))
			return fmt.Errorf("检查菜单是否存在失败: %w", err)
		}

		// 检查是否有子菜单
		var count int64
		if err := tx.Model(&system.Menu{}).Where("parent_id = ? AND is_deleted = ?", id, 0).Count(&count).Error; err != nil {
			global.KUBEGALE_LOG.Error("检查子菜单失败", zap.Error(err), zap.Int("ID", id))
			return fmt.Errorf("检查子菜单失败: %w", err)
		}
		if count > 0 {
			global.KUBEGALE_LOG.Warn("存在子菜单,不能删除", zap.Int("ID", id), zap.Int64("childCount", count))
			return errors.New("存在子菜单,不能删除")
		}

		// 软删除菜单
		updates := map[string]interface{}{
			"is_deleted":  1,
			"update_time": time.Now().Unix(),
		}
		result := tx.Model(&system.Menu{}).Where("id = ? AND is_deleted = ?", id, 0).Updates(updates)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("删除菜单失败", zap.Error(result.Error), zap.Int("ID", id))
			return fmt.Errorf("删除菜单失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			global.KUBEGALE_LOG.Warn("菜单不存在或已被删除", zap.Int("ID", id))
			return errors.New("菜单不存在或已被删除")
		}

		global.KUBEGALE_LOG.Info("菜单删除成功", zap.Int("id", id), zap.String("name", menu.Name))
		return nil
	})
}

// ListMenuTree 获取菜单树形结构
func (m *MenuService) ListMenuTree() ([]*system.Menu, error) {
	// 预分配合适的初始容量
	menus := make([]*system.Menu, 0, 50)

	// 查询所有未删除的菜单
	if err := global.KUBEGALE_DB.
		Select("id, name, parent_id, path, component, route_name, hidden, redirect, meta, create_time, update_time").
		Where("is_deleted = ?", 0).
		Find(&menus).Error; err != nil {
		global.KUBEGALE_LOG.Error("查询菜单列表失败", zap.Error(err))
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}

	global.KUBEGALE_LOG.Info("查询菜单列表成功", zap.Int("count", len(menus)))

	// 预分配map容量
	menuMap := make(map[int]*system.Menu, len(menus))
	rootMenus := make([]*system.Menu, 0, len(menus)/3) // 假设大约1/3的菜单是根菜单

	// 第一次遍历,建立ID到菜单的映射
	for _, menu := range menus {
		if menu == nil {
			continue
		}

		menu.Children = make([]*system.Menu, 0, 4) // 预分配子菜单切片,假设平均4个子菜单
		menuMap[menu.ID] = menu
	}

	// 第二次遍历,构建树形结构
	for _, menu := range menus {
		if menu == nil {
			continue
		}
		if menu.ParentID == 0 {
			rootMenus = append(rootMenus, menu)
		} else {
			if parent, exists := menuMap[menu.ParentID]; exists {
				parent.Children = append(parent.Children, menu)
			} else {
				// 如果找不到父节点,作为根节点处理
				global.KUBEGALE_LOG.Warn("菜单的父节点不存在，作为根节点处理", 
					zap.Int("menuID", menu.ID), 
					zap.Int("parentID", menu.ParentID))
				rootMenus = append(rootMenus, menu)
			}
		}
	}

	// 对根菜单和所有子菜单进行排序
	sortMenus(rootMenus)

	global.KUBEGALE_LOG.Info("构建菜单树成功", zap.Int("rootCount", len(rootMenus)))
	return rootMenus, nil
}

// sortMenus 递归对菜单及其子菜单进行排序
func sortMenus(menus []*system.Menu) {
	// 按照Meta.Order排序
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].Meta.Order < menus[j].Meta.Order
	})

	// 递归排序子菜单
	for _, menu := range menus {
		if len(menu.Children) > 0 {
			sortMenus(menu.Children)
		}
	}
}

// UpdateUserMenu 更新用户菜单关联
func (m *MenuService) UpdateUserMenu(userId int, menuIds []int) error {
	if userId <= 0 {
		global.KUBEGALE_LOG.Error("无效的用户ID", zap.Int("userId", userId))
		return errors.New("无效的用户ID")
	}
	if len(menuIds) == 0 {
		global.KUBEGALE_LOG.Error("无效的菜单ID", zap.Ints("menuIds", menuIds))
		return errors.New("无效的菜单ID")
	}

	// 检查用户是否存在且未删除
	var user system.User
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", userId, 0).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Error("用户不存在或已删除", zap.Int("userId", userId))
			return errors.New("用户不存在或已删除")
		}
		global.KUBEGALE_LOG.Error("查询用户失败", zap.Error(err))
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查所有菜单是否存在且未删除
	var count int64
	if err := global.KUBEGALE_DB.Model(&system.Menu{}).Where("id IN ? AND is_deleted = ?", menuIds, 0).Count(&count).Error; err != nil {
		global.KUBEGALE_LOG.Error("查询菜单失败", zap.Error(err))
		return fmt.Errorf("查询菜单失败: %w", err)
	}
	if int(count) != len(menuIds) {
		global.KUBEGALE_LOG.Error("部分菜单不存在或已删除", zap.Ints("menuIds", menuIds), zap.Int64("existCount", count))
		return errors.New("部分菜单不存在或已删除")
	}

	// 使用事务保证数据一致性
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 构建批量插入的数据
		userMenus := make([]map[string]interface{}, 0, len(menuIds))
		for _, menuId := range menuIds {
			userMenus = append(userMenus, map[string]interface{}{
				"user_id":     userId,
				"menu_id":     menuId,
				"create_time": time.Now().Unix(),
				"update_time": time.Now().Unix(),
			})
		}

		// 先删除已有的关联
		if err := tx.Table("user_menus").Where("user_id = ?", userId).Delete(nil).Error; err != nil {
			global.KUBEGALE_LOG.Error("删除已有关联失败", zap.Error(err))
			return fmt.Errorf("删除已有关联失败: %w", err)
		}

		// 批量创建新的关联
		if len(userMenus) > 0 {
			if err := tx.Table("user_menus").Create(userMenus).Error; err != nil {
				global.KUBEGALE_LOG.Error("添加用户菜单关联失败", zap.Error(err))
				return fmt.Errorf("添加用户菜单关联失败: %w", err)
			}
		}

		global.KUBEGALE_LOG.Info("更新用户菜单关联成功", 
			zap.Int("userId", userId), 
			zap.Int("menuCount", len(menuIds)))
		return nil
	})
}