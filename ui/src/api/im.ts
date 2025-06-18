import request from '@/utils/request'
import type { ApiResponse, NotificationCardContent } from '@/types/im'

// 获取通知卡片内容
export function getCardContent(params: { notification_id: number }) {
  return request<ApiResponse<NotificationCardContent>>({
    url: '/notification/getCardContent',
    method: 'get',
    params
  })
} 