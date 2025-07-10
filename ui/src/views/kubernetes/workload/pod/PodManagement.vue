<template>
  <PageLayout title="Pod 管理">
    <template #subtitle>
      管理 Kubernetes Pod 资源，包括查看、创建、删除等操作
    </template>
    
    <template #actions>
      <ResourceActions
        :selected-items="selectedPods"
        :refreshing="loading"
        @bulk-delete="handleBulkDelete"
        @refresh="handleRefresh"
        @create="handleCreate"
      />
    </template>

    <template #search>
      <SearchFilter
        v-model="searchParams"
        :namespaces="namespaces"
        :status-options="podStatusOptions"
        @search="handleSearch"
        @reset="handleReset"
      />
    </template>

    <K8sTable
      :data="pods"
      :loading="loading"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      :show-selection="true"
      :show-actions="true"
      @selection-change="handleSelectionChange"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
      @refresh="handleRefresh"
    >
      <el-table-column prop="name" label="名称" min-width="200">
        <template #default="{ row }">
          <el-link type="primary" @click="handleShowDetail(row)">
            {{ row.name }}
          </el-link>
        </template>
      </el-table-column>
      
      <el-table-column prop="namespace" label="命名空间" width="120" />
      
      <el-table-column prop="status.phase" label="状态" width="120">
        <template #default="{ row }">
          <K8sStatusBadge 
            :status="row.status?.phase || 'Unknown'"
            :show-icon="true"
          />
        </template>
      </el-table-column>
      
      <el-table-column prop="ready" label="就绪状态" width="100">
        <template #default="{ row }">
          {{ getReadyStatus(row) }}
        </template>
      </el-table-column>
      
      <el-table-column prop="restarts" label="重启次数" width="100">
        <template #default="{ row }">
          {{ getRestartCount(row) }}
        </template>
      </el-table-column>
      
      <el-table-column prop="node" label="节点" min-width="150">
        <template #default="{ row }">
          {{ row.spec?.nodeName || 'Unscheduled' }}
        </template>
      </el-table-column>
      
      <el-table-column prop="age" label="运行时间" width="120">
        <template #default="{ row }">
          {{ formatAge(row.metadata?.creationTimestamp) }}
        </template>
      </el-table-column>

      <template #actions="{ row }">
        <el-button size="small" @click="handleShowDetail(row)">详情</el-button>
        <el-button size="small" @click="handleShowLogs(row)">日志</el-button>
        <el-button size="small" type="warning" @click="handleRestart(row)">重启</el-button>
        <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
      </template>
    </K8sTable>

    <!-- Pod 详情抽屉 -->
    <ResourceDrawer
      v-model="showDetail"
      title="Pod 详情"
      :basic-info="detailBasicInfo"
      :labels="currentPod?.metadata?.labels"
      :annotations="currentPod?.metadata?.annotations"
      :yaml-content="yamlContent"
      @edit="handleEdit"
      @delete="handleDeleteFromDetail"
    >
      <!-- 容器信息 -->
      <div class="resource-section">
        <h3>容器信息</h3>
        <el-table :data="containers" size="small">
          <el-table-column prop="name" label="容器名称" />
          <el-table-column prop="image" label="镜像" show-overflow-tooltip />
          <el-table-column label="状态">
            <template #default="{ row }">
              <K8sStatusBadge 
                :status="getContainerStatus(row.name)"
                :show-icon="true"
                size="small"
              />
            </template>
          </el-table-column>
          <el-table-column label="重启次数" width="100">
            <template #default="{ row }">
              {{ getContainerRestartCount(row.name) }}
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 事件信息 -->
      <div class="resource-section" v-if="events.length > 0">
        <h3>相关事件</h3>
        <el-table :data="events" size="small" max-height="300">
          <el-table-column prop="type" label="类型" width="80">
            <template #default="{ row }">
              <K8sStatusBadge 
                :status="row.type"
                size="small"
              />
            </template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" width="150" />
          <el-table-column prop="message" label="消息" show-overflow-tooltip />
          <el-table-column prop="firstTimestamp" label="首次时间" width="120">
            <template #default="{ row }">
              {{ formatAge(row.firstTimestamp) }}
            </template>
          </el-table-column>
        </el-table>
      </div>
    </ResourceDrawer>

    <!-- Pod 日志对话框 -->
    <el-dialog
      v-model="showLogs"
      title="Pod 日志"
      width="80%"
      :before-close="handleCloseLog"
    >
      <div class="log-controls">
        <el-select 
          v-model="selectedContainer" 
          placeholder="选择容器"
          @change="handleContainerChange"
        >
          <el-option 
            v-for="container in logContainers" 
            :key="container" 
            :label="container" 
            :value="container"
          />
        </el-select>
        <el-switch 
          v-model="followLogs" 
          active-text="实时日志"
          @change="handleFollowChange"
        />
        <el-button @click="handleClearLogs">清空</el-button>
        <el-button @click="handleRefreshLogs">刷新</el-button>
      </div>
      
      <div class="log-content">
        <pre ref="logContainer">{{ logContent }}</pre>
      </div>
    </el-dialog>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageLayout from '@/components/layout/PageLayout.vue'
