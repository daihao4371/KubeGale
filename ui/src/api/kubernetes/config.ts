import service from '@/api/request'
import type {
  // 配置资源相关类型
  ConfigMap,
  Secret,
  ResourceQuota,
  LimitRange,
  HorizontalPodAutoscaler,
  PodDisruptionBudget,
  VerticalPodAutoscaler,
  PriorityClass,
  GetConfigListRequest,
  GetConfigDetailRequest,
  CreateConfigRequest,
  UpdateConfigRequest,
  DeleteConfigRequest,
  CreateConfigMapFromFileRequest,
  CreateSecretFromFileRequest,
  CreateTLSSecretRequest,
  CreateDockerSecretRequest,
  ScaleHPARequest,
  ConfigListResponse,
  ConfigDetailResponse,
  ConfigMetricsResponse
} from '@/types/kubernetes/config'

export const configApi = {
  // ==================== ConfigMap 相关 API ====================
  
  // 获取 ConfigMap 列表
  getConfigMapList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<ConfigMap>>('/kubernetes/configmap', { params })
  },

  // 获取 ConfigMap 详情
  getConfigMapDetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<ConfigMap>>('/kubernetes/configmap/detail', { params })
  },

  // 创建 ConfigMap
  createConfigMap(data: CreateConfigRequest) {
    return service.post<ConfigMap>('/kubernetes/configmap', data)
  },

  // 从文件创建 ConfigMap
  createConfigMapFromFile(data: CreateConfigMapFromFileRequest) {
    const formData = new FormData()
    formData.append('cluster_id', data.cluster_id.toString())
    formData.append('namespace', data.namespace)
    formData.append('name', data.name)
    data.files.forEach(file => {
      formData.append('files', file)
    })
    if (data.labels) {
      formData.append('labels', JSON.stringify(data.labels))
    }
    if (data.annotations) {
      formData.append('annotations', JSON.stringify(data.annotations))
    }
    return service.post<ConfigMap>('/kubernetes/configmap/from-file', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // 更新 ConfigMap
  updateConfigMap(data: UpdateConfigRequest) {
    return service.put<ConfigMap>('/kubernetes/configmap', data)
  },

  // 删除 ConfigMap
  deleteConfigMap(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/configmap', { data })
  },

  // ==================== Secret 相关 API ====================
  
  // 获取 Secret 列表
  getSecretList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<Secret>>('/kubernetes/secret', { params })
  },

  // 获取 Secret 详情
  getSecretDetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<Secret>>('/kubernetes/secret/detail', { params })
  },

  // 创建 Secret
  createSecret(data: CreateConfigRequest) {
    return service.post<Secret>('/kubernetes/secret', data)
  },

  // 从文件创建 Secret
  createSecretFromFile(data: CreateSecretFromFileRequest) {
    const formData = new FormData()
    formData.append('cluster_id', data.cluster_id.toString())
    formData.append('namespace', data.namespace)
    formData.append('name', data.name)
    formData.append('type', data.type)
    data.files.forEach(file => {
      formData.append('files', file)
    })
    if (data.labels) {
      formData.append('labels', JSON.stringify(data.labels))
    }
    if (data.annotations) {
      formData.append('annotations', JSON.stringify(data.annotations))
    }
    return service.post<Secret>('/kubernetes/secret/from-file', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // 创建 TLS Secret
  createTLSSecret(data: CreateTLSSecretRequest) {
    const formData = new FormData()
    formData.append('cluster_id', data.cluster_id.toString())
    formData.append('namespace', data.namespace)
    formData.append('name', data.name)
    formData.append('cert_file', data.cert_file)
    formData.append('key_file', data.key_file)
    if (data.labels) {
      formData.append('labels', JSON.stringify(data.labels))
    }
    if (data.annotations) {
      formData.append('annotations', JSON.stringify(data.annotations))
    }
    return service.post<Secret>('/kubernetes/secret/tls', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // 创建 Docker Secret
  createDockerSecret(data: CreateDockerSecretRequest) {
    return service.post<Secret>('/kubernetes/secret/docker', data)
  },

  // 更新 Secret
  updateSecret(data: UpdateConfigRequest) {
    return service.put<Secret>('/kubernetes/secret', data)
  },

  // 删除 Secret
  deleteSecret(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/secret', { data })
  },

  // ==================== ResourceQuota 相关 API ====================
  
  // 获取 ResourceQuota 列表
  getResourceQuotaList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<ResourceQuota>>('/kubernetes/resourcequota', { params })
  },

  // 获取 ResourceQuota 详情
  getResourceQuotaDetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<ResourceQuota>>('/kubernetes/resourcequota/detail', { params })
  },

  // 创建 ResourceQuota
  createResourceQuota(data: CreateConfigRequest) {
    return service.post<ResourceQuota>('/kubernetes/resourcequota', data)
  },

  // 更新 ResourceQuota
  updateResourceQuota(data: UpdateConfigRequest) {
    return service.put<ResourceQuota>('/kubernetes/resourcequota', data)
  },

  // 删除 ResourceQuota
  deleteResourceQuota(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/resourcequota', { data })
  },

  // ==================== LimitRange 相关 API ====================
  
  // 获取 LimitRange 列表
  getLimitRangeList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<LimitRange>>('/kubernetes/limitrange', { params })
  },

  // 获取 LimitRange 详情
  getLimitRangeDetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<LimitRange>>('/kubernetes/limitrange/detail', { params })
  },

  // 创建 LimitRange
  createLimitRange(data: CreateConfigRequest) {
    return service.post<LimitRange>('/kubernetes/limitrange', data)
  },

  // 更新 LimitRange
  updateLimitRange(data: UpdateConfigRequest) {
    return service.put<LimitRange>('/kubernetes/limitrange', data)
  },

  // 删除 LimitRange
  deleteLimitRange(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/limitrange', { data })
  },

  // ==================== HorizontalPodAutoscaler 相关 API ====================
  
  // 获取 HPA 列表
  getHPAList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<HorizontalPodAutoscaler>>('/kubernetes/hpa', { params })
  },

  // 获取 HPA 详情
  getHPADetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<HorizontalPodAutoscaler>>('/kubernetes/hpa/detail', { params })
  },

  // 创建 HPA
  createHPA(data: CreateConfigRequest) {
    return service.post<HorizontalPodAutoscaler>('/kubernetes/hpa', data)
  },

  // 更新 HPA
  updateHPA(data: UpdateConfigRequest) {
    return service.put<HorizontalPodAutoscaler>('/kubernetes/hpa', data)
  },

  // 删除 HPA
  deleteHPA(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/hpa', { data })
  },

  // 缩放 HPA
  scaleHPA(data: ScaleHPARequest) {
    return service.patch('/kubernetes/hpa/scale', data)
  },

  // ==================== PodDisruptionBudget 相关 API ====================
  
  // 获取 PDB 列表
  getPDBList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<PodDisruptionBudget>>('/kubernetes/pdb', { params })
  },

  // 获取 PDB 详情
  getPDBDetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<PodDisruptionBudget>>('/kubernetes/pdb/detail', { params })
  },

  // 创建 PDB
  createPDB(data: CreateConfigRequest) {
    return service.post<PodDisruptionBudget>('/kubernetes/pdb', data)
  },

  // 更新 PDB
  updatePDB(data: UpdateConfigRequest) {
    return service.put<PodDisruptionBudget>('/kubernetes/pdb', data)
  },

  // 删除 PDB
  deletePDB(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/pdb', { data })
  },

  // ==================== VerticalPodAutoscaler 相关 API ====================
  
  // 获取 VPA 列表
  getVPAList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<VerticalPodAutoscaler>>('/kubernetes/vpa', { params })
  },

  // 获取 VPA 详情
  getVPADetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<VerticalPodAutoscaler>>('/kubernetes/vpa/detail', { params })
  },

  // 创建 VPA
  createVPA(data: CreateConfigRequest) {
    return service.post<VerticalPodAutoscaler>('/kubernetes/vpa', data)
  },

  // 更新 VPA
  updateVPA(data: UpdateConfigRequest) {
    return service.put<VerticalPodAutoscaler>('/kubernetes/vpa', data)
  },

  // 删除 VPA
  deleteVPA(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/vpa', { data })
  },

  // ==================== PriorityClass 相关 API ====================
  
  // 获取 PriorityClass 列表
  getPriorityClassList(params: GetConfigListRequest) {
    return service.get<ConfigListResponse<PriorityClass>>('/kubernetes/priorityclass', { params })
  },

  // 获取 PriorityClass 详情
  getPriorityClassDetail(params: GetConfigDetailRequest) {
    return service.get<ConfigDetailResponse<PriorityClass>>('/kubernetes/priorityclass/detail', { params })
  },

  // 创建 PriorityClass
  createPriorityClass(data: CreateConfigRequest) {
    return service.post<PriorityClass>('/kubernetes/priorityclass', data)
  },

  // 更新 PriorityClass
  updatePriorityClass(data: UpdateConfigRequest) {
    return service.put<PriorityClass>('/kubernetes/priorityclass', data)
  },

  // 删除 PriorityClass
  deletePriorityClass(data: DeleteConfigRequest) {
    return service.delete('/kubernetes/priorityclass', { data })
  },

  // ==================== 配置监控和统计 API ====================
  
  // 获取配置指标
  getConfigMetrics(params: { cluster_id: number; namespace?: string }) {
    return service.get<ConfigMetricsResponse>('/kubernetes/config/metrics', { params })
  },

  // 获取配置使用情况
  getConfigUsage(params: { cluster_id: number; namespace?: string; resource_type?: string }) {
    return service.get('/kubernetes/config/usage', { params })
  },

  // 配置对比
  compareConfigs(data: {
    cluster_id: number
    namespace: string
    resource_type: string
    source_name: string
    target_name: string
  }) {
    return service.post('/kubernetes/config/compare', data)
  },

  // 配置模板
  getConfigTemplates(params: { cluster_id: number; resource_type: string }) {
    return service.get('/kubernetes/config/templates', { params })
  },

  // 配置验证
  validateConfig(data: { cluster_id: number; namespace: string; resource_type: string; content: any }) {
    return service.post('/kubernetes/config/validate', data)
  },

  // 配置导入导出
  exportConfigs(data: { cluster_id: number; namespace?: string; resource_types: string[] }) {
    return service.post('/kubernetes/config/export', data, { responseType: 'blob' })
  },

  importConfigs(data: { cluster_id: number; namespace: string; file: File }) {
    const formData = new FormData()
    formData.append('cluster_id', data.cluster_id.toString())
    formData.append('namespace', data.namespace)
    formData.append('file', data.file)
    return service.post('/kubernetes/config/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}