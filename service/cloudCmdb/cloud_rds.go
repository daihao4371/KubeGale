package cloudCmdb

import (
	"KubeGale/global"
	model "KubeGale/model/cloudCmdb"
	cloudcmdbreq "KubeGale/model/cloudCmdb/cloudcmdb"
	"KubeGale/model/common/request"
	"KubeGale/utils/cloudCmdb/aliyun"
	"KubeGale/utils/cloudCmdb/aws"
	"KubeGale/utils/cloudCmdb/huawei"
	"KubeGale/utils/cloudCmdb/tencent"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CloudRDSService struct{}

// List RDS列表
func (r *CloudRDSService) List(slb model.RDS, info cloudcmdbreq.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(model.RDS{})
	var slbList []model.RDS

	if info.Keyword != "" && info.Field != "" {
		db = db.Where(info.Field+" LIKE ?", "%"+info.Keyword+"%")
	}

	if slb.InstanceId != "" {
		db = db.Where("instance_id LIKE ?", "%"+slb.InstanceId+"%")
	}

	if slb.Name != "" {
		db = db.Where("name = ?", slb.Name)
	}

	if slb.Region != "" {
		db = db.Where("region = ?", slb.Region)
	}

	err = db.Count(&total).Error

	if err != nil {
		return slbList, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			orderMap := make(map[string]bool, 5)
			orderMap["id"] = true
			orderMap["instance_id"] = true
			orderMap["name"] = true
			orderMap["status"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return slbList, total, err
			}

			err = db.Order(OrderStr).Find(&slbList).Error
		} else {
			err = db.Order("id").Find(&slbList).Error
		}
	}

	return slbList, total, err
}

// UpdateRDS 更新RDS信息
func (r *CloudRDSService) UpdateRDS(list []model.RDS) {
	db := global.KUBEGALE_DB.Model(model.RDS{})

	for _, machine := range list {
		// 开始事务
		tx := db.Begin()

		// 更新所有存在的记录，忽略不存在的记录
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"instance_id",
				"private_addr",
				"public_addr",
				"region",
				"region_name",
				"status",
				"creation_time",
				"cloud_platform_id",
			}),
		}).Create(&machine).Error; err != nil {
			global.KUBEGALE_LOG.Error("rds  messages update fail!", zap.Error(err))
			tx.Rollback()
		}

		// 插入不存在的记录
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&machine).Error; err != nil {
			global.KUBEGALE_LOG.Error("rds messages insert fail!", zap.Error(err))
			tx.Rollback()
		}

		// 提交事务
		tx.Commit()
	}
}

// AliyunSyncRDS 阿里云同步RDS
func (r *CloudRDSService) AliyunSyncLoadRDS(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("aliyun ecs list get fail: %s", err))
				}
			}()

			rds := aliyun.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("aliyun LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				r.UpdateRDS(list)
			}

		}(region)
	}

	return err
}

// TencentSyncRDS 腾讯云同步RDS
func (r *CloudRDSService) TencentSyncRDS(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("Tencent slb list get fail: %s", err))
				}
			}()

			ecs := tencent.NewRDS()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("Tencent LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				r.UpdateRDS(list)
			}

		}(region)
	}
	return err
}

// HuaweiSyncRDS 华为云同步RDS
func (r *CloudRDSService) HuaweiSyncRDS(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("huawei rds list get fail: %s", err))
				}
			}()

			rds := huawei.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("huawei rds list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				r.UpdateRDS(list)
			}

		}(region)
	}
	return err
}

// AwsSyncRDS 亚马逊云同步RDS
func (r *CloudRDSService) AwsSyncRDS(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("aws rds list get fail: %s", err))
				}
			}()

			rds := aws.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("aws RDS list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				r.UpdateRDS(list)
			}
		}(region)
	}

	return err
}

// SyncRDS 同步各个厂商的RDS
func (r *CloudRDSService) SyncRDS(id int) (err error) {
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloud model.CloudPlatform
	if err := db.Where("id = ?", id).First(&cloud).Error; err != nil {
		return err
	}

	if cloud.Platform == "aliyun" {
		if err = r.AliyunSyncLoadRDS(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "tencent" {
		if err = r.TencentSyncRDS(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "huawei" {
		if err = r.HuaweiSyncRDS(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "aws" {
		if err = r.AwsSyncRDS(cloud); err != nil {
			return err
		}
	}

	return err
}

// RDSTree rds目录树
func (r *CloudRDSService) RDSTree(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list interface{}, err error) {
	info.PageSize, info.Page = 1000, 1
	var platformTree []model.PlatformTree
	var platform CloudPlatformService
	platformList, _, err := platform.List(cloud, info, order, desc)
	if err != nil {
		return nil, err
	}

	for _, pt := range platformList {
		var rdslist []model.RDS
		var regions []model.Regions
		if err := global.KUBEGALE_DB.Table("cloud_rds").Select("DISTINCT region, region_name").Where("cloud_platform_id = ?", pt.ID).Find(&rdslist).Error; err != nil {
			global.KUBEGALE_LOG.Error("ecs machine DISTINCT fail!", zap.Error(err))
			return nil, err
		}

		if len(rdslist) > 0 {
			for _, vmRegion := range rdslist {
				regions = append(regions, model.Regions{
					ID:         vmRegion.Region,
					Name:       vmRegion.RegionName,
					RegionId:   vmRegion.Region,
					RegionName: vmRegion.RegionName,
				})
			}
		}

		platformTree = append(platformTree, model.PlatformTree{ID: pt.ID, Name: pt.Name, Region: regions})
	}

	return platformTree, err
}
