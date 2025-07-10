<template>
  <el-dialog
    v-model="dialogVisible"
    title="Pod 详情"
    width="80%"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <div v-if="pod" class="pod-detail">
      <!-- 基本信息 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
          </div>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="名称">{{ pod.metadata.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ pod.metadata.namespace }}</el-descriptions-item>
          <el-descriptions-item label="UID">{{ pod.metadata.uid }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getPodStatusType()">
              {{ pod.status.phase }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="节点">{{ (pod.spec as any).nodeName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="IP">{{ pod.status.podIP || '-' }}</el-descriptions-item>
          <el-descriptions-item label="主机IP">{{ pod.status.hostIP || '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(pod.metadata.creationTimestamp) }}</el-descriptions-item>
          <el-descriptions-item label="启动时间">{{ pod.status.startTime ? formatDate(pod.status.startTime) : '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 标签和注解 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>标签和注解</span>
          </div>
        </template>
        <el-tabs>
          <el-tab-pane label="标签" name="labels">
            <div v-if="pod.metadata.labels && Object.keys(pod.metadata.labels).length > 0">
              <el-tag
                v-for="(value, key) in pod.metadata.labels"
                :key="key"
                class="mr-2 mb-2"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
            <el-empty v-else description="无标签" />
          </el-tab-pane>
          <el-tab-pane label="注解" name="annotations">
            <div v-if="pod.metadata.annotations && Object.keys(pod.metadata.annotations).length > 0">
              <el-tag
                v-for="(value, key) in pod.metadata.annotations"
                :key="key"
                class="mr-2 mb-2"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
            <el-empty v-else description="无注解" />
          </el-tab-pane>
        </el-tabs>
      </el-card>

      <!-- 容器信息 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>容器信息</span>
          </div>
        </template>
        <div v-for="(container, index) in pod.spec.containers" :key="index" class="container-info">
          <h4>容器 #{{ index + 1 }}: {{ container.name }}</h4>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="镜像">{{ container.image }}</el-descriptions-item>
            <el-descriptions-item label="镜像拉取策略">{{ (container as any).imagePullPolicy || 'IfNotPresent' }}</el-descriptions-item>
            <el-descriptions-item label="端口">
              <div v-if="container.ports && container.ports.length > 0">
                <div v-for="port in container.ports" :key="port.containerPort">
                  {{ port.containerPort }}{{ port.protocol ? '/' + port.protocol : '' }}{{ port.name ? ' (' + port.name + ')' : '' }}
                </div>
              </div>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item label="环境变量">
              <div v-if="container.env && container.env.length > 0">
                <div v-for="env in container.env" :key="env.name">
                  {{ env.name }}={{ env.value || '***' }}
                </div>
              </div>
              <span v-else>-</span>
            </el-descriptions-item>
          </el-descriptions>
          
          <!-- 容器状态 -->
          <div v-if="getContainerStatus(container.name)" class="container-status">
            <h5>状态信息</h5>
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="就绪状态">
                <el-tag :type="getContainerStatus(container.name)?.ready ? 'success' : 'danger'">
                  {{ getContainerStatus(container.name)?.ready ? '就绪' : '未就绪' }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="重启次数">{{ getContainerStatus(container.name)?.restartCount || 0 }}</el-descriptions-item>
              <el-descriptions-item label="状态">
                <el-tag :type="getContainerStateType(getContainerStatus(container.name))">
                  {{ getContainerStateText(getContainerStatus(container.name)) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="镜像ID">{{ getContainerStatus(container.name)?.imageID || '-' }}</el-descriptions-item>
            </el-descriptions>
          </div>
        </div>
      </el-card>

      <!-- 条件状态 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>条件状态</span>
          </div>
        </template>
        <el-table :data="pod.status.conditions" border stripe>
          <el-table-column prop="type" label="类型" width="150" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'True' ? 'success' : 'danger'">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="lastTransitionTime" label="最后转换时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.lastTransitionTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" />
          <el-table-column prop="message" label="消息" show-overflow-tooltip />
        </el-table>
      </el-card>
    </div>

    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Pod, ContainerStatus } from '@/types/kubernetes/workload'

// Props
interface Props {
  visible: boolean
  pod?: Pod | null
  clusterId: string | number
  namespace: string
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  pod: null,
  clusterId: '',
  namespace: ''
})

// Emits
const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 方法
const getPodStatusType = () => {
  if (!props.pod) return 'info'
  
  const phase = props.pod.status.phase
  switch (phase) {
    case 'Running':
      return 'success'
    case 'Pending':
      return 'warning'
    case 'Succeeded':
      return 'info'
    case 'Failed':
      return 'danger'
    default:
      return 'info'
  }
}

const getContainerStatus = (containerName: string): ContainerStatus | undefined => {
  if (!props.pod?.status.containerStatuses) return undefined
  return props.pod.status.containerStatuses.find(container => container.name === containerName)
}

const getContainerStateType = (containerStatus?: ContainerStatus) => {
  if (!containerStatus) return 'info'
  
  if (containerStatus.state.running) return 'success'
  if (containerStatus.state.waiting) return 'warning'
  if (containerStatus.state.terminated) return 'danger'
  return 'info'
}

const getContainerStateText = (containerStatus?: ContainerStatus) => {
  if (!containerStatus) return '未知'
  
  if (containerStatus.state.running) return '运行中'
  if (containerStatus.state.waiting) return `等待中 (${containerStatus.state.waiting.reason})`
  if (containerStatus.state.terminated) return `已终止 (${containerStatus.state.terminated.reason})`
  return '未知'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const handleClose = () => {
  dialogVisible.value = false
}
</script>

<style scoped>
.pod-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.detail-card {
  margin-bottom: 20px;
}

.detail-card:last-child {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.container-info {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

.container-info:last-child {
  margin-bottom: 0;
}

.container-info h4 {
  margin: 0 0 10px 0;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.container-status {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #e4e7ed;
}

.container-status h5 {
  margin: 0 0 10px 0;
  color: var(--el-text-color-primary);
  font-size: 13px;
}

.mr-2 {
  margin-right: 8px;
}

.mb-2 {
  margin-bottom: 8px;
}
</style> 