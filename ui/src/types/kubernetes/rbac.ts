// RBAC 资源基础类型
export interface RBACBase {
  metadata: {
    name: string
    namespace?: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
}

// Role 相关类型
export interface Role extends RBACBase {
  rules: PolicyRule[]
}

// ClusterRole 相关类型
export interface ClusterRole extends Omit<RBACBase, 'namespace'> {
  rules: PolicyRule[]
  aggregationRule?: {
    clusterRoleSelectors: Array<{
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }>
  }
}

// PolicyRule 类型
export interface PolicyRule {
  verbs: string[]
  apiGroups?: string[]
  resources?: string[]
  resourceNames?: string[]
  nonResourceURLs?: string[]
}

// RoleBinding 相关类型
export interface RoleBinding extends RBACBase {
  subjects?: Subject[]
  roleRef: RoleRef
}

// ClusterRoleBinding 相关类型
export interface ClusterRoleBinding extends Omit<RBACBase, 'namespace'> {
  subjects?: Subject[]
  roleRef: RoleRef
}

// Subject 类型
export interface Subject {
  kind: string // User, Group, ServiceAccount
  name: string
  namespace?: string
  apiGroup?: string
}

// RoleRef 类型
export interface RoleRef {
  kind: string // Role or ClusterRole
  name: string
  apiGroup: string
}

// ServiceAccount 相关类型
export interface ServiceAccount extends RBACBase {
  secrets?: Array<{
    name: string
  }>
  imagePullSecrets?: Array<{
    name: string
  }>
  automountServiceAccountToken?: boolean
}

// PodSecurityPolicy 相关类型
export interface PodSecurityPolicy extends Omit<RBACBase, 'namespace'> {
  spec: {
    privileged: boolean
    allowPrivilegeEscalation?: boolean
    defaultAllowPrivilegeEscalation?: boolean
    requiredDropCapabilities?: string[]
    allowedCapabilities?: string[]
    defaultAddCapabilities?: string[]
    volumes?: string[]
    hostNetwork?: boolean
    hostIPC?: boolean
    hostPID?: boolean
    hostPorts?: Array<{
      min: number
      max: number
    }>
    runAsUser?: {
      rule: string
      ranges?: Array<{
        min: number
        max: number
      }>
    }
    runAsGroup?: {
      rule: string
      ranges?: Array<{
        min: number
        max: number
      }>
    }
    supplementalGroups?: {
      rule: string
      ranges?: Array<{
        min: number
        max: number
      }>
    }
    fsGroup?: {
      rule: string
      ranges?: Array<{
        min: number
        max: number
      }>
    }
    seLinux?: {
      rule: string
      seLinuxOptions?: {
        level?: string
        role?: string
        type?: string
        user?: string
      }
    }
    readOnlyRootFilesystem?: boolean
    allowedHostPaths?: Array<{
      pathPrefix: string
      readOnly?: boolean
    }>
    allowedFlexVolumes?: Array<{
      driver: string
    }>
    allowedCSIDrivers?: Array<{
      name: string
    }>
    allowedUnsafeSysctls?: string[]
    forbiddenSysctls?: string[]
    allowedProcMountTypes?: string[]
    runtimeClass?: {
      allowedRuntimeClassNames: string[]
      defaultRuntimeClassName?: string
    }
  }
}

// SecurityContextConstraints (OpenShift) 相关类型
export interface SecurityContextConstraints extends Omit<RBACBase, 'namespace'> {
  allowHostDirVolumePlugin: boolean
  allowHostIPC: boolean
  allowHostNetwork: boolean
  allowHostPID: boolean
  allowHostPorts: boolean
  allowPrivilegedContainer: boolean
  allowedCapabilities?: string[]
  allowedFlexVolumes?: Array<{
    driver: string
  }>
  allowedUnsafeSysctls?: string[]
  defaultAddCapabilities?: string[]
  defaultAllowPrivilegeEscalation?: boolean
  forbiddenSysctls?: string[]
  fsGroup?: {
    type: string
    ranges?: Array<{
      min: number
      max: number
    }>
  }
  groups?: string[]
  priority?: number
  readOnlyRootFilesystem: boolean
  requiredDropCapabilities?: string[]
  runAsUser?: {
    type: string
    uid?: number
    uidRangeMin?: number
    uidRangeMax?: number
  }
  seLinuxContext?: {
    type: string
    seLinuxOptions?: {
      level?: string
      role?: string
      type?: string
      user?: string
    }
  }
  supplementalGroups?: {
    type: string
    ranges?: Array<{
      min: number
      max: number
    }>
  }
  users?: string[]
  volumes?: string[]
}

// API 请求类型
export interface GetRBACListRequest {
  cluster_id: number
  namespace?: string
  page?: number
  pageSize?: number
  keyword?: string
  labelSelector?: string
  fieldSelector?: string
}

export interface GetRBACDetailRequest {
  cluster_id: number
  namespace?: string
  name: string
}

export interface CreateRBACRequest {
  cluster_id: number
  namespace?: string
  content: any
}

export interface UpdateRBACRequest {
  cluster_id: number
  namespace?: string
  name: string
  content: any
}

export interface DeleteRBACRequest {
  cluster_id: number
  namespace?: string
  name: string
}

// RBAC 特殊操作请求
export interface CreateRoleBindingRequest {
  cluster_id: number
  namespace: string
  name: string
  role_name: string
  role_kind: string // Role or ClusterRole
  subjects: Subject[]
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

export interface CreateClusterRoleBindingRequest {
  cluster_id: number
  name: string
  cluster_role_name: string
  subjects: Subject[]
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

export interface CreateServiceAccountRequest {
  cluster_id: number
  namespace: string
  name: string
  automount_service_account_token?: boolean
  image_pull_secrets?: string[]
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

export interface BindRoleToSubjectRequest {
  cluster_id: number
  namespace?: string
  role_name: string
  role_kind: string
  subject: Subject
}

export interface UnbindRoleFromSubjectRequest {
  cluster_id: number
  namespace?: string
  role_binding_name: string
  subject: Subject
}

// 获取用户权限请求
export interface GetUserPermissionsRequest {
  cluster_id: number
  user_name: string
  user_kind: string // User, Group, ServiceAccount
  namespace?: string
}

// 获取资源权限请求
export interface GetResourcePermissionsRequest {
  cluster_id: number
  namespace?: string
  api_group: string
  resource: string
  resource_name?: string
}

// API 响应类型
export interface RBACListResponse<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface RBACDetailResponse<T = any> {
  items: T
}

// 权限查询响应
export interface UserPermissionsResponse {
  user: {
    name: string
    kind: string
    namespace?: string
  }
  permissions: Array<{
    namespace?: string
    role_name: string
    role_kind: string
    rules: PolicyRule[]
  }>
  effective_permissions: Array<{
    api_group: string
    resource: string
    verbs: string[]
    resource_names?: string[]
    namespaces?: string[]
  }>
}

export interface ResourcePermissionsResponse {
  resource: {
    api_group: string
    resource: string
    resource_name?: string
    namespace?: string
  }
  allowed_users: Array<{
    name: string
    kind: string
    namespace?: string
    verbs: string[]
    source: {
      role_name: string
      role_kind: string
      binding_name: string
    }
  }>
}

// RBAC 统计信息
export interface RBACMetrics {
  role_count: number
  cluster_role_count: number
  role_binding_count: number
  cluster_role_binding_count: number
  service_account_count: number
  pod_security_policy_count: number
  security_context_constraints_count: number
}

export interface RBACMetricsResponse {
  metrics: RBACMetrics
  timestamp: string
}