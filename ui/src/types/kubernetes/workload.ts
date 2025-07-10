// 工作负载基础类型
export interface WorkloadBase {
  metadata: {
    name: string
    namespace: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
}

// Deployment 相关类型
export interface Deployment extends WorkloadBase {
  spec: {
    replicas: number
    selector: {
      matchLabels: Record<string, string>
    }
    template: {
      metadata: {
        labels: Record<string, string>
      }
      spec: PodSpec
    }
    strategy: {
      type: string
      rollingUpdate?: {
        maxUnavailable: string
        maxSurge: string
      }
    }
  }
  status: {
    replicas: number
    updatedReplicas: number
    readyReplicas: number
    availableReplicas: number
    unavailableReplicas: number
    conditions: Array<{
      type: string
      status: string
      lastTransitionTime: string
      reason?: string
      message?: string
    }>
  }
}

// Pod 相关类型
export interface Pod extends WorkloadBase {
  spec: PodSpec
  status: PodStatus
}

export interface PodStatus {
  phase: string
  podIP?: string
  hostIP?: string
  qosClass?: string
  startTime?: string
  containerStatuses?: ContainerStatus[]
  initContainerStatuses?: ContainerStatus[]
  conditions: Array<{
    type: string
    status: string
    lastTransitionTime: string
    reason?: string
    message?: string
  }>
}

export interface PodSpec {
  containers: Container[]
  initContainers?: Container[]
  volumes?: Volume[]
  restartPolicy: string
  terminationGracePeriodSeconds?: number
  dnsPolicy: string
  nodeName?: string
  nodeSelector?: Record<string, string>
  serviceAccountName?: string
  securityContext?: SecurityContext
  imagePullSecrets?: Array<{
    name: string
  }>
}

export interface Container {
  name: string
  image: string
  command?: string[]
  args?: string[]
  ports?: Array<{
    name?: string
    containerPort: number
    protocol: string
  }>
  env?: Array<{
    name: string
    value?: string
    valueFrom?: {
      fieldRef?: {
        fieldPath: string
      }
      configMapKeyRef?: {
        name: string
        key: string
      }
      secretKeyRef?: {
        name: string
        key: string
      }
    }
  }>
  resources?: {
    limits?: Record<string, string>
    requests?: Record<string, string>
  }
  volumeMounts?: Array<{
    name: string
    mountPath: string
    readOnly?: boolean
  }>
  livenessProbe?: Probe
  readinessProbe?: Probe
  startupProbe?: Probe
}

export interface ContainerStatus {
  name: string
  state: {
    running?: {
      startedAt: string
    }
    waiting?: {
      reason: string
      message?: string
    }
    terminated?: {
      exitCode: number
      reason: string
      startedAt: string
      finishedAt: string
    }
  }
  ready: boolean
  restartCount: number
  image: string
  imageID: string
}

export interface Probe {
  httpGet?: {
    path: string
    port: number
    scheme: string
  }
  tcpSocket?: {
    port: number
  }
  exec?: {
    command: string[]
  }
  initialDelaySeconds?: number
  periodSeconds?: number
  timeoutSeconds?: number
  failureThreshold?: number
  successThreshold?: number
}

export interface Volume {
  name: string
  emptyDir?: {}
  hostPath?: {
    path: string
    type: string
  }
  persistentVolumeClaim?: {
    claimName: string
  }
  configMap?: {
    name: string
    items?: Array<{
      key: string
      path: string
    }>
  }
  secret?: {
    secretName: string
    items?: Array<{
      key: string
      path: string
    }>
  }
}

export interface SecurityContext {
  runAsUser?: number
  runAsGroup?: number
  fsGroup?: number
  runAsNonRoot?: boolean
  capabilities?: {
    add?: string[]
    drop?: string[]
  }
}

// StatefulSet 相关类型
export interface StatefulSet extends WorkloadBase {
  spec: {
    replicas: number
    selector: {
      matchLabels: Record<string, string>
    }
    template: {
      metadata: {
        labels: Record<string, string>
      }
      spec: PodSpec
    }
    serviceName: string
    volumeClaimTemplates?: Array<{
      metadata: {
        name: string
      }
      spec: {
        accessModes: string[]
        resources: {
          requests: Record<string, string>
        }
        storageClassName?: string
      }
    }>
  }
  status: {
    replicas: number
    readyReplicas: number
    currentReplicas: number
    updatedReplicas: number
    currentRevision: string
    updateRevision: string
    collisionCount?: number
  }
}

// DaemonSet 相关类型
export interface DaemonSet extends WorkloadBase {
  spec: {
    selector: {
      matchLabels: Record<string, string>
    }
    template: {
      metadata: {
        labels: Record<string, string>
      }
      spec: PodSpec
    }
    updateStrategy: {
      type: string
      rollingUpdate?: {
        maxUnavailable: string
      }
    }
  }
  status: {
    currentNumberScheduled: number
    numberMisscheduled: number
    desiredNumberScheduled: number
    numberReady: number
    updatedNumberScheduled: number
    numberAvailable: number
    numberUnavailable: number
  }
}

// Job 相关类型
export interface Job extends WorkloadBase {
  spec: {
    template: {
      metadata: {
        labels: Record<string, string>
      }
      spec: PodSpec
    }
    parallelism?: number
    completions?: number
    backoffLimit?: number
    activeDeadlineSeconds?: number
  }
  status: {
    startTime?: string
    completionTime?: string
    active: number
    succeeded: number
    failed: number
  }
}

// CronJob 相关类型
export interface CronJob extends WorkloadBase {
  spec: {
    schedule: string
    jobTemplate: {
      metadata: {
        labels: Record<string, string>
      }
      spec: {
        template: {
          metadata: {
            labels: Record<string, string>
          }
          spec: PodSpec
        }
        parallelism?: number
        completions?: number
        backoffLimit?: number
        activeDeadlineSeconds?: number
      }
    }
    concurrencyPolicy?: string
    suspend?: boolean
    successfulJobsHistoryLimit?: number
    failedJobsHistoryLimit?: number
  }
  status: {
    lastScheduleTime?: string
    lastSuccessfulTime?: string
    active: Array<{
      name: string
      uid: string
    }>
  }
}

// ReplicaSet 相关类型
export interface ReplicaSet extends WorkloadBase {
  spec: {
    replicas: number
    selector: {
      matchLabels: Record<string, string>
    }
    template: {
      metadata: {
        labels: Record<string, string>
      }
      spec: PodSpec
    }
  }
  status: {
    replicas: number
    fullyLabeledReplicas: number
    readyReplicas: number
    availableReplicas: number
    observedGeneration: number
  }
}

// API 请求类型
export interface GetWorkloadListRequest {
  cluster_id: number
  namespace: string
  page?: number
  pageSize?: number
  keyword?: string
}

export interface GetWorkloadDetailRequest {
  cluster_id: number
  namespace: string
  name: string
}

export interface CreateWorkloadRequest {
  cluster_id: number
  namespace: string
  content: any
}

export interface UpdateWorkloadRequest {
  cluster_id: number
  namespace: string
  name: string
  content: any
}

export interface DeleteWorkloadRequest {
  cluster_id: number
  namespace: string
  name: string
}

export interface ScaleWorkloadRequest {
  cluster_id: number
  namespace: string
  name: string
  replicas: number
}

export interface RollbackDeploymentRequest {
  cluster_id: number
  namespace: string
  name: string
  revision: number
}

// API 响应类型
export interface WorkloadListResponse<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface WorkloadDetailResponse<T = any> {
  items: T
}

// Pod 相关 API 类型
export interface PodLogRequest {
  cluster_id: number
  namespace: string
  pod_name: string
  container?: string
  tail_lines?: number
  follow?: boolean
}

export interface PodLogResponse {
  logs: string
}

export interface PodTerminalRequest {
  cluster_id: number
  namespace: string
  pod_name: string
  container?: string
}

export interface PodTerminalResponse {
  terminal_url: string
  token: string
} 