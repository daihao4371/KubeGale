import service from '@/api/request'
import type {
  // 存储资源相关类型
  PersistentVolume,
  PersistentVolumeClaim,
  StorageClass,
  VolumeSnapshot,
  VolumeSnapshotClass,
  CSIDriver,
  CSINode,
  GetStorageListRequest,
  GetStorageDetailRequest,
  CreateStorageRequest,
  UpdateStorageRequest,
  DeleteStorageRequest,
  ExpandPVCRequest,
  CreateSnapshotRequest,
  RestoreFromSnapshotRequest,
  StorageListResponse,
  StorageDetailResponse,
  StorageMetricsResponse
} from '@/types/kubernetes/storage'

export const storageApi = {
  // ==================== PersistentVolume 相关 API ====================
  
  // 获取 PV 列表
  getPVList(params: GetStorageListRequest) {
    return service.get<StorageListResponse<PersistentVolume>>('/kubernetes/pv', { params })
  },

  // 获取 PV 详情
  getPVDetail(params: GetStorageDetailRequest) {
    return service.get<StorageDetailResponse<PersistentVolume>>('/kubernetes/pvDetails', { params })
  },

  // 创建 PV
  createPV(data: CreateStorageRequest) {
    return service.post<PersistentVolume>('/kubernetes/pv', data)
  },

  // 更新 PV
  updatePV(data: UpdateStorageRequest) {
    return service.put<PersistentVolume>('/kubernetes/pv', data)
  },

  // 删除 PV
  deletePV(data: DeleteStorageRequest) {
    return service.delete('/kubernetes/pv', { data })
  },

  // ==================== PersistentVolumeClaim 相关 API ====================
  
  // 获取 PVC 列表
  getPVCList(params: GetStorageListRequest) {
    return service.get<StorageListResponse<PersistentVolumeClaim>>('/kubernetes/pvc', { params })
  },

  // 获取 PVC 详情
  getPVCDetail(params: GetStorageDetailRequest) {
    return service.get<StorageDetailResponse<PersistentVolumeClaim>>('/kubernetes/pvcDetails', { params })
  },

  // 创建 PVC
  createPVC(data: CreateStorageRequest) {
    return service.post<PersistentVolumeClaim>('/kubernetes/pvc', data)
  },

  // 更新 PVC
  updatePVC(data: UpdateStorageRequest) {
    return service.put<PersistentVolumeClaim>('/kubernetes/pvc', data)
  },

  // 删除 PVC
  deletePVC(data: DeleteStorageRequest) {
    return service.delete('/kubernetes/pvc', { data })
  },

  // 扩容 PVC
  expandPVC(data: ExpandPVCRequest) {
    return service.patch('/kubernetes/pvc/expand', data)
  },

  // ==================== StorageClass 相关 API ====================
  
  // 获取 StorageClass 列表
  getStorageClassList(params: GetStorageListRequest) {
    return service.get<StorageListResponse<StorageClass>>('/kubernetes/storageclass', { params })
  },

  // 获取 StorageClass 详情
  getStorageClassDetail(params: GetStorageDetailRequest) {
    return service.get<StorageDetailResponse<StorageClass>>('/kubernetes/storageclass/detail', { params })
  },

  // 创建 StorageClass
  createStorageClass(data: CreateStorageRequest) {
    return service.post<StorageClass>('/kubernetes/storageclass', data)
  },

  // 更新 StorageClass
  updateStorageClass(data: UpdateStorageRequest) {
    return service.put<StorageClass>('/kubernetes/storageclass', data)
  },

  // 删除 StorageClass
  deleteStorageClass(data: DeleteStorageRequest) {
    return service.delete('/kubernetes/storageclass', { data })
  },

  // 设置默认 StorageClass
  setDefaultStorageClass(data: { cluster_id: number; name: string }) {
    return service.patch('/kubernetes/storageclass/default', data)
  },

  // ==================== VolumeSnapshot 相关 API ====================
  
  // 获取快照列表
  getVolumeSnapshotList(params: GetStorageListRequest) {
    return service.get<StorageListResponse<VolumeSnapshot>>('/kubernetes/volumesnapshot', { params })
  },

  // 获取快照详情
  getVolumeSnapshotDetail(params: GetStorageDetailRequest) {
    return service.get<StorageDetailResponse<VolumeSnapshot>>('/kubernetes/volumesnapshot/detail', { params })
  },

  // 创建快照
  createVolumeSnapshot(data: CreateSnapshotRequest) {
    return service.post<VolumeSnapshot>('/kubernetes/volumesnapshot', data)
  },

  // 删除快照
  deleteVolumeSnapshot(data: DeleteStorageRequest) {
    return service.delete('/kubernetes/volumesnapshot', { data })
  },

  // 从快照恢复
  restoreFromSnapshot(data: RestoreFromSnapshotRequest) {
    return service.post('/kubernetes/volumesnapshot/restore', data)
  },

  // ==================== VolumeSnapshotClass 相关 API ====================
  
  // 获取快照类列表
  getVolumeSnapshotClassList(params: GetStorageListRequest) {
    return service.get<StorageListResponse<VolumeSnapshotClass>>('/kubernetes/volumesnapshotclass', { params })
  },

  // 获取快照类详情
  getVolumeSnapshotClassDetail(params: GetStorageDetailRequest) {
    return service.get<StorageDetailResponse<VolumeSnapshotClass>>('/kubernetes/volumesnapshotclass/detail', { params })
  },

  // 创建快照类
  createVolumeSnapshotClass(data: CreateStorageRequest) {
    return service.post<VolumeSnapshotClass>('/kubernetes/volumesnapshotclass', data)
  },

  // 更新快照类
  updateVolumeSnapshotClass(data: UpdateStorageRequest) {
    return service.put<VolumeSnapshotClass>('/kubernetes/volumesnapshotclass', data)
  },

  // 删除快照类
  deleteVolumeSnapshotClass(data: DeleteStorageRequest) {
    return service.delete('/kubernetes/volumesnapshotclass', { data })
  },

  // ==================== CSI 相关 API ====================
  
  // 获取 CSIDriver 列表
  getCSIDriverList(params: GetStorageListRequest) {
    return service.get<StorageListResponse<CSIDriver>>('/kubernetes/csidriver', { params })
  },

  // 获取 CSIDriver 详情
  getCSIDriverDetail(params: GetStorageDetailRequest) {
    return service.get<StorageDetailResponse<CSIDriver>>('/kubernetes/csidriver/detail', { params })
  },

  // 获取 CSINode 列表
  getCSINodeList(params: GetStorageListRequest) {
    return service.get<StorageListResponse<CSINode>>('/kubernetes/csinode', { params })
  },

  // 获取 CSINode 详情
  getCSINodeDetail(params: GetStorageDetailRequest) {
    return service.get<StorageDetailResponse<CSINode>>('/kubernetes/csinode/detail', { params })
  },

  // ==================== 存储监控和统计 API ====================
  
  // 获取存储指标
  getStorageMetrics(params: { cluster_id: number; namespace?: string }) {
    return service.get<StorageMetricsResponse>('/kubernetes/storage/metrics', { params })
  },

  // 获取存储使用趋势
  getStorageUsageTrend(params: {
    cluster_id: number
    namespace?: string
    start_time: string
    end_time: string
    granularity: string
  }) {
    return service.get('/kubernetes/storage/usage-trend', { params })
  },

  // 获取存储性能指标
  getStoragePerformance(params: {
    cluster_id: number
    storage_class?: string
    pv_name?: string
    start_time: string
    end_time: string
  }) {
    return service.get('/kubernetes/storage/performance', { params })
  },

  // 获取存储容量规划
  getStorageCapacityPlanning(params: { cluster_id: number; namespace?: string }) {
    return service.get('/kubernetes/storage/capacity-planning', { params })
  },

  // 存储健康检查
  checkStorageHealth(data: { cluster_id: number; namespace?: string }) {
    return service.post('/kubernetes/storage/health-check', data)
  },

  // 存储性能测试
  performStorageBenchmark(data: {
    cluster_id: number
    namespace: string
    storage_class: string
    test_type: string // sequential-read, sequential-write, random-read, random-write
    file_size: string
    duration: number
  }) {
    return service.post('/kubernetes/storage/benchmark', data)
  },

  // 获取存储故障信息
  getStorageIssues(params: { cluster_id: number; namespace?: string; severity?: string }) {
    return service.get('/kubernetes/storage/issues', { params })
  },

  // 存储资源清理
  cleanupOrphanedStorage(data: { cluster_id: number; namespace?: string; dry_run?: boolean }) {
    return service.post('/kubernetes/storage/cleanup', data)
  }
}