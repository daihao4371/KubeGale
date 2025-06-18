import request from '@/utils/request'
import type { ApiResponse } from '@/types/cmdb'

interface UserInfo {
  id: string
  username: string
  email?: string
  role?: string
}

export function getUserInfo(userId: string) {
  return request<ApiResponse<UserInfo>>({
    url: `/user/${userId}`,
    method: 'get'
  })
} 