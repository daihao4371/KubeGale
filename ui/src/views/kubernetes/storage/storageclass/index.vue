<template>
  <PageLayout>
    <template #header>
      <h1>StorageClass管理</h1>
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
        :data="storageClasses"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="provisioner" label="供应商" width="200" />
        <el-table-column label="回收策略" width="120">
          <template #default="{ row }">
            {{ row.reclaimPolicy || 'Delete' }}
          </template>
        </el-table-column>
        <el-table-column label="允许扩容" width="120">
          <template #default="{ row }">
            {{ row.allowVolumeExpansion ? '是' : '否' }}
          </template>
        </el-table-column>
        <el-table-column label="卷绑定模式" width="150">
          <template #default="{ row }">
            {{ row.volumeBindingMode || 'Immediate' }}
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
    title="StorageClass详情"
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
  name: 'StorageClassManagement'
})

// 响应式数据
const loading = ref(false)
const storageClasses = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  status: '',
  keyword: ''
})

// 方法
const fetchStorageClasses = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取StorageClass列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    storageClasses.value = []
  } catch (error) {
    console.error('获取StorageClass列表失败:', error)
    ElMessage.error('获取StorageClass列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  ElMessage.info('创建StorageClass功能开发中')
}

const handleRefresh = () => {
  fetchStorageClasses()
}

const handleSearch = () => {
  fetchStorageClasses()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑StorageClass功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除StorageClass功能开发中')
}

// 生命周期
onMounted(() => {
  fetchStorageClasses()
})
</script>