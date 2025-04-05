package system

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/system"
	systemReq "KubeGale/model/system/request"
	"encoding/json"
)

type OperationRecordService struct{}

var OperationRecordServiceApp = new(OperationRecordService)

// 创建记录
func (operationRecordService *OperationRecordService) CreateSysOperationRecord(sysOperationRecord system.SysOperationRecord) (err error) {
	// 如果有用户ID，查询用户信息并设置操作人姓名
	if sysOperationRecord.UserID > 0 {
		var user system.User
		if err := global.KUBEGALE_DB.Where("id = ?", sysOperationRecord.UserID).First(&user).Error; err == nil {
			// 将用户信息设置到操作记录的User字段
			sysOperationRecord.User = user

			// 在数据库中创建记录时，User字段不会被保存
			// 但我们可以在这里设置操作人信息，方便后续查询时使用
			global.KUBEGALE_LOG.Info("记录操作，操作人: " + user.Username + "(" + user.RealName + ")")
		} else {
			global.KUBEGALE_LOG.Warn("查询用户信息失败，无法关联操作人")
		}
	}

	// 创建操作记录
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
	// 预加载用户信息
	err = global.KUBEGALE_DB.Where("id = ?", id).Preload("User").First(&sysOperationRecord).Error

	// 设置操作人信息
	if err == nil && sysOperationRecord.User.ID > 0 {
		sysOperationRecord.OperatorName = sysOperationRecord.User.Username
		sysOperationRecord.OperatorRealName = sysOperationRecord.User.RealName
	}

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

	// 查询操作记录并预加载用户信息
	err = db.Order("id desc").Limit(limit).Offset(offset).Preload("User").Find(&sysOperationRecords).Error
	if err != nil {
		return nil, 0, err
	}

	// 处理操作人信息
	var result []map[string]interface{}
	for _, record := range sysOperationRecords {
		recordMap := make(map[string]interface{})

		// 将记录转换为map
		recordBytes, _ := json.Marshal(record)
		json.Unmarshal(recordBytes, &recordMap)

		// 添加操作人信息
		if record.User.ID > 0 {
			recordMap["operator_name"] = record.User.Username
			recordMap["operator_real_name"] = record.User.RealName
		}

		result = append(result, recordMap)
	}

	return result, total, err
}