import K8sTable from '@/components/K8sTable.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import ResourceActions from '@/components/kubernetes/ResourceActions.vue'
import ResourceDrawer from '@/components/kubernetes/ResourceDrawer.vue'
import SearchFilter from '@/components/kubernetes/SearchFilter.vue'
import type { Pod, PodStatus, Container } from '@/types/kubernetes/workload'
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'PodManagement'
})

// 响应式数据
const pods = ref<Pod[]>([])
const selectedPods = ref<Pod[]>([])
const currentPod = ref<Pod | null>(null)
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const namespaces = ref<string[]>([])
const events = ref<any[]>([])
const containers = ref<Container[]>([])

// 搜索参数
const searchParams = reactive({
  name: '',
  namespace: '',
  status: '',
  labels: ''
})

// 详情相关
const showDetail = ref(false)
const yamlContent = ref('')

// 日志相关
const showLogs = ref(false)
const logContent = ref('')
const selectedContainer = ref('')
const logContainers = ref<string[]>([])
const followLogs = ref(false)
const logContainer = ref<HTMLElement>()

// 状态选项
const podStatusOptions = [
  { label: '运行中', value: 'Running' },
  { label: '等待中', value: 'Pending' },
  { label: '成功', value: 'Succeeded' },
  { label: '失败', value: 'Failed' },
  { label: '未知', value: 'Unknown' }
]

// 计算属性
const detailBasicInfo = computed(() => {
  if (!currentPod.value) return {}
  
  const pod = currentPod.value
  return {
    name: { label: '名称', value: pod.metadata?.name || '-' },
    namespace: { label: '命名空间', value: pod.metadata?.namespace || '-' },
    status: { 
      label: '状态', 
      value: pod.status?.phase || 'Unknown',
      class: `status-${(pod.status?.phase || 'unknown').toLowerCase()}`
    },
    node: { label: '节点', value: pod.spec?.nodeName || 'Unscheduled' },
    podIP: { label: 'Pod IP', value: pod.status?.podIP || '-' },
    hostIP: { label: '主机 IP', value: pod.status?.hostIP || '-' },
    qosClass: { label: 'QoS 类别', value: pod.status?.qosClass || '-' },
    restartPolicy: { label: '重启策略', value: pod.spec?.restartPolicy || '-' },
    creationTime: { 
      label: '创建时间', 
      value: pod.metadata?.creationTimestamp ? 
        new Date(pod.metadata.creationTimestamp).toLocaleString() : '-'
    }
  }
})

// 方法
const handleSearch = () => {
  currentPage.value = 1
  fetchPods()
}

const handleReset = () => {
  Object.assign(searchParams, {
    name: '',
    namespace: '',
    status: '',
    labels: ''
  })
  handleSearch()
}

const handleRefresh = () => {
  fetchPods()
}

const handleCreate = () => {
  // TODO: 实现创建 Pod 功能
  ElMessage.info('创建功能正在开发中')
}

const handleSelectionChange = (selection: Pod[]) => {
  selectedPods.value = selection
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchPods()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchPods()
}

const handleShowDetail = async (row: Pod) => {
  currentPod.value = row
  showDetail.value = true
  
  // 获取容器信息
  containers.value = row.spec?.containers || []
  
  // 获取 YAML
  try {
    // TODO: 调用 API 获取 YAML
    yamlContent.value = `# Pod YAML 配置\napiVersion: v1\nkind: Pod\nmetadata:\n  name: ${row.metadata?.name}\n  namespace: ${row.metadata?.namespace}`
  } catch (error) {
    console.error('获取 YAML 失败:', error)
  }
  
  // 获取相关事件
  try {
    // TODO: 调用 API 获取事件
    events.value = []
  } catch (error) {
    console.error('获取事件失败:', error)
  }
}

const handleShowLogs = (row: Pod) => {
  currentPod.value = row
  logContainers.value = row.spec?.containers?.map(c => c.name) || []
  if (logContainers.value.length > 0) {
    selectedContainer.value = logContainers.value[0]
  }
  showLogs.value = true
  fetchLogs()
}

const handleRestart = async (row: Pod) => {
  try {
    await ElMessageBox.confirm(
      `确定要重启 Pod "${row.metadata?.name}" 吗？`,
      '确认重启',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    // TODO: 调用重启 API
    ElMessage.success('重启成功')
    await fetchPods()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重启失败:', error)
    }
  }
}

