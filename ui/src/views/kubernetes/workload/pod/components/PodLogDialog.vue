<template>
  <el-dialog
    v-model="dialogVisible"
    title="Pod 日志"
    width="80%"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <div v-if="pod" class="pod-log">
      <!-- 日志配置 -->
      <el-card class="config-card">
        <el-form :model="logConfig" inline>
          <el-form-item label="容器">
            <el-select v-model="logConfig.container" placeholder="选择容器" @change="loadLogs">
              <el-option
                v-for="container in pod.spec.containers"
                :key="container.name"
                :label="container.name"
                :value="container.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="行数">
            <el-input-number
              v-model="logConfig.tail_lines"
              :min="1"
              :max="10000"
              placeholder="显示行数"
              @change="loadLogs"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadLogs" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button @click="clearLogs">
              <el-icon><Delete /></el-icon>
              清空
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 日志内容 -->
      <el-card class="log-card">
        <template #header>
          <div class="card-header">
            <span>日志内容</span>
            <div class="header-actions">
              <el-button size="small" @click="copyLogs">
                <el-icon><CopyDocument /></el-icon>
                复制
              </el-button>
              <el-button size="small" @click="downloadLogs">
                <el-icon><Download /></el-icon>
                下载
              </el-button>
            </div>
          </div>
        </template>
        
        <div v-loading="loading" class="log-content">
          <pre v-if="logs" class="log-text">{{ logs }}</pre>
          <el-empty v-else description="暂无日志" />
        </div>
      </el-card>
    </div>

    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Delete, CopyDocument, Download } from '@element-plus/icons-vue'
import type { Pod } from '@/types/kubernetes/workload'
import { getPodLogs } from '@/api/kubernetes/pods'

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

// 响应式数据
const loading = ref(false)
const logs = ref('')

// 日志配置
const logConfig = reactive({
  container: '',
  tail_lines: 100
})

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 方法
const loadLogs = async () => {
  if (!props.pod || !props.clusterId || !props.namespace) return

  loading.value = true
  try {
    const response = await getPodLogs({
      cluster_id: Number(props.clusterId),
      namespace: props.namespace,
      pod_name: props.pod.metadata.name,
      container: logConfig.container || undefined,
      tail_lines: logConfig.tail_lines,
      follow: false
    })
    
    logs.value = response.logs || ''
  } catch (error) {
    console.error('加载日志失败:', error)
    ElMessage.error('加载日志失败')
    logs.value = ''
  } finally {
    loading.value = false
  }
}

const clearLogs = () => {
  logs.value = ''
}

const copyLogs = async () => {
  if (!logs.value) {
    ElMessage.warning('暂无日志内容')
    return
  }

  try {
    await navigator.clipboard.writeText(logs.value)
    ElMessage.success('日志已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

const downloadLogs = () => {
  if (!logs.value) {
    ElMessage.warning('暂无日志内容')
    return
  }

  const blob = new Blob([logs.value], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${props.pod?.metadata.name || 'pod'}-logs.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  ElMessage.success('日志下载成功')
}

const handleClose = () => {
  dialogVisible.value = false
}

// 监听对话框显示
const watchDialogVisible = () => {
  if (dialogVisible.value && props.pod) {
    // 设置默认容器
    if (props.pod.spec.containers.length > 0) {
      logConfig.container = props.pod.spec.containers[0].name
    }
    loadLogs()
  }
}

// 生命周期
onMounted(() => {
  watchDialogVisible()
})
</script>

<style scoped>
.pod-log {
  max-height: 70vh;
  overflow-y: auto;
}

.config-card {
  margin-bottom: 20px;
}

.log-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.log-content {
  min-height: 400px;
  max-height: 500px;
  overflow-y: auto;
}

.log-text {
  margin: 0;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
}
</style> 