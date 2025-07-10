// 配置资源基础类型
export interface ConfigBase {
  metadata: {
    name: string
    namespace: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
}

// ConfigMap 相关类型
export interface ConfigMap extends ConfigBase {
  data?: Record<string, string>
  binaryData?: Record<string, string>
  immutable?: boolean
}

// Secret 相关类型
export interface Secret extends ConfigBase {
  type: string
  data?: Record<string, string>
  stringData?: Record<string, string>
  immutable?: boolean
}

// ResourceQuota 相关类型
export interface ResourceQuota extends ConfigBase {
  spec: {
    hard?: Record<string, string>
    scopes?: string[]
    scopeSelector?: {
      matchExpressions: Array<{
        scopeName: string
        operator: string
        values?: string[]
      }>
    }
  }
  status: {
    hard?: Record<string, string>
    used?: Record<string, string>
  }
}

// LimitRange 相关类型
export interface LimitRange extends ConfigBase {
  spec: {
    limits: Array<{
      type: string
      max?: Record<string, string>
      min?: Record<string, string>
      default?: Record<string, string>
      defaultRequest?: Record<string, string>
      maxLimitRequestRatio?: Record<string, string>
    }>
  }
}

// HorizontalPodAutoscaler 相关类型
export interface HorizontalPodAutoscaler extends ConfigBase {
  spec: {
    scaleTargetRef: {
      apiVersion: string
      kind: string
      name: string
    }
    minReplicas?: number
    maxReplicas: number
    targetCPUUtilizationPercentage?: number
    metrics?: Array<{
      type: string
      resource?: {
        name: string
        target: {
          type: string
          averageUtilization?: number
          averageValue?: string
          value?: string
        }
      }
      pods?: {
        metric: {
          name: string
          selector?: {
            matchLabels?: Record<string, string>
          }
        }
        target: {
          type: string
          averageValue: string
        }
      }
      object?: {
        metric: {
          name: string
          selector?: {
            matchLabels?: Record<string, string>
          }
        }
        target: {
          type: string
          value: string
        }
        describedObject: {
          apiVersion: string
          kind: string
          name: string
        }
      }
      external?: {
        metric: {
          name: string
          selector?: {
            matchLabels?: Record<string, string>
          }
        }
        target: {
          type: string
          value?: string
          averageValue?: string
        }
      }
    }>
    behavior?: {
      scaleDown?: ScalingPolicy
      scaleUp?: ScalingPolicy
    }
  }
  status: {
    observedGeneration?: number
    lastScaleTime?: string
    currentReplicas: number
    desiredReplicas: number
    currentMetrics?: Array<{
      type: string
      resource?: {
        name: string
        current: {
          averageUtilization?: number
          averageValue?: string
        }
      }
    }>
    conditions: Array<{
      type: string
      status: string
      lastTransitionTime: string
      reason?: string
      message?: string
    }>
  }
}

interface ScalingPolicy {
  stabilizationWindowSeconds?: number
  selectPolicy?: string
  policies?: Array<{
    type: string
    value: number
    periodSeconds: number
  }>
}

// PodDisruptionBudget 相关类型
export interface PodDisruptionBudget extends ConfigBase {
  spec: {
    minAvailable?: number | string
    maxUnavailable?: number | string
    selector?: {
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }
    unhealthyPodEvictionPolicy?: string
  }
  status: {
    observedGeneration?: number
    disruptionsAllowed: number
    currentHealthy: number
    desiredHealthy: number
    expectedPods: number
    conditions?: Array<{
      type: string
      status: string
      lastTransitionTime: string
      reason?: string
      message?: string
      observedGeneration?: number
    }>
    disruptedPods?: Record<string, string>
  }
}

// VerticalPodAutoscaler 相关类型
export interface VerticalPodAutoscaler extends ConfigBase {
  spec: {
    targetRef: {
      apiVersion: string
      kind: string
      name: string
    }
    updatePolicy?: {
      updateMode: string
    }
    resourcePolicy?: {
      containerPolicies?: Array<{
        containerName?: string
        mode?: string
        minAllowed?: Record<string, string>
        maxAllowed?: Record<string, string>
        controlledResources?: string[]
        controlledValues?: string
      }>
    }
  }
  status: {
    conditions?: Array<{
      type: string
      status: string
      lastTransitionTime: string
      reason?: string
      message?: string
    }>
    recommendation?: {
      containerRecommendations?: Array<{
        containerName: string
        target: Record<string, string>
        lowerBound?: Record<string, string>
        upperBound?: Record<string, string>
        uncappedTarget?: Record<string, string>
      }>
    }
  }
}

// Priority 相关类型
export interface PriorityClass extends Omit<ConfigBase, 'namespace'> {
  value: number
  globalDefault?: boolean
  description?: string
  preemptionPolicy?: string
}

// API 请求类型
export interface GetConfigListRequest {
  cluster_id: number
  namespace?: string
  page?: number
  pageSize?: number
  keyword?: string
  labelSelector?: string
  fieldSelector?: string
}

export interface GetConfigDetailRequest {
  cluster_id: number
  namespace: string
  name: string
}

export interface CreateConfigRequest {
  cluster_id: number
  namespace: string
  content: any
}

export interface UpdateConfigRequest {
  cluster_id: number
  namespace: string
  name: string
  content: any
}

export interface DeleteConfigRequest {
  cluster_id: number
  namespace: string
  name: string
}

// ConfigMap 和 Secret 特殊操作
export interface CreateConfigMapFromFileRequest {
  cluster_id: number
  namespace: string
  name: string
  files: File[]
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

export interface CreateSecretFromFileRequest {
  cluster_id: number
  namespace: string
  name: string
  type: string
  files: File[]
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

export interface CreateTLSSecretRequest {
  cluster_id: number
  namespace: string
  name: string
  cert_file: File
  key_file: File
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

export interface CreateDockerSecretRequest {
  cluster_id: number
  namespace: string
  name: string
  docker_server: string
  docker_username: string
  docker_password: string
  docker_email?: string
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

// HPA 操作请求
export interface ScaleHPARequest {
  cluster_id: number
  namespace: string
  name: string
  min_replicas: number
  max_replicas: number
  target_cpu_percentage?: number
}

// API 响应类型
export interface ConfigListResponse<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface ConfigDetailResponse<T = any> {
  items: T
}

// 配置统计信息
export interface ConfigMetrics {
  configmap_count: number
  secret_count: number
  resource_quota_count: number
  limit_range_count: number
  hpa_count: number
  pdb_count: number
  vpa_count: number
  priority_class_count: number
}

export interface ConfigMetricsResponse {
  metrics: ConfigMetrics
  timestamp: string
}