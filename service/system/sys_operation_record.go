package system

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/system"
	systemReq "KubeGale/model/system/request"
)

type OperationRecordService struct{}

var OperationRecordServiceApp = new(OperationRecordService)

// 创建记录
func (operationRecordService *OperationRecordService) CreateSysOperationRecord(sysOperationRecord system.SysOperationRecord) (err error) {
	err = global.KUBEGALE_DB.Create(&sysOperationRecord).Error
	return err
}

// 批量删除记录
func (operationRecordService *OperationRecordService) DeleteSysOperationRecordByIds(ids request.IdsReq) (err error) {
	err = global.KUBEGALE_DB.Delete(&[]system.SysOperationRecord{}, "id in (?)", ids.Ids).Error
	return err
}

// 删除操作记录
func (operationRecordService *OperationRecordService) DeleteSysOperationRecord(sysOperationRecord system.SysOperationRecord) (err error) {
	err = global.KUBEGALE_DB.Delete(&sysOperationRecord).Error
	return err
}

// 根据id获取单条操作记录
func (operationRecordService *OperationRecordService) GetSysOperationRecord(id uint) (sysOperationRecord system.SysOperationRecord, err error) {
	err = global.KUBEGALE_DB.Where("id = ?", id).First(&sysOperationRecord).Error
	return
}

// 分页获取操作记录列表
func (operationRecordService *OperationRecordService) GetSysOperationRecordInfoList(info systemReq.SysOperationRecordSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.KUBEGALE_DB.Model(&system.SysOperationRecord{})
	var sysOperationRecords []system.SysOperationRecord
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Method != "" {
		db = db.Where("method = ?", info.Method)
	}
	if info.Path != "" {
		db = db.Where("path LIKE ?", "%"+info.Path+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	return sysOperationRecords, total, err
}
