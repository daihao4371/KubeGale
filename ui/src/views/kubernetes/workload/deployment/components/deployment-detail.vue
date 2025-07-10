<template>
  <el-dialog
    v-model="dialogVisible"
    title="Deployment 详情"
    width="900px"
    :before-close="handleClose"
  >
    <div v-if="deployment" class="deployment-detail">
      <!-- 基本信息 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="名称">
            {{ deployment.metadata.name }}
          </el-descriptions-item>
          <el-descriptions-item label="命名空间">
            {{ deployment.metadata.namespace }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(deployment.metadata.creationTimestamp) }}
          </el-descriptions-item>
          <el-descriptions-item label="UID">
            {{ deployment.metadata.uid }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 状态信息 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>状态信息</span>
          </div>
        </template>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="副本数">
            <el-tag :type="getReplicaStatusType()">
              {{ deployment.status.replicas }}/{{ deployment.spec.replicas }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="就绪副本">
            {{ deployment.status.readyReplicas }}
          </el-descriptions-item>
          <el-descriptions-item label="可用副本">
            {{ deployment.status.availableReplicas }}
          </el-descriptions-item>
          <el-descriptions-item label="更新副本">
            {{ deployment.status.updatedReplicas }}
          </el-descriptions-item>
          <el-descriptions-item label="不可用副本">
            {{ deployment.status.unavailableReplicas }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType()">
              {{ getStatusText() }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 容器信息 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>容器信息</span>
          </div>
        </template>
        <div
          v-for="(container, index) in deployment.spec.template.spec.containers"
          :key="index"
          class="container-info"
        >
          <h4>{{ container.name }}</h4>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="镜像">
              {{ container.image }}
            </el-descriptions-item>
            <el-descriptions-item label="端口">
              <div v-if="container.ports && container.ports.length > 0">
                <el-tag
                  v-for="port in container.ports"
                  :key="port.name || port.containerPort"
                  size="small"
                  class="mr-2"
                >
                  {{ port.name || 'unnamed' }}:{{ port.containerPort }}/{{ port.protocol }}
                </el-tag>
              </div>
              <span v-else>无</span>
            </el-descriptions-item>
            <el-descriptions-item label="环境变量">
              <div v-if="container.env && container.env.length > 0">
                <el-tag
                  v-for="env in container.env"
                  :key="env.name"
                  size="small"
                  class="mr-2"
                >
                  {{ env.name }}={{ env.value }}
                </el-tag>
              </div>
              <span v-else>无</span>
            </el-descriptions-item>
            <el-descriptions-item label="资源限制">
              <div v-if="container.resources">
                <div v-if="container.resources.requests">
                  <strong>请求:</strong> 
                  CPU: {{ container.resources.requests.cpu || '未设置' }}, 
                  内存: {{ container.resources.requests.memory || '未设置' }}
                </div>
                <div v-if="container.resources.limits">
                  <strong>限制:</strong> 
                  CPU: {{ container.resources.limits.cpu || '未设置' }}, 
                  内存: {{ container.resources.limits.memory || '未设置' }}
                </div>
              </div>
              <span v-else>未设置</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </el-card>

      <!-- 更新策略 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>更新策略</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="策略类型">
            {{ deployment.spec.strategy.type }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="deployment.spec.strategy.type === 'RollingUpdate'"
            label="最大不可用"
          >
            {{ deployment.spec.strategy.rollingUpdate?.maxUnavailable }}
          </el-descriptions-item>
          <el-descriptions-item
            v-if="deployment.spec.strategy.type === 'RollingUpdate'"
            label="最大超出"
          >
            {{ deployment.spec.strategy.rollingUpdate?.maxSurge }}
          </el-descriptions-item>
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
            <div v-if="deployment.metadata.labels && Object.keys(deployment.metadata.labels).length > 0">
              <el-tag
                v-for="(value, key) in deployment.metadata.labels"
                :key="key"
                class="mr-2 mb-2"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
            <el-empty v-else description="无标签" />
          </el-tab-pane>
          <el-tab-pane label="注解" name="annotations">
            <div v-if="deployment.metadata.annotations && Object.keys(deployment.metadata.annotations).length > 0">
              <el-tag
                v-for="(value, key) in deployment.metadata.annotations"
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

      <!-- 条件状态 -->
      <el-card class="detail-card">
        <template #header>
          <div class="card-header">
            <span>条件状态</span>
          </div>
        </template>
        <el-table :data="deployment.status.conditions" border stripe>
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
import type { Deployment } from '@/types/kubernetes/workload'

// Props
interface Props {
  visible: boolean
  deployment?: Deployment | null
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  deployment: null
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
const getReplicaStatusType = () => {
  if (!props.deployment) return 'info'
  
  const { replicas, readyReplicas } = props.deployment.status
  if (readyReplicas === replicas) return 'success'
  if (readyReplicas > 0) return 'warning'
  return 'danger'
}

const getStatusType = () => {
  if (!props.deployment) return 'info'
  
  const { replicas, readyReplicas, availableReplicas } = props.deployment.status
  if (availableReplicas === replicas) return 'success'
  if (readyReplicas > 0) return 'warning'
  return 'danger'
}

const getStatusText = () => {
  if (!props.deployment) return '未知'
  
  const { replicas, readyReplicas, availableReplicas } = props.deployment.status
  if (availableReplicas === replicas) return '运行中'
  if (readyReplicas > 0) return '部分就绪'
  return '未就绪'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const handleClose = () => {
  dialogVisible.value = false
}
</script>

<style scoped>
.deployment-detail {
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

.mr-2 {
  margin-right: 8px;
}

.mb-2 {
  margin-bottom: 8px;
}
</style> 