package cmdb

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	cmdbReq "KubeGale/model/cmdb/request"
	"fmt"

	"gorm.io/gorm"
)

type CmdbProjectsService struct{}

// CreateCmdbProjects 创建cmdbProjects表记录
func (cmdbProjectsService *CmdbProjectsService) CreateCmdbProjects(cmdbProjects *cmdb.CmdbProjects) (err error) {
	err = global.KUBEGALE_DB.Create(cmdbProjects).Error
	return err
}

// DeleteCmdbProjects 删除cmdbProjects表记录
func (cmdbProjectsService *CmdbProjectsService) DeleteCmdbProjects(ID string, userID uint) (err error) {
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 先检查记录是否存在
		var project cmdb.CmdbProjects
		if err := tx.Where("id = ?", ID).First(&project).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("项目不存在")
			}
			return err
		}

		// 执行物理删除
		if err := tx.Unscoped().Delete(&cmdb.CmdbProjects{}, "id = ?", ID).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteCmdbProjectsByIds 批量删除cmdbProjects表记录
func (cmdbProjectsService *CmdbProjectsService) DeleteCmdbProjectsByIds(IDs []string, deleted_by uint) (err error) {
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 检查所有记录是否存在
		var count int64
		if err := tx.Model(&cmdb.CmdbProjects{}).Where("id IN ?", IDs).Count(&count).Error; err != nil {
			return err
		}
		if int(count) != len(IDs) {
			return fmt.Errorf("部分项目不存在")
		}

		// 执行批量物理删除
		if err := tx.Unscoped().Delete(&cmdb.CmdbProjects{}, "id IN ?", IDs).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateCmdbProjects 更新cmdbProjects表记录
func (cmdbProjectsService *CmdbProjectsService) UpdateCmdbProjects(cmdbProjects cmdb.CmdbProjects) (err error) {
	err = global.KUBEGALE_DB.Model(&cmdb.CmdbProjects{}).Where("id = ?", cmdbProjects.ID).Updates(&cmdbProjects).Error
	return err
}

// GetCmdbProjects 根据ID获取cmdbProjects表记录
func (cmdbProjectsService *CmdbProjectsService) GetCmdbProjects(ID string) (cmdbProjects cmdb.CmdbProjects, err error) {
	err = global.KUBEGALE_DB.Where("id = ?", ID).First(&cmdbProjects).Error
	return
}

// GetCmdbProjectsInfoList 分页获取cmdbProjects表记录
func (cmdbProjectsService *CmdbProjectsService) GetCmdbProjectsInfoList(info cmdbReq.CmdbProjectsSearch) (list []cmdb.CmdbProjects, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.KUBEGALE_DB.Model(&cmdb.CmdbProjects{})
	var cmdbProjectss []cmdb.CmdbProjects
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	// 只查询未删除的记录
	db = db.Unscoped().Where("deleted_at IS NULL")

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&cmdbProjectss).Error
	return cmdbProjectss, total, err
}
