<template>
  <PageLayout>
    <template #header>
      <h1>Job管理</h1>
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
        :data="jobs"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <K8sStatusBadge :status="getJobStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="完成情况" width="120">
          <template #default="{ row }">
            {{ row.status?.succeeded || 0 }}/{{ row.spec?.completions || 1 }}
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
    title="Job详情"
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
  name: 'JobManagement'
})

// 响应式数据
const loading = ref(false)
const jobs = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchJobs = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取Job列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    jobs.value = []
  } catch (error) {
    console.error('获取Job列表失败:', error)
    ElMessage.error('获取Job列表失败')
  } finally {
    loading.value = false
  }
}

const getJobStatus = (job: any): string => {
  if (job.status?.succeeded) {
    return 'Succeeded'
  } else if (job.status?.failed) {
    return 'Failed'
  } else if (job.status?.active) {
    return 'Running'
  } else {
    return 'Pending'
  }
}

const handleCreate = () => {
  ElMessage.info('创建Job功能开发中')
}

const handleRefresh = () => {
  fetchJobs()
}

const handleSearch = () => {
  fetchJobs()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑Job功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除Job功能开发中')
}

// 生命周期
onMounted(() => {
  fetchJobs()
})
</script>