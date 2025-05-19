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

// List 厂商列表
// 功能：获取云平台厂商列表，支持分页、排序和条件筛选
// 参数：
//   - cloud: 云平台查询条件
//   - info: 分页信息
//   - order: 排序字段
//   - desc: 是否降序排序
//
// 返回：
//   - list: 云平台列表
//   - total: 总记录数
//   - err: 错误信息
func (p *CloudPlatformService) List(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list []model.CloudPlatform, total int64, err error) {
	// 计算分页参数
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloudList []model.CloudPlatform

	// 添加名称模糊查询条件
	if cloud.Name != "" {
		db = db.Where("name LIKE ?", "%"+cloud.Name+"%")
	}

	// 获取总记录数
	err = db.Count(&total).Error

	if err != nil {
		return cloudList, total, err
	} else {
		// 设置分页
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

			// 执行排序查询
			err = db.Order(OrderStr).Find(&cloudList).Error
		} else {
			// 默认按ID排序
			err = db.Order("id").Find(&cloudList).Error
		}
	}

	return cloudList, total, err
}

// GetCloudPlatformById 获取单个厂商信息
// 功能：根据ID获取云平台厂商信息及其关联的区域信息
// 参数：
//   - id: 云平台ID
//
// 返回：
//   - cloud: 云平台信息
//   - regions: 关联的区域信息列表
//   - err: 错误信息
func (p *CloudPlatformService) GetCloudPlatformById(id int) (cloud model.CloudPlatform, regions []model.CloudRegions, err error) {
	// 获取云平台基本信息
	if err = global.KUBEGALE_DB.Where("id = ?", id).First(&cloud).Error; err != nil {
		return cloud, regions, nil
	}

	// 获取关联的区域信息
	if err = global.KUBEGALE_DB.Select("id, name").Where("cloud_platform_id = ?", id).Find(&regions).Error; err != nil {
		return cloud, regions, nil
	}

	return cloud, regions, nil
}

// CreateCloudPlatform 创建厂商信息
// 功能：创建新的云平台厂商信息，检查名称是否重复
// 参数：
//   - cloud: 云平台信息
//
// 返回：
//   - err: 错误信息
func (p *CloudPlatformService) CreateCloudPlatform(cloud model.CloudPlatform) (err error) {
	// 检查厂商名称是否已存在
	if !errors.Is(global.KUBEGALE_DB.Where("name = ?", cloud.Name).First(&model.CloudPlatform{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同厂商")
	}
	// 创建新厂商记录
	return global.KUBEGALE_DB.Create(&cloud).Error
}

// UpdateCloudPlatform 更新厂商信息
// 功能：更新指定ID的云平台厂商信息
// 参数：
//   - cloud: 云平台信息
//
// 返回：
//   - err: 错误信息
func (p *CloudPlatformService) UpdateCloudPlatform(cloud model.CloudPlatform) (err error) {
	// 先查询记录是否存在，然后更新
	return global.KUBEGALE_DB.Where("id = ?", cloud.ID).First(&model.CloudPlatform{}).Updates(&cloud).Error
}

// DeleteCloudPlatform 删除厂商信息
// 功能：删除指定ID的云平台厂商信息
// 参数：
//   - req: 包含ID的请求参数
//
// 返回：
//   - err: 错误信息
func (p *CloudPlatformService) DeleteCloudPlatform(req request.GetById) (err error) {
	var cloud model.CloudPlatform
	// 检查记录是否存在
	if err = global.KUBEGALE_DB.Where("id = ?", req.ID).First(&cloud).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 执行删除操作
	if err = global.KUBEGALE_DB.Delete(&cloud).Error; err != nil {
		return err
	}

	return err
}

// DeleteCloudPlatformByIds 批量删除厂商信息
// 功能：批量删除多个云平台厂商信息
// 参数：
//   - ids: 包含多个ID的请求参数
//
// 返回：
//   - err: 错误信息
func (p *CloudPlatformService) DeleteCloudPlatformByIds(ids request.IdsReq) (err error) {
	// 批量删除指定ID的记录
	return global.KUBEGALE_DB.Delete(&[]model.CloudPlatform{}, "id in ?", ids.Ids).Error
}
