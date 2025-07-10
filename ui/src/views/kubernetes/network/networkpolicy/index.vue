<template>
  <PageLayout>
    <template #header>
      <h1>NetworkPolicy管理</h1>
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
        :data="networkPolicies"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="Pod选择器" width="200">
          <template #default="{ row }">
            {{ getPodSelector(row) }}
          </template>
        </el-table-column>
        <el-table-column label="策略类型" width="150">
          <template #default="{ row }">
            {{ getPolicyTypes(row) }}
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
    title="NetworkPolicy详情"
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
  name: 'NetworkPolicyManagement'
})

// 响应式数据
const loading = ref(false)
const networkPolicies = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchNetworkPolicies = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取NetworkPolicy列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    networkPolicies.value = []
  } catch (error) {
    console.error('获取NetworkPolicy列表失败:', error)
    ElMessage.error('获取NetworkPolicy列表失败')
  } finally {
    loading.value = false
  }
}

const getPodSelector = (networkPolicy: any): string => {
  const selector = networkPolicy.spec?.podSelector?.matchLabels
  if (!selector || Object.keys(selector).length === 0) return 'All Pods'
  
  return Object.entries(selector).map(([key, value]) => `${key}=${value}`).join(', ')
}

const getPolicyTypes = (networkPolicy: any): string => {
  const types = networkPolicy.spec?.policyTypes || []
  return types.join(', ') || '-'
}

const handleCreate = () => {
  ElMessage.info('创建NetworkPolicy功能开发中')
}

const handleRefresh = () => {
  fetchNetworkPolicies()
}

const handleSearch = () => {
  fetchNetworkPolicies()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑NetworkPolicy功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除NetworkPolicy功能开发中')
}

// 生命周期
onMounted(() => {
  fetchNetworkPolicies()
})
</script>