<template>
  <PageLayout>
    <template #header>
      <h1>PersistentVolume管理</h1>
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
        :data="persistentVolumes"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <K8sStatusBadge :status="row.status?.phase || 'Unknown'" />
          </template>
        </el-table-column>
        <el-table-column label="容量" width="100">
          <template #default="{ row }">
            {{ row.spec?.capacity?.storage || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="spec.storageClassName" label="存储类" width="120" />
        <el-table-column label="访问模式" width="150">
          <template #default="{ row }">
            {{ getAccessModes(row) }}
          </template>
        </el-table-column>
        <el-table-column label="回收策略" width="120">
          <template #default="{ row }">
            {{ row.spec?.persistentVolumeReclaimPolicy || '-' }}
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
    title="PersistentVolume详情"
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
  name: 'PersistentVolumeManagement'
})

// 响应式数据
const loading = ref(false)
const persistentVolumes = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  status: '',
  keyword: ''
})

// 方法
const fetchPersistentVolumes = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取PersistentVolume列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    persistentVolumes.value = []
  } catch (error) {
    console.error('获取PersistentVolume列表失败:', error)
    ElMessage.error('获取PersistentVolume列表失败')
  } finally {
    loading.value = false
  }
}

const getAccessModes = (pv: any): string => {
  const modes = pv.spec?.accessModes || []
  return modes.join(', ') || '-'
}

const handleCreate = () => {
  ElMessage.info('创建PersistentVolume功能开发中')
}

const handleRefresh = () => {
  fetchPersistentVolumes()
}

const handleSearch = () => {
  fetchPersistentVolumes()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑PersistentVolume功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除PersistentVolume功能开发中')
}

// 生命周期
onMounted(() => {
  fetchPersistentVolumes()
})
</script>