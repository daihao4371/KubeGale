import request from '@/api/request'
import type { AddUserReq, UpdateUserReq, ChangePasswordReq } from '@/types/cmdb'

// 用户角色设置接口
export interface SetUserAuthoritiesReq {
  ID: number
  authorityIds: number[]
}

// 获取用户列表
export const getUserList = (data: { page: number; pageSize: number }) => {
  return request({
    url: '/user/getUserList',
    method: 'post',
    data
  })
}

// 添加用户
export const addUser = (data: AddUserReq) => {
  return request({
    url: '/user/admin_register',
    method: 'post',
    data
  })
}

// 编辑用户
export const updateUser = (data: UpdateUserReq) => {
  return request({
    url: '/user/setUserInfo',
    method: 'put',
    data
  })
}

// 设置用户角色
export const setUserAuthorities = (data: SetUserAuthoritiesReq) => {
  return request({
    url: '/user/setUserAuthorities',
    method: 'post',
    data
  })
}

// 删除用户
export const deleteUser = (id: number) => {
  return request({
    url: '/user/deleteUser',
    method: 'delete',
    data: { id }
  })
}

// 重设密码
export const resetPassword = (data: { ID: number; newPassword: string }) => {
  return request({
    url: '/user/resetPassword',
    method: 'post',
    data
  })
}

// 获取用户信息
export const getUserInfo = () => {
  return request({
    url: '/user/getUserInfo',
    method: 'get'
  })
}

// 修改密码
export const changePassword = (data: ChangePasswordReq) => {
  return request({
    url: '/user/changePassword',
    method: 'post',
    data
  })
}

