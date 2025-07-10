// 监控相关类型
export interface MonitoringBase {
  metadata: {
    name: string
    namespace?: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
}

// Prometheus 相关类型
export interface Prometheus extends MonitoringBase {
  spec: {
    replicas?: number
    version?: string
    serviceAccountName?: string
    serviceMonitorSelector?: {
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }
    ruleSelector?: {
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }
    resources?: {
      requests?: Record<string, string>
      limits?: Record<string, string>
    }
    storage?: {
      volumeClaimTemplate?: {
        spec: {
          accessModes: string[]
          resources: {
            requests: {
              storage: string
            }
          }
          storageClassName?: string
        }
      }
    }
    retention?: string
    retentionSize?: string
    logLevel?: string
    scrapeInterval?: string
    evaluationInterval?: string
    externalUrl?: string
    routePrefix?: string
    alerting?: {
      alertmanagers: Array<{
        namespace: string
        name: string
        port: string
        pathPrefix?: string
      }>
    }
    additionalScrapeConfigs?: {
      name: string
      key: string
    }
    remoteWrite?: Array<{
      url: string
      remoteTimeout?: string
      writeRelabelConfigs?: Array<{
        sourceLabels?: string[]
        separator?: string
        regex?: string
        targetLabel?: string
        replacement?: string
        action?: string
      }>
      basicAuth?: {
        username: {
          name: string
          key: string
        }
        password: {
          name: string
          key: string
        }
      }
      bearerToken?: string
      bearerTokenFile?: string
      tlsConfig?: {
        caFile?: string
        certFile?: string
        keyFile?: string
        serverName?: string
        insecureSkipVerify?: boolean
      }
    }>
  }
  status?: {
    availableReplicas: number
    paused: boolean
    replicas: number
    updatedReplicas: number
  }
}

// ServiceMonitor 相关类型
export interface ServiceMonitor extends MonitoringBase {
  spec: {
    selector: {
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }
    namespaceSelector?: {
      matchNames?: string[]
      any?: boolean
    }
    endpoints: Array<{
      port?: string
      targetPort?: string | number
      path?: string
      scheme?: string
      params?: Record<string, string[]>
      interval?: string
      scrapeTimeout?: string
      honorLabels?: boolean
      honorTimestamps?: boolean
      metricRelabelings?: Array<{
        sourceLabels?: string[]
        separator?: string
        regex?: string
        targetLabel?: string
        replacement?: string
        action?: string
      }>
      relabelings?: Array<{
        sourceLabels?: string[]
        separator?: string
        regex?: string
        targetLabel?: string
        replacement?: string
        action?: string
      }>
      basicAuth?: {
        username?: {
          name: string
          key: string
        }
        password?: {
          name: string
          key: string
        }
      }
      bearerTokenFile?: string
      bearerTokenSecret?: {
        name: string
        key: string
      }
      tlsConfig?: {
        caFile?: string
        certFile?: string
        keyFile?: string
        serverName?: string
        insecureSkipVerify?: boolean
      }
    }>
    jobLabel?: string
    targetLabels?: string[]
    podTargetLabels?: string[]
    sampleLimit?: number
    targetLimit?: number
    labelLimit?: number
    labelNameLengthLimit?: number
    labelValueLengthLimit?: number
  }
}

// PodMonitor 相关类型
export interface PodMonitor extends MonitoringBase {
  spec: {
    selector: {
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }
    namespaceSelector?: {
      matchNames?: string[]
      any?: boolean
    }
    podMetricsEndpoints: Array<{
      port?: string
      targetPort?: string | number
      path?: string
      scheme?: string
      params?: Record<string, string[]>
      interval?: string
      scrapeTimeout?: string
      honorLabels?: boolean
      honorTimestamps?: boolean
      metricRelabelings?: Array<{
        sourceLabels?: string[]
        separator?: string
        regex?: string
        targetLabel?: string
        replacement?: string
        action?: string
      }>
      relabelings?: Array<{
        sourceLabels?: string[]
        separator?: string
        regex?: string
        targetLabel?: string
        replacement?: string
        action?: string
      }>
    }>
    jobLabel?: string
    podTargetLabels?: string[]
    sampleLimit?: number
    targetLimit?: number
    labelLimit?: number
    labelNameLengthLimit?: number
    labelValueLengthLimit?: number
  }
}

// PrometheusRule 相关类型
export interface PrometheusRule extends MonitoringBase {
  spec: {
    groups: Array<{
      name: string
      interval?: string
      rules: Array<{
        record?: string
        alert?: string
        expr: string
        for?: string
        labels?: Record<string, string>
        annotations?: Record<string, string>
      }>
    }>
  }
}

// Alertmanager 相关类型
export interface Alertmanager extends MonitoringBase {
  spec: {
    replicas?: number
    version?: string
    serviceAccountName?: string
    configSecret?: string
    secrets?: string[]
    configMaps?: string[]
    storage?: {
      volumeClaimTemplate?: {
        spec: {
          accessModes: string[]
          resources: {
            requests: {
              storage: string
            }
          }
          storageClassName?: string
        }
      }
    }
    retention?: string
    logLevel?: string
    resources?: {
      requests?: Record<string, string>
      limits?: Record<string, string>
    }
    routePrefix?: string
    externalUrl?: string
    web?: {
      httpConfig?: {
        http2?: boolean
        headers?: Record<string, string>
      }
      tlsConfig?: {
        keySecret?: {
          name: string
          key: string
        }
        cert?: {
          secret?: {
            name: string
            key: string
          }
        }
        clientAuthType?: string
        clientCA?: {
          secret?: {
            name: string
            key: string
          }
        }
      }
    }
  }
  status?: {
    availableReplicas: number
    paused: boolean
    replicas: number
    updatedReplicas: number
  }
}

// Grafana 相关类型
export interface GrafanaDashboard extends MonitoringBase {
  spec: {
    json?: string
    jsonnet?: string
    gzipJson?: string
    url?: string
    configMapRef?: {
      name: string
      key: string
    }
    datasources?: Array<{
      inputName: string
      datasourceName: string
    }>
    plugins?: Array<{
      name: string
      version: string
    }>
    folder?: string
    customFolderName?: string
  }
  status?: {
    lastResyncTime?: string
    uid?: string
    hash?: string
  }
}

// 指标查询相关类型
export interface MetricsQueryRequest {
  cluster_id: number
  query: string
  start?: string
  end?: string
  step?: string
  timeout?: string
}

export interface MetricsQueryResponse {
  status: string
  data: {
    resultType: string
    result: Array<{
      metric: Record<string, string>
      value?: [number, string]
      values?: Array<[number, string]>
    }>
  }
  warnings?: string[]
}

export interface MetricsRangeQueryRequest {
  cluster_id: number
  query: string
  start: string
  end: string
  step: string
  timeout?: string
}

// 集群监控指标
export interface ClusterMetrics {
  cluster_id: number
  timestamp: string
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
  nodes: {
    total: number
    ready: number
    not_ready: number
  }
  pods: {
    total: number
    running: number
    pending: number
    failed: number
    succeeded: number
  }
  services: {
    total: number
  }
  ingresses: {
    total: number
  }
  pvs: {
    total: number
    available: number
    bound: number
    released: number
    failed: number
  }
  pvcs: {
    total: number
    pending: number
    bound: number
    lost: number
  }
}

// 节点监控指标
export interface NodeMetrics {
  node_name: string
  cluster_id: number
  timestamp: string
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
  network: {
    rx_bytes: number
    tx_bytes: number
    rx_packets: number
    tx_packets: number
  }
  pods: {
    total: number
    running: number
  }
  conditions: Array<{
    type: string
    status: string
    reason?: string
    message?: string
  }>
}

// Pod 监控指标
export interface PodMetrics {
  pod_name: string
  namespace: string
  cluster_id: number
  timestamp: string
  cpu: {
    used_cores: number
    limit_cores?: number
    request_cores?: number
    usage_percentage?: number
  }
  memory: {
    used_bytes: number
    limit_bytes?: number
    request_bytes?: number
    usage_percentage?: number
  }
  network?: {
    rx_bytes: number
    tx_bytes: number
  }
  containers: Array<{
    name: string
    cpu: {
      used_cores: number
      limit_cores?: number
      request_cores?: number
    }
    memory: {
      used_bytes: number
      limit_bytes?: number
      request_bytes?: number
    }
  }>
}

// API 请求类型
export interface GetMonitoringListRequest {
  cluster_id: number
  namespace?: string
  page?: number
  pageSize?: number
  keyword?: string
  labelSelector?: string
  fieldSelector?: string
}

export interface GetMonitoringDetailRequest {
  cluster_id: number
  namespace?: string
  name: string
}

export interface CreateMonitoringRequest {
  cluster_id: number
  namespace?: string
  content: any
}

export interface UpdateMonitoringRequest {
  cluster_id: number
  namespace?: string
  name: string
  content: any
}

export interface DeleteMonitoringRequest {
  cluster_id: number
  namespace?: string
  name: string
}

// API 响应类型
export interface MonitoringListResponse<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface MonitoringDetailResponse<T = any> {
  items: T
}

export interface ClusterMetricsResponse {
  metrics: ClusterMetrics
}

export interface NodeMetricsResponse {
  metrics: NodeMetrics[]
}

export interface PodMetricsResponse {
  metrics: PodMetrics[]
}