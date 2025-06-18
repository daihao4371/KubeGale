<template>
  <div class="region-tree-select">
    <el-popover
      v-model:visible="popoverVisible"
      placement="bottom-start"
      :width="300"
      trigger="click"
    >
      <template #reference>
        <el-input
          v-model="selectedLabel"
          :placeholder="placeholder"
          readonly
          clearable
          @clear="handleClear"
        >
          <template #prefix>
            <el-icon><Location /></el-icon>
          </template>
        </el-input>
      </template>
      
      <div class="region-tree-container">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索区域"
          clearable
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-tree
          ref="treeRef"
          :data="treeData"
          :props="defaultProps"
          :filter-node-method="filterNode"
          node-key="id"
          highlight-current
          @node-click="handleNodeClick"
        >
          <template #default="{ data }">
            <div class="custom-tree-node">
              <span>{{ data.region_name }}</span>
              <el-tag size="small" type="info">{{ data.region_id }}</el-tag>
            </div>
          </template>
        </el-tree>
      </div>
    </el-popover>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Location, Search } from '@element-plus/icons-vue'
import type { Region } from '@/api/cloudCmdb/cloud_rds'

const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const popoverVisible = ref(false)
const searchKeyword = ref('')
const treeRef = ref()
const treeData = ref<Region[]>([])

const defaultProps = {
  children: 'children',
  label: 'region_name'
}

// 计算选中的标签
const selectedLabel = computed(() => {
  const selected = treeData.value.find(item => item.region_id === props.modelValue)
  return selected ? selected.region_name : ''
})

// 监听搜索关键词变化
watch(searchKeyword, (val) => {
  treeRef.value?.filter(val)
})

// 过滤节点
const filterNode = (value: string, data: Region) => {
  if (!value) return true
  return (
    data.region_name.toLowerCase().includes(value.toLowerCase()) ||
    data.region_id.toLowerCase().includes(value.toLowerCase())
  )
}

// 处理节点点击
const handleNodeClick = (data: Region) => {
  emit('update:modelValue', data.region_id)
  popoverVisible.value = false
}

// 处理清除
const handleClear = () => {
  emit('update:modelValue', '')
}

// 设置树形数据
const setTreeData = (data: Region[]) => {
  treeData.value = data
}

defineExpose({
  setTreeData
})
</script>

<style scoped>
.region-tree-select {
  width: 100%;
}

.region-tree-container {
  height: 300px;
  display: flex;
  flex-direction: column;
}

.search-input {
  margin-bottom: 8px;
}

.custom-tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding-right: 8px;
}

:deep(.el-tree-node__content) {
  height: 32px;
}

:deep(.el-tree-node__content:hover) {
  background-color: #f5f7fa;
}
</style> 