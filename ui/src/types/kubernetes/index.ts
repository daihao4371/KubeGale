// 集群管理相关类型
export * from './cluster'

// 工作负载相关类型
export * from './workload'

// 网络资源相关类型
export * from './network'

// 存储资源相关类型
export * from './storage'

// 配置管理相关类型
export * from './config'

// RBAC 权限相关类型
export * from './rbac'

// 监控相关类型
export * from './monitoring'

// 节点管理相关类型
export * from './node'

// 公共类型定义

// 通用请求参数
export interface BaseRequest {
  cluster_id: number
  namespace?: string
}

// 分页信息
export interface PageInfo {
  page: number
  pageSize: number
  keyword?: string
}

// 通用列表请求
export interface ListRequest extends BaseRequest {
  labelSelector?: string
  fieldSelector?: string
  pageInfo?: PageInfo
}

// 通用列表响应
export interface ListResponse<T = any> {
  items: T[]
  total: number
  page?: number
  pageSize?: number
}

// 通用详情请求
export interface DetailRequest extends BaseRequest {
  name: string
}

// 通用详情响应
export interface DetailResponse<T = any> {
  items: T
}

// 通用创建请求
export interface CreateRequest extends BaseRequest {
  content: any
}

// 通用更新请求
export interface UpdateRequest extends BaseRequest {
  name: string
  content: any
}

// 通用删除请求
export interface DeleteRequest extends BaseRequest {
  name: string
}

// Kubernetes 资源状态
export type ResourceStatus = 'Running' | 'Pending' | 'Succeeded' | 'Failed' | 'Unknown'

// Kubernetes 资源条件
export interface ResourceCondition {
  type: string
  status: string
  lastTransitionTime: string
  lastUpdateTime?: string
  reason?: string
  message?: string
}

// Kubernetes 标签选择器
export interface LabelSelector {
  matchLabels?: Record<string, string>
  matchExpressions?: Array<{
    key: string
    operator: string // In, NotIn, Exists, DoesNotExist
    values?: string[]
  }>
}

// Kubernetes 资源引用
export interface ResourceRef {
  apiVersion: string
  kind: string
  name: string
  namespace?: string
  uid?: string
}

// 操作结果
export interface OperationResult {
  success: boolean
  message?: string
  data?: any
  error?: string
}

// 批量操作结果
export interface BatchOperationResult {
  total: number
  success: number
  failed: number
  results: Array<{
    name: string
    success: boolean
    message?: string
    error?: string
  }>
}

// WebSocket 消息类型
export interface WebSocketMessage {
  type: string
  data: any
  timestamp: string
}

// 终端消息
export interface TerminalMessage extends WebSocketMessage {
  type: 'data' | 'resize' | 'ping' | 'pong' | 'error' | 'close'
  data: {
    input?: string
    output?: string
    cols?: number
    rows?: number
    error?: string
  }
}

// 日志消息
export interface LogMessage extends WebSocketMessage {
  type: 'log' | 'error' | 'close'
  data: {
    content: string
    timestamp: string
    source?: string
  }
}

// 事件消息
export interface EventMessage extends WebSocketMessage {
  type: 'event'
  data: {
    object: any
    event_type: string // ADDED, MODIFIED, DELETED
    resource_version: string
  }
}

// API 错误类型
export interface ApiError {
  code: number
  message: string
  details?: any
  timestamp: string
}

// 文件上传相关
export interface FileUpload {
  file: File
  path: string
  progress?: number
  status?: 'pending' | 'uploading' | 'success' | 'error'
  error?: string
}

// 文件下载相关
export interface FileDownload {
  path: string
  filename: string
  size?: number
  progress?: number
  status?: 'pending' | 'downloading' | 'success' | 'error'
  error?: string
}

// 搜索类型
export interface SearchOption {
  label: string
  value: string
  type?: string
  description?: string
}

// 过滤器类型
export interface FilterOption {
  key: string
  label: string
  type: 'text' | 'select' | 'date' | 'number'
  options?: Array<{
    label: string
    value: any
  }>
  placeholder?: string
}

// 排序类型
export interface SortOption {
  key: string
  order: 'asc' | 'desc'
  label?: string
}

// 表格列配置
export interface TableColumn {
  key: string
  title: string
  width?: number | string
  sortable?: boolean
  filterable?: boolean
  resizable?: boolean
  align?: 'left' | 'center' | 'right'
  fixed?: 'left' | 'right'
  render?: (value: any, record: any, index: number) => any
}

// 表格配置
export interface TableConfig {
  columns: TableColumn[]
  rowKey: string
  bordered?: boolean
  striped?: boolean
  size?: 'small' | 'medium' | 'large'
  loading?: boolean
  pagination?: {
    current: number
    pageSize: number
    total: number
    showSizeChanger?: boolean
    showQuickJumper?: boolean
    showTotal?: boolean
  }
  selection?: {
    type: 'checkbox' | 'radio'
    selectedRowKeys: any[]
    onChange: (selectedRowKeys: any[], selectedRows: any[]) => void
  }
}

// 图表配置
export interface ChartConfig {
  type: 'line' | 'bar' | 'pie' | 'area' | 'scatter'
  title?: string
  width?: number
  height?: number
  data: any[]
  xField?: string
  yField?: string
  seriesField?: string
  color?: string | string[]
  legend?: boolean
  tooltip?: boolean
  grid?: boolean
  animation?: boolean
}

// 布局配置
export interface LayoutConfig {
  sidebar: {
    collapsed: boolean
    width: number
  }
  header: {
    height: number
    fixed: boolean
  }
  footer: {
    height: number
    fixed: boolean
  }
  theme: 'light' | 'dark'
  language: string
}

// 用户配置
export interface UserPreference {
  theme: 'light' | 'dark'
  language: string
  timezone: string
  date_format: string
  time_format: string
  page_size: number
  auto_refresh: boolean
  refresh_interval: number
  sidebar_collapsed: boolean
  table_density: 'compact' | 'normal' | 'loose'
  chart_theme: string
}