<template>
  <PageLayout>
    <template #header>
      <h1>Endpoint管理</h1>
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
        :data="endpoints"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="端点数量" width="120">
          <template #default="{ row }">
            {{ getEndpointCount(row) }}
          </template>
        </el-table-column>
        <el-table-column label="地址" width="300">
          <template #default="{ row }">
            {{ getEndpointAddresses(row) }}
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
    title="Endpoint详情"
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
  name: 'EndpointManagement'
})

// 响应式数据
const loading = ref(false)
const endpoints = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchEndpoints = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取Endpoint列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    endpoints.value = []
  } catch (error) {
    console.error('获取Endpoint列表失败:', error)
    ElMessage.error('获取Endpoint列表失败')
  } finally {
    loading.value = false
  }
}

const getEndpointCount = (endpoint: any): number => {
  const subsets = endpoint.subsets || []
  return subsets.reduce((count: number, subset: any) => {
    return count + (subset.addresses?.length || 0)
  }, 0)
}

const getEndpointAddresses = (endpoint: any): string => {
  const subsets = endpoint.subsets || []
  if (subsets.length === 0) return '-'
  
  const addresses: string[] = []
  subsets.forEach((subset: any) => {
    const ports = subset.ports || []
    const ips = subset.addresses || []
    
    ips.forEach((addr: any) => {
      if (ports.length > 0) {
        ports.forEach((port: any) => {
          addresses.push(`${addr.ip}:${port.port}`)
        })
      } else {
        addresses.push(addr.ip)
      }
    })
  })
  
  return addresses.slice(0, 3).join(', ') + (addresses.length > 3 ? '...' : '')
}

const handleCreate = () => {
  ElMessage.info('创建Endpoint功能开发中')
}

const handleRefresh = () => {
  fetchEndpoints()
}

const handleSearch = () => {
  fetchEndpoints()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑Endpoint功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除Endpoint功能开发中')
}

// 生命周期
onMounted(() => {
  fetchEndpoints()
})
</script>