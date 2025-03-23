package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ApiService struct{}

var ApiServiceApp = new(ApiService)

// CreateApi 创建新的API记录
func (a *ApiService) CreateApi(api *system.Api) error {
	// 参数验证
	if api == nil {
		global.KUBEGALE_LOG.Warn("API对象不能为空")
		return errors.New("API对象不能为空")
	}

	// 检查必填字段
	if api.Name == "" {
		global.KUBEGALE_LOG.Warn("API名称不能为空")
		return errors.New("API名称不能为空")
	}
	if api.Path == "" {
		global.KUBEGALE_LOG.Warn("API路径不能为空")
		return errors.New("API路径不能为空")
	}
	if api.Method <= 0 || api.Method > 4 {
		global.KUBEGALE_LOG.Warn("无效的HTTP方法", zap.Int("method", api.Method))
		return errors.New("无效的HTTP方法")
	}

	if api.Category <= 0 {
		api.Category = 1 // 默认为系统API
	}

	// 设置时间戳
	api.CreateTime = time.Now().Unix()
	api.UpdateTime = time.Now().Unix()
	api.IsDeleted = 0

	// 使用事务创建API
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 检查API路径和方法是否已存在
		var count int64
		if err := tx.Model(&system.Api{}).Where("path = ? AND method = ? AND is_deleted = ?", api.Path, api.Method, 0).Count(&count).Error; err != nil {
			global.KUBEGALE_LOG.Error("检查API是否存在失败", zap.Error(err))
			return fmt.Errorf("检查API是否存在失败: %w", err)
		}
		if count > 0 {
			global.KUBEGALE_LOG.Warn("相同路径和方法的API已存在", zap.String("path", api.Path), zap.Int("method", api.Method))
			return errors.New("相同路径和方法的API已存在")
		}

		// 创建API
		if err := tx.Create(api).Error; err != nil {
			global.KUBEGALE_LOG.Error("创建API失败", zap.Error(err))
			return fmt.Errorf("创建API失败: %w", err)
		}

		global.KUBEGALE_LOG.Info("API创建成功", zap.Int("id", api.ID), zap.String("name", api.Name), zap.String("path", api.Path))
		return nil
	})
}

// GetApiById 根据ID获取API记录
func (a *ApiService) GetApiById(id int) (*system.Api, error) {
	// 参数验证
	if id <= 0 {
		global.KUBEGALE_LOG.Warn("无效的API ID", zap.Int("ID", id))
		return nil, errors.New("无效的API ID")
	}

	var api system.Api
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", id, 0).First(&api).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Warn("API不存在", zap.Int("ID", id))
			return nil, errors.New("API不存在")
		}
		global.KUBEGALE_LOG.Error("查询API失败", zap.Error(err), zap.Int("ID", id))
		return nil, fmt.Errorf("查询API失败: %w", err)
	}

	global.KUBEGALE_LOG.Info("获取API成功", zap.Int("id", api.ID), zap.String("name", api.Name), zap.String("path", api.Path))
	return &api, nil
}

// UpdateApi 更新API记录
func (a *ApiService) UpdateApi(api *system.Api) error {
	// 参数验证
	if api == nil {
		global.KUBEGALE_LOG.Warn("API对象不能为空")
		return errors.New("API对象不能为空")
	}
	if api.ID <= 0 {
		global.KUBEGALE_LOG.Warn("无效的API ID", zap.Int("ID", api.ID))
		return errors.New("无效的API ID")
	}
	if api.Name == "" {
		global.KUBEGALE_LOG.Warn("API名称不能为空")
		return errors.New("API名称不能为空")
	}
	if api.Path == "" {
		global.KUBEGALE_LOG.Warn("API路径不能为空")
		return errors.New("API路径不能为空")
	}
	if api.Method <= 0 || api.Method > 4 {
		global.KUBEGALE_LOG.Warn("无效的HTTP方法", zap.Int("method", api.Method))
		return errors.New("无效的HTTP方法")
	}

	// 获取旧的API记录
	oldApi, err := a.GetApiById(api.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取原API记录失败", zap.Error(err), zap.Int("ID", api.ID))
		return fmt.Errorf("获取原API记录失败: %w", err)
	}

	// 使用事务保证数据一致性
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 如果路径或方法发生变化，检查是否与其他API冲突
		if oldApi.Path != api.Path || oldApi.Method != api.Method {
			var count int64
			if err := tx.Model(&system.Api{}).Where("path = ? AND method = ? AND id != ? AND is_deleted = ?",
				api.Path, api.Method, api.ID, 0).Count(&count).Error; err != nil {
				global.KUBEGALE_LOG.Error("检查API是否存在失败", zap.Error(err))
				return fmt.Errorf("检查API是否存在失败: %w", err)
			}
			if count > 0 {
				global.KUBEGALE_LOG.Warn("相同路径和方法的API已存在", zap.String("path", api.Path), zap.Int("method", api.Method))
				return errors.New("相同路径和方法的API已存在")
			}
		}

		// 准备更新的字段
		updates := map[string]interface{}{
			"name":        api.Name,
			"path":        api.Path,
			"method":      api.Method,
			"description": api.Description,
			"version":     api.Version,
			"category":    api.Category,
			"is_public":   api.IsPublic,
			"update_time": time.Now().Unix(),
		}

		// 更新API记录
		result := tx.Model(&system.Api{}).Where("id = ? AND is_deleted = ?", api.ID, 0).Updates(updates)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("更新API失败", zap.Error(result.Error), zap.Int("ID", api.ID))
			return fmt.Errorf("更新API失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			global.KUBEGALE_LOG.Warn("API不存在或已被删除", zap.Int("ID", api.ID))
			return errors.New("API不存在或已被删除")
		}

		// 记录更新成功的日志
		global.KUBEGALE_LOG.Info("API更新成功",
			zap.Int("id", api.ID),
			zap.String("name", api.Name),
			zap.String("path", api.Path),
			zap.Int("method", api.Method))
		return nil
	})
}

