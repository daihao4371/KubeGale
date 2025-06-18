import request from '@/utils/request'
import type { ApiResponse } from '@/types/cloudCmdb'
import type { CloudProvider, CloudProviderListResponse, CloudProviderDetailResponse } from '@/types/cloudCmdb'

interface ListParams {
  page: number
  pageSize: number
  name?: string
  type?: string
}

interface IdParams {
  id: number
}

interface IdsParams {
  ids: number[]
}

// 获取云厂商列表
export const cloudplatformlist = (params: ListParams): Promise<ApiResponse<CloudProviderListResponse>> => {
  return request<ApiResponse<CloudProviderListResponse>>({
    url: '/cloud_platform/list',
    method: 'post',
    data: params
  })
}

// 获取云厂商详情
export const cloudplatformById = (data: { id: number }): Promise<ApiResponse<CloudProviderDetailResponse>> => {
  return request<ApiResponse<CloudProviderDetailResponse>>({
    url: '/cloud_platform/getById',
    method: 'post',
    data
  })
}

// 创建云厂商
export const cloudplatformCreate = (data: Omit<CloudProvider, 'id' | 'created_at'>): Promise<ApiResponse<null>> => {
  return request<ApiResponse<null>>({
    url: '/cloud_platform/create',
    method: 'post',
    data
  })
}

// 更新云厂商
export const cloudplatformUpdate = (data: {
  id: number
  name: string
  type: string
  accessKey: string
  secretKey: string
}): Promise<ApiResponse<null>> => {
  return request<ApiResponse<null>>({
    url: '/cloud_platform/update',
    method: 'put',
    data
  })
}

// 删除云厂商
export const cloudplatformDelete = (data: { id: number }): Promise<ApiResponse<null>> => {
  return request<ApiResponse<null>>({
    url: '/cloud_platform/delete',
    method: 'delete',
    data
  })
}

// 批量删除云厂商
export const cloudplatformDeleteByIds = (data: { ids: number[] }): Promise<ApiResponse<null>> => {
  return request<ApiResponse<null>>({
    url: '/cloud_platform/deleteByIds',
    method: 'delete',
    data
  })
}