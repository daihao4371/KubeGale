<template>
  <PageLayout>
    <template #header>
      <h1>CronJob管理</h1>
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
        :data="cronjobs"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column prop="spec.schedule" label="调度规则" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <K8sStatusBadge :status="getCronJobStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="上次调度" width="150">
          <template #default="{ row }">
            {{ row.status?.lastScheduleTime ? formatAge(row.status.lastScheduleTime) : '-' }}
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
    title="CronJob详情"
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
  name: 'CronJobManagement'
})

// 响应式数据
const loading = ref(false)
const cronjobs = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchCronJobs = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取CronJob列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    cronjobs.value = []
  } catch (error) {
    console.error('获取CronJob列表失败:', error)
    ElMessage.error('获取CronJob列表失败')
  } finally {
    loading.value = false
  }
}

const getCronJobStatus = (cronjob: any): string => {
  if (cronjob.spec?.suspend) {
    return 'Suspended'
  } else if (cronjob.status?.active?.length > 0) {
    return 'Running'
  } else {
    return 'Ready'
  }
}

const handleCreate = () => {
  ElMessage.info('创建CronJob功能开发中')
}

const handleRefresh = () => {
  fetchCronJobs()
}

const handleSearch = () => {
  fetchCronJobs()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑CronJob功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除CronJob功能开发中')
}

// 生命周期
onMounted(() => {
  fetchCronJobs()
})
</script>