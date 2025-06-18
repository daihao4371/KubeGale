import request from '@/utils/request'
import type { ApiResponse, RdsInstance, RdsListResponse } from '@/types/cloudCmdb'

// RDS列表请求参数
export interface RdsListParams {
  page: number
  pageSize: number
  name?: string
  instance_id?: string
  region?: string
  status?: string
}

// 获取RDS列表
export const getRdsList = (params: RdsListParams): Promise<ApiResponse<RdsListResponse>> => {
  return request<ApiResponse<RdsListResponse>>({
    url: '/rds/list',
    method: 'post',
    data: params
  })
}

// 创建RDS实例
export const createRds = (data: Omit<RdsInstance, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<null>> => {
  return request<ApiResponse<null>>({
    url: '/cloud_rds/create',
    method: 'post',
    data
  })
}

// 获取单个RDS实例信息
export const getRdsInstance = (params: { name?: string; instanceId?: string }): Promise<ApiResponse<RdsInstance>> => {
  return request<ApiResponse<RdsInstance>>({
    url: '/rds/get',
    method: 'post',
    data: params
  })
}

// 删除RDS实例
export const deleteRds = (id: number): Promise<ApiResponse<null>> => {
  return request<ApiResponse<null>>({
    url: '/cloud_rds/delete',
    method: 'delete',
    data: { id }
  })
}

// 同步RDS数据
export const syncRds = (data: { id: number }) => {
  return request<ApiResponse<null>>({
    url: '/rds/sync',
    method: 'post',
    data
  })
}

// 区域接口
export interface Region {
  id: string
  name: string
  region_id: string
  region_name: string
}

// 云平台树形结构响应接口
export interface CloudPlatformTreeItem {
  id: number
  name: string
  region: string | null
}

export interface CloudPlatformTreeResponse {
  list: CloudPlatformTreeItem[]
  total: number
  page: number
  pageSize: number
}

// 获取云平台树形结构
export const getCloudPlatformTree = (): Promise<ApiResponse<CloudPlatformTreeResponse>> => {
  return request<ApiResponse<CloudPlatformTreeResponse>>({
    url: '/rds/tree',
    method: 'post',
    data: {
      page: 1,
      pageSize: 1000
    }
  })
}

// 平台-区域树类型
export interface RegionNode {
  id: string
  name: string
  region_id: string
  region_name: string
}

export interface PlatformRegionTreeItem {
  id: number
  name: string
  region: RegionNode[]
}

// 获取平台-区域树
export const getPlatformRegionTree = () => {
  return request<{ code: number, data: PlatformRegionTreeItem[], msg: string }>({
    url: '/cloud_region/tree',
    method: 'get'
  })
}
