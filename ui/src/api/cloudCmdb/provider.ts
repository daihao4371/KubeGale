import request from '@/utils/request'

interface ApiResponse<T> {
  code: number
  data: T
  msg: string
}

interface CloudPlatformList {
  list: CloudPlatform[]
  total: number
  page: number
  pageSize: number
}

interface CloudPlatform {
  id: number
  name: string
  platform: string
  access_key_id: string
  access_key_secret: string
  created_at: string
  updated_at: string
}

// 获取云厂商列表
export function getCloudPlatformList(data: {
  page: number
  pageSize: number
  name?: string
  type?: string
}): Promise<ApiResponse<CloudPlatformList>> {
  return request({
    url: '/cloud_platform/list',
    method: 'post',
    data
  })
}

// 新增云厂商
export function addCloudPlatform(data: {
  name: string
  platform: string
  access_key_id: string
  access_key_secret: string
}) {
  return request({
    url: '/cloud_platform/add',
    method: 'post',
    data
  })
}

// 编辑云厂商
export function updateCloudPlatform(data: {
  id: number
  name: string
  platform: string
  access_key_id: string
  access_key_secret: string
}) {
  return request({
    url: '/cloud_platform/update',
    method: 'post',
    data
  })
}

// 删除云厂商
export function deleteCloudPlatform(id: number) {
  return request({
    url: '/cloud_platform/delete',
    method: 'post',
    data: { id }
  })
}

// 同步区域
export function syncRegions(id: number) {
  return request({
    url: '/cloud_platform/sync_regions',
    method: 'post',
    data: { id }
  })
} 