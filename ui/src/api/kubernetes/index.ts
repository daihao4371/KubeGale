// Kubernetes API 统一入口

// 集群管理 API
export { clusterApi } from './cluster'

// 工作负载 API
export { workloadApi } from './workload'

// 网络资源 API
export { networkApi } from './network'

// 存储资源 API
export { storageApi } from './storage'

// 配置管理 API
export { configApi } from './config'

// RBAC 权限 API
export { rbacApi } from './rbac'

// 节点管理 API
export { nodeApi } from './node'

// 整合所有 Kubernetes API
export const kubernetesApi = {
  cluster: clusterApi,
  workload: workloadApi,
  network: networkApi,
  storage: storageApi,
  config: configApi,
  rbac: rbacApi,
  node: nodeApi
}

// 默认导出
export default kubernetesApi