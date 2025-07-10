// Node 相关类型
export interface Node {
  metadata: {
    name: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
  spec: {
    podCIDR?: string
    podCIDRs?: string[]
    providerID?: string
    taints?: Array<{
      key: string
      value?: string
      effect: string
      timeAdded?: string
    }>
    unschedulable?: boolean
  }
  status: {
    capacity: Record<string, string>
    allocatable: Record<string, string>
    phase?: string
    conditions: Array<{
      type: string
      status: string
      lastHeartbeatTime: string
      lastTransitionTime: string
      reason?: string
      message?: string
    }>
    addresses: Array<{
      type: string
      address: string
    }>
    daemonEndpoints?: {
      kubeletEndpoint?: {
        Port: number
      }
    }
    nodeInfo?: {
      machineID: string
      systemUUID: string
      bootID: string
      kernelVersion: string
      osImage: string
      containerRuntimeVersion: string
      kubeletVersion: string
      kubeProxyVersion: string
      operatingSystem: string
      architecture: string
    }
    images?: Array<{
      names: string[]
      sizeBytes: number
    }>
    volumesInUse?: string[]
    volumesAttached?: Array<{
      name: string
      devicePath: string
    }>
    config?: {
      active?: {
        configMap?: {
          name: string
          namespace: string
          uid: string
          kubeletConfigKey: string
          resourceVersion: string
        }
      }
      assigned?: {
        configMap?: {
          name: string
          namespace: string
          uid: string
          kubeletConfigKey: string
          resourceVersion: string
        }
      }
      lastKnownGood?: {
        configMap?: {
          name: string
          namespace: string
          uid: string
          kubeletConfigKey: string
          resourceVersion: string
        }
      }
      error?: string
    }
  }
}

// 节点统计信息
export interface NodeStats {
  node_name: string
  cluster_id: number
  cpu: {
    total_cores: number
    used_cores: number
    usage_percentage: number
  }
  memory: {
    total_bytes: number
    used_bytes: number
    usage_percentage: number
  }
  storage: {
    total_bytes: number
    used_bytes: number
    usage_percentage: number
  }
  pods: {
    total: number
    running: number
    capacity: number
  }
  conditions: Array<{
    type: string
    status: string
    reason?: string
    message?: string
  }>
  taints: Array<{
    key: string
    value?: string
    effect: string
  }>
  addresses: Array<{
    type: string
    address: string
  }>
  node_info: {
    kernel_version: string
    os_image: string
    container_runtime_version: string
    kubelet_version: string
    kube_proxy_version: string
    operating_system: string
    architecture: string
  }
  age: string
  ready: boolean
  schedulable: boolean
}

// 节点管理操作
export interface NodeTaintRequest {
  cluster_id: number
  node_name: string
  taints: Array<{
    key: string
    value?: string
    effect: string
  }>
}

export interface NodeLabelRequest {
  cluster_id: number
  node_name: string
  labels: Record<string, string>
}

export interface NodeAnnotationRequest {
  cluster_id: number
  node_name: string
  annotations: Record<string, string>
}

export interface NodeSchedulableRequest {
  cluster_id: number
  node_name: string
  schedulable: boolean
}

export interface NodeDrainRequest {
  cluster_id: number
  node_name: string
  ignore_daemonsets?: boolean
  delete_emptydir_data?: boolean
  force?: boolean
  grace_period?: number
  timeout?: number
}

export interface NodeCordonRequest {
  cluster_id: number
  node_name: string
  reason?: string
}

export interface NodeUncordonRequest {
  cluster_id: number
  node_name: string
}

// 节点终端连接
export interface NodeTerminalRequest {
  cluster_id: number
  node_name: string
  container_runtime?: string // docker, containerd, cri-o
  shell?: string // bash, sh, zsh
}

export interface NodeTerminalResponse {
  terminal_url: string
  token: string
  expires_at: string
}

// 节点日志
export interface NodeLogRequest {
  cluster_id: number
  node_name: string
  log_type: string // kubelet, kube-proxy, system
  lines?: number
  since?: string
  follow?: boolean
}

export interface NodeLogResponse {
  logs: string
  timestamp: string
}

// 节点事件
export interface NodeEvent {
  type: string
  reason: string
  message: string
  source: {
    component: string
    host: string
  }
  first_timestamp: string
  last_timestamp: string
  count: number
  involved_object: {
    kind: string
    name: string
    namespace?: string
    uid: string
    api_version: string
    resource_version: string
  }
}

// API 请求类型
export interface GetNodeListRequest {
  cluster_id: number
  page?: number
  pageSize?: number
  keyword?: string
  labelSelector?: string
  fieldSelector?: string
}

export interface GetNodeDetailRequest {
  cluster_id: number
  name: string
}

export interface UpdateNodeRequest {
  cluster_id: number
  name: string
  content: any
}

export interface DeleteNodeRequest {
  cluster_id: number
  name: string
}

// API 响应类型
export interface NodeListResponse {
  items: Node[]
  total: number
  page: number
  pageSize: number
}

export interface NodeDetailResponse {
  items: Node
}

export interface NodeStatsResponse {
  items: NodeStats[]
  total: number
}

export interface NodeEventsResponse {
  events: NodeEvent[]
  total: number
}

// 节点监控指标
export interface NodeMetricsRequest {
  cluster_id: number
  node_name?: string
  start_time?: string
  end_time?: string
  step?: string
}

export interface NodeMetricsData {
  node_name: string
  metrics: Array<{
    timestamp: number
    cpu_usage: number
    memory_usage: number
    disk_usage: number
    network_rx: number
    network_tx: number
    pod_count: number
  }>
}

export interface NodeMetricsResponse {
  data: NodeMetricsData[]
  time_range: {
    start: string
    end: string
    step: string
  }
}

// 节点容量规划
export interface NodeCapacityPlanning {
  node_name: string
  cluster_id: number
  current_capacity: {
    cpu: string
    memory: string
    storage: string
    pods: string
  }
  current_allocatable: {
    cpu: string
    memory: string
    storage: string
    pods: string
  }
  current_usage: {
    cpu: number
    memory: number
    storage: number
    pods: number
  }
  recommended_capacity: {
    cpu: string
    memory: string
    storage: string
    pods: string
  }
  scaling_recommendations: Array<{
    resource: string
    current_value: string
    recommended_value: string
    reason: string
    priority: string
  }>
}

export interface NodeCapacityPlanningResponse {
  planning: NodeCapacityPlanning[]
  cluster_summary: {
    total_nodes: number
    total_capacity: {
      cpu: string
      memory: string
      storage: string
      pods: string
    }
    total_allocatable: {
      cpu: string
      memory: string
      storage: string
      pods: string
    }
    total_usage: {
      cpu: number
      memory: number
      storage: number
      pods: number
    }
  }
}