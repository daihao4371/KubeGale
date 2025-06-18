import request from '@/utils/request'
import type { Project } from '@/views/cmdb/project/project'

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

// 创建项目
export const createProject = (data: {
  name: string
  description: string
  manager: string
}) => {
  return request<ApiResponse>({
    url: '/cmdb/projects',
    method: 'post',
    data
  })
}

// 获取项目列表
export const getProjectList = (params: {
  page: number
  pageSize: number
  name?: string
  description?: string
}) => {
  return request<ApiResponse<ListResponse<Project>>>({
    url: '/cmdb/projects',
    method: 'get',
    params
  })
}

// 更新项目
export const updateProject = (data: {
  id: number
  name: string
  description: string
  manager: string
}) => {
  return request<ApiResponse>({
    url: '/cmdb/projects',
    method: 'put',
    data
  })
}

// 删除项目
export const deleteProject = (data: { id: number }) => {
  return request<ApiResponse>({
    url: '/cmdb/projects',
    method: 'delete',
    data
  })
}

// 批量删除项目
export const batchDeleteProjects = (data: { ids: number[] }) => {
  return request<ApiResponse>({
    url: '/cmdb/projectsByIds',
    method: 'delete',
    data
  })
}

// 获取项目详情
export const getProjectDetail = (id: number) => {
  return request<ApiResponse<Project>>({
    url: `/cmdb/projectsById`,
    method: 'get',
    params: { id }
  })
}

