<template>
  <div class="yaml-editor-container">
    <div class="editor-header">
      <h3>YAML 编辑器</h3>
      <div class="header-actions">
        <el-button size="small" @click="handleFormat">格式化</el-button>
        <el-button size="small" @click="handleValidate">验证</el-button>
      </div>
    </div>
    
    <div class="editor-content">
      <el-input
        v-model="yamlContent"
        type="textarea"
        :rows="20"
        placeholder="请输入YAML内容..."
        @input="handleInput"
      />
    </div>
    
    <div class="editor-footer">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleSave">保存</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'

interface Props {
  modelValue?: string
  readonly?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'save', value: string): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  readonly: false
})

const emit = defineEmits<Emits>()

const yamlContent = ref(props.modelValue)

watch(() => props.modelValue, (newValue) => {
  yamlContent.value = newValue
})

const handleInput = (value: string) => {
  emit('update:modelValue', value)
}

const handleFormat = () => {
  try {
    // 简单的YAML格式化逻辑
    const lines = yamlContent.value.split('\n')
    const formatted = lines.map(line => line.trim()).filter(line => line).join('\n')
    yamlContent.value = formatted
    emit('update:modelValue', formatted)
  } catch (error) {
    console.error('YAML格式化失败:', error)
  }
}

const handleValidate = () => {
  try {
    // 简单的YAML验证逻辑
    const lines = yamlContent.value.split('\n')
    let indentLevel = 0
    let isValid = true
    
    for (const line of lines) {
      if (line.trim() === '') continue
      
      const currentIndent = line.search(/\S/)
      if (currentIndent === -1) continue
      
      if (currentIndent > indentLevel + 2) {
        isValid = false
        break
      }
      indentLevel = currentIndent
    }
    
    if (isValid) {
      ElMessage.success('YAML格式正确')
    } else {
      ElMessage.error('YAML格式错误')
    }
  } catch (error) {
    ElMessage.error('YAML验证失败')
  }
}

const handleSave = () => {
  emit('save', yamlContent.value)
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.yaml-editor-container {
  padding: 20px;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.editor-header h3 {
  margin: 0;
  color: var(--el-text-color-primary);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.editor-content {
  margin-bottom: 20px;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 