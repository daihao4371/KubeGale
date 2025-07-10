<template>
  <PageLayout>
    <template #header>
      <h1>ResourceQuota管理</h1>
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
        :data="resourceQuotas"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="CPU限制" width="120">
          <template #default="{ row }">
            {{ getCpuQuota(row) }}
          </template>
        </el-table-column>
        <el-table-column label="内存限制" width="120">
          <template #default="{ row }">
            {{ getMemoryQuota(row) }}
          </template>
        </el-table-column>
        <el-table-column label="Pod数量限制" width="120">
          <template #default="{ row }">
            {{ getPodQuota(row) }}
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
    title="ResourceQuota详情"
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
  name: 'ResourceQuotaManagement'
})

// 响应式数据
const loading = ref(false)
const resourceQuotas = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchResourceQuotas = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取ResourceQuota列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    resourceQuotas.value = []
  } catch (error) {
    console.error('获取ResourceQuota列表失败:', error)
    ElMessage.error('获取ResourceQuota列表失败')
  } finally {
    loading.value = false
  }
}

const getCpuQuota = (quota: any): string => {
  const used = quota.status?.used?.['requests.cpu'] || quota.status?.used?.['cpu'] || '0'
  const hard = quota.spec?.hard?.['requests.cpu'] || quota.spec?.hard?.['cpu']
  
  if (!hard) return '-'
  return `${used}/${hard}`
}

const getMemoryQuota = (quota: any): string => {
  const used = quota.status?.used?.['requests.memory'] || quota.status?.used?.['memory'] || '0'
  const hard = quota.spec?.hard?.['requests.memory'] || quota.spec?.hard?.['memory']
  
  if (!hard) return '-'
  return `${used}/${hard}`
}

const getPodQuota = (quota: any): string => {
  const used = quota.status?.used?.['count/pods'] || '0'
  const hard = quota.spec?.hard?.['count/pods']
  
  if (!hard) return '-'
  return `${used}/${hard}`
}

const handleCreate = () => {
  ElMessage.info('创建ResourceQuota功能开发中')
}

const handleRefresh = () => {
  fetchResourceQuotas()
}

const handleSearch = () => {
  fetchResourceQuotas()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑ResourceQuota功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除ResourceQuota功能开发中')
}

// 生命周期
onMounted(() => {
  fetchResourceQuotas()
})
</script>