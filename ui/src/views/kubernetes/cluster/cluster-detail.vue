<template>
  <div class="cluster-detail-page" v-if="cluster">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>集群详情: {{ cluster.name }}</span>
          <el-button class="button" type="primary" @click="goBack">返回列表</el-button>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="集群ID">{{ cluster.id }}</el-descriptions-item>
        <el-descriptions-item label="集群名称">{{ cluster.name }}</el-descriptions-item>
        <el-descriptions-item label="API Server地址">{{ cluster.api_address }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(cluster.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="集群类型">
          <el-tag :type="cluster.kube_type === 1 ? 'success' : 'info'">
            {{ formatKubeType(cluster.kube_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="Prometheus URL">{{ cluster.prometheus_url || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Prometheus认证">
            {{ formatPrometheusAuthType(cluster.prometheus_auth_type) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="cluster.prometheus_auth_type === 1" label="Prometheus 用户">
            {{ cluster.prometheus_user || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-tabs v-model="activeTab" class="detail-tabs">
        <el-tab-pane label="KubeConfig" name="kubeconfig" v-if="cluster.kube_type === 1">
          <pre class="kubeconfig-display">{{ cluster.kube_config }}</pre>
        </el-tab-pane>
        <el-tab-pane label="用户管理" name="users">
          <p>用户管理功能将在此处实现 (cluster-users.vue integration point)</p>
          <!-- <ClusterUsers :clusterId="cluster.id" v-if="cluster.id" /> -->
        </el-tab-pane>
        <el-tab-pane label="角色管理" name="roles">
          <p>角色管理功能将在此处实现 (cluster-roles.vue integration point)</p>
          <!-- <ClusterRoles :clusterId="cluster.id" v-if="cluster.id" /> -->
        </el-tab-pane>
        <!-- Add more tabs as needed, e.g., API Groups, Credentials -->
      </el-tabs>

    </el-card>
  </div>
  <div v-else-if="loading">
    <p>加载中...</p>
  </div>
  <div v-else>
    <p>未找到集群信息。</p>
    <el-button type="primary" @click="goBack">返回列表</el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getClustersById } from '@/api/kubernetes/cluster/k8sCluster';
import type { ClusterData } from '@/api/kubernetes/cluster/k8sCluster';
import { ElMessage } from 'element-plus';
// import ClusterUsers from './cluster-users.vue'; // Placeholder for later
// import ClusterRoles from './cluster-roles.vue'; // Placeholder for later

const route = useRoute();
const router = useRouter();

const cluster = ref<ClusterData | null>(null);
const loading = ref(true);
const activeTab = ref('kubeconfig'); // Default tab

const formatKubeType = (type: number | undefined) => {
  if (type === undefined) return '未知';
  return type === 1 ? 'KubeConfig' : type === 2 ? 'Agent' : '未知';
};

const formatPrometheusAuthType = (type: number | undefined) => {
    if (type === undefined) return '未知';
    switch (type) {
        case 0: return 'None';
        case 1: return 'Basic Auth';
        default: return '未知';
    }
};

const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '';
  return new Date(dateTime).toLocaleString();
};

const fetchClusterDetails = async (id: string) => {
  loading.value = true;
  try {
    const response = await getClustersById({ id });
    if (response.code === 0 && response.data && response.data.cluster) {
      cluster.value = response.data.cluster;
      if(cluster.value.kube_type !== 1) { // If not KubeConfig, switch tab
        activeTab.value = 'users';
      }
    } else {
      ElMessage.error(response.message || '获取集群详情失败');
      cluster.value = null;
    }
  } catch (error) {
    console.error('Error fetching cluster details:', error);
    ElMessage.error('获取集群详情失败');
    cluster.value = null;
  } finally {
    loading.value = false;
  }
};

const goBack = () => {
  router.push({ name: 'ClusterList' }); // Assuming 'ClusterList' is the route name for the list page
};

onMounted(() => {
  const clusterId = route.params.id as string;
  if (clusterId) {
    fetchClusterDetails(clusterId);
  } else {
    ElMessage.error('未提供集群ID');
    loading.value = false;
    router.push({ name: 'ClusterList' }); // Redirect if no ID
  }
});
</script>

<style scoped>
.cluster-detail-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.kubeconfig-display {
  background-color: #f5f5f5;
  padding: 15px;
  border-radius: 4px;
  white-space: pre-wrap; /* Handles line breaks */
  word-break: break-all; /* Handles long lines without spaces */
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #e0e0e0;
}
.detail-tabs {
  margin-top: 20px;
}
</style>
