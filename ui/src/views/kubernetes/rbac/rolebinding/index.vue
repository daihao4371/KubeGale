<template>
  <PageLayout>
    <template #header>
      <h1>RoleBinding管理</h1>
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
        :data="roleBindings"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="角色" width="150">
          <template #default="{ row }">
            {{ getRoleRef(row) }}
          </template>
        </el-table-column>
        <el-table-column label="主体数量" width="120">
          <template #default="{ row }">
            {{ getSubjectCount(row) }}
          </template>
        </el-table-column>
        <el-table-column label="主体类型" width="150">
          <template #default="{ row }">
            {{ getSubjectTypes(row) }}
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
    title="RoleBinding详情"
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
  name: 'RoleBindingManagement'
})

// 响应式数据
const loading = ref(false)
const roleBindings = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchRoleBindings = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取RoleBinding列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    roleBindings.value = []
  } catch (error) {
    console.error('获取RoleBinding列表失败:', error)
    ElMessage.error('获取RoleBinding列表失败')
  } finally {
    loading.value = false
  }
}

const getRoleRef = (roleBinding: any): string => {
  const roleRef = roleBinding.roleRef
  if (!roleRef) return '-'
  
  return `${roleRef.kind}/${roleRef.name}`
}

const getSubjectCount = (roleBinding: any): number => {
  return roleBinding.subjects?.length || 0
}

const getSubjectTypes = (roleBinding: any): string => {
  const subjects = roleBinding.subjects || []
  const types = new Set<string>()
  
  subjects.forEach((subject: any) => {
    if (subject.kind) {
      types.add(subject.kind)
    }
  })
  
  return Array.from(types).join(', ') || '-'
}

const handleCreate = () => {
  ElMessage.info('创建RoleBinding功能开发中')
}

const handleRefresh = () => {
  fetchRoleBindings()
}

const handleSearch = () => {
  fetchRoleBindings()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑RoleBinding功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除RoleBinding功能开发中')
}

// 生命周期
onMounted(() => {
  fetchRoleBindings()
})
</script>