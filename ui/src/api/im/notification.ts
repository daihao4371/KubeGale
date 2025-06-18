import request from '@/utils/request'
import type {
  ApiResponse,
  CreateFeiShuParams,
  UpdateNotificationParams,
  DeleteNotificationParams,
  TestNotificationParams,
  CardContentParams,
  GetNotificationListParams,
  GetNotificationDetailParams,
  NotificationListResponse,
  CreateUpdateResponse,
  NotificationItem,
  DingTalkCardContent,
  NotificationCardContent,
  FeiShuNotificationListResponse
} from '@/types/im'

interface UpdateFeiShuParams {
  id: number
  name: string
  notification_policy: string
  robot_url: string
  send_daily_stats: boolean
  card_content: {
    alert_level: string
    alert_name: string
    notification_policy: string
    alert_content: string
    alert_time: string
    notified_users: string[]
    last_similar_alert: string
    alert_handler: string
    claim_alert: boolean
    resolve_alert: boolean
    mute_alert: boolean
    unresolved_alert: boolean
  }
}

// 创建飞书通知
export const createFeiShu = (data: CreateFeiShuParams) => {
  return request<ApiResponse<CreateUpdateResponse>>({
    url: '/notification/createFeiShu',
    method: 'post',
    data
  })
}

// 更新飞书通知
export const updateFeiShu = (data: UpdateFeiShuParams) => {
  return request<ApiResponse>({
    url: '/notification/updateFeiShu',
    method: 'put',
    data
  })
}

// 删除通知配置
export const deleteNotification = (params: { id: number; type: string }) => {
  return request<ApiResponse>({
    url: '/notification/deleteNotification',
    method: 'delete',
    params
  })
}

// 测试通知发送
export const testNotification = (data: TestNotificationParams) => {
  return request<ApiResponse>({
    url: '/notification/testNotification',
    method: 'post',
    data
  })
}

// 创建卡片内容
export const createCardContent = (data: CardContentParams) => {
  return request<ApiResponse<CreateUpdateResponse>>({
    url: '/notification/createCardContent',
    method: 'post',
    data
  })
}

// 更新卡片内容
export const updateCardContent = (data: CardContentParams) => {
  return request<ApiResponse<CreateUpdateResponse>>({
    url: '/notification/updateCardContent',
    method: 'put',
    data
  })
}

// 获取通知列表
export const getNotificationList = (params: {
  page: number
  page_size: number
}) => {
  return request<ApiResponse<NotificationListResponse>>({
    url: '/notification/getNotificationList',
    method: 'post',
    data: params
  })
}

// 根据ID获取通知配置
export const getNotificationById = (params: GetNotificationDetailParams) => {
  return request<ApiResponse<NotificationItem>>({
    url: '/notification/getNotificationById',
    method: 'get',
    params
  })
}

// 根据通知ID获取卡片内容
export const getCardContent = (params: { notification_id: number }) => {
  return request<ApiResponse<NotificationCardContent>>({
    url: '/notification/getCardContent',
    method: 'get',
    params
  })
}

