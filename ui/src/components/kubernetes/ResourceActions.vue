<template>
  <div class="resource-actions">
    <div class="bulk-actions" v-if="selectedCount > 0">
      <span class="selection-info">已选择 {{ selectedCount }} 项</span>
      <el-button 
        v-if="showBulkDelete"
        type="danger" 
        size="small"
        @click="handleBulkDelete"
      >
        批量删除
      </el-button>
      <slot name="bulk-actions" :selected="selectedItems" />
    </div>
    
    <div class="item-actions">
      <el-button 
        v-if="showRefresh"
        :icon="Refresh" 
        @click="handleRefresh"
        :loading="refreshing"
      >
        刷新
      </el-button>
      <el-button 
        v-if="showCreate"
        type="primary" 
        :icon="Plus"
        @click="handleCreate"
      >
        {{ createText }}
      </el-button>
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Plus, Refresh } from '@element-plus/icons-vue'

interface Props {
  selectedItems?: any[]
  showBulkDelete?: boolean
  showRefresh?: boolean
  showCreate?: boolean
  createText?: string
  refreshing?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  selectedItems: () => [],
  showBulkDelete: true,
  showRefresh: true,
  showCreate: true,
  createText: '新建',
  refreshing: false,
})

const emit = defineEmits<{
  bulkDelete: [items: any[]]
  refresh: []
  create: []
}>()

const selectedCount = computed(() => props.selectedItems.length)

const handleBulkDelete = () => {
  emit('bulkDelete', props.selectedItems)
}

const handleRefresh = () => {
  emit('refresh')
}

const handleCreate = () => {
  emit('create')
}
</script>

<style scoped>
.resource-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 0 20px;
  min-height: 40px;
}

.bulk-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.selection-info {
  font-size: 14px;
  color: #409eff;
  font-weight: 500;
}

.item-actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 768px) {
  .resource-actions {
    flex-direction: column;
    gap: 12px;
    padding: 0 15px;
  }
  
  .bulk-actions,
  .item-actions {
    width: 100%;
    justify-content: center;
  }
}
</style>