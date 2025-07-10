<template>
  <div class="cluster-detail">
    <el-card class="detail-card">
      <template #header>
        <div class="card-header">
          <span class="title">集群详情</span>
          <div class="actions">
            <el-button @click="goBack" icon="ArrowLeft">返回</el-button>
            <el-button type="primary" @click="refreshData" :loading="loading" icon="Refresh">刷新</el-button>
          </div>
        </div>
      </template>

      <div v-if="currentCluster" class="detail-content">
        <!-- 集群基本信息 -->
        <div class="basic-info">
          <h3>基本信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="集群名称">{{ currentCluster.name }}</el-descriptions-item>
            <el-descriptions-item label="集群别名">{{ currentCluster.alias }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <K8sStatusBadge :status="currentCluster.status" />
            </el-descriptions-item>
            <el-descriptions-item label="版本">{{ currentCluster.version }}</el-descriptions-item>
            <el-descriptions-item label="所在城市">{{ currentCluster.city }}</el-descriptions-item>
            <el-descriptions-item label="所在区域">{{ currentCluster.district }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(currentCluster.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(currentCluster.updated_at) }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ currentCluster.description || '暂无描述' }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 标签页 -->
        <el-tabs v-model="activeTab" class="detail-tabs" @tab-click="handleTabClick">
          <el-tab-pane label="用户管理" name="users">
            <ClusterUsersTab :cluster-id="clusterId" />
          </el-tab-pane>
          <el-tab-pane label="角色管理" name="roles">
            <ClusterRolesTab :cluster-id="clusterId" />
          </el-tab-pane>
          <el-tab-pane label="命名空间" name="namespaces">
            <ClusterNamespacesTab :cluster-id="clusterId" />
          </el-tab-pane>
          <el-tab-pane label="API资源" name="apiGroups">
            <ClusterApiGroupsTab :cluster-id="clusterId" />
          </el-tab-pane>
        </el-tabs>
      </div>

      <div v-else-if="loading" class="loading">
        <el-skeleton :rows="10" animated />
      </div>

      <div v-else class="empty">
        <el-empty description="集群信息不存在" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useClusterStore } from '@/stores/kubernetes/cluster'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import ClusterUsersTab from './components/ClusterUsersTab.vue'
import ClusterRolesTab from './components/ClusterRolesTab.vue'
import ClusterNamespacesTab from './components/ClusterNamespacesTab.vue'
import ClusterApiGroupsTab from './components/ClusterApiGroupsTab.vue'
import { formatDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const clusterStore = useClusterStore()

// 响应式数据
const activeTab = ref('users')
const loading = ref(false)

// 计算属性
const clusterId = computed(() => Number(route.params.id))
const currentCluster = computed(() => clusterStore.currentCluster)

// 方法
const fetchClusterDetail = async () => {
  if (!clusterId.value) return
  
  try {
    loading.value = true
    await clusterStore.fetchClusterDetail(clusterId.value)
  } catch (error) {
    console.error('获取集群详情失败:', error)
  } finally {
    loading.value = false
  }
}

const refreshData = async () => {
  await fetchClusterDetail()
}

const handleTabClick = (tab: any) => {
  console.log('切换到标签页:', tab.name)
}

const goBack = () => {
  router.back()
}

// 生命周期
onMounted(() => {
  fetchClusterDetail()
})
</script>

<style scoped>
.cluster-detail {
  padding: 20px;
}

.detail-card {
  min-height: calc(100vh - 120px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 18px;
  font-weight: 600;
}

.actions {
  display: flex;
  gap: 10px;
}

.detail-content {
  margin-top: 20px;
}

.basic-info {
  margin-bottom: 30px;
}

.basic-info h3 {
  margin-bottom: 15px;
  color: #303133;
  font-weight: 600;
}

.detail-tabs {
  margin-top: 20px;
}

.loading {
  padding: 40px 0;
}

.empty {
  padding: 60px 0;
  text-align: center;
}
</style> 