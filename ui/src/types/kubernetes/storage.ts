// 存储资源基础类型
export interface StorageBase {
  metadata: {
    name: string
    namespace?: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
}

// PersistentVolume 相关类型
export interface PersistentVolume extends StorageBase {
  spec: {
    capacity: {
      storage: string
    }
    accessModes: string[]
    persistentVolumeReclaimPolicy: string
    storageClassName?: string
    volumeMode?: string
    mountOptions?: string[]
    nodeAffinity?: {
      required?: {
        nodeSelectorTerms: Array<{
          matchExpressions?: Array<{
            key: string
            operator: string
            values?: string[]
          }>
          matchFields?: Array<{
            key: string
            operator: string
            values?: string[]
          }>
        }>
      }
    }
    // Volume 源配置
    hostPath?: {
      path: string
      type?: string
    }
    nfs?: {
      server: string
      path: string
      readOnly?: boolean
    }
    iscsi?: {
      targetPortal: string
      iqn: string
      lun: number
      fsType?: string
      readOnly?: boolean
    }
    cephfs?: {
      monitors: string[]
      path?: string
      user?: string
      secretFile?: string
      secretRef?: {
        name: string
      }
      readOnly?: boolean
    }
    rbd?: {
      monitors: string[]
      image: string
      fsType?: string
      pool?: string
      user?: string
      keyring?: string
      secretRef?: {
        name: string
      }
      readOnly?: boolean
    }
    csi?: {
      driver: string
      volumeHandle: string
      readOnly?: boolean
      fsType?: string
      volumeAttributes?: Record<string, string>
      controllerPublishSecretRef?: {
        name: string
        namespace: string
      }
      nodeStageSecretRef?: {
        name: string
        namespace: string
      }
      nodePublishSecretRef?: {
        name: string
        namespace: string
      }
      controllerExpandSecretRef?: {
        name: string
        namespace: string
      }
    }
  }
  status: {
    phase: string
    message?: string
    reason?: string
  }
}

// PersistentVolumeClaim 相关类型
export interface PersistentVolumeClaim extends StorageBase {
  spec: {
    accessModes: string[]
    resources: {
      requests: {
        storage: string
      }
      limits?: {
        storage: string
      }
    }
    storageClassName?: string
    volumeMode?: string
    selector?: {
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }
    volumeName?: string
    dataSource?: {
      name: string
      kind: string
      apiGroup?: string
    }
    dataSourceRef?: {
      name: string
      kind: string
      apiGroup?: string
      namespace?: string
    }
  }
  status: {
    phase: string
    accessModes?: string[]
    capacity?: {
      storage: string
    }
    conditions?: Array<{
      type: string
      status: string
      lastTransitionTime: string
      reason?: string
      message?: string
    }>
    allocatedResources?: {
      storage: string
    }
    resizeStatus?: string
  }
}

// StorageClass 相关类型
export interface StorageClass extends Omit<StorageBase, 'namespace'> {
  provisioner: string
  parameters?: Record<string, string>
  reclaimPolicy?: string
  allowVolumeExpansion?: boolean
  volumeBindingMode?: string
  allowedTopologies?: Array<{
    matchLabelExpressions?: Array<{
      key: string
      values: string[]
    }>
  }>
  mountOptions?: string[]
}

// VolumeSnapshot 相关类型
export interface VolumeSnapshot extends StorageBase {
  spec: {
    source: {
      persistentVolumeClaimName?: string
      volumeSnapshotContentName?: string
    }
    volumeSnapshotClassName?: string
  }
  status: {
    boundVolumeSnapshotContentName?: string
    creationTime?: string
    readyToUse?: boolean
    restoreSize?: string
    error?: {
      time: string
      message: string
    }
  }
}

// VolumeSnapshotClass 相关类型
export interface VolumeSnapshotClass extends Omit<StorageBase, 'namespace'> {
  driver: string
  parameters?: Record<string, string>
  deletionPolicy: string
}

// CSIDriver 相关类型
export interface CSIDriver extends Omit<StorageBase, 'namespace'> {
  spec: {
    attachRequired?: boolean
    podInfoOnMount?: boolean
    volumeLifecycleModes?: string[]
    storageCapacity?: boolean
    fsGroupPolicy?: string
    tokenRequests?: Array<{
      audience: string
      expirationSeconds?: number
    }>
    requiresRepublish?: boolean
    seLinuxMount?: boolean
  }
}

// CSINode 相关类型
export interface CSINode extends Omit<StorageBase, 'namespace'> {
  spec: {
    drivers: Array<{
      name: string
      nodeID: string
      topologyKeys?: string[]
      allocatable?: {
        count: number
      }
    }>
  }
}

// API 请求类型
export interface GetStorageListRequest {
  cluster_id: number
  namespace?: string
  page?: number
  pageSize?: number
  keyword?: string
  labelSelector?: string
  fieldSelector?: string
}

export interface GetStorageDetailRequest {
  cluster_id: number
  namespace?: string
  name: string
}

export interface CreateStorageRequest {
  cluster_id: number
  namespace?: string
  content: any
}

export interface UpdateStorageRequest {
  cluster_id: number
  namespace?: string
  name: string
  content: any
}

export interface DeleteStorageRequest {
  cluster_id: number
  namespace?: string
  name: string
}

// 存储相关操作请求
export interface ExpandPVCRequest {
  cluster_id: number
  namespace: string
  name: string
  storage: string
}

export interface CreateSnapshotRequest {
  cluster_id: number
  namespace: string
  pvc_name: string
  snapshot_name: string
  snapshot_class?: string
}

export interface RestoreFromSnapshotRequest {
  cluster_id: number
  namespace: string
  snapshot_name: string
  pvc_name: string
  storage_class?: string
  access_modes: string[]
  storage: string
}

// API 响应类型
export interface StorageListResponse<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface StorageDetailResponse<T = any> {
  items: T
}

// 存储统计信息
export interface StorageMetrics {
  total_capacity: string
  used_capacity: string
  available_capacity: string
  usage_percentage: number
  pv_count: number
  pvc_count: number
  storage_class_count: number
}

export interface StorageMetricsResponse {
  metrics: StorageMetrics
  timestamp: string
}