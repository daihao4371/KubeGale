import request from '@/utils/request'
import type { Host } from '@/types/cmdb'

interface ApiResponse<T = Record<string, never>> {
  code: number
  data: T
  msg: string
}

interface ListResponse<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

interface AuthenticationResponse {
  sessionID: string
  wsUrl: string
}

// 获取主机列表
export const getHostList = (params: {
  project?: number
  page: number
  pageSize: number
}) => {
  return request<ApiResponse<ListResponse<Host>>>({
    url: '/cmdb/hostsList',
    method: 'post',
    data: params
  })
}

// 创建主机
export const createHost = (data: {
  name: string
  serverHost: string
  port: number
  username: string
  password: string
  project: number
  note?: string
}) => {
  return request<ApiResponse>({
    url: '/cmdb/hosts',
    method: 'post',
    data
  })
}

// 更新主机
export const updateHost = (data: {
  ID: number
  name?: string
  serverHost?: string
  port?: number
  username?: string
  password?: string
  project?: number
  note?: string
}) => {
  return request<ApiResponse>({
    url: '/cmdb/hosts',
    method: 'put',
    data
  })
}

// 删除主机
export const deleteHost = (data: { id: number }) => {
  return request<ApiResponse>({
    url: '/cmdb/hosts',
    method: 'delete',
    data
  })
}

// 批量删除主机
export const batchDeleteHosts = (data: { ids: number[] }) => {
  return request<ApiResponse>({
    url: '/cmdb/hostsByIds',
    method: 'delete',
    data
  })
}

// 获取主机详情
export const getHostDetail = (id: number) => {
  return request<ApiResponse<Host>>({
    url: '/cmdb/hostsById',
    method: 'get',
    params: { id }
  })
}

// 批量导入主机
export const importHosts = (data: FormData) => {
  return request<ApiResponse>({
    url: '/cmdb/hosts/import',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// SSH 认证
export const authenticateHost = (data: {
  name: string
  serverHost: string
  port: number
  username: string
  password: string
  project: number
  note: string
}) => {
  return request<ApiResponse<AuthenticationResponse>>({
    url: '/cmdb/hosts/authentication',
    method: 'post',
    data
  })
}

// 获取终端 WebSocket 连接
export const getTerminalWebSocket = (wsUrl: string) => {
  return new WebSocket(wsUrl)
}
