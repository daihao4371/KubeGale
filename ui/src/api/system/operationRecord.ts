import request from '../request'
import type { OperationRecordParams, SysOperationRecord } from '@/types/system'

// 获取操作记录列表
export const getOperationRecordList = (params: OperationRecordParams) => {
  return request({
    url: '/sysOperationRecord/getSysOperationRecordList',
    method: 'get',
    params
  })
}

// 删除操作记录
export const deleteOperationRecord = (ID: number) => {
  return request({
    url: '/sysOperationRecord/deleteSysOperationRecord',
    method: 'delete',
    data: { ID }
  })
}

// 批量删除操作记录
export const batchDeleteOperationRecord = (ids: number[]) => {
  return request({
    url: '/sysOperationRecord/deleteSysOperationRecordByIds',
    method: 'delete',
    data: { ids }
  })
}


