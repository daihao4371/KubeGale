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
func (r CloudRegionService) TencentRegion(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	region := tencent.NewRegion()
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            *instance.RegionName,
			RegionId:        "tencent-" + *instance.Region,
			RegionName:      *instance.RegionName,
			CloudPlatformId: cloud.ID,
		})
	}

	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// AliyunRegion 阿里云同步Region
func (r CloudRegionService) AliyunRegion(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	region := aliyun.NewRegion()
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            instance.LocalName,
			RegionId:        "aliyun-" + instance.RegionId,
			RegionName:      instance.LocalName,
			CloudPlatformId: cloud.ID,
		})
	}

	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// HuaweiRegion 华为云同步Region
func (r CloudRegionService) HuaweiRegion(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	region := huawei.NewRegion()
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            instance.Name,
			RegionId:        "huawei-" + instance.RegionId,
			RegionName:      instance.Name,
			CloudPlatformId: cloud.ID,
		})
	}

	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// AwsRegion 亚马逊云同步Region
func (r CloudRegionService) AwsRegion(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	region := aws.NewRegion()
	list, err := region.List(cloud.AccessKeyId, cloud.AccessKeySecret)
	if err != nil {
		return err
	}

	for _, instance := range list {
		regions = append(regions, model.CloudRegions{
			Name:            *instance.RegionName,
			RegionId:        "aws-" + *instance.RegionName,
			RegionName:      *instance.RegionName,
			CloudPlatformId: cloud.ID,
		})
	}

	if len(regions) > 0 {
		r.UpdateRegions(regions)
	}

	return err
}

// UpdateRegions 更新Region信息
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
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"region_id",
			}),
		}).Create(&region).Error; err != nil {
			global.KUBEGALE_LOG.Error("region messages update fail!", zap.Error(err))
			tx.Rollback()
		}

		// 插入不存在的记录
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&region).Error; err != nil {
			global.KUBEGALE_LOG.Error("region messages insert fail!", zap.Error(err))
			tx.Rollback()
		}

		// 提交事务
		tx.Commit()
	}
}

// SyncRegion 同步云平台所有Region
func (r CloudRegionService) SyncRegion(id int) (err error) {
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloud model.CloudPlatform
	if err := db.Where("id = ?", id).First(&cloud).Error; err != nil {
		return err
	}

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
