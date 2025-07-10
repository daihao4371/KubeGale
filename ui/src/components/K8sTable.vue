<template>
  <div class="k8s-table">
    <!-- 表格工具栏 -->
    <div v-if="showToolbar" class="table-toolbar">
      <div class="toolbar-left">
        <slot name="toolbar-left" />
      </div>
      <div class="toolbar-right">
        <el-button
          v-if="showRefresh"
          :icon="Refresh"
          @click="handleRefresh"
          :loading="refreshing"
          circle
        />
        <el-button
          v-if="showFullscreen"
          :icon="FullScreen"
          @click="handleFullscreen"
          circle
        />
        <slot name="toolbar-right" />
      </div>
    </div>

    <!-- 表格主体 -->
    <el-table
      ref="tableRef"
      :data="data"
      :loading="loading"
      :height="height"
      :max-height="maxHeight"
      :border="border"
      :stripe="stripe"
      :size="size"
      :empty-text="emptyText"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
      @filter-change="handleFilterChange"
      @row-click="handleRowClick"
      @row-dblclick="handleRowDblClick"
    >
      <!-- 多选列 -->
      <el-table-column
        v-if="showSelection"
        type="selection"
        width="55"
        align="center"
        fixed="left"
      />
      
      <!-- 序号列 -->
      <el-table-column
        v-if="showIndex"
        type="index"
        label="序号"
        width="60"
        align="center"
        fixed="left"
      />
      
      <!-- 自定义列 -->
      <slot />
      
      <!-- 操作列 -->
      <el-table-column
        v-if="showActions"
        label="操作"
        :width="actionWidth"
        fixed="right"
        align="center"
      >
        <template #default="scope">
          <slot name="actions" :row="scope.row" :index="scope.$index" />
        </template>
      </el-table-column>
    </el-table>
    
    <!-- 分页器 -->
    <el-pagination
      v-if="showPagination && total > 0"
      :current-page="currentPage"
      :page-size="pageSize"
      :total="total"
      :page-sizes="pageSizes"
      :layout="paginationLayout"
      :background="paginationBackground"
      class="pagination"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Refresh, FullScreen } from '@element-plus/icons-vue'

interface Props {
  data: any[]
  loading?: boolean
  height?: string | number
  maxHeight?: string | number
  border?: boolean
  stripe?: boolean
  size?: 'large' | 'default' | 'small'
  emptyText?: string
  showSelection?: boolean
  showIndex?: boolean
  showActions?: boolean
  showPagination?: boolean
  showToolbar?: boolean
  showRefresh?: boolean
  showFullscreen?: boolean
  actionWidth?: number
  total?: number
  currentPage?: number
  pageSize?: number
  pageSizes?: number[]
  paginationLayout?: string
  paginationBackground?: boolean
  refreshing?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  border: true,
  stripe: true,
  size: 'default',
  emptyText: '暂无数据',
  showSelection: false,
  showIndex: false,
  showActions: false,
  showPagination: true,
  showToolbar: false,
  showRefresh: true,
  showFullscreen: false,
  actionWidth: 150,
  total: 0,
  currentPage: 1,
  pageSize: 20,
  pageSizes: () => [10, 20, 50, 100],
  paginationLayout: 'total, sizes, prev, pager, next, jumper',
  paginationBackground: true,
  refreshing: false,
})

const emit = defineEmits<{
  selectionChange: [selection: any[]]
  sizeChange: [size: number]
  currentChange: [page: number]
  sortChange: [sort: { column: any; prop: string; order: string }]
  filterChange: [filters: any]
  rowClick: [row: any, column: any, event: Event]
  rowDblclick: [row: any, column: any, event: Event]
  refresh: []
  fullscreen: []
}>()

const tableRef = ref()

const handleSelectionChange = (selection: any[]) => {
  emit('selectionChange', selection)
}

const handleSizeChange = (size: number) => {
  emit('sizeChange', size)
}

const handleCurrentChange = (page: number) => {
  emit('currentChange', page)
}

const handleSortChange = (sort: { column: any; prop: string; order: string }) => {
  emit('sortChange', sort)
}

const handleFilterChange = (filters: any) => {
  emit('filterChange', filters)
}

const handleRowClick = (row: any, column: any, event: Event) => {
  emit('rowClick', row, column, event)
}

const handleRowDblClick = (row: any, column: any, event: Event) => {
  emit('rowDblclick', row, column, event)
}

const handleRefresh = () => {
  emit('refresh')
}

const handleFullscreen = () => {
  emit('fullscreen')
}

// 暴露表格实例方法
defineExpose({
  tableRef,
  clearSelection: () => tableRef.value?.clearSelection(),
  toggleRowSelection: (row: any, selected?: boolean) => 
    tableRef.value?.toggleRowSelection(row, selected),
  toggleAllSelection: () => tableRef.value?.toggleAllSelection(),
  setCurrentRow: (row: any) => tableRef.value?.setCurrentRow(row),
  clearSort: () => tableRef.value?.clearSort(),
  clearFilter: (columnKeys?: string[]) => tableRef.value?.clearFilter(columnKeys),
  doLayout: () => tableRef.value?.doLayout(),
  sort: (prop: string, order: string) => tableRef.value?.sort(prop, order),
})
</script>

<style scoped>
.k8s-table {
  width: 100%;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  background: #fafbfc;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pagination {
  margin-top: 20px;
  padding: 0 20px 20px;
  display: flex;
  justify-content: flex-end;
}

/* 表格内容样式优化 */
.k8s-table :deep(.el-table) {
  border-radius: 0;
}

.k8s-table :deep(.el-table__header) {
  background: #f8f9fa;
}

.k8s-table :deep(.el-table th) {
  background: #f8f9fa;
  color: #495057;
  font-weight: 600;
  border-bottom: 1px solid #dee2e6;
}

.k8s-table :deep(.el-table td) {
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.k8s-table :deep(.el-table__row:hover) {
  background-color: #f8f9fa;
}

.k8s-table :deep(.el-table__empty-text) {
  color: #6c757d;
}

/* 操作按钮样式 */
.k8s-table :deep(.el-button--small) {
  padding: 4px 8px;
  font-size: 12px;
}

.k8s-table :deep(.el-button + .el-button) {
  margin-left: 4px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .table-toolbar {
    flex-direction: column;
    gap: 12px;
    padding: 12px 15px;
  }
  
  .toolbar-left,
  .toolbar-right {
    width: 100%;
    justify-content: center;
  }
  
  .pagination {
    padding: 0 15px 15px;
    justify-content: center;
  }
  
  .k8s-table :deep(.el-pagination) {
    flex-wrap: wrap;
    justify-content: center;
  }
}

/* 加载状态优化 */
.k8s-table :deep(.el-loading-mask) {
  border-radius: 8px;
}

/* 选择列样式优化 */
.k8s-table :deep(.el-table__column--selection .el-checkbox) {
  display: flex;
  justify-content: center;
}
</style> 