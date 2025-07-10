<template>
  <PageLayout>
    <template #header>
      <h1>Role管理</h1>
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
        :data="roles"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
        <el-table-column label="规则数量" width="120">
          <template #default="{ row }">
            {{ getRuleCount(row) }}
          </template>
        </el-table-column>
        <el-table-column label="资源类型" width="200">
          <template #default="{ row }">
            {{ getResourceTypes(row) }}
          </template>
        </el-table-column>
        <el-table-column label="操作权限" width="150">
          <template #default="{ row }">
            {{ getVerbs(row) }}
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
    title="Role详情"
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
  name: 'RoleManagement'
})

// 响应式数据
const loading = ref(false)
const roles = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  namespace: '',
  status: '',
  keyword: ''
})

// 方法
const fetchRoles = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取Role列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    roles.value = []
  } catch (error) {
    console.error('获取Role列表失败:', error)
    ElMessage.error('获取Role列表失败')
  } finally {
    loading.value = false
  }
}

const getRuleCount = (role: any): number => {
  return role.rules?.length || 0
}

const getResourceTypes = (role: any): string => {
  const rules = role.rules || []
  const resourceTypes = new Set<string>()
  
  rules.forEach((rule: any) => {
    if (rule.resources) {
      rule.resources.forEach((resource: string) => {
        resourceTypes.add(resource)
      })
    }
  })
  
  const types = Array.from(resourceTypes).slice(0, 3)
  return types.join(', ') + (resourceTypes.size > 3 ? '...' : '')
}

const getVerbs = (role: any): string => {
  const rules = role.rules || []
  const verbs = new Set<string>()
  
  rules.forEach((rule: any) => {
    if (rule.verbs) {
      rule.verbs.forEach((verb: string) => {
        verbs.add(verb)
      })
    }
  })
  
  const verbList = Array.from(verbs).slice(0, 3)
  return verbList.join(', ') + (verbs.size > 3 ? '...' : '')
}

const handleCreate = () => {
  ElMessage.info('创建Role功能开发中')
}

const handleRefresh = () => {
  fetchRoles()
}

const handleSearch = () => {
  fetchRoles()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑Role功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除Role功能开发中')
}

// 生命周期
onMounted(() => {
  fetchRoles()
})
</script>