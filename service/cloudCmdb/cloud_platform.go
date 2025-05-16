package cloudCmdb

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"KubeGale/model/common/request"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CloudPlatformService struct{}

// List  厂商列表
func (p *CloudPlatformService) List(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list []model.CloudPlatform, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloudList []model.CloudPlatform

	if cloud.Name != "" {
		db = db.Where("name LIKE ?", "%"+cloud.Name+"%")
	}

	err = db.Count(&total).Error

	if err != nil {
		return cloudList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["name"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return cloudList, total, err
			}

			err = db.Order(OrderStr).Find(&cloudList).Error
		} else {
			err = db.Order("id").Find(&cloudList).Error
		}
	}

	return cloudList, total, err
}

// GetCloudPlatformById 获取单个厂商信息
func (p *CloudPlatformService) GetCloudPlatformById(id int) (cloud model.CloudPlatform, regions []model.CloudRegions, err error) {
	if err = global.KUBEGALE_DB.Where("id = ?", id).First(&cloud).Error; err != nil {
		return cloud, regions, nil
	}

	if err = global.KUBEGALE_DB.Select("id, name").Where("cloud_platform_id = ?", id).Find(&regions).Error; err != nil {
		return cloud, regions, nil
	}

	return cloud, regions, nil
}

// CreateCloudPlatform 创建厂商信息
func (p *CloudPlatformService) CreateCloudPlatform(cloud model.CloudPlatform) (err error) {
	if !errors.Is(global.KUBEGALE_DB.Where("name = ?", cloud.Name).First(&model.CloudPlatform{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同厂商")
	}
	return global.KUBEGALE_DB.Create(&cloud).Error
}

// UpdateCloudPlatform 更新厂商信息
func (p *CloudPlatformService) UpdateCloudPlatform(cloud model.CloudPlatform) (err error) {
	return global.KUBEGALE_DB.Where("id = ?", cloud.ID).First(&model.CloudPlatform{}).Updates(&cloud).Error
}

// DeleteCloudPlatform 删除厂商信息
func (p *CloudPlatformService) DeleteCloudPlatform(req request.GetById) (err error) {
	var cloud model.CloudPlatform
	if err = global.KUBEGALE_DB.Where("id = ?", req.ID).First(&cloud).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err = global.KUBEGALE_DB.Delete(&cloud).Error; err != nil {
		return err
	}

	return err
}

// DeleteCloudPlatformByIds 批量删除厂商信息
func (p *CloudPlatformService) DeleteCloudPlatformByIds(ids request.IdsReq) (err error) {
	return global.KUBEGALE_DB.Delete(&[]model.CloudPlatform{}, "id in ?", ids.Ids).Error
}
