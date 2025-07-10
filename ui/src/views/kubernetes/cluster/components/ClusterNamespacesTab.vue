<template>
  <div class="cluster-namespaces-tab">
    <div class="tab-header">
      <h4>命名空间</h4>
      <el-button type="primary" @click="showCreateDialog = true" icon="Plus">创建命名空间</el-button>
    </div>
    
    <K8sTable
      :data="clusterNamespaces"
      :loading="loading"
      :columns="columns"
      :show-pagination="true"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @refresh="fetchNamespaces"
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
  { prop: 'metadata.name', label: '命名空间名称' },
  { prop: 'status.phase', label: '状态' },
  { prop: 'metadata.creationTimestamp', label: '创建时间' }
]

const clusterNamespaces = computed(() => clusterStore.clusterNamespaces)

const fetchNamespaces = async () => {
  await clusterStore.fetchClusterNamespaces(props.clusterId)
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchNamespaces()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchNamespaces()
}

onMounted(() => {
  fetchNamespaces()
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