import service from '@/api/request'
import type {
  // 节点相关类型
  Node,
  NodeStats,
  NodeTaintRequest,
  NodeLabelRequest,
  NodeAnnotationRequest,
  NodeSchedulableRequest,
  NodeDrainRequest,
  NodeCordonRequest,
  NodeUncordonRequest,
  NodeTerminalRequest,
  NodeTerminalResponse,
  NodeLogRequest,
  NodeLogResponse,
  NodeEvent,
  GetNodeListRequest,
  GetNodeDetailRequest,
  UpdateNodeRequest,
  DeleteNodeRequest,
  NodeListResponse,
  NodeDetailResponse,
  NodeStatsResponse,
  NodeEventsResponse,
  NodeMetricsRequest,
  NodeMetricsResponse,
  NodeCapacityPlanningResponse
} from '@/types/kubernetes/node'

export const nodeApi = {
  // ==================== 节点基本管理 API ====================
  
  // 获取节点列表
  getNodeList(params: GetNodeListRequest) {
    return service.get<NodeListResponse>('/kubernetes/node', { params })
  },

  // 获取节点详情
  getNodeDetail(params: GetNodeDetailRequest) {
    return service.get<NodeDetailResponse>('/kubernetes/node/detail', { params })
  },

  // 更新节点
  updateNode(data: UpdateNodeRequest) {
    return service.put<Node>('/kubernetes/node', data)
  },

  // 删除节点
  deleteNode(data: DeleteNodeRequest) {
    return service.delete('/kubernetes/node', { data })
  },

  // 获取节点统计信息
  getNodeStats(params: GetNodeListRequest) {
    return service.get<NodeStatsResponse>('/kubernetes/node/stats', { params })
  },

  // ==================== 节点标签和注释管理 API ====================
  
  // 更新节点标签
  updateNodeLabels(data: NodeLabelRequest) {
    return service.patch('/kubernetes/node/labels', data)
  },

  // 更新节点注释
  updateNodeAnnotations(data: NodeAnnotationRequest) {
    return service.patch('/kubernetes/node/annotations', data)
  },

  // ==================== 节点污点管理 API ====================
  
  // 更新节点污点
  updateNodeTaints(data: NodeTaintRequest) {
    return service.patch('/kubernetes/node/taints', data)
  },

  // 添加节点污点
  addNodeTaint(data: {
    cluster_id: number
    node_name: string
    key: string
    value?: string
    effect: string
  }) {
    return service.post('/kubernetes/node/taint', data)
  },

  // 移除节点污点
  removeNodeTaint(data: {
    cluster_id: number
    node_name: string
    key: string
    effect?: string
  }) {
    return service.delete('/kubernetes/node/taint', { data })
  },

  // ==================== 节点调度管理 API ====================
  
  // 设置节点可调度性
  setNodeSchedulable(data: NodeSchedulableRequest) {
    return service.patch('/kubernetes/node/schedulable', data)
  },

  // 禁止调度 (Cordon)
  cordonNode(data: NodeCordonRequest) {
    return service.post('/kubernetes/node/cordon', data)
  },

  // 允许调度 (Uncordon)
  uncordonNode(data: NodeUncordonRequest) {
    return service.post('/kubernetes/node/uncordon', data)
  },

  // 排空节点 (Drain)
  drainNode(data: NodeDrainRequest) {
    return service.post('/kubernetes/node/drain', data)
  },

  // 获取排空状态
  getDrainStatus(params: { cluster_id: number; node_name: string; drain_id: string }) {
    return service.get('/kubernetes/node/drain-status', { params })
  },

  // 取消排空
  cancelDrain(data: { cluster_id: number; node_name: string; drain_id: string }) {
    return service.post('/kubernetes/node/cancel-drain', data)
  },

  // ==================== 节点终端和日志 API ====================
  
  // 节点终端连接
  getNodeTerminal(data: NodeTerminalRequest) {
    return service.post<NodeTerminalResponse>('/kubernetes/node/terminal', data)
  },

  // 获取节点日志
  getNodeLog(params: NodeLogRequest) {
    return service.get<NodeLogResponse>('/kubernetes/node/logs', { params })
  },

  // 获取系统日志
  getSystemLog(params: {
    cluster_id: number
    node_name: string
    service_name: string // kubelet, docker, containerd, etc.
    lines?: number
    since?: string
  }) {
    return service.get('/kubernetes/node/system-logs', { params })
  },

  // ==================== 节点事件和监控 API ====================
  
  // 获取节点事件
  getNodeEvents(params: {
    cluster_id: number
    node_name?: string
    namespace?: string
    since?: string
    page?: number
    pageSize?: number
  }) {
    return service.get<NodeEventsResponse>('/kubernetes/node/events', { params })
  },

  // 获取节点指标
  getNodeMetrics(params: NodeMetricsRequest) {
    return service.get<NodeMetricsResponse>('/kubernetes/node/metrics', { params })
  },

  // 获取节点实时指标
  getNodeRealtimeMetrics(params: { cluster_id: number; node_name?: string }) {
    return service.get('/kubernetes/node/realtime-metrics', { params })
  },

  // 获取节点Top进程
  getNodeTopProcesses(params: {
    cluster_id: number
    node_name: string
    limit?: number
    sort_by?: string // cpu, memory, pid, etc.
  }) {
    return service.get('/kubernetes/node/top-processes', { params })
  },

  // ==================== 节点容量和资源管理 API ====================
  
  // 获取节点容量规划
  getNodeCapacityPlanning(params: { cluster_id: number; node_name?: string }) {
    return service.get<NodeCapacityPlanningResponse>('/kubernetes/node/capacity-planning', { params })
  },

  // 获取节点资源分配
  getNodeResourceAllocation(params: {
    cluster_id: number
    node_name?: string
    resource_type?: string // cpu, memory, storage, pods
  }) {
    return service.get('/kubernetes/node/resource-allocation', { params })
  },

  // 节点资源预订
  reserveNodeResources(data: {
    cluster_id: number
    node_name: string
    cpu?: string
    memory?: string
    storage?: string
  }) {
    return service.post('/kubernetes/node/reserve-resources', data)
  },

  // ==================== 节点维护管理 API ====================
  
  // 节点健康检查
  checkNodeHealth(data: { cluster_id: number; node_name?: string }) {
    return service.post('/kubernetes/node/health-check', data)
  },

  // 节点诊断
  diagnoseNode(data: {
    cluster_id: number
    node_name: string
    check_types: string[] // network, storage, kubelet, docker, etc.
  }) {
    return service.post('/kubernetes/node/diagnose', data)
  },

  // 节点重启
  rebootNode(data: {
    cluster_id: number
    node_name: string
    force?: boolean
    delay?: number
  }) {
    return service.post('/kubernetes/node/reboot', data)
  },

  // 节点维护模式
  setNodeMaintenanceMode(data: {
    cluster_id: number
    node_name: string
    enabled: boolean
    reason?: string
  }) {
    return service.post('/kubernetes/node/maintenance-mode', data)
  },

  // ==================== 节点网络管理 API ====================
  
  // 获取节点网络信息
  getNodeNetworkInfo(params: { cluster_id: number; node_name: string }) {
    return service.get('/kubernetes/node/network-info', { params })
  },

  // 网络连通性测试
  testNodeConnectivity(data: {
    cluster_id: number
    source_node: string
    target_node: string
    test_type: string // ping, traceroute, telnet
    target_host?: string
    target_port?: number
  }) {
    return service.post('/kubernetes/node/connectivity-test', data)
  },

  // 获取节点网络统计
  getNodeNetworkStats(params: {
    cluster_id: number
    node_name: string
    interface?: string
    start_time?: string
    end_time?: string
  }) {
    return service.get('/kubernetes/node/network-stats', { params })
  },

  // ==================== 节点存储管理 API ====================
  
  // 获取节点存储信息
  getNodeStorageInfo(params: { cluster_id: number; node_name: string }) {
    return service.get('/kubernetes/node/storage-info', { params })
  },

  // 获取节点磁盘使用情况
  getNodeDiskUsage(params: {
    cluster_id: number
    node_name: string
    path?: string
  }) {
    return service.get('/kubernetes/node/disk-usage', { params })
  },

  // 清理节点存储
  cleanupNodeStorage(data: {
    cluster_id: number
    node_name: string
    cleanup_types: string[] // docker-images, docker-containers, kubelet-pods, logs
    dry_run?: boolean
  }) {
    return service.post('/kubernetes/node/storage-cleanup', data)
  },

  // ==================== 节点集群管理 API ====================
  
  // 添加节点到集群
  addNodeToCluster(data: {
    cluster_id: number
    node_config: {
      hostname: string
      ip: string
      ssh_user: string
      ssh_password?: string
      ssh_key?: string
      role: string // master, worker
    }
  }) {
    return service.post('/kubernetes/node/add', data)
  },

  // 从集群移除节点
  removeNodeFromCluster(data: {
    cluster_id: number
    node_name: string
    force?: boolean
    cleanup?: boolean
  }) {
    return service.post('/kubernetes/node/remove', data)
  },

  // 获取节点加入命令
  getNodeJoinCommand(params: { cluster_id: number; role: string }) {
    return service.get('/kubernetes/node/join-command', { params })
  },

  // 更新节点证书
  updateNodeCertificate(data: {
    cluster_id: number
    node_name: string
    certificate_type: string // kubelet, kube-proxy
  }) {
    return service.post('/kubernetes/node/update-certificate', data)
  }
}