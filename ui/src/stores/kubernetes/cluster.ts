import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { clusterApi } from '@/api/kubernetes/cluster'
import type { 
  Cluster, 
  CreateClusterRequest, 
  UpdateClusterRequest, 
  ClusterUser, 
  ClusterRole, 
  ApiGroupOption,
  ClusterSearchParams,
  CreateClusterUserRequest,
  UpdateClusterUserRequest,
  DeleteClusterUserRequest,
  CreateClusterRoleRequest,
  UpdateClusterRoleRequest,
  DeleteClusterRoleRequest,
  Namespace,
  ClusterDetail
} from '@/types/kubernetes/cluster'
import { ElMessage } from 'element-plus'

export const useClusterStore = defineStore('cluster', () => {
  // 状态
  const clusters = ref<Cluster[]>([])
  const currentCluster = ref<ClusterDetail | null>(null)
  const clusterUsers = ref<ClusterUser[]>([])
  const clusterRoles = ref<ClusterRole[]>([])
  const clusterNamespaces = ref<Namespace[]>([])
  const clusterApiGroups = ref<ApiGroupOption[]>([])
  const loading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 计算属性
  const hasClusters = computed(() => clusters.value.length > 0)
  const currentClusterId = computed(() => currentCluster.value?.id)

  // 获取集群列表
  const fetchClusters = async (params?: ClusterSearchParams) => {
    try {
      loading.value = true
      const response = await clusterApi.getClusterList({
        page: currentPage.value,
        pageSize: pageSize.value,
        ...params
      })
      clusters.value = response.list
      total.value = response.total
      return response
    } catch (error) {
      ElMessage.error('获取集群列表失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取集群详情
  const fetchClusterDetail = async (id: number) => {
    try {
      loading.value = true
      const response = await clusterApi.getClusterById(id)
      currentCluster.value = response
      return response
    } catch (error) {
      ElMessage.error('获取集群详情失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 创建集群
  const createCluster = async (data: CreateClusterRequest) => {
    try {
      loading.value = true
      const response = await clusterApi.createCluster(data)
      ElMessage.success('集群创建成功')
      await fetchClusters()
      return response
    } catch (error) {
      ElMessage.error('集群创建失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新集群
  const updateCluster = async (data: UpdateClusterRequest) => {
    try {
      loading.value = true
      const response = await clusterApi.updateCluster(data)
      ElMessage.success('集群更新成功')
      await fetchClusters()
      if (currentCluster.value?.id === data.id) {
        await fetchClusterDetail(data.id)
      }
      return response
    } catch (error) {
      ElMessage.error('集群更新失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除集群
  const deleteCluster = async (id: number) => {
    try {
      loading.value = true
      await clusterApi.deleteCluster(id)
      ElMessage.success('集群删除成功')
      await fetchClusters()
      if (currentCluster.value?.id === id) {
        currentCluster.value = null
      }
    } catch (error) {
      ElMessage.error('集群删除失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 批量删除集群
  const deleteClusters = async (ids: number[]) => {
    try {
      loading.value = true
      await clusterApi.deleteClusters(ids)
      ElMessage.success('批量删除成功')
      await fetchClusters()
      if (currentCluster.value && ids.includes(currentCluster.value.id)) {
        currentCluster.value = null
      }
    } catch (error) {
      ElMessage.error('批量删除失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 创建集群凭据
  const createCredential = async (id: number) => {
    try {
      loading.value = true
      await clusterApi.createCredential(id)
      ElMessage.success('凭据创建成功')
    } catch (error) {
      ElMessage.error('凭据创建失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取集群用户
  const fetchClusterUsers = async (clusterId: number) => {
    try {
      loading.value = true
      const response = await clusterApi.getClusterUserById({ id: clusterId })
      clusterUsers.value = response
      return response
    } catch (error) {
      ElMessage.error('获取集群用户失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取集群角色
  const fetchClusterRoles = async (clusterId: number, roleType: string = 'cluster') => {
    try {
      loading.value = true
      const response = await clusterApi.getClusterRoles({ cluster_id: clusterId, role_type: roleType })
      clusterRoles.value = response
      return response
    } catch (error) {
      ElMessage.error('获取集群角色失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取集群API组
  const fetchClusterApiGroups = async (clusterId: number, apiType: string = 'cluster') => {
    try {
      loading.value = true
      const response = await clusterApi.getClusterApiGroups({ cluster_id: clusterId, api_type: apiType })
      clusterApiGroups.value = response
      return response
    } catch (error) {
      ElMessage.error('获取集群API组失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取集群命名空间
  const fetchClusterNamespaces = async (clusterId: number) => {
    try {
      loading.value = true
      const response = await clusterApi.getClusterListNamespace({ id: clusterId })
      clusterNamespaces.value = response.namespaces
      return response.namespaces
    } catch (error) {
      ElMessage.error('获取集群命名空间失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 创建集群角色
  const createClusterRole = async (data: CreateClusterRoleRequest) => {
    try {
      loading.value = true
      await clusterApi.createClusterRole(data)
      ElMessage.success('集群角色创建成功')
      await fetchClusterRoles(data.cluster_id)
    } catch (error) {
      ElMessage.error('集群角色创建失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新集群角色
  const updateClusterRole = async (data: UpdateClusterRoleRequest) => {
    try {
      loading.value = true
      await clusterApi.updateClusterRole(data)
      ElMessage.success('集群角色更新成功')
      await fetchClusterRoles(data.cluster_id)
    } catch (error) {
      ElMessage.error('集群角色更新失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除集群角色
  const deleteClusterRole = async (data: DeleteClusterRoleRequest) => {
    try {
      loading.value = true
      await clusterApi.deleteClusterRole(data)
      ElMessage.success('集群角色删除成功')
      await fetchClusterRoles(data.cluster_id)
    } catch (error) {
      ElMessage.error('集群角色删除失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 创建集群用户
  const createClusterUser = async (data: CreateClusterUserRequest) => {
    try {
      loading.value = true
      await clusterApi.createClusterUser(data)
      ElMessage.success('集群用户创建成功')
      await fetchClusterUsers(data.cluster_id)
    } catch (error) {
      ElMessage.error('集群用户创建失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新集群用户
  const updateClusterUser = async (data: UpdateClusterUserRequest) => {
    try {
      loading.value = true
      await clusterApi.updateClusterUser(data)
      ElMessage.success('集群用户更新成功')
      await fetchClusterUsers(data.cluster_id)
    } catch (error) {
      ElMessage.error('集群用户更新失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除集群用户
  const deleteClusterUser = async (data: DeleteClusterUserRequest) => {
    try {
      loading.value = true
      await clusterApi.deleteClusterUser(data)
      ElMessage.success('集群用户删除成功')
      await fetchClusterUsers(data.cluster_id)
    } catch (error) {
      ElMessage.error('集群用户删除失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 重置状态
  const resetState = () => {
    clusters.value = []
    currentCluster.value = null
    clusterUsers.value = []
    clusterRoles.value = []
    clusterNamespaces.value = []
    clusterApiGroups.value = []
    loading.value = false
    total.value = 0
    currentPage.value = 1
    pageSize.value = 10
  }

  return {
    // 状态
    clusters,
    currentCluster,
    clusterUsers,
    clusterRoles,
    clusterNamespaces,
    clusterApiGroups,
    loading,
    total,
    currentPage,
    pageSize,
    
    // 计算属性
    hasClusters,
    currentClusterId,
    
    // 方法
    fetchClusters,
    fetchClusterDetail,
    createCluster,
    updateCluster,
    deleteCluster,
    deleteClusters,
    createCredential,
    fetchClusterUsers,
    fetchClusterRoles,
    fetchClusterApiGroups,
    fetchClusterNamespaces,
    createClusterRole,
    updateClusterRole,
    deleteClusterRole,
    createClusterUser,
    updateClusterUser,
    deleteClusterUser,
    resetState
  }
}) 