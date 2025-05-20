package cloudCmdb

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	"KubeGale/utils/cloudCmdb/aliyun"
	"KubeGale/utils/cloudCmdb/aws"
	"KubeGale/utils/cloudCmdb/huawei"
	"KubeGale/utils/cloudCmdb/tencent"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CloudRegionService struct{}

// TencentRegion 腾讯云同步Region
// 功能：同步腾讯云平台的所有区域信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r CloudRegionService) TencentRegion(cloud model.CloudPlatform) (err error) {
	// 初始化区域列表
	var regions []model.CloudRegions
	// 创建腾讯云区域实例
	region := tencent.NewRegion()
	// 获取区域列表
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	// 转换区域信息格式
	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            *instance.RegionName,
			RegionId:        "tencent-" + *instance.Region,
			RegionName:      *instance.RegionName,
			CloudPlatformId: cloud.ID,
		})
	}

	// 更新数据库：如果有数据则进行更新
	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// AliyunRegion 阿里云同步Region
// 功能：同步阿里云平台的所有区域信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r CloudRegionService) AliyunRegion(cloud model.CloudPlatform) (err error) {
	// 初始化区域列表
	var regions []model.CloudRegions
	// 创建阿里云区域实例
	region := aliyun.NewRegion()
	// 获取区域列表
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	// 转换区域信息格式
	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            instance.LocalName,
			RegionId:        "aliyun-" + instance.RegionId,
			RegionName:      instance.LocalName,
			CloudPlatformId: cloud.ID,
		})
	}

	// 更新数据库：如果有数据则进行更新
	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// HuaweiRegion 华为云同步Region
// 功能：同步华为云平台的所有区域信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r CloudRegionService) HuaweiRegion(cloud model.CloudPlatform) (err error) {
	// 初始化区域列表
	var regions []model.CloudRegions
	// 创建华为云区域实例
	region := huawei.NewRegion()
	// 获取区域列表
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	// 转换区域信息格式
	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            instance.Name,
			RegionId:        "huawei-" + instance.RegionId,
			RegionName:      instance.Name,
			CloudPlatformId: cloud.ID,
		})
	}

	// 更新数据库：如果有数据则进行更新
	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// AwsRegion 亚马逊云同步Region
// 功能：同步亚马逊云平台的所有区域信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r CloudRegionService) AwsRegion(cloud model.CloudPlatform) (err error) {
	// 初始化区域列表
	var regions []model.CloudRegions
	// 创建AWS区域实例
	region := aws.NewRegion()
	// 获取区域列表
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	// 转换区域信息格式
	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            *instance.RegionName,
			RegionId:        "aws-" + *instance.RegionName,
			RegionName:      *instance.RegionName,
			CloudPlatformId: cloud.ID,
		})
	}

	// 更新数据库：如果有数据则进行更新
	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// UpdateRegions 更新Region信息
// 功能：批量更新或插入区域信息到数据库
// 参数：
//   - list: 区域信息列表
func (r *CloudRegionService) UpdateRegions(list []model.CloudRegions) {
	db := global.KUBEGALE_DB.Model(model.CloudRegions{})

	for _, region := range list {
		// 开始事务
		tx := db.Begin()

		tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"region_id",
			}),
		}).Create(&region)

		// 更新所有存在的记录，忽略不存在的记录
		// 使用 OnConflict 子句处理冲突情况
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"region_id",
			}),
		}).Create(&region).Error; err != nil {
			global.KUBEGALE_LOG.Error("region messages update fail!", zap.Error(err))
			tx.Rollback()
		}

		// 插入不存在的记录
		// 使用 OnConflict 子句的 DoNothing 选项处理插入冲突
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&region).Error; err != nil {
			global.KUBEGALE_LOG.Error("region messages insert fail!", zap.Error(err))
			tx.Rollback()
		}

		// 提交事务：确认所有操作成功完成
		tx.Commit()
	}
}

// SyncRegion 同步云平台所有Region
// 功能：根据云平台类型选择对应的同步方法
// 参数：
//   - id: 云平台ID
//
// 返回：
//   - err: 错误信息
func (r CloudRegionService) SyncRegion(id int) (err error) {
	// 获取云平台信息
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloud model.CloudPlatform
	if err := db.Where("id = ?", id).First(&cloud).Error; err != nil {
		return err
	}

	// 根据云平台类型选择对应的同步方法
	if cloud.Platform == "aliyun" {
		go func(cloud model.CloudPlatform) {
			if err = r.AliyunRegion(cloud); err != nil {
				global.KUBEGALE_LOG.Error("aliyun region aync fail!", zap.Error(err))
			}
		}(cloud)
	}

	if cloud.Platform == "tencent" {
		go func(cloud model.CloudPlatform) {
			if err = r.TencentRegion(cloud); err != nil {
				global.KUBEGALE_LOG.Error("tencent region aync fail!", zap.Error(err))
			}
		}(cloud)
	}

	if cloud.Platform == "huawei" {
		go func(cloud model.CloudPlatform) {
			if err = r.HuaweiRegion(cloud); err != nil {
				global.KUBEGALE_LOG.Error("huawei region aync fail!", zap.Error(err))
			}
		}(cloud)
	}

	if cloud.Platform == "aws" {
		go func(cloud model.CloudPlatform) {
			if err = r.AwsRegion(cloud); err != nil {
				global.KUBEGALE_LOG.Error("aws region aync fail!", zap.Error(err))
			}
		}(cloud)
	}

	return
}

// GetRegionTree 获取区域树形结构
// 功能：获取所有云平台及其区域的树形结构
// 返回：
//   - list: 区域树形结构列表
//   - err: 错误信息
func (r *CloudRegionService) GetRegionTree() (list []model.PlatformTree, err error) {
	// 获取所有云平台
	var platforms []model.CloudPlatform
	if err := global.KUBEGALE_DB.Find(&platforms).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取云平台列表失败", zap.Error(err))
		return nil, err
	}

	// 遍历每个云平台
	for _, platform := range platforms {
		// 获取该云平台下的所有区域
		var regions []model.CloudRegions
		if err := global.KUBEGALE_DB.Where("cloud_platform_id = ?", platform.ID).Find(&regions).Error; err != nil {
			global.KUBEGALE_LOG.Error("获取区域列表失败",
				zap.Error(err),
				zap.String("platform", platform.Name))
			continue
		}

		// 转换区域格式
		var regionList []model.Regions
		for _, region := range regions {
			regionList = append(regionList, model.Regions{
				ID:         region.RegionId,
				Name:       region.Name,
				RegionId:   region.RegionId,
				RegionName: region.RegionName,
			})
		}

		// 添加到树形结构
		list = append(list, model.PlatformTree{
			ID:     platform.ID,
			Name:   platform.Name,
			Region: regionList,
		})
	}

	global.KUBEGALE_LOG.Info("获取区域树形结构完成",
		zap.Int("platformCount", len(list)))

	return list, nil
}
