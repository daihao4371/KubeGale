<template>
  <PageLayout>
    <template #header>
      <h1>ClusterRole管理</h1>
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
        :data="clusterRoles"
        :loading="loading"
        @row-click="handleRowClick"
      >
        <el-table-column prop="metadata.name" label="名称" min-width="150" />
        <el-table-column label="规则数量" width="120">
          <template #default="{ row }">
            {{ getRuleCount(row) }}
          </template>
        </el-table-column>
        <el-table-column label="资源类型" width="250">
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
    title="ClusterRole详情"
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
  name: 'ClusterRoleManagement'
})

// 响应式数据
const loading = ref(false)
const clusterRoles = ref([])
const drawerVisible = ref(false)
const selectedResource = ref(null)

const filters = reactive({
  status: '',
  keyword: ''
})

// 方法
const fetchClusterRoles = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取ClusterRole列表
    await new Promise(resolve => setTimeout(resolve, 1000))
    clusterRoles.value = []
  } catch (error) {
    console.error('获取ClusterRole列表失败:', error)
    ElMessage.error('获取ClusterRole列表失败')
  } finally {
    loading.value = false
  }
}

const getRuleCount = (clusterRole: any): number => {
  return clusterRole.rules?.length || 0
}

const getResourceTypes = (clusterRole: any): string => {
  const rules = clusterRole.rules || []
  const resourceTypes = new Set<string>()
  
  rules.forEach((rule: any) => {
    if (rule.resources) {
      rule.resources.forEach((resource: string) => {
        resourceTypes.add(resource)
      })
    }
  })
  
  const types = Array.from(resourceTypes).slice(0, 4)
  return types.join(', ') + (resourceTypes.size > 4 ? '...' : '')
}

const getVerbs = (clusterRole: any): string => {
  const rules = clusterRole.rules || []
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
  ElMessage.info('创建ClusterRole功能开发中')
}

const handleRefresh = () => {
  fetchClusterRoles()
}

const handleSearch = () => {
  fetchClusterRoles()
}

const handleRowClick = (row: any) => {
  selectedResource.value = row
  drawerVisible.value = true
}

const handleEdit = (row: any) => {
  ElMessage.info('编辑ClusterRole功能开发中')
}

const handleDelete = (row: any) => {
  ElMessage.info('删除ClusterRole功能开发中')
}

// 生命周期
onMounted(() => {
  fetchClusterRoles()
})
</script>