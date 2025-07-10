<template>
  <div class="cluster-users-tab">
    <div class="tab-header">
      <h4>用户管理</h4>
      <el-button type="primary" @click="showCreateDialog = true" icon="Plus">添加用户</el-button>
    </div>
    
    <K8sTable
      :data="clusterUsers"
      :loading="loading"
      :columns="columns"
      :show-pagination="true"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @refresh="fetchUsers"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useClusterStore } from '@/stores/kubernetes/cluster'
import K8sTable from '@/components/K8sTable.vue'

const props = defineProps<{
  clusterId: number
}>()

const clusterStore = useClusterStore()
const loading = ref(false)
const showCreateDialog = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const columns = [
  { prop: 'userName', label: '用户名' },
  { prop: 'nickName', label: '昵称' },
  { prop: 'created_at', label: '创建时间' }
]

const clusterUsers = computed(() => clusterStore.clusterUsers)

const fetchUsers = async () => {
  await clusterStore.fetchClusterUsers(props.clusterId)
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchUsers()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchUsers()
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.tab-header h4 {
  margin: 0;
  color: #303133;
}
</style> 