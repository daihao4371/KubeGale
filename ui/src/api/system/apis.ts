import request from '@/api/request'

// 获取所有API
export const getAllApis = () => {
  return request({
    url: '/api/getAllApis',
    method: 'post',
    data: {}
  })
}

// 获取API分组
export const getApiGroups = () => {
  return request({
    url: '/api/getApiGroups',
    method: 'get'
  })
}

// 创建API
export const createApi = (data: {
  path: string
  description: string
  apiGroup: string
  method: string
}) => {
  return request({
    url: '/api/createApi',
    method: 'post',
    data
  })
}

// 根据ID获取API
export const getApiById = (data: { ID: number }) => {
  return request({
    url: '/api/getApiById',
    method: 'post',
    data
  })
}

// 更新API
export const updateApi = (data: {
  ID: number
  path: string
  description: string
  apiGroup: string
  method: string
}) => {
  return request({
    url: '/api/updateApi',
    method: 'post',
    data
  })
}

// 删除API
export const deleteApi = (data: { ID: number }) => {
  return request({
    url: '/api/deleteApi',
    method: 'post',
    data
  })
}

// 批量删除API
export const deleteApisByIds = (data: { ids: number[] }) => {
  return request({
    url: '/api/deleteApisByIds',
    method: 'delete',
    data
  })
}

// 刷新Casbin缓存
export const freshCasbin = () => {
  return request({
    url: '/api/freshCasbin',
    method: 'get'
  })
}
