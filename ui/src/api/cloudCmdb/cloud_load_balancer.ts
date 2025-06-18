import request from '@/utils/request'

// 负载均衡接口参数类型
interface LoadBalancerParams {
  page?: number
  pageSize?: number
  name?: string
  provider?: string
}

// 负载均衡表单数据类型
interface LoadBalancerForm {
  name: string
  provider: string
  region: string
}

// 负载均衡列表请求参数
export interface LoadBalancerListParams {
  loadBalancer: {
    name?: string
    instanceId?: string
    region?: string
  }
  page: number
  pageSize: number
  keyword?: string
  field?: string
  orderKey?: 'id' | 'instance_id' | 'name' | 'status'
  desc?: boolean
}

// 负载均衡列表项类型
export interface LoadBalancerItem {
  id: number
  name: string
  instance_id: string
  private_addr: string
  public_addr: string
  bandwidth: string
  region: string
  region_name: string
  status: string
  creation_time: string
  expired_time: string
  cloud_platform_id: number
  cloud_platform: {
    id: number
    name: string
    access_key_id: string
    access_key_secret: string
    platform: string
    created_at: string
    updated_at: string
  }
  created_at: string
  updated_at: string
}

// 负载均衡列表响应类型
export interface LoadBalancerListResponse {
  code: number
  data: {
    list: LoadBalancerItem[]
    total: number
    page: number
    pageSize: number
  }
  msg: string
}

// 负载均衡树形结构响应类型
export interface LoadBalancerTreeResponse {
  code: number
  data: {
    list: Array<{
      id: number
      name: string
      region: string | null
    }>
    total: number
    page: number
    pageSize: number
  }
  msg: string
}

// 同步响应类型
export interface SyncResponse {
  code: number
  data: Record<string, never>
  msg: string
}

// 获取负载均衡列表
export const getLoadBalancerList = (params: LoadBalancerListParams) => {
  return request<LoadBalancerListResponse>({
    url: '/loadBalancer/list',
    method: 'post',
    data: params
  })
}

// 获取负载均衡树形结构
export function getLoadBalancerTree() {
  return request<LoadBalancerTreeResponse>({
    url: '/loadBalancer/tree',
    method: 'post'
  })
}

// 创建负载均衡
export const createLoadBalancer = (data: LoadBalancerForm) => {
  return request({
    url: '/api/v1/cloud/loadbalancer',
    method: 'post',
    data
  })
}

// 更新负载均衡
export const updateLoadBalancer = (id: number, data: LoadBalancerForm) => {
  return request({
    url: `/api/v1/cloud/loadbalancer/${id}`,
    method: 'put',
    data
  })
}

// 删除负载均衡
export const deleteLoadBalancer = (id: number) => {
  return request({
    url: `/api/v1/cloud/loadbalancer/${id}`,
    method: 'delete'
  })
}

// 批量删除负载均衡
export const batchDeleteLoadBalancer = (ids: number[]) => {
  return request({
    url: '/api/v1/cloud/loadbalancer/batch',
    method: 'delete',
    data: { ids }
  })
}

// 同步负载均衡数据
export const syncLoadBalancer = (id: number) => {
  return request<SyncResponse>({
    url: '/loadBalancer/sync',
    method: 'post',
    data: { id }
  })
}
