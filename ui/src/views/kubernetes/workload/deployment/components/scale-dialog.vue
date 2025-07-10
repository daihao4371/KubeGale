<template>
  <el-dialog
    v-model="dialogVisible"
    title="扩缩容 Deployment"
    width="500px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      v-loading="loading"
    >
      <el-form-item label="Deployment">
        <el-input
          :value="deployment?.metadata.name"
          disabled
          placeholder="Deployment 名称"
        />
      </el-form-item>
      
      <el-form-item label="当前副本数">
        <el-input
          :value="deployment?.spec.replicas"
          disabled
          placeholder="当前副本数"
        />
      </el-form-item>
      
      <el-form-item label="目标副本数" prop="replicas">
        <el-input-number
          v-model="form.replicas"
          :min="0"
          :max="100"
          placeholder="请输入目标副本数"
          style="width: 100%"
        />
      </el-form-item>
      
      <el-form-item label="状态预览">
        <div class="status-preview">
          <el-tag :type="getStatusType()">
            {{ getStatusText() }}
          </el-tag>
          <span class="status-desc">
            {{ getStatusDescription() }}
          </span>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">
        确认扩缩容
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { deploymentApi } from '@/api/kubernetes/workload'
import type { Deployment, ScaleWorkloadRequest } from '@/types/kubernetes/workload'

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
  'success': []
}>()

// 响应式数据
const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 表单数据
const form = reactive({
  replicas: 1
})

// 表单验证规则
const rules: FormRules = {
  replicas: [
    { required: true, message: '请输入目标副本数', trigger: 'blur' },
    { type: 'number', min: 0, max: 100, message: '副本数必须在 0-100 之间', trigger: 'blur' }
  ]
}

// 方法
const getStatusType = () => {
  if (!props.deployment || !form.replicas) return 'info'
  
  const current = props.deployment.spec.replicas
  const target = form.replicas
  
  if (target > current) return 'success'
  if (target < current) return 'warning'
  return 'info'
}

const getStatusText = () => {
  if (!props.deployment || !form.replicas) return '未知'
  
  const current = props.deployment.spec.replicas
  const target = form.replicas
  
  if (target > current) return '扩容'
  if (target < current) return '缩容'
  return '无变化'
}

const getStatusDescription = () => {
  if (!props.deployment || !form.replicas) return ''
  
  const current = props.deployment.spec.replicas
  const target = form.replicas
  
  if (target > current) {
    return `将副本数从 ${current} 增加到 ${target}`
  }
  if (target < current) {
    return `将副本数从 ${current} 减少到 ${target}`
  }
  return '副本数保持不变'
}

const handleSubmit = async () => {
  if (!formRef.value || !props.deployment) return

  try {
    await formRef.value.validate()
    submitting.value = true

    const requestData: ScaleWorkloadRequest = {
      cluster_id: 0, // TODO: 从上下文获取集群ID
      namespace: props.deployment.metadata.namespace,
      name: props.deployment.metadata.name,
      replicas: form.replicas
    }

    await deploymentApi.scaleDeployment(requestData)
    ElMessage.success('扩缩容操作已提交')
    emit('success')
  } catch (error) {
    console.error('扩缩容失败:', error)
    ElMessage.error('扩缩容失败')
  } finally {
    submitting.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
}

// 监听对话框显示状态，初始化表单
watch(() => props.visible, (visible) => {
  if (visible && props.deployment) {
    form.replicas = props.deployment.spec.replicas
  }
})
</script>

<style scoped>
.status-preview {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-desc {
  color: var(--el-text-color-regular);
  font-size: 14px;
}
</style> 