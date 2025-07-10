// 集群基础信息
export interface Cluster {
  id: number
  uuid: string
  name: string
  alias: string
  status: string
  version: string
  city: string
  district: string
  description?: string
  kube_type: number // 1: KubeConfig, 2: Token
  kube_config: string
  kube_ca_crt: string
  api_address: string
  prometheus_url: string
  prometheus_auth_type: number
  prometheus_user: string
  prometheus_pwd: string
  created_at: string
  updated_at: string
}

// 创建集群请求
export interface CreateClusterRequest {
  name: string
  kube_type: number
  kube_config?: string
  kube_ca_crt?: string
  api_address: string
  prometheus_url?: string
  prometheus_auth_type?: number
  prometheus_user?: string
  prometheus_pwd?: string
}

// 更新集群请求
export interface UpdateClusterRequest {
  id: number
  name?: string
  kube_type?: number
  kube_config?: string
  kube_ca_crt?: string
  api_address?: string
  prometheus_url?: string
  prometheus_auth_type?: number
  prometheus_user?: string
  prometheus_pwd?: string
}

// 集群用户信息
export interface ClusterUser {
  id: number
  uuid: string
  userName: string
  nickName: string
  kube_config: string
  cluster_roles: string
  namespace_roles: string
  cluster_id: number
  created_at: string
  updated_at: string
}

// 集群角色信息
export interface ClusterRole {
  metadata: {
    name: string
    creationTimestamp: string
  }
  rules: Array<{
    apiGroups: string[]
    resources: string[]
    verbs: string[]
  }>
}

// 集群API组信息
export interface ApiGroupOption {
  group: string
  resources: ApiResourceOption[]
}

export interface ApiResourceOption {
  resource: string
  verbs: string[]
}

// 集群搜索参数
export interface ClusterSearchParams {
  name?: string
  kube_type?: number
  page?: number
  pageSize?: number
}

// 分页响应
export interface PageResponse<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 命名空间角色信息
export interface NamespaceRoles {
  namespace: string
  roles: string[]
}

// 创建集群用户请求
export interface CreateClusterUserRequest {
  cluster_id: number
  cluster_roles: string[]
  namespace_roles: NamespaceRoles[]
  user_uuids: string[]
}

// 更新集群用户请求
export interface UpdateClusterUserRequest {
  cluster_id: number
  cluster_roles: string[]
  namespace_roles: NamespaceRoles[]
  uuid: string
}

// 删除集群用户请求
export interface DeleteClusterUserRequest {
  cluster_id: number
  user_uuids: string[]
}

// 集群角色类型请求
export interface ClusterRoleTypeRequest {
  role_type: string
  cluster_id: number
}

// 集群API组请求
export interface ClusterApiGroupsRequest {
  api_type: string
  cluster_id: number
}

// 角色数据
export interface RoleData {
  cluster_id: number
  name: string
  annotations: Record<string, string>
  labels: Record<string, string>
  rules: Array<{
    apiGroups: string[]
    resources: string[]
    verbs: string[]
  }>
}

// 创建集群角色请求
export interface CreateClusterRoleRequest {
  cluster_id: number
  role: RoleData
}

// 更新集群角色请求
export interface UpdateClusterRoleRequest {
  cluster_id: number
  role: RoleData
}

// 删除集群角色请求
export interface DeleteClusterRoleRequest {
  cluster_id: number
  name: string
}

// 命名空间信息
export interface Namespace {
  metadata: {
    name: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
  status: {
    phase: string
  }
}

// 集群用户命名空间响应
export interface ClusterUserNamespaceResponse {
  namespaces: string[]
}

// 集群命名空间列表响应
export interface ClusterListNamespaceResponse {
  namespaces: Namespace[]
}

// 集群详情信息
export interface ClusterDetail extends Cluster {
  users?: ClusterUser[]
  namespaces?: Namespace[]
  roles?: ClusterRole[]
  apiGroups?: ApiGroupOption[]
} 