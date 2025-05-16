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
// 功能：获取RDS实例列表，支持分页、排序和条件筛选
// 参数：
//   - slb: RDS查询条件
//   - info: 分页信息
//   - order: 排序字段
//   - desc: 是否降序排序
//
// 返回：
//   - list: RDS实例列表
//   - total: 总记录数
//   - err: 错误信息
func (r *CloudRDSService) List(slb model.RDS, info cloudcmdbreq.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	// 计算分页参数
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(model.RDS{})
	var slbList []model.RDS

	// 添加关键字搜索条件
	if info.Keyword != "" && info.Field != "" {
		db = db.Where(info.Field+" LIKE ?", "%"+info.Keyword+"%")
	}

	// 添加实例ID查询条件
	if slb.InstanceId != "" {
		db = db.Where("instance_id LIKE ?", "%"+slb.InstanceId+"%")
	}

	// 添加名称查询条件
	if slb.Name != "" {
		db = db.Where("name = ?", slb.Name)
	}

	// 添加区域查询条件
	if slb.Region != "" {
		db = db.Where("region = ?", slb.Region)
	}

	// 获取总记录数
	err = db.Count(&total).Error

	if err != nil {
		return slbList, total, err
	} else {
		// 设置分页
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

			// 执行排序查询
			err = db.Order(OrderStr).Find(&slbList).Error
		} else {
			// 默认按ID排序
			err = db.Order("id").Find(&slbList).Error
		}
	}

	return slbList, total, err
}

// UpdateRDS 更新RDS信息
// 功能：批量更新或插入RDS实例信息到数据库
// 参数：
//   - list: RDS实例信息列表
func (r *CloudRDSService) UpdateRDS(list []model.RDS) {
	db := global.KUBEGALE_DB.Model(model.RDS{})

	for _, machine := range list {
		// 开始事务：确保数据一致性
		tx := db.Begin()

		// 更新所有存在的记录，忽略不存在的记录
		// 使用 OnConflict 子句处理冲突情况
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
			global.KUBEGALE_LOG.Error("rds messages update fail!", zap.Error(err))
			tx.Rollback()
		}

		// 插入不存在的记录
		// 使用 OnConflict 子句的 DoNothing 选项处理插入冲突
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&machine).Error; err != nil {
			global.KUBEGALE_LOG.Error("rds messages insert fail!", zap.Error(err))
			tx.Rollback()
		}

		// 提交事务：确认所有操作成功完成
		tx.Commit()
	}
}

// AliyunSyncRDS 阿里云同步RDS
// 功能：同步阿里云平台的所有区域下的RDS实例信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r *CloudRDSService) AliyunSyncLoadRDS(cloud model.CloudPlatform) (err error) {
	// 获取云平台所有区域信息
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	// 并发处理每个区域
	for _, region := range regions {
		go func(region model.CloudRegions) {
			// 错误恢复机制：捕获并记录可能的panic
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("aliyun ecs list get fail: %s", err))
				}
			}()

			// 创建RDS实例并获取列表
			rds := aliyun.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("aliyun LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			// 更新数据库：如果有数据则进行更新
			if len(list) > 0 {
				r.UpdateRDS(list)
			}
		}(region)
	}

	return err
}

// TencentSyncRDS 腾讯云同步RDS
// 功能：同步腾讯云平台的所有区域下的RDS实例信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r *CloudRDSService) TencentSyncRDS(cloud model.CloudPlatform) (err error) {
	// 获取云平台所有区域信息
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	// 并发处理每个区域
	for _, region := range regions {
		go func(region model.CloudRegions) {
			// 错误恢复机制：捕获并记录可能的panic
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("Tencent slb list get fail: %s", err))
				}
			}()

			// 创建RDS实例并获取列表
			ecs := tencent.NewRDS()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("Tencent LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			// 更新数据库：如果有数据则进行更新
			if len(list) > 0 {
				r.UpdateRDS(list)
			}
		}(region)
	}
	return err
}

// HuaweiSyncRDS 华为云同步RDS
// 功能：同步华为云平台的所有区域下的RDS实例信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r *CloudRDSService) HuaweiSyncRDS(cloud model.CloudPlatform) (err error) {
	// 获取云平台所有区域信息
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	// 并发处理每个区域
	for _, region := range regions {
		go func(region model.CloudRegions) {
			// 错误恢复机制：捕获并记录可能的panic
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("huawei rds list get fail: %s", err))
				}
			}()

			// 创建RDS实例并获取列表
			rds := huawei.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("huawei rds list get fail: ", zap.Error(err))
				return
			}

			// 更新数据库：如果有数据则进行更新
			if len(list) > 0 {
				r.UpdateRDS(list)
			}
		}(region)
	}
	return err
}

// AwsSyncRDS 亚马逊云同步RDS
// 功能：同步亚马逊云平台的所有区域下的RDS实例信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (r *CloudRDSService) AwsSyncRDS(cloud model.CloudPlatform) (err error) {
	// 获取云平台所有区域信息
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	// 并发处理每个区域
	for _, region := range regions {
		go func(region model.CloudRegions) {
			// 错误恢复机制：捕获并记录可能的panic
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("aws rds list get fail: %s", err))
				}
			}()

			// 创建RDS实例并获取列表
			rds := aws.NewRDS()
			list, err := rds.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("aws RDS list get fail: ", zap.Error(err))
				return
			}

			// 更新数据库：如果有数据则进行更新
			if len(list) > 0 {
				r.UpdateRDS(list)
			}
		}(region)
	}

	return err
}

// SyncRDS 同步各个厂商的RDS
// 功能：根据云平台类型选择对应的同步方法
// 参数：
//   - id: 云平台ID
//
// 返回：
//   - err: 错误信息
func (r *CloudRDSService) SyncRDS(id int) (err error) {
	// 获取云平台信息
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloud model.CloudPlatform
	if err := db.Where("id = ?", id).First(&cloud).Error; err != nil {
		return err
	}

	// 根据云平台类型选择对应的同步方法
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

// RDSTree RDS目录树
// 功能：获取RDS实例的目录树结构，按云平台和区域组织
// 参数：
//   - cloud: 云平台信息
//   - info: 分页信息
//   - order: 排序字段
//   - desc: 是否降序排序
//
// 返回：
//   - list: 目录树结构
//   - err: 错误信息
func (r *CloudRDSService) RDSTree(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list interface{}, err error) {
	// 设置分页参数
	info.PageSize, info.Page = 1000, 1
	var platformTree []model.PlatformTree
	var platform CloudPlatformService
	platformList, _, err := platform.List(cloud, info, order, desc)
	if err != nil {
		return nil, err
	}

	// 构建目录树结构
	for _, pt := range platformList {
		var rdslist []model.RDS
		var regions []model.Regions
		// 获取区域信息
		if err := global.KUBEGALE_DB.Table("cloud_rds").Select("DISTINCT region, region_name").Where("cloud_platform_id = ?", pt.ID).Find(&rdslist).Error; err != nil {
			global.KUBEGALE_LOG.Error("ecs machine DISTINCT fail!", zap.Error(err))
			return nil, err
		}

		// 处理区域信息
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

		// 添加到目录树
		platformTree = append(platformTree, model.PlatformTree{ID: pt.ID, Name: pt.Name, Region: regions})
	}

	return platformTree, err
}
