<template>
  <PageLayout>
    <template #header>
      <h1>StatefulSet管理</h1>
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
        :data="statefulsets"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <K8sStatusBadge :status="getStatefulSetStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column label="就绪副本" width="120">
          <template #default="{ row }">
            {{ row.status?.readyReplicas || 0 }}/{{ row.spec?.replicas || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="运行时间" width="120">
          <template #default="{ row }">
            {{ formatAge(row.metadata?.creationTimestamp) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" @click="handleScale(row)">扩缩容</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </K8sTable>
    </template>
  </PageLayout>
  
  <ResourceDrawer
    v-model="drawerVisible"
    :resource="selectedResource"
    title="StatefulSet详情"
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
  name: 'StatefulSetManagement'
})

// 响应式数据
const loading = ref(false)
const statefulsets = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchStatefulSets = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取StatefulSet列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    statefulsets.value = []
  } catch (error) {
    console.error('获取StatefulSet列表失败:', error)
    ElMessage.error('获取StatefulSet列表失败')
  } finally {
    loading.value = false
  }
}

const getStatefulSetStatus = (statefulset: any): string => {
  const readyReplicas = statefulset.status?.readyReplicas || 0
  const replicas = statefulset.spec?.replicas || 0
  
  if (readyReplicas === replicas && replicas > 0) {
    return 'Running'
  } else if (readyReplicas > 0) {
    return 'Pending'
  } else {
    return 'Failed'
  }
}

const handleCreate = () => {
  ElMessage.info('创建StatefulSet功能开发中')
}

const handleRefresh = () => {
  fetchStatefulSets()
}

const handleSearch = () => {
  fetchStatefulSets()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑StatefulSet功能开发中')
}

const handleScale = (row: any) => {
  ElMessage.info('扩缩容StatefulSet功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除StatefulSet功能开发中')
}

// 生命周期
onMounted(() => {
  fetchStatefulSets()
})
</script>