const handleDelete = async (row: Pod) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Pod "${row.metadata?.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    // TODO: 调用删除 API
    ElMessage.success('删除成功')
    await fetchPods()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const handleBulkDelete = async (pods: Pod[]) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${pods.length} 个 Pod 吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    // TODO: 调用批量删除 API
    ElMessage.success('批量删除成功')
    await fetchPods()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
    }
  }
}

const handleEdit = () => {
  // TODO: 实现编辑功能
  ElMessage.info('编辑功能正在开发中')
}

const handleDeleteFromDetail = () => {
  if (currentPod.value) {
    showDetail.value = false
    handleDelete(currentPod.value)
  }
}

const handleContainerChange = () => {
  fetchLogs()
}

const handleFollowChange = () => {
  if (followLogs.value) {
    // TODO: 开启实时日志
  } else {
    // TODO: 关闭实时日志
  }
}

const handleClearLogs = () => {
  logContent.value = ''
}

const handleRefreshLogs = () => {
  fetchLogs()
}

const handleCloseLog = () => {
  followLogs.value = false
  logContent.value = ''
  showLogs.value = false
}

// 工具函数
const getReadyStatus = (pod: Pod): string => {
  const conditions = pod.status?.conditions || []
  const readyCondition = conditions.find(c => c.type === 'Ready')
  return readyCondition?.status === 'True' ? '1/1' : '0/1'
}

const getRestartCount = (pod: Pod): number => {
  const containerStatuses = pod.status?.containerStatuses || []
  return containerStatuses.reduce((total, status) => total + (status.restartCount || 0), 0)
}

const getContainerStatus = (containerName: string): string => {
  if (!currentPod.value) return 'Unknown'
  
  const containerStatuses = currentPod.value.status?.containerStatuses || []
  const status = containerStatuses.find(s => s.name === containerName)
  
  if (status?.state?.running) return 'Running'
  if (status?.state?.waiting) return status.state.waiting.reason || 'Waiting'
  if (status?.state?.terminated) return status.state.terminated.reason || 'Terminated'
  
  return 'Unknown'
}

const getContainerRestartCount = (containerName: string): number => {
  if (!currentPod.value) return 0
  
  const containerStatuses = currentPod.value.status?.containerStatuses || []
  const status = containerStatuses.find(s => s.name === containerName)
  
  return status?.restartCount || 0
}

const fetchPods = async () => {
  try {
    loading.value = true
    
    // TODO: 调用实际的 API
    // const response = await podApi.getPods({
    //   page: currentPage.value,
    //   pageSize: pageSize.value,
    //   ...searchParams
    // })
    
    // 模拟数据
    pods.value = [
      {
        metadata: {
          name: 'nginx-pod-1',
          namespace: 'default',
          creationTimestamp: '2024-01-01T10:00:00Z',
          labels: { app: 'nginx' }
        },
        spec: {
          containers: [{ name: 'nginx', image: 'nginx:1.20' }],
          nodeName: 'node-1',
          restartPolicy: 'Always'
        },
        status: {
          phase: 'Running',
          podIP: '10.244.1.5',
          hostIP: '192.168.1.100',
          qosClass: 'BestEffort',
          conditions: [{ type: 'Ready', status: 'True' }],
          containerStatuses: [{ 
            name: 'nginx', 
            restartCount: 0,
            state: { running: { startedAt: '2024-01-01T10:00:30Z' } }
          }]
        }
      }
    ]
    total.value = 1
    
  } catch (error) {
    ElMessage.error('获取 Pod 列表失败')
    console.error('获取 Pod 列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchNamespaces = async () => {
  try {
    // TODO: 调用实际的 API
    namespaces.value = ['default', 'kube-system', 'kube-public']
  } catch (error) {
    console.error('获取命名空间失败:', error)
  }
}

const fetchLogs = async () => {
  if (!currentPod.value || !selectedContainer.value) return
  
  try {
    // TODO: 调用实际的日志 API
    logContent.value = `[${new Date().toISOString()}] Pod ${currentPod.value.metadata?.name} 容器 ${selectedContainer.value} 日志开始...\n`
    logContent.value += `[${new Date().toISOString()}] 容器启动成功\n`
    logContent.value += `[${new Date().toISOString()}] 应用程序正在运行\n`
    
    // 自动滚动到底部
    await nextTick()
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  } catch (error) {
    ElMessage.error('获取日志失败')
    console.error('获取日志失败:', error)
  }
}

// 生命周期
onMounted(async () => {
  await Promise.all([
    fetchPods(),
    fetchNamespaces()
  ])
})
</script>

<style scoped>
.log-controls {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.log-content {
  height: 400px;
  overflow: auto;
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 12px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

.resource-section {
  margin-bottom: 24px;
}

.resource-section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}
</style>