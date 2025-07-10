<template>
  <el-drawer
    v-model="visible"
    :title="title"
    :size="size"
    :direction="direction"
    :before-close="handleClose"
    class="resource-drawer"
  >
    <div class="drawer-content">
      <!-- 资源基本信息 -->
      <div class="resource-info">
        <h3>基本信息</h3>
        <div class="info-grid">
          <div class="info-item" v-for="(item, key) in basicInfo" :key="key">
            <label>{{ item.label }}:</label>
            <span :class="item.class">{{ item.value }}</span>
          </div>
        </div>
      </div>

      <!-- 标签和注解 -->
      <div v-if="labels && Object.keys(labels).length > 0" class="resource-section">
        <h3>标签</h3>
        <div class="tags-container">
          <el-tag 
            v-for="(value, key) in labels" 
            :key="key"
            size="small"
            class="tag-item"
          >
            {{ key }}: {{ value }}
          </el-tag>
        </div>
      </div>

      <div v-if="annotations && Object.keys(annotations).length > 0" class="resource-section">
        <h3>注解</h3>
        <div class="annotations-container">
          <div 
            v-for="(value, key) in annotations" 
            :key="key"
            class="annotation-item"
          >
            <strong>{{ key }}:</strong> {{ value }}
          </div>
        </div>
      </div>

      <!-- 自定义内容 -->
      <div v-if="$slots.default" class="resource-section">
        <slot />
      </div>

      <!-- YAML 配置 -->
      <div v-if="showYaml && yamlContent" class="resource-section">
        <h3>YAML 配置</h3>
        <el-input
          :model-value="yamlContent"
          type="textarea"
          :rows="20"
          readonly
          class="yaml-content"
        />
      </div>
    </div>

    <template #footer>
      <div class="drawer-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button v-if="showEdit" type="primary" @click="handleEdit">编辑</el-button>
        <el-button v-if="showDelete" type="danger" @click="handleDelete">删除</el-button>
        <slot name="footer" />
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'

interface BasicInfoItem {
  label: string
  value: string
  class?: string
}

interface Props {
  modelValue: boolean
  title: string
  size?: string | number
  direction?: 'ltr' | 'rtl' | 'ttb' | 'btt'
  basicInfo?: Record<string, BasicInfoItem>
  labels?: Record<string, string>
  annotations?: Record<string, string>
  yamlContent?: string
  showYaml?: boolean
  showEdit?: boolean
  showDelete?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: '60%',
  direction: 'rtl',
  basicInfo: () => ({}),
  labels: () => ({}),
  annotations: () => ({}),
  yamlContent: '',
  showYaml: true,
  showEdit: false,
  showDelete: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  edit: []
  delete: []
  close: []
}>()

const visible = ref(false)

watch(
  () => props.modelValue,
  (newValue) => {
    visible.value = newValue
  },
  { immediate: true }
)

watch(visible, (newValue) => {
  emit('update:modelValue', newValue)
})

const handleClose = () => {
  visible.value = false
  emit('close')
}

const handleEdit = () => {
  emit('edit')
}

const handleDelete = () => {
  emit('delete')
}
</script>

<style scoped>
.resource-drawer :deep(.el-drawer__header) {
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.resource-drawer :deep(.el-drawer__body) {
  padding: 0;
}

.drawer-content {
  padding: 20px;
}

.resource-info {
  margin-bottom: 24px;
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

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f2f5;
}

.info-item label {
  min-width: 100px;
  color: #606266;
  font-weight: 500;
}

.info-item span {
  flex: 1;
  color: #2c3e50;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  word-break: break-all;
}

.annotations-container {
  max-height: 200px;
  overflow-y: auto;
}

.annotation-item {
  padding: 8px 0;
  border-bottom: 1px solid #f0f2f5;
  word-break: break-all;
}

.annotation-item strong {
  color: #606266;
}

.yaml-content {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e4e7ed;
}

/* 状态样式 */
.status-running {
  color: #67c23a;
}

.status-pending {
  color: #e6a23c;
}

.status-failed {
  color: #f56c6c;
}

.status-error {
  color: #f56c6c;
}

.status-success {
  color: #67c23a;
}

.status-warning {
  color: #e6a23c;
}

@media (max-width: 768px) {
  .resource-drawer {
    width: 100% !important;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .drawer-footer {
    flex-direction: column;
  }
}
</style>