import service from '@/api/request'
import type {
  // RBAC 资源相关类型
  Role,
  ClusterRole,
  RoleBinding,
  ClusterRoleBinding,
  ServiceAccount,
  PodSecurityPolicy,
  SecurityContextConstraints,
  GetRBACListRequest,
  GetRBACDetailRequest,
  CreateRBACRequest,
  UpdateRBACRequest,
  DeleteRBACRequest,
  CreateRoleBindingRequest,
  CreateClusterRoleBindingRequest,
  CreateServiceAccountRequest,
  BindRoleToSubjectRequest,
  UnbindRoleFromSubjectRequest,
  GetUserPermissionsRequest,
  GetResourcePermissionsRequest,
  RBACListResponse,
  RBACDetailResponse,
  UserPermissionsResponse,
  ResourcePermissionsResponse,
  RBACMetricsResponse
} from '@/types/kubernetes/rbac'

export const rbacApi = {
  // ==================== Role 相关 API ====================
  
  // 获取 Role 列表
  getRoleList(params: GetRBACListRequest) {
    return service.get<RBACListResponse<Role>>('/kubernetes/role', { params })
  },

  // 获取 Role 详情
  getRoleDetail(params: GetRBACDetailRequest) {
    return service.get<RBACDetailResponse<Role>>('/kubernetes/role/detail', { params })
  },

  // 创建 Role
  createRole(data: CreateRBACRequest) {
    return service.post<Role>('/kubernetes/role', data)
  },

  // 更新 Role
  updateRole(data: UpdateRBACRequest) {
    return service.put<Role>('/kubernetes/role', data)
  },

  // 删除 Role
  deleteRole(data: DeleteRBACRequest) {
    return service.delete('/kubernetes/role', { data })
  },

  // ==================== ClusterRole 相关 API ====================
  
  // 获取 ClusterRole 列表
  getClusterRoleList(params: GetRBACListRequest) {
    return service.get<RBACListResponse<ClusterRole>>('/kubernetes/clusterrole', { params })
  },

  // 获取 ClusterRole 详情
  getClusterRoleDetail(params: GetRBACDetailRequest) {
    return service.get<RBACDetailResponse<ClusterRole>>('/kubernetes/clusterrole/detail', { params })
  },

  // 创建 ClusterRole
  createClusterRole(data: CreateRBACRequest) {
    return service.post<ClusterRole>('/kubernetes/clusterrole', data)
  },

  // 更新 ClusterRole
  updateClusterRole(data: UpdateRBACRequest) {
    return service.put<ClusterRole>('/kubernetes/clusterrole', data)
  },

  // 删除 ClusterRole
  deleteClusterRole(data: DeleteRBACRequest) {
    return service.delete('/kubernetes/clusterrole', { data })
  },

  // ==================== RoleBinding 相关 API ====================
  
  // 获取 RoleBinding 列表
  getRoleBindingList(params: GetRBACListRequest) {
    return service.get<RBACListResponse<RoleBinding>>('/kubernetes/rolebinding', { params })
  },

  // 获取 RoleBinding 详情
  getRoleBindingDetail(params: GetRBACDetailRequest) {
    return service.get<RBACDetailResponse<RoleBinding>>('/kubernetes/rolebinding/detail', { params })
  },

  // 创建 RoleBinding
  createRoleBinding(data: CreateRoleBindingRequest) {
    return service.post<RoleBinding>('/kubernetes/rolebinding', data)
  },

  // 更新 RoleBinding
  updateRoleBinding(data: UpdateRBACRequest) {
    return service.put<RoleBinding>('/kubernetes/rolebinding', data)
  },

  // 删除 RoleBinding
  deleteRoleBinding(data: DeleteRBACRequest) {
    return service.delete('/kubernetes/rolebinding', { data })
  },

  // ==================== ClusterRoleBinding 相关 API ====================
  
  // 获取 ClusterRoleBinding 列表
  getClusterRoleBindingList(params: GetRBACListRequest) {
    return service.get<RBACListResponse<ClusterRoleBinding>>('/kubernetes/clusterrolebinding', { params })
  },

  // 获取 ClusterRoleBinding 详情
  getClusterRoleBindingDetail(params: GetRBACDetailRequest) {
    return service.get<RBACDetailResponse<ClusterRoleBinding>>('/kubernetes/clusterrolebinding/detail', { params })
  },

  // 创建 ClusterRoleBinding
  createClusterRoleBinding(data: CreateClusterRoleBindingRequest) {
    return service.post<ClusterRoleBinding>('/kubernetes/clusterrolebinding', data)
  },

  // 更新 ClusterRoleBinding
  updateClusterRoleBinding(data: UpdateRBACRequest) {
    return service.put<ClusterRoleBinding>('/kubernetes/clusterrolebinding', data)
  },

  // 删除 ClusterRoleBinding
  deleteClusterRoleBinding(data: DeleteRBACRequest) {
    return service.delete('/kubernetes/clusterrolebinding', { data })
  },

  // ==================== ServiceAccount 相关 API ====================
  
  // 获取 ServiceAccount 列表
  getServiceAccountList(params: GetRBACListRequest) {
    return service.get<RBACListResponse<ServiceAccount>>('/kubernetes/serviceaccount', { params })
  },

  // 获取 ServiceAccount 详情
  getServiceAccountDetail(params: GetRBACDetailRequest) {
    return service.get<RBACDetailResponse<ServiceAccount>>('/kubernetes/serviceaccount/detail', { params })
  },

  // 创建 ServiceAccount
  createServiceAccount(data: CreateServiceAccountRequest) {
    return service.post<ServiceAccount>('/kubernetes/serviceaccount', data)
  },

  // 更新 ServiceAccount
  updateServiceAccount(data: UpdateRBACRequest) {
    return service.put<ServiceAccount>('/kubernetes/serviceaccount', data)
  },

  // 删除 ServiceAccount
  deleteServiceAccount(data: DeleteRBACRequest) {
    return service.delete('/kubernetes/serviceaccount', { data })
  },

  // ==================== PodSecurityPolicy 相关 API ====================
  
  // 获取 PodSecurityPolicy 列表
  getPodSecurityPolicyList(params: GetRBACListRequest) {
    return service.get<RBACListResponse<PodSecurityPolicy>>('/kubernetes/podsecuritypolicy', { params })
  },

  // 获取 PodSecurityPolicy 详情
  getPodSecurityPolicyDetail(params: GetRBACDetailRequest) {
    return service.get<RBACDetailResponse<PodSecurityPolicy>>('/kubernetes/podsecuritypolicy/detail', { params })
  },

  // 创建 PodSecurityPolicy
  createPodSecurityPolicy(data: CreateRBACRequest) {
    return service.post<PodSecurityPolicy>('/kubernetes/podsecuritypolicy', data)
  },

  // 更新 PodSecurityPolicy
  updatePodSecurityPolicy(data: UpdateRBACRequest) {
    return service.put<PodSecurityPolicy>('/kubernetes/podsecuritypolicy', data)
  },

  // 删除 PodSecurityPolicy
  deletePodSecurityPolicy(data: DeleteRBACRequest) {
    return service.delete('/kubernetes/podsecuritypolicy', { data })
  },

  // ==================== SecurityContextConstraints 相关 API (OpenShift) ====================
  
  // 获取 SCC 列表
  getSCCList(params: GetRBACListRequest) {
    return service.get<RBACListResponse<SecurityContextConstraints>>('/kubernetes/scc', { params })
  },

  // 获取 SCC 详情
  getSCCDetail(params: GetRBACDetailRequest) {
    return service.get<RBACDetailResponse<SecurityContextConstraints>>('/kubernetes/scc/detail', { params })
  },

  // 创建 SCC
  createSCC(data: CreateRBACRequest) {
    return service.post<SecurityContextConstraints>('/kubernetes/scc', data)
  },

  // 更新 SCC
  updateSCC(data: UpdateRBACRequest) {
    return service.put<SecurityContextConstraints>('/kubernetes/scc', data)
  },

  // 删除 SCC
  deleteSCC(data: DeleteRBACRequest) {
    return service.delete('/kubernetes/scc', { data })
  },

  // ==================== RBAC 权限管理 API ====================
  
  // 绑定角色到主体
  bindRoleToSubject(data: BindRoleToSubjectRequest) {
    return service.post('/kubernetes/rbac/bind', data)
  },

  // 解绑角色从主体
  unbindRoleFromSubject(data: UnbindRoleFromSubjectRequest) {
    return service.post('/kubernetes/rbac/unbind', data)
  },

  // 获取用户权限
  getUserPermissions(params: GetUserPermissionsRequest) {
    return service.get<UserPermissionsResponse>('/kubernetes/rbac/user-permissions', { params })
  },

  // 获取资源权限
  getResourcePermissions(params: GetResourcePermissionsRequest) {
    return service.get<ResourcePermissionsResponse>('/kubernetes/rbac/resource-permissions', { params })
  },

  // 权限检查
  checkPermissions(data: {
    cluster_id: number
    user_name: string
    user_kind: string
    namespace?: string
    api_group: string
    resource: string
    verb: string
    resource_name?: string
  }) {
    return service.post('/kubernetes/rbac/can-i', data)
  },

  // ==================== RBAC 监控和统计 API ====================
  
  // 获取 RBAC 指标
  getRBACMetrics(params: { cluster_id: number; namespace?: string }) {
    return service.get<RBACMetricsResponse>('/kubernetes/rbac/metrics', { params })
  },

  // 获取权限矩阵
  getPermissionMatrix(params: { cluster_id: number; namespace?: string }) {
    return service.get('/kubernetes/rbac/permission-matrix', { params })
  },

  // 安全审计
  getSecurityAudit(params: { cluster_id: number; namespace?: string }) {
    return service.get('/kubernetes/rbac/security-audit', { params })
  },

  // 权限风险评估
  assessPermissionRisks(data: { cluster_id: number; namespace?: string }) {
    return service.post('/kubernetes/rbac/risk-assessment', data)
  },

  // 最小权限建议
  getLeastPrivilegeRecommendations(params: {
    cluster_id: number
    namespace?: string
    user_name?: string
    service_account?: string
  }) {
    return service.get('/kubernetes/rbac/least-privilege', { params })
  },

  // 权限对比
  comparePermissions(data: {
    cluster_id: number
    source_user: { name: string; kind: string; namespace?: string }
    target_user: { name: string; kind: string; namespace?: string }
  }) {
    return service.post('/kubernetes/rbac/compare-permissions', data)
  },

  // 权限模拟
  simulatePermissions(data: {
    cluster_id: number
    user_name: string
    user_kind: string
    namespace?: string
    actions: Array<{
      api_group: string
      resource: string
      verb: string
      resource_name?: string
    }>
  }) {
    return service.post('/kubernetes/rbac/simulate', data)
  },

  // 批量权限操作
  batchPermissionOperations(data: {
    cluster_id: number
    operations: Array<{
      action: 'grant' | 'revoke'
      user_name: string
      user_kind: string
      namespace?: string
      role_name: string
      role_kind: string
    }>
  }) {
    return service.post('/kubernetes/rbac/batch-operations', data)
  }
}