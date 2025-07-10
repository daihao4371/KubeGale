import service from '@/api/request'
import type {
  // 工作负载相关类型
  Pod,
  Deployment,
  StatefulSet,
  DaemonSet,
  Job,
  CronJob,
  ReplicaSet,
  GetWorkloadListRequest,
  GetWorkloadDetailRequest,
  CreateWorkloadRequest,
  UpdateWorkloadRequest,
  DeleteWorkloadRequest,
  ScaleWorkloadRequest,
  RollbackDeploymentRequest,
  WorkloadListResponse,
  WorkloadDetailResponse,
  PodLogRequest,
  PodLogResponse,
  PodTerminalRequest,
  PodTerminalResponse
} from '@/types/kubernetes/workload'

export const workloadApi = {
  // ==================== Pod 相关 API ====================
  
  // 获取 Pod 列表
  getPodList(params: GetWorkloadListRequest) {
    return service.get<WorkloadListResponse<Pod>>('/kubernetes/pods', { params })
  },

  // 获取 Pod 详情
  getPodDetail(params: GetWorkloadDetailRequest) {
    return service.get<WorkloadDetailResponse<Pod>>('/kubernetes/podDetails', { params })
  },

  // 创建 Pod
  createPod(data: CreateWorkloadRequest) {
    return service.post<Pod>('/kubernetes/pods', data)
  },

  // 更新 Pod
  updatePod(data: UpdateWorkloadRequest) {
    return service.put<Pod>('/kubernetes/pods', data)
  },

  // 删除 Pod
  deletePod(data: DeleteWorkloadRequest) {
    return service.delete('/kubernetes/pods', { data })
  },

  // 获取 Pod 指标
  getPodMetrics(params: GetWorkloadListRequest) {
    return service.get('/kubernetes/pods/metrics', { params })
  },

  // 获取 Pod 日志
  getPodLog(params: PodLogRequest) {
    return service.get<PodLogResponse>('/kubernetes/pods/logs', { params })
  },

  // Pod 终端连接
  getPodTerminal(data: PodTerminalRequest) {
    return service.post<PodTerminalResponse>('/kubernetes/pods/terminal', data)
  },

  // 列出 Pod 文件
  listPodFiles(data: { cluster_id: number; namespace: string; pod_name: string; path: string }) {
    return service.post('/kubernetes/pods/listFiles', data)
  },

  // 上传文件到 Pod
  uploadFileToPod(data: FormData) {
    return service.post('/kubernetes/pods/uploadFile', data, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // 从 Pod 下载文件
  downloadFileFromPod(params: { cluster_id: number; namespace: string; pod_name: string; file_path: string }) {
    return service.get('/kubernetes/pods/downloadFile', { params, responseType: 'blob' })
  },

  // 删除 Pod 文件
  deletePodFiles(data: { cluster_id: number; namespace: string; pod_name: string; file_paths: string[] }) {
    return service.post('/kubernetes/pods/deleteFiles', data)
  },

  // ==================== Deployment 相关 API ====================
  
  // 获取 Deployment 列表
  getDeploymentList(params: GetWorkloadListRequest) {
    return service.get<WorkloadListResponse<Deployment>>('/kubernetes/deployment', { params })
  },

  // 获取 Deployment 详情
  getDeploymentDetail(params: GetWorkloadDetailRequest) {
    return service.get<WorkloadDetailResponse<Deployment>>('/kubernetes/deployment/detail', { params })
  },

  // 创建 Deployment
  createDeployment(data: CreateWorkloadRequest) {
    return service.post<Deployment>('/kubernetes/deployment', data)
  },

  // 更新 Deployment
  updateDeployment(data: UpdateWorkloadRequest) {
    return service.put<Deployment>('/kubernetes/deployment', data)
  },

  // 删除 Deployment
  deleteDeployment(data: DeleteWorkloadRequest) {
    return service.delete('/kubernetes/deployment', { data })
  },

  // 缩放 Deployment
  scaleDeployment(data: ScaleWorkloadRequest) {
    return service.patch('/kubernetes/deployment/scale', data)
  },

  // 回滚 Deployment
  rollbackDeployment(data: RollbackDeploymentRequest) {
    return service.patch('/kubernetes/deployment', data)
  },

  // ==================== StatefulSet 相关 API ====================
  
  // 获取 StatefulSet 列表
  getStatefulSetList(params: GetWorkloadListRequest) {
    return service.get<WorkloadListResponse<StatefulSet>>('/kubernetes/statefulset', { params })
  },

  // 获取 StatefulSet 详情
  getStatefulSetDetail(params: GetWorkloadDetailRequest) {
    return service.get<WorkloadDetailResponse<StatefulSet>>('/kubernetes/statefulset/detail', { params })
  },

  // 创建 StatefulSet
  createStatefulSet(data: CreateWorkloadRequest) {
    return service.post<StatefulSet>('/kubernetes/statefulset', data)
  },

  // 更新 StatefulSet
  updateStatefulSet(data: UpdateWorkloadRequest) {
    return service.put<StatefulSet>('/kubernetes/statefulset', data)
  },

  // 删除 StatefulSet
  deleteStatefulSet(data: DeleteWorkloadRequest) {
    return service.delete('/kubernetes/statefulset', { data })
  },

  // 缩放 StatefulSet
  scaleStatefulSet(data: ScaleWorkloadRequest) {
    return service.patch('/kubernetes/statefulset/scale', data)
  },

  // ==================== DaemonSet 相关 API ====================
  
  // 获取 DaemonSet 列表
  getDaemonSetList(params: GetWorkloadListRequest) {
    return service.get<WorkloadListResponse<DaemonSet>>('/kubernetes/daemonset', { params })
  },

  // 获取 DaemonSet 详情
  getDaemonSetDetail(params: GetWorkloadDetailRequest) {
    return service.get<WorkloadDetailResponse<DaemonSet>>('/kubernetes/daemonset/detail', { params })
  },

  // 创建 DaemonSet
  createDaemonSet(data: CreateWorkloadRequest) {
    return service.post<DaemonSet>('/kubernetes/daemonset', data)
  },

  // 更新 DaemonSet
  updateDaemonSet(data: UpdateWorkloadRequest) {
    return service.put<DaemonSet>('/kubernetes/daemonset', data)
  },

  // 删除 DaemonSet
  deleteDaemonSet(data: DeleteWorkloadRequest) {
    return service.delete('/kubernetes/daemonset', { data })
  },

  // ==================== Job 相关 API ====================
  
  // 获取 Job 列表
  getJobList(params: GetWorkloadListRequest) {
    return service.get<WorkloadListResponse<Job>>('/kubernetes/job', { params })
  },

  // 获取 Job 详情
  getJobDetail(params: GetWorkloadDetailRequest) {
    return service.get<WorkloadDetailResponse<Job>>('/kubernetes/job/detail', { params })
  },

  // 创建 Job
  createJob(data: CreateWorkloadRequest) {
    return service.post<Job>('/kubernetes/job', data)
  },

  // 更新 Job
  updateJob(data: UpdateWorkloadRequest) {
    return service.put<Job>('/kubernetes/job', data)
  },

  // 删除 Job
  deleteJob(data: DeleteWorkloadRequest) {
    return service.delete('/kubernetes/job', { data })
  },

  // ==================== CronJob 相关 API ====================
  
  // 获取 CronJob 列表
  getCronJobList(params: GetWorkloadListRequest) {
    return service.get<WorkloadListResponse<CronJob>>('/kubernetes/cronjob', { params })
  },

  // 获取 CronJob 详情
  getCronJobDetail(params: GetWorkloadDetailRequest) {
    return service.get<WorkloadDetailResponse<CronJob>>('/kubernetes/cronjob/detail', { params })
  },

  // 创建 CronJob
  createCronJob(data: CreateWorkloadRequest) {
    return service.post<CronJob>('/kubernetes/cronjob', data)
  },

  // 更新 CronJob
  updateCronJob(data: UpdateWorkloadRequest) {
    return service.put<CronJob>('/kubernetes/cronjob', data)
  },

  // 删除 CronJob
  deleteCronJob(data: DeleteWorkloadRequest) {
    return service.delete('/kubernetes/cronjob', { data })
  },

  // 暂停 CronJob
  suspendCronJob(data: { cluster_id: number; namespace: string; name: string; suspend: boolean }) {
    return service.patch('/kubernetes/cronjob/suspend', data)
  },

  // ==================== ReplicaSet 相关 API ====================
  
  // 获取 ReplicaSet 列表
  getReplicaSetList(params: GetWorkloadListRequest) {
    return service.get<WorkloadListResponse<ReplicaSet>>('/kubernetes/replicaset', { params })
  },

  // 获取 ReplicaSet 详情
  getReplicaSetDetail(params: GetWorkloadDetailRequest) {
    return service.get<WorkloadDetailResponse<ReplicaSet>>('/kubernetes/replicaset/detail', { params })
  },

  // 创建 ReplicaSet
  createReplicaSet(data: CreateWorkloadRequest) {
    return service.post<ReplicaSet>('/kubernetes/replicaset', data)
  },

  // 更新 ReplicaSet
  updateReplicaSet(data: UpdateWorkloadRequest) {
    return service.put<ReplicaSet>('/kubernetes/replicaset', data)
  },

  // 删除 ReplicaSet
  deleteReplicaSet(data: DeleteWorkloadRequest) {
    return service.delete('/kubernetes/replicaset', { data })
  },

  // 缩放 ReplicaSet
  scaleReplicaSet(data: ScaleWorkloadRequest) {
    return service.patch('/kubernetes/replicaset/scale', data)
  }
}