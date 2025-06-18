// 钉钉通知卡片内容
export interface DingTalkCardContent {
  alert_level: string
  alert_name: string
  notification_policy: string
  alert_content: string
  notified_users: string
  last_similar_alert: string
  alert_handler: string
  claim_alert: boolean
  resolve_alert: boolean
  mute_alert: boolean
  unresolved_alert: boolean
}

// 飞书通知卡片内容
export interface FeiShuCardContent {
  alert_level: string
  alert_name: string
  alert_content: string
  notified_users: string
  alert_handler: string
}

// 创建钉钉通知请求参数
export interface CreateDingTalkParams {
  id?: number
  name: string
  notificationPolicy: string
  robotURL: string
  signatureKey?: string
  card_content: DingTalkCardContent
}

// 创建飞书通知请求参数
export interface CreateFeiShuParams {
  name: string
  type: 'feishu'
  enabled: boolean
  webhook_url: string
  description: string
  tags: string[]
  notify_events: string[]
  receivers: string[]
  send_daily_stats: boolean
}

// 更新通知请求参数
export interface UpdateNotificationParams {
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
    notified_users: string
    alert_handler: string
    claim_alert: boolean
    resolve_alert: boolean
    mute_alert: boolean
    unresolved_alert: boolean
    alert_time: string
  }
}

// 删除通知请求参数
export interface DeleteNotificationParams {
  id: number
  type: 'feishu'
}

// 测试通知请求参数
export interface TestNotificationParams {
  id: number
  type: 'feishu'
  message?: string
}

// 卡片内容请求参数
export interface CardContentParams {
  notification_id: number
  card_content: FeiShuCardContent
}

// 获取通知列表请求参数
export interface GetNotificationListParams {
  page?: number
  pageSize?: number
  orderKey?: string
  desc?: boolean
}

// 通知配置详情
export interface NotificationConfig {
  id: number
  name: string
  type: string
  notification_policy: string
  send_daily_stats: boolean
  created_at: string
  updated_at: string
  robot_url: string
}

// 通知卡片内容
export interface NotificationCardContent {
  id: number
  notification_id: number
  alert_level: string
  alert_name: string
  notification_policy: string
  alert_content: string
  alert_time: string
  notified_users: string
  last_similar_alert: string
  alert_handler: string
  claim_alert: boolean
  resolve_alert: boolean
  mute_alert: boolean
  unresolved_alert: boolean
}

// 通知详情响应
export interface NotificationDetailResponse {
  config: NotificationConfig
  card_content: NotificationCardContent
}

// 获取通知详情请求参数
export interface GetNotificationDetailParams {
  id: number
  type: 'feishu'
}

// 通知列表项
export interface NotificationItem {
  id: number
  name: string
  type: string
  notification_policy: string
  robot_url: string
  created_at: string
  updated_at: string
  send_daily_stats: boolean
  config?: NotificationConfig
  card_content?: NotificationCardContent
}

// 通知列表响应
export interface NotificationListResponse {
  total: number
  list: NotificationItem[]
  page: number
  pageSize: number
}

// 飞书通知列表项
export interface FeiShuNotificationItem {
  id: number
  name: string
  type: string
  notification_policy: string
  robot_url: string
  created_at: string
  updated_at: string
  send_daily_stats: boolean
}

// 飞书通知列表响应
export interface FeiShuNotificationListResponse {
  total: number
  list: FeiShuNotificationItem[]
  page: number
  pageSize: number
}

// 通用响应格式
export interface ApiResponse<T = Record<string, never>> {
  code: number
  data: T
  msg: string
}

// 创建/更新响应
export interface CreateUpdateResponse {
  error?: string
}

// 更新钉钉通知请求参数
export interface UpdateDingTalkParams {
  id: number
  name: string
  notification_policy: string
  send_daily_stats: boolean
  signature_key: string
  robot_url: string
  type: 'dingtalk' | 'feishu'
  card_content: {
    alert_level: string
    alert_name: string
    notification_policy: string
    alert_content: string
    notified_users: string
    last_similar_alert?: string
    alert_handler: string
    claim_alert: boolean
    resolve_alert: boolean
    mute_alert: boolean
    unresolved_alert: boolean
  }
}
