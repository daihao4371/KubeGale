import service from '@/api/request'
import type { 
  Cluster, 
  CreateClusterRequest, 
  UpdateClusterRequest, 
  ClusterUser, 
  ClusterRole, 
  ApiGroupOption,
  ClusterSearchParams,
  PageResponse,
  CreateClusterUserRequest,
  UpdateClusterUserRequest,
  DeleteClusterUserRequest,
  ClusterRoleTypeRequest,
  ClusterApiGroupsRequest,
  CreateClusterRoleRequest,
  UpdateClusterRoleRequest,
  DeleteClusterRoleRequest,
  ClusterUserNamespaceResponse,
  ClusterListNamespaceResponse,
  ClusterDetail
} from '@/types/kubernetes/cluster'

export const clusterApi = {
  // 获取集群列表
  getClusterList(params?: ClusterSearchParams) {
    return service.get<PageResponse<Cluster>>('/kubernetes/clusterList', { params })
  },

  // 获取集群详情
  getClusterById(id: number) {
    return service.get<ClusterDetail>(`/kubernetes/findK8sCluster`, { params: { id } })
  },

  // 根据名称获取集群
  getClusterByName(name: string) {
    return service.get<ClusterDetail>(`/kubernetes/findK8sClusterByName`, { params: { name } })
  },

  // 创建集群
  createCluster(data: CreateClusterRequest) {
    return service.post<Cluster>('/kubernetes/cluster', data)
  },

  // 更新集群
  updateCluster(data: UpdateClusterRequest) {
    return service.put<Cluster>('/kubernetes/cluster', data)
  },

  // 删除集群
  deleteCluster(id: number) {
    return service.delete(`/kubernetes/cluster`, { data: { id } })
  },

  // 批量删除集群
  deleteClusters(ids: number[]) {
    return service.delete(`/kubernetes/clusterByIds`, { data: { ids } })
  },

  // 创建集群凭据
  createCredential(id: number) {
    return service.post('/kubernetes/credential', { id })
  },

  // 获取集群用户
  getClusterUserById(data: { id: number }) {
    return service.post<ClusterUser[]>('/kubernetes/getUserById', data)
  },

  // 获取集群角色
  getClusterRoles(data: ClusterRoleTypeRequest) {
    return service.post<ClusterRole[]>('/kubernetes/getClusterRoles', data)
  },

  // 获取集群资源分组
  getClusterApiGroups(data: ClusterApiGroupsRequest) {
    return service.post<ApiGroupOption[]>('/kubernetes/getClusterApiGroups', data)
  },

  // 创建集群角色
  createClusterRole(data: CreateClusterRoleRequest) {
    return service.post('/kubernetes/createClusterRole', data)
  },

  // 更新集群角色
  updateClusterRole(data: UpdateClusterRoleRequest) {
    return service.put('/kubernetes/updateClusterRole', data)
  },

  // 删除集群角色
  deleteClusterRole(data: DeleteClusterRoleRequest) {
    return service.delete('/kubernetes/deleteClusterRole', { data })
  },

  // 创建集群用户授权
  createClusterUser(data: CreateClusterUserRequest) {
    return service.post('/kubernetes/createClusterUser', data)
  },

  // 更新集群用户授权
  updateClusterUser(data: UpdateClusterUserRequest) {
    return service.put('/kubernetes/updateClusterUser', data)
  },

  // 删除集群用户
  deleteClusterUser(data: DeleteClusterUserRequest) {
    return service.delete('/kubernetes/deleteClusterUser', { data })
  },

  // 获取集群用户命名空间
  getClusterUserNamespace(data: { id: number }) {
    return service.post<ClusterUserNamespaceResponse>('/kubernetes/getClusterUserNamespace', data)
  },

  // 获取集群命名空间列表
  getClusterListNamespace(data: { id: number }) {
    return service.post<ClusterListNamespaceResponse>('/kubernetes/getClusterListNamespace', data)
  }
} 