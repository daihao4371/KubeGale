<template>
  <div class="platform-region-tree-select">
    <el-input
      v-model="searchKeyword"
      placeholder="搜索云平台/区域"
      clearable
      class="search-input"
    />
    <el-tree
      ref="treeRef"
      :data="filteredTreeData"
      :props="defaultProps"
      node-key="id"
      highlight-current
      :expand-on-click-node="false"
      @node-click="handleNodeClick"
      :filter-node-method="filterNode"
      default-expand-all
    >
      <template #default="{ data }">
        <span>
          <el-tag v-if="!data.region_id" size="small" type="info">{{ data.name }}</el-tag>
          <span v-else>{{ data.region_name }}</span>
        </span>
      </template>
    </el-tree>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getPlatformRegionTree } from '@/api/cloudCmdb/cloud_rds'
import type { PlatformRegionTreeItem, RegionNode } from '@/api/cloudCmdb/cloud_rds'

const emit = defineEmits<{
  (e: 'select', payload: { platform: PlatformRegionTreeItem, region: RegionNode | null }): void
}>()

const treeRef = ref()
const treeData = ref<PlatformRegionTreeItem[]>([])
const searchKeyword = ref('')

// el-tree 配置
const defaultProps = {
  children: 'children',
  label: (data: PlatformRegionTreeItem | RegionNode) => 'region_name' in data ? data.region_name : data.name
}

// 加载数据并转换结构
const loadTreeData = async () => {
  const res = await getPlatformRegionTree()
  if (res.code === 0) {
    // 转换 region 为 children
    treeData.value = res.data.map(platform => ({
      ...platform,
      children: (platform.region || []).map(region => ({
        ...region,
        parentPlatform: platform
      }))
    }))
  }
}
onMounted(loadTreeData)

// 搜索过滤
const filteredTreeData = computed(() => {
  if (!searchKeyword.value) return treeData.value
  const keyword = searchKeyword.value.toLowerCase()
  return treeData.value
    .map(platform => {
      const p = platform as PlatformRegionTreeItem & { children?: RegionNode[] }
      const matchedChildren = (p.children || []).filter(
        (region: RegionNode) =>
          region.region_name.toLowerCase().includes(keyword) ||
          region.region_id.toLowerCase().includes(keyword)
      )
      if (
        p.name.toLowerCase().includes(keyword) ||
        matchedChildren.length
      ) {
        return {
          ...p,
          children: matchedChildren
        }
      }
      return null
    })
    .filter(Boolean)
})

// 过滤节点
const filterNode = (value: string, data: PlatformRegionTreeItem | RegionNode) => {
  if (!value) return true
  if ('region_name' in data) {
    return (
      data.region_name.toLowerCase().includes(value.toLowerCase()) ||
      data.region_id.toLowerCase().includes(value.toLowerCase())
    )
  }
  return data.name.toLowerCase().includes(value.toLowerCase())
}

// 处理节点点击
const handleNodeClick = (data: PlatformRegionTreeItem | (RegionNode & { parentPlatform: PlatformRegionTreeItem })) => {
  if ('region_id' in data && data.parentPlatform) {
    // 区域节点
    emit('select', { platform: data.parentPlatform, region: data })
  } else if (!('region_id' in data)) {
    // 平台节点
    emit('select', { platform: data, region: null })
  }
}
</script>

<style scoped>
.platform-region-tree-select {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  /* background-color: var(--bg-color-card); /* Or transparent if part of a larger card */
  /* padding: var(--spacing-md); */ /* Add padding if it's a standalone card-like element */
  /* border-radius: var(--border-radius-md); */ /* Add radius if standalone */
}

.search-input {
  margin-bottom: var(--spacing-md); /* Use theme variable */
}

/* Styling for el-tree items. Element Plus uses CSS variables for this too.
   Refer to Element Plus documentation for their specific CSS variables for el-tree
   if deeper customization than hover/active is needed.
   Example: --el-tree-node-hover-bg-color
*/
:deep(.el-tree-node__content) {
  padding-top: var(--spacing-xs);
  padding-bottom: var(--spacing-xs);
  border-radius: var(--border-radius-sm);
}

:deep(.el-tree-node__content:hover) {
  background-color: var(--color-primary-light-hover); /* Use a light hover from theme */
}

:deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: var(--color-primary-light) !important; /* Use primary for selection */
  color: var(--text-color-on-primary); /* Text color on primary background */
  font-weight: var(--font-weight-bold);
}

:deep(.el-tree-node.is-current > .el-tree-node__content .el-tag),
:deep(.el-tree-node.is-current > .el-tree-node__content span) { /* Ensure text and tag color on selection */
  color: var(--text-color-on-primary);
}
:deep(.el-tree-node.is-current > .el-tree-node__content .el-tag .el-tag__close){
  color: var(--text-color-on-primary);
}


/* Default state text color for tree items if needed, usually inherited */
:deep(.el-tree-node__label) {
  color: var(--text-color-primary);
  font-size: var(--font-size-base);
}
:deep(.el-tag){
  background-color: var(--bg-color-page) !important;
  color: var(--text-color-secondary) !important;
  border-color: var(--border-color-light) !important;
}
:deep(.el-tree-node.is-current > .el-tree-node__content .el-tag){
  background-color: var(--color-primary) !important;
  color: var(--text-color-on-primary) !important;
  border-color: var(--color-primary-dark) !important;
}

/* Search input styling - should mostly inherit from global el-input styles in main.css */
/* .search-input :deep(.el-input__wrapper) {
  background-color: var(--bg-color-input);
  border-radius: var(--border-radius-md);
  box-shadow: 0 0 0 1px var(--border-color-light) inset;
}
.search-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--color-primary) inset;
}
.search-input :deep(.el-input__inner) {
  color: var(--text-color-primary);
} */
</style>
