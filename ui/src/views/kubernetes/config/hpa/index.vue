<template>
  <PageLayout>
    <template #header>
      <h1>HPA管理</h1>
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
        :data="hpas"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="目标" width="150">
          <template #default="{ row }">
            {{ getScaleTargetRef(row) }}
          </template>
        </el-table-column>
        <el-table-column label="副本数" width="120">
          <template #default="{ row }">
            {{ getCurrentReplicas(row) }}
          </template>
        </el-table-column>
        <el-table-column label="副本范围" width="120">
          <template #default="{ row }">
            {{ getReplicaRange(row) }}
          </template>
        </el-table-column>
        <el-table-column label="CPU使用率" width="120">
          <template #default="{ row }">
            {{ getCpuUtilization(row) }}
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
    title="HPA详情"
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
  name: 'HPAManagement'
})

// 响应式数据
const loading = ref(false)
const hpas = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchHPAs = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取HPA列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    hpas.value = []
  } catch (error) {
    console.error('获取HPA列表失败:', error)
    ElMessage.error('获取HPA列表失败')
  } finally {
    loading.value = false
  }
}

const getScaleTargetRef = (hpa: any): string => {
  const target = hpa.spec?.scaleTargetRef
  if (!target) return '-'
  
  return `${target.kind}/${target.name}`
}

const getCurrentReplicas = (hpa: any): string => {
  const current = hpa.status?.currentReplicas || 0
  const desired = hpa.status?.desiredReplicas || 0
  
  return `${current}/${desired}`
}

const getReplicaRange = (hpa: any): string => {
  const min = hpa.spec?.minReplicas || 1
  const max = hpa.spec?.maxReplicas || 1
  
  return `${min}-${max}`
}

const getCpuUtilization = (hpa: any): string => {
  const current = hpa.status?.currentCPUUtilizationPercentage
  const target = hpa.spec?.targetCPUUtilizationPercentage
  
  if (current === undefined || target === undefined) return '-'
  
  return `${current}%/${target}%`
}

const handleCreate = () => {
  ElMessage.info('创建HPA功能开发中')
}

const handleRefresh = () => {
  fetchHPAs()
}

const handleSearch = () => {
  fetchHPAs()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑HPA功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除HPA功能开发中')
}

// 生命周期
onMounted(() => {
  fetchHPAs()
})
</script>