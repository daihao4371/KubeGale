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

type CloudVirtualMachineService struct{}

// List 云主机列表
func (vm CloudVirtualMachineService) List(machine model.VirtualMachine, info cloudcmdbreq.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.KUBEGALE_DB.Model(model.VirtualMachine{})
	var vmList []model.VirtualMachine

	if info.Keyword != "" && info.Field != "" {
		db = db.Where(info.Field+" LIKE ?", "%"+info.Keyword+"%")
	}

	if machine.Name != "" {
		db = db.Where("name LIKE ?", "%"+machine.Name+"%")
	}

	if machine.InstanceId != "" {
		db = db.Where("instance_id = ?", machine.InstanceId)
	}

	if machine.Region != "" {
		db = db.Where("region = ?", machine.Region)
	}

	err = db.Count(&total).Error

	if err != nil {
		return vmList, total, err
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
			orderMap["os_type"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", order)
				return vmList, total, err
			}

			err = db.Order(OrderStr).Find(&vmList).Error
		} else {
			err = db.Order("id").Find(&vmList).Error
		}
	}

	return vmList, total, err
}

// TencentSyncMachine 腾讯云同步主机
func (vm *CloudVirtualMachineService) TencentSyncMachine(cloud model.CloudPlatform) (err error) {
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		return err
	}

	for _, region := range regions {
		go func(region model.CloudRegions) {
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("tencent ecs list get fail: %s", err))
				}
			}()

			ecs := tencent.NewECS()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("tencent ecs list get  fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				vm.UpdateMachine(list)
			}

		}(region)
	}

	return err
}

// HuaweiSyncMachine 华为云同步主机
func (vm *CloudVirtualMachineService) HuaweiSyncMachine(cloud model.CloudPlatform) (err error) {
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

			ecs := huawei.NewECS()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("huawei ecs list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				vm.UpdateMachine(list)
			}

		}(region)
	}

	return err
}

// AliyunSyncMachine 阿里云同步主机
func (vm *CloudVirtualMachineService) AliyunSyncMachine(cloud model.CloudPlatform) (err error) {
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

			ecs := aliyun.NewECS()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error("aliyun ecs list get fail: ", zap.Error(err))
				return
			}

			if len(list) > 0 {
				vm.UpdateMachine(list)
			}

		}(region)
	}

	return err
}

// AwsSyncMachine 亚马逊云同步主机
// 功能：同步亚马逊云平台的所有区域下的主机信息
// 参数：
//   - cloud: 云平台信息，包含访问凭证和平台ID
//
// 返回：
//   - err: 错误信息
func (vm *CloudVirtualMachineService) AwsSyncMachine(cloud model.CloudPlatform) (err error) {
	// 获取云平台所有区域信息
	var regions []model.CloudRegions
	if err = global.KUBEGALE_DB.Where("cloud_platform_id = ?", cloud.ID).Find(&regions).Error; err != nil {
		global.KUBEGALE_LOG.Error("aws get regions fail!", zap.Error(err))
		return err
	}

	// 并发处理每个区域
	for _, region := range regions {
		go func(region model.CloudRegions) {
			// 错误恢复机制：捕获并记录可能的panic
			defer func() {
				if err := recover(); err != nil {
					global.KUBEGALE_LOG.Error(fmt.Sprintf("aws ecs list get fail in region %s: %s", region.RegionId, err))
				}
			}()

			// 创建 ECS 实例并获取列表
			// 使用云平台凭证和区域信息初始化 ECS 客户端
			ecs := aws.NewECS()
			list, err := ecs.List(cloud.ID, region, cloud.AccessKeyId, cloud.AccessKeySecret)
			if err != nil {
				global.KUBEGALE_LOG.Error(fmt.Sprintf("aws ecs list get fail in region %s: ", region.RegionId), zap.Error(err))
				return
			}

			// 更新数据库：如果有数据则进行更新
			if len(list) > 0 {
				vm.UpdateMachine(list)
			}
		}(region)
	}

	return err
}

// UpdateMachine 更新主机信息
// 功能：批量更新或插入主机信息到数据库
// 参数：
//   - list: 主机信息列表
func (vm *CloudVirtualMachineService) UpdateMachine(list []model.VirtualMachine) {
	db := global.KUBEGALE_DB.Model(model.VirtualMachine{})

	for _, machine := range list {
		// 开始事务：确保数据一致性
		tx := db.Begin()

		// 更新所有存在的记录，忽略不存在的记录
		// 使用 OnConflict 子句处理冲突情况
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"instance_id",
				"cpu",
				"memory",
				"os",
				"os_type",
				"private_addr",
				"public_addr",
				"region",
				"region_name",
				"status",
				"creation_time",
				"expired_time",
				"cloud_platform_id",
			}),
		}).Create(&machine).Error; err != nil {
			global.KUBEGALE_LOG.Error("ecs machine messages update fail!", zap.Error(err))
			tx.Rollback()
		}

		// 插入不存在的记录
		// 使用 OnConflict 子句的 DoNothing 选项处理插入冲突
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&machine).Error; err != nil {
			global.KUBEGALE_LOG.Error("ecs machine messages insert fail!", zap.Error(err))
			tx.Rollback()
		}

		// 提交事务：确认所有操作成功完成
		tx.Commit()
	}
}

// SyncMachine 同步各个厂商的云主机
// 功能：根据云平台类型选择对应的同步方法
// 参数：
//   - id: 云平台ID
//
// 返回：
//   - err: 错误信息
func (vm *CloudVirtualMachineService) SyncMachine(id int) (err error) {
	// 获取云平台信息
	db := global.KUBEGALE_DB.Model(model.CloudPlatform{})
	var cloud model.CloudPlatform
	if err := db.Where("id = ?", id).First(&cloud).Error; err != nil {
		return err
	}

	// 根据云平台类型选择对应的同步方法
	if cloud.Platform == "aliyun" {
		if err = vm.AliyunSyncMachine(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "tencent" {
		if err = vm.TencentSyncMachine(cloud); err != nil {
			return err
		}
	}

	if cloud.Platform == "huawei" {
		if err = vm.HuaweiSyncMachine(cloud); err != nil {
			return err
		}
	}

	// 修复：将 huawei 改为 aws
	if cloud.Platform == "aws" {
		if err = vm.AwsSyncMachine(cloud); err != nil {
			return err
		}
	}

	return err
}

// MachineTree  主机厂商目录树
func (vm *CloudVirtualMachineService) MachineTree(cloud model.CloudPlatform, info request.PageInfo, order string, desc bool) (list interface{}, err error) {
	info.PageSize, info.Page = 1000, 1
	var platformTree []model.PlatformTree
	var platform CloudPlatformService
	platformList, _, err := platform.List(cloud, info, order, desc)
	if err != nil {
		return nil, err
	}

	for _, pt := range platformList {
		var vmlist []model.VirtualMachine
		var regions []model.Regions
		if err := global.KUBEGALE_DB.Table("cloud_virtual_machine").Select("DISTINCT region, region_name").Where("cloud_platform_id = ?", pt.ID).Find(&vmlist).Error; err != nil {
			global.KUBEGALE_LOG.Error("ecs machine DISTINCT fail!", zap.Error(err))
			return nil, err
		}

		if len(vmlist) > 0 {
			for _, vmRegion := range vmlist {
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
