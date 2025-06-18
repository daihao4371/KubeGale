import request from '@/utils/request'
import type { ApiResponse } from '@/types/cmdb'

interface ExecuteCommandParams {
  hosts: string[]
  users: string[]
  ports: number[]
  command: string
  language: 'python' | 'shell'
}

interface ExecuteCommandResponse {
  allHosts: string[]
  successHosts: string[]
  failureHosts: string[] | null
  executionLogs: Array<{
    host: string
    output: string
    error?: string
  }>
  status: 'success' | 'partial_failure' | 'failed'
}

export interface ExecutionLog {
  ID: number
  UserId: number
  Command: string
  AllHosts: string
  SuccessHosts: string
  FailureHosts: string
  ExecutionLogs: string
  Status: string
  CreatedAt: string
}

export interface BatchExecuteCommandReq {
  hosts: string[]
  users: string[]
  ports: number[]
  command: string
  language: 'shell' | 'python'
}

export interface BatchExecuteCommandRes {
  allHosts: string[]
  successHosts: string[] | null
  failureHosts: string[] | null
  executionLogs: {
    host: string
    output: string
  }[]
  status: string
}

// 执行批量命令
export const executeCommand = (data: ExecuteCommandParams) => {
  return request<ApiResponse<ExecuteCommandResponse>>({
    url: '/cmdb/batchOperations/execute',
    method: 'post',
    data
  })
}

// 获取执行日志列表
export const getExecutionLogs = (params: { page: number; pageSize: number }) => {
  return request<ApiResponse<ExecutionLog[]>>({
    url: '/cmdb/batchOperations/execLogs',
    method: 'get',
    params
  })
}

export const batchExecuteCommand = (data: BatchExecuteCommandReq) => {
  return request<ApiResponse<BatchExecuteCommandRes>>({
    url: '/cmdb/batchOperations/execute',
    method: 'post',
    data
  })
}
