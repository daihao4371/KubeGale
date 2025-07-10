import service from '@/api/request'
import type {
  // 网络资源相关类型
  Service,
  Ingress,
  Endpoint,
  NetworkPolicy,
  GetNetworkListRequest,
  GetNetworkDetailRequest,
  CreateNetworkRequest,
  UpdateNetworkRequest,
  DeleteNetworkRequest,
  NetworkListResponse,
  NetworkDetailResponse
} from '@/types/kubernetes/network'

export const networkApi = {
  // ==================== Service 相关 API ====================
  
  // 获取 Service 列表
  getServiceList(params: GetNetworkListRequest) {
    return service.get<NetworkListResponse<Service>>('/kubernetes/service', { params })
  },

  // 获取 Service 详情
  getServiceDetail(params: GetNetworkDetailRequest) {
    return service.get<NetworkDetailResponse<Service>>('/kubernetes/service/detail', { params })
  },

  // 创建 Service
  createService(data: CreateNetworkRequest) {
    return service.post<Service>('/kubernetes/service', data)
  },

  // 更新 Service
  updateService(data: UpdateNetworkRequest) {
    return service.put<Service>('/kubernetes/service', data)
  },

  // 删除 Service
  deleteService(data: DeleteNetworkRequest) {
    return service.delete('/kubernetes/service', { data })
  },

  // ==================== Ingress 相关 API ====================
  
  // 获取 Ingress 列表
  getIngressList(params: GetNetworkListRequest) {
    return service.get<NetworkListResponse<Ingress>>('/kubernetes/ingress', { params })
  },

  // 获取 Ingress 详情
  getIngressDetail(params: GetNetworkDetailRequest) {
    return service.get<NetworkDetailResponse<Ingress>>('/kubernetes/ingress/detail', { params })
  },

  // 创建 Ingress
  createIngress(data: CreateNetworkRequest) {
    return service.post<Ingress>('/kubernetes/ingress', data)
  },

  // 更新 Ingress
  updateIngress(data: UpdateNetworkRequest) {
    return service.put<Ingress>('/kubernetes/ingress', data)
  },

  // 删除 Ingress
  deleteIngress(data: DeleteNetworkRequest) {
    return service.delete('/kubernetes/ingress', { data })
  },

  // ==================== Endpoint 相关 API ====================
  
  // 获取 Endpoint 列表
  getEndpointList(params: GetNetworkListRequest) {
    return service.get<NetworkListResponse<Endpoint>>('/kubernetes/endpoint', { params })
  },

  // 获取 Endpoint 详情
  getEndpointDetail(params: GetNetworkDetailRequest) {
    return service.get<NetworkDetailResponse<Endpoint>>('/kubernetes/endpoint/detail', { params })
  },

  // 创建 Endpoint
  createEndpoint(data: CreateNetworkRequest) {
    return service.post<Endpoint>('/kubernetes/endpoint', data)
  },

  // 更新 Endpoint
  updateEndpoint(data: UpdateNetworkRequest) {
    return service.put<Endpoint>('/kubernetes/endpoint', data)
  },

  // 删除 Endpoint
  deleteEndpoint(data: DeleteNetworkRequest) {
    return service.delete('/kubernetes/endpoint', { data })
  },

  // ==================== NetworkPolicy 相关 API ====================
  
  // 获取 NetworkPolicy 列表
  getNetworkPolicyList(params: GetNetworkListRequest) {
    return service.get<NetworkListResponse<NetworkPolicy>>('/kubernetes/networkpolicy', { params })
  },

  // 获取 NetworkPolicy 详情
  getNetworkPolicyDetail(params: GetNetworkDetailRequest) {
    return service.get<NetworkDetailResponse<NetworkPolicy>>('/kubernetes/networkpolicy/detail', { params })
  },

  // 创建 NetworkPolicy
  createNetworkPolicy(data: CreateNetworkRequest) {
    return service.post<NetworkPolicy>('/kubernetes/networkpolicy', data)
  },

  // 更新 NetworkPolicy
  updateNetworkPolicy(data: UpdateNetworkRequest) {
    return service.put<NetworkPolicy>('/kubernetes/networkpolicy', data)
  },

  // 删除 NetworkPolicy
  deleteNetworkPolicy(data: DeleteNetworkRequest) {
    return service.delete('/kubernetes/networkpolicy', { data })
  },

  // ==================== 网络诊断和监控 API ====================
  
  // 获取网络拓扑
  getNetworkTopology(params: { cluster_id: number; namespace?: string }) {
    return service.get('/kubernetes/network/topology', { params })
  },

  // 网络连通性测试
  testNetworkConnectivity(data: {
    cluster_id: number
    source_namespace: string
    source_pod: string
    target_namespace: string
    target_pod: string
    port: number
    protocol: string
  }) {
    return service.post('/kubernetes/network/connectivity-test', data)
  },

  // 获取网络策略规则
  getNetworkPolicyRules(params: { cluster_id: number; namespace: string; pod_name: string }) {
    return service.get('/kubernetes/network/policy-rules', { params })
  },

  // 获取网络流量统计
  getNetworkTrafficStats(params: {
    cluster_id: number
    namespace?: string
    start_time?: string
    end_time?: string
    granularity?: string
  }) {
    return service.get('/kubernetes/network/traffic-stats', { params })
  },

  // DNS 查询测试
  testDNSLookup(data: {
    cluster_id: number
    namespace: string
    pod_name?: string
    domain: string
    record_type?: string
  }) {
    return service.post('/kubernetes/network/dns-lookup', data)
  },

  // 获取服务网格信息 (Istio)
  getServiceMeshInfo(params: { cluster_id: number; namespace?: string }) {
    return service.get('/kubernetes/network/service-mesh', { params })
  },

  // 获取负载均衡器状态
  getLoadBalancerStatus(params: { cluster_id: number; namespace: string; service_name: string }) {
    return service.get('/kubernetes/network/loadbalancer-status', { params })
  }
}