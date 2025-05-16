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

type CloudLoadBalancerService struct{}

// List 负载均衡列表
// 功能：获取负载均衡器列表，支持分页、排序和条件筛选
// 参数：
//   - slb: 负载均衡器查询条件
//   - info: 分页信息
//   - order: 排序字段
//   - desc: 是否降序排序
//
// 返回：
//   - list: 负载均衡器列表
//   - total: 总记录数
//   - err: 错误信息
func (l *CloudLoadBalancerService) List(slb model.LoadBalancer, info cloudcmdbreq.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(model.LoadBalancer{})
	var slbList []model.LoadBalancer

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

// UpdateLoadBalancer 更新负载均衡信息
// 功能：批量更新或插入负载均衡器信息到数据库
// 参数：
//   - list: 负载均衡器信息列表
func (l *CloudLoadBalancerService) UpdateLoadBalancer(list []model.LoadBalancer) {
	db := global.KUBEGALE_DB.Model(model.LoadBalancer{})

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
				"bandwidth",
				"region",
				"region_name",
				"status",
				"creation_time",
				"cloud_platform_id",
			}),
		}).Create(&machine).Error; err != nil {
			global.KUBEGALE_LOG.Error("LoadBalancer  messages update fail!", zap.Error(err))
			tx.Rollback()
		}

		// 插入不存在的记录
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&machine).Error; err != nil {
			global.KUBEGALE_LOG.Error("LoadBalancer messages insert fail!", zap.Error(err))
			tx.Rollback()
		}

		// 提交事务
		tx.Commit()
	}
}

// AliyunSyncLoadBalancer 阿里云同步负载均衡
// 功能：同步阿里云平台的所有区域下的负载均衡器信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (l *CloudLoadBalancerService) AliyunSyncLoadBalancer(cloud model.CloudPlatform) (err error) {
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

			ecs := aliyun.NewLoadBalancer()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("aliyun LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				l.UpdateLoadBalancer(list)
			}

		}(region)
	}

	return err
}

// TencentSyncLoadBalancer 腾讯云同步负载均衡
// 功能：同步腾讯云平台的所有区域下的负载均衡器信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (l *CloudLoadBalancerService) TencentSyncLoadBalancer(cloud model.CloudPlatform) (err error) {
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

			ecs := tencent.NewLoadBalancer()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("Tencent LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				l.UpdateLoadBalancer(list)
			}

		}(region)
	}
	return err
}

// HuaweiSyncLoadBalancer 华为云同步负载均衡
// 功能：同步华为云平台的所有区域下的负载均衡器信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (l *CloudLoadBalancerService) HuaweiSyncLoadBalancer(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("huawei ecs list get fail: %s", err))
				}
			}()

			ecs := huawei.NewLoadBalancer()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("huawei LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				l.UpdateLoadBalancer(list)
			}

		}(region)
	}
	return err
}

// AwsSyncLoadBalancer 亚马逊云同步负载均衡
// 功能：同步亚马逊云平台的所有区域下的负载均衡器信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (l *CloudLoadBalancerService) AwsSyncLoadBalancer(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("aws load balancer list get fail: %s", err))
				}
			}()

			lb := aws.NewLoadBalancer()
			list, err := lb.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("aws LoadBalancer list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				l.UpdateLoadBalancer(list)
			}
		}(region)
	}

	return err
}

// SyncLoadBalancer 同步各个厂商的负载均衡
// 功能：根据云平台类型选择对应的同步方法
// 参数：
//   - id: 云平台ID
//
// 返回：
//   - err: 错误信息
func (l *CloudLoadBalancerService) SyncLoadBalancer(id int) (err error) {
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloud model.CloudPlatform
	if err := db.Where("id = ?", id).First(&cloud).Error; err != nil {
		return err
	}

	if cloud.Platform == "aliyun" {
		if err = l.AliyunSyncLoadBalancer(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "tencent" {
		if err = l.TencentSyncLoadBalancer(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "huawei" {
		if err = l.HuaweiSyncLoadBalancer(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "aws" {
		if err = l.AwsSyncLoadBalancer(cloud); err != nil {
			return err
		}
	}

	return err
}

// LoadBalancerTree 负载均衡目录树
// 功能：获取负载均衡器的目录树结构，按云平台和区域组织
// 参数：
//   - cloud: 云平台信息
//   - info: 分页信息
//   - order: 排序字段
//   - desc: 是否降序排序
//
// 返回：
//   - list: 目录树结构
//   - err: 错误信息
func (l *CloudLoadBalancerService) LoadBalancerTree(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list interface{}, err error) {
	info.PageSize, info.Page = 1000, 1
	var platformTree []model.PlatformTree
	var platform CloudPlatformService
	platformList, _, err := platform.List(cloud, info, order, desc)
	if err != nil {
		return nil, err
	}

	for _, pt := range platformList {
		var slblist []model.LoadBalancer
		var regions []model.Regions
		if err := global.KUBEGALE_DB.Table("cloud_load_balancer").Select("DISTINCT region, region_name").Where("cloud_platform_id = ?", pt.ID).Find(&slblist).Error; err != nil {
			global.KUBEGALE_LOG.Error("ecs machine DISTINCT fail!", zap.Error(err))
			return nil, err
		}

		if len(slblist) > 0 {
			for _, vmRegion := range slblist {
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
