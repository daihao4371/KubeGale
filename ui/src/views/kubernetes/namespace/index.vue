<template>
  <PageLayout>
    <template #header>
      <h1>命名空间管理</h1>
    </template>
    
    <template #actions>
      <ResourceActions
        @create="handleCreate"
        @refresh="handleRefresh"
      />
    </template>
    
    <template #search>
      <SearchFilter
        v-model:status="filters.status"
        v-model:keyword="filters.keyword"
        @search="handleSearch"
        :show-namespace="false"
      />
    </template>
    
    <template #default>
      <K8sTable
        :data="namespaces"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <K8sStatusBadge :status="row.status?.phase || 'Unknown'" />
          </template>
        </el-table-column>
        <el-table-column label="标签" width="200">
          <template #default="{ row }">
            {{ getLabels(row) }}
          </template>
        </el-table-column>
        <el-table-column label="运行时间" width="120">
          <template #default="{ row }">
            {{ formatAge(row.metadata?.creationTimestamp) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)" :disabled="isSystemNamespace(row)">删除</el-button>
          </template>
        </el-table-column>
      </K8sTable>
    </template>
  </PageLayout>
  
  <ResourceDrawer
    v-model="drawerVisible"
    :resource="selectedResource"
    title="命名空间详情"
  />
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import PageLayout from '@/components/layout/PageLayout.vue'
import ResourceActions from '@/components/kubernetes/ResourceActions.vue'
import SearchFilter from '@/components/kubernetes/SearchFilter.vue'
import ResourceDrawer from '@/components/kubernetes/ResourceDrawer.vue'
import K8sTable from '@/components/K8sTable.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'NamespaceManagement'
})

// 响应式数据
const loading = ref(false)
const namespaces = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  status: '',
  keyword: ''
})

// 方法
const fetchNamespaces = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取命名空间列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    namespaces.value = []
  } catch (error) {
    console.error('获取命名空间列表失败:', error)
    ElMessage.error('获取命名空间列表失败')
  } finally {
    loading.value = false
  }
}

const getLabels = (namespace: any): string => {
  const labels = namespace.metadata?.labels || {}
  const labelStrings = Object.entries(labels)
    .slice(0, 2)
    .map(([key, value]) => `${key}=${value}`)
  
  return labelStrings.length > 0 ? labelStrings.join(', ') + (Object.keys(labels).length > 2 ? '...' : '') : '-'
}

const isSystemNamespace = (namespace: any): boolean => {
  const systemNamespaces = ['default', 'kube-system', 'kube-public', 'kube-node-lease']
  return systemNamespaces.includes(namespace.metadata?.name)
}

const handleCreate = () => {
  ElMessage.info('创建命名空间功能开发中')
}

const handleRefresh = () => {
  fetchNamespaces()
}

const handleSearch = () => {
  fetchNamespaces()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑命名空间功能开发中')
}

const handleDelete = (row: any) => {
  if (isSystemNamespace(row)) {
    ElMessage.warning('系统命名空间不允许删除')
    return
  }
  ElMessage.info('删除命名空间功能开发中')
}

// 生命周期
onMounted(() => {
  fetchNamespaces()
})
</script> 