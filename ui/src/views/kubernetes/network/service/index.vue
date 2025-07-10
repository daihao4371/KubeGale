<template>
  <PageLayout>
    <template #header>
      <h1>Service管理</h1>
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
        :data="services"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column prop="spec.type" label="类型" width="120" />
        <el-table-column label="集群IP" width="150">
          <template #default="{ row }">
            {{ row.spec?.clusterIP || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="端口" width="120">
          <template #default="{ row }">
            {{ getServicePorts(row) }}
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
    title="Service详情"
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
  name: 'ServiceManagement'
})

// 响应式数据
const loading = ref(false)
const services = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchServices = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取Service列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    services.value = []
  } catch (error) {
    console.error('获取Service列表失败:', error)
    ElMessage.error('获取Service列表失败')
  } finally {
    loading.value = false
  }
}

const getServicePorts = (service: any): string => {
  const ports = service.spec?.ports || []
  if (ports.length === 0) return '-'
  
  return ports.map((port: any) => `${port.port}${port.protocol ? `/${port.protocol}` : ''}`).join(', ')
}

const handleCreate = () => {
  ElMessage.info('创建Service功能开发中')
}

const handleRefresh = () => {
  fetchServices()
}

const handleSearch = () => {
  fetchServices()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑Service功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除Service功能开发中')
}

// 生命周期
onMounted(() => {
  fetchServices()
})
</script>