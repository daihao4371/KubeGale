<template>
  <PageLayout>
    <template #header>
      <h1>ConfigMap管理</h1>
    </template>
    
    <template #actions>
      <ResourceActions
        @create="handleCreate"
        @refresh="handleRefresh"
      />
    </template>
    
    <template #search>
      <SearchFilter
        v-model:namespace="filters.namespace"
        v-model:status="filters.status"
        v-model:keyword="filters.keyword"
        @search="handleSearch"
      />
    </template>
    
    <template #default>
      <K8sTable
        :data="configMaps"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="数据项数量" width="120">
          <template #default="{ row }">
            {{ getDataCount(row) }}
          </template>
        </el-table-column>
        <el-table-column label="数据大小" width="120">
          <template #default="{ row }">
            {{ getDataSize(row) }}
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
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </K8sTable>
    </template>
  </PageLayout>
  
  <ResourceDrawer
    v-model="drawerVisible"
    :resource="selectedResource"
    title="ConfigMap详情"
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
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'ConfigMapManagement'
})

// 响应式数据
const loading = ref(false)
const configMaps = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchConfigMaps = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取ConfigMap列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    configMaps.value = []
  } catch (error) {
    console.error('获取ConfigMap列表失败:', error)
    ElMessage.error('获取ConfigMap列表失败')
  } finally {
    loading.value = false
  }
}

const getDataCount = (configMap: any): number => {
  const data = configMap.data || {}
  const binaryData = configMap.binaryData || {}
  return Object.keys(data).length + Object.keys(binaryData).length
}

const getDataSize = (configMap: any): string => {
  const data = configMap.data || {}
  const binaryData = configMap.binaryData || {}
  
  let totalSize = 0
  
  // 计算文本数据大小
  Object.values(data).forEach((value: any) => {
    if (typeof value === 'string') {
      totalSize += new Blob([value]).size
    }
  })
  
  // 计算二进制数据大小（Base64编码）
  Object.values(binaryData).forEach((value: any) => {
    if (typeof value === 'string') {
      totalSize += value.length * 0.75 // Base64 编码大约是原始数据的 4/3
    }
  })
  
  if (totalSize < 1024) {
    return `${totalSize} B`
  } else if (totalSize < 1024 * 1024) {
    return `${(totalSize / 1024).toFixed(1)} KB`
  } else {
    return `${(totalSize / (1024 * 1024)).toFixed(1)} MB`
  }
}

const handleCreate = () => {
  ElMessage.info('创建ConfigMap功能开发中')
}

const handleRefresh = () => {
  fetchConfigMaps()
}

const handleSearch = () => {
  fetchConfigMaps()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑ConfigMap功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除ConfigMap功能开发中')
}

// 生命周期
onMounted(() => {
  fetchConfigMaps()
})
</script>