// DeleteApi 软删除API记录
func (a *ApiService) DeleteApi(id int) error {
	// 参数验证
	if id <= 0 {
		global.KUBEGALE_LOG.Warn("无效的API ID", zap.Int("ID", id))
		return errors.New("无效的API ID")
	}

	// 检查API是否存在
	var api system.Api
	if err := global.KUBEGALE_DB.Where("id = ? AND is_deleted = ?", id, 0).First(&api).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.KUBEGALE_LOG.Warn("API不存在或已被删除", zap.Int("ID", id))
			return errors.New("API不存在或已被删除")
		}
		global.KUBEGALE_LOG.Error("查询API失败", zap.Error(err), zap.Int("ID", id))
		return fmt.Errorf("查询API失败: %w", err)
	}

	// 使用事务保证数据一致性
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 软删除API记录
		updates := map[string]interface{}{
			"is_deleted":  1,
			"update_time": time.Now().Unix(),
		}

		result := tx.Model(&system.Api{}).Where("id = ? AND is_deleted = ?", id, 0).Updates(updates)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("删除API失败", zap.Error(result.Error), zap.Int("ID", id))
			return fmt.Errorf("删除API失败: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			global.KUBEGALE_LOG.Warn("API不存在或已被删除", zap.Int("ID", id))
			return errors.New("API不存在或已被删除")
		}

		// 可选：清除与该API相关的权限关联
		if err := tx.Exec("DELETE FROM role_apis WHERE api_id = ?", id).Error; err != nil {
			global.KUBEGALE_LOG.Error("清除API权限关联失败", zap.Error(err), zap.Int("apiID", id))
			return fmt.Errorf("清除API权限关联失败: %w", err)
		}

		if err := tx.Exec("DELETE FROM user_apis WHERE api_id = ?", id).Error; err != nil {
			global.KUBEGALE_LOG.Error("清除用户API关联失败", zap.Error(err), zap.Int("apiID", id))
			return fmt.Errorf("清除用户API关联失败: %w", err)
		}

		global.KUBEGALE_LOG.Info("API删除成功", zap.Int("id", id), zap.String("name", api.Name), zap.String("path", api.Path))
		return nil
	})
}

// ListApis 分页获取API列表
func (a *ApiService) ListApis(page, pageSize int) ([]*system.Api, int64, error) {
	// 参数验证
	if page <= 0 {
		page = 1
		global.KUBEGALE_LOG.Warn("页码参数无效，使用默认值", zap.Int("page", page))
	}
	if pageSize <= 0 {
		pageSize = 10
		global.KUBEGALE_LOG.Warn("每页数量参数无效，使用默认值", zap.Int("pageSize", pageSize))
	}
	if pageSize > 100 {
		pageSize = 100
		global.KUBEGALE_LOG.Warn("每页数量过大，使用最大值", zap.Int("pageSize", pageSize))
	}

	var apis []*system.Api
	var total int64

	// 构建基础查询
	db := global.KUBEGALE_DB.Model(&system.Api{}).Where("is_deleted = ?", 0)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取API总数失败", zap.Error(err))
		return nil, 0, fmt.Errorf("获取API总数失败: %w", err)
	}

	// 如果没有记录，直接返回空列表
	if total == 0 {
		global.KUBEGALE_LOG.Info("没有找到符合条件的API记录")
		return []*system.Api{}, 0, nil
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&apis).Error; err != nil {
		global.KUBEGALE_LOG.Error("查询API列表失败", zap.Error(err),
			zap.Int("page", page), zap.Int("pageSize", pageSize))
		return nil, 0, fmt.Errorf("查询API列表失败: %w", err)
	}

	global.KUBEGALE_LOG.Info("获取API列表成功",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int64("total", total),
		zap.Int("count", len(apis)))
	return apis, total, nil
}
