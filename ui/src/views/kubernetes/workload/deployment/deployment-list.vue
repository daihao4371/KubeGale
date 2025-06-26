<template>
  <div class="deployment-list-page">
    <!-- Search and Filter Bar -->
    <div class="dycloud-search-box">
      <el-form :inline="true" :model="searchQuery" ref="searchFormRef">
        <el-form-item label="集群">
          <el-select v-model="searchQuery.clusterId" placeholder="选择集群" @change="onClusterChange" clearable>
            <!-- Populate with actual cluster list -->
            <el-option v-for="cluster in clusterList" :key="cluster.id" :label="cluster.name" :value="cluster.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="命名空间">
          <el-select v-model="searchQuery.namespace" placeholder="选择命名空间" @change="onNamespaceChange" clearable :disabled="!searchQuery.clusterId">
            <!-- Populate with namespaces for selected cluster -->
            <el-option v-for="ns in namespaceList" :key="ns.name" :label="ns.name" :value="ns.name"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="searchQuery.name" placeholder="输入Deployment名称"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleSearch" type="primary" icon="search" :disabled="!searchQuery.clusterId || !searchQuery.namespace">查询</el-button>
          <el-button @click="handleResetSearch" icon="refresh">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- Action Buttons -->
    <div class="dycloud-table-operator-area">
      <el-button @click="handleCreateDeployment" type="primary" icon="plus" :disabled="!searchQuery.clusterId || !searchQuery.namespace">创建Deployment</el-button>
    </div>

    <!-- Deployments Table -->
    <div class="dycloud-table-box">
      <el-table :data="deployments" v-loading="loading" style="width: 100%">
        <el-table-column prop="metadata.name" label="名称" min-width="150"></el-table-column>
        <el-table-column label="命名空间" prop="metadata.namespace" min-width="120"></el-table-column>
        <el-table-column label="副本数(Ready/Target)" min-width="150">
          <template #default="scope">
            <span>{{ scope.row.status?.readyReplicas || 0 }}/{{ scope.row.spec?.replicas || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="镜像" min-width="200">
          <template #default="scope">
            <!-- Display first container's image as an example -->
            <div v-if="scope.row.spec?.template?.spec?.containers?.length > 0">
              <el-tooltip :content="scope.row.spec.template.spec.containers[0].image" placement="top">
                 <span>{{ truncateText(scope.row.spec.template.spec.containers[0].image, 30) }}</span>
              </el-tooltip>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="160">
          <template #default="scope">
            <span>{{ formatDateTime(scope.row.metadata?.creationTimestamp) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="scope">
            <el-button @click="handleViewDetails(scope.row)" type="primary" link icon="view">详情</el-button>
            <el-button @click="handleEditDeployment(scope.row)" type="primary" link icon="edit">编辑</el-button>
            <el-button @click="handleScaleDeployment(scope.row)" type="primary" link>扩缩容</el-button>
            <el-button @click="handleRollbackDeployment(scope.row)" type="primary" link>回滚</el-button>
            <el-button @click="handleDeleteDeployment(scope.row)" type="danger" link icon="delete">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Pagination -->
    <div class="dycloud-pagination">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="pagination.page"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pagination.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total">
      </el-pagination>
    </div>

    <!-- Create/Edit Dialog (Placeholder - will use a separate form component) -->
    <el-dialog :title="dialogTitle" v-model="dialogFormVisible" width="70%" :close-on-click-modal="false">
      <DeploymentForm
        v-if="dialogFormVisible"
        :cluster-id="searchQuery.clusterId"
        :namespace-id="searchQuery.namespace"
        :deployment="selectedDeployment"
        @submit="onFormSubmitted"
        @cancel="dialogFormVisible = false" />
    </el-dialog>

     <!-- Scale Dialog (Placeholder) -->
    <el-dialog title="调整副本数" v-model="scaleDialogVisible" width="30%">
      <el-form :model="scaleForm" label-width="80px" v-if="selectedDeployment">
        <p>Deployment: {{ selectedDeployment.metadata?.name }}</p>
        <el-form-item label="副本数">
          <el-input-number v-model="scaleForm.replicas" :min="0" :max="100"></el-input-number>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scaleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitScaleDeployment">确定</el-button>
      </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRouter } from 'vue-router';
import { getClustersList, getAllClusterNamespaces, K8sCluster } from '@/api/kubernetes/cluster/k8sCluster'; // Assuming Namespace type might be here or in its own file
// TODO: Create and import Deployment API functions
// import { getDeployments, deleteDeployment, scaleDeployment } from '@/api/kubernetes/workload/deployment';
import DeploymentForm from './deployment-form.vue'; // To be created
// TODO: Define Deployment type (likely from Kubernetes client types or custom)
// For now, using 'any' as a placeholder for Kubernetes objects
type K8sDeployment = any;
type K8sNamespace = { name: string; [key: string]: any };


const router = useRouter();

const searchQuery = reactive({
  clusterId: '',
  namespace: '',
  name: '',
});

const clusterList = ref<K8sCluster[]>([]);
const namespaceList = ref<K8sNamespace[]>([]); // { name: string } for now

const deployments = ref<K8sDeployment[]>([]);
const loading = ref(false);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});

const dialogFormVisible = ref(false);
const dialogTitle = ref('');
const selectedDeployment = ref<K8sDeployment | null>(null);

const scaleDialogVisible = ref(false);
const scaleForm = reactive({ replicas: 0 });

// --- Utility Functions ---
const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '';
  return new Date(dateTime).toLocaleString();
};

const truncateText = (text: string, maxLength: number) => {
  if (!text) return '';
  return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
};


// --- Fetching Data ---
const fetchClusterList = async () => {
  try {
    const res = await getClustersList({ page: 1, pageSize: 1000 }); // Fetch all clusters for dropdown
    if (res.code === 0) {
      clusterList.value = res.data.list;
    } else {
      ElMessage.error(res.message || '获取集群列表失败');
    }
  } catch (error) {
    ElMessage.error('获取集群列表接口错误');
  }
};

const fetchNamespacesForCluster = async (clusterId: string) => {
  if (!clusterId) {
    namespaceList.value = [];
    searchQuery.namespace = ''; // Reset namespace if cluster changes
    return;
  }
  try {
    // Assuming getAllClusterNamespaces returns a list of namespace objects
    const res = await getAllClusterNamespaces(clusterId);
    if (res.code === 0) {
      namespaceList.value = res.data.items || []; // Ensure items is an array
       if (namespaceList.value.length > 0 && !searchQuery.namespace) {
         // searchQuery.namespace = 'default'; // Or the first one, or leave empty
       }
    } else {
      ElMessage.error(res.message || '获取命名空间列表失败');
      namespaceList.value = [];
    }
  } catch (error) {
    ElMessage.error('获取命名空间接口错误');
    namespaceList.value = [];
  }
};

const fetchDeployments = async () => {
  if (!searchQuery.clusterId || !searchQuery.namespace) {
    // deployments.value = []; // Clear deployments if cluster/ns not selected
    // pagination.total = 0;
    // ElMessage.info('请先选择集群和命名空间');
    return;
  }
  loading.value = true;
  try {
    // const response = await getDeployments({
    //   clusterId: searchQuery.clusterId,
    //   namespace: searchQuery.namespace,
    //   name: searchQuery.name,
    //   page: pagination.page,
    //   pageSize: pagination.pageSize,
    // });
    // MOCK DATA for now
    await new Promise(resolve => setTimeout(resolve, 500));
    const mockData = {
      code: 0,
      data: {
        items: [
          { metadata: { name: 'nginx-deployment', namespace: 'default', creationTimestamp: new Date().toISOString() }, spec: { replicas: 2, template: { spec: { containers: [{ name: 'nginx', image: 'nginx:latest' }]}}}, status: { readyReplicas: 2 } },
          { metadata: { name: 'myapp-deployment', namespace: 'default', creationTimestamp: new Date().toISOString() }, spec: { replicas: 3, template: { spec: { containers: [{ name: 'myapp', image: 'myregistry/myapp:1.2.3-very-long-image-tag-that-needs-truncation' }]}}}, status: { readyReplicas: 1 } },
        ],
        total: 2,
      }
    };
    const response = mockData;
    // END MOCK DATA

    if (response.code === 0) {
      deployments.value = response.data.items;
      pagination.total = response.data.total;
    } else {
      // ElMessage.error(response.message || '获取Deployments失败');
      deployments.value = [];
      pagination.total = 0;
    }
  } catch (error) {
    console.error('Error fetching deployments:', error);
    ElMessage.error('获取Deployments接口错误');
    deployments.value = [];
    pagination.total = 0;
  } finally {
    loading.value = false;
  }
};

// --- Event Handlers ---
watch(() => searchQuery.clusterId, (newClusterId) => {
  searchQuery.namespace = ''; // Reset namespace when cluster changes
  namespaceList.value = [];
  deployments.value = [];
  pagination.total = 0;
  if (newClusterId) {
    fetchNamespacesForCluster(newClusterId);
  }
});


const onClusterChange = () => {
    // Namespace will be fetched by the watcher.
    // We can optionally trigger an immediate fetch or wait for user to select namespace.
    // For now, let user select namespace explicitly.
    deployments.value = [];
    pagination.total = 0;
};

const onNamespaceChange = () => {
    // Fetch deployments if both cluster and namespace are selected
    if (searchQuery.clusterId && searchQuery.namespace) {
        pagination.page = 1;
        fetchDeployments();
    } else {
        deployments.value = [];
        pagination.total = 0;
    }
};


const handleSearch = () => {
  pagination.page = 1;
  fetchDeployments();
};

const handleResetSearch = () => {
  searchQuery.clusterId = '';
  searchQuery.namespace = '';
  searchQuery.name = '';
  namespaceList.value = [];
  deployments.value = [];
  pagination.page = 1;
  pagination.total = 0;
  // fetchDeployments(); // Or clear table
};

const handleCreateDeployment = () => {
  if (!searchQuery.clusterId || !searchQuery.namespace) {
    ElMessage.warning('请先选择集群和命名空间');
    return;
  }
  selectedDeployment.value = null;
  dialogTitle.value = '创建Deployment';
  dialogFormVisible.value = true;
};

const handleEditDeployment = (deployment: K8sDeployment) => {
  selectedDeployment.value = JSON.parse(JSON.stringify(deployment)); // Deep copy
  dialogTitle.value = `编辑Deployment: ${deployment.metadata.name}`;
  dialogFormVisible.value = true;
};

const handleViewDetails = (deployment: K8sDeployment) => {
  // router.push({ name: 'DeploymentDetail', params: { clusterId: searchQuery.clusterId, namespace: deployment.metadata.namespace, name: deployment.metadata.name } });
  ElMessage.info(`详情: ${deployment.metadata.name} (to be implemented)`);
};


const handleScaleDeployment = (deployment: K8sDeployment) => {
  selectedDeployment.value = deployment;
  scaleForm.replicas = deployment.spec?.replicas || 0;
  scaleDialogVisible.value = true;
};

const submitScaleDeployment = async () => {
  if (!selectedDeployment.value) return;
  // const { metadata, spec } = selectedDeployment.value;
  // try {
  //   await scaleDeployment(searchQuery.clusterId, metadata.namespace, metadata.name, scaleForm.replicas);
  //   ElMessage.success('副本数调整成功');
  //   fetchDeployments();
  //   scaleDialogVisible.value = false;
  // } catch (error) {
  //   ElMessage.error('调整副本数失败');
  // }
  ElMessage.info(`模拟调整副本数: ${selectedDeployment.value.metadata.name} to ${scaleForm.replicas}`);
  // Update mock data
  const index = deployments.value.findIndex(d => d.metadata.name === selectedDeployment.value?.metadata.name);
  if (index !== -1) {
    deployments.value[index].spec.replicas = scaleForm.replicas;
    // simulate some ready replicas
    deployments.value[index].status.readyReplicas = Math.min(scaleForm.replicas, deployments.value[index].status.readyReplicas || 0);
  }
  scaleDialogVisible.value = false;
};

const handleRollbackDeployment = (deployment: K8sDeployment) => {
  ElMessage.info(`回滚: ${deployment.metadata.name} (to be implemented)`);
  // This will likely involve fetching revision history and allowing user to select one.
};

const handleDeleteDeployment = (deployment: K8sDeployment) => {
  ElMessageBox.confirm(`确定要删除Deployment "${deployment.metadata.name}" 吗?`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    // try {
    //   await deleteDeployment(searchQuery.clusterId, deployment.metadata.namespace, deployment.metadata.name);
    //   ElMessage.success('删除成功');
    //   fetchDeployments();
    // } catch (error) {
    //   ElMessage.error('删除失败');
    // }
    ElMessage.info(`模拟删除: ${deployment.metadata.name}`);
    deployments.value = deployments.value.filter(d => d.metadata.name !== deployment.metadata.name);
    pagination.total = deployments.value.length;

  }).catch(() => ElMessage.info('已取消删除'));
};

const onFormSubmitted = () => {
  dialogFormVisible.value = false;
  fetchDeployments();
};

const handleSizeChange = (val: number) => {
  pagination.pageSize = val;
  fetchDeployments();
};

const handleCurrentChange = (val: number) => {
  pagination.page = val;
  fetchDeployments();
};

// --- Lifecycle Hooks ---
onMounted(() => {
  fetchClusterList();
  // Initial fetch if cluster/namespace are pre-selected (e.g., from query params or local storage)
  if (searchQuery.clusterId) { // Could be loaded from store/localStorage
      fetchNamespacesForCluster(searchQuery.clusterId).then(() => {
          if (searchQuery.namespace) { // Could be loaded from store/localStorage
              fetchDeployments();
          }
      });
  }
});

</script>

<style scoped>
.deployment-list-page {
  padding: 20px;
}
.dycloud-search-box {
  margin-bottom: 16px;
  padding: 16px;
  background-color: #f7f7f7;
  border-radius: 4px;
}
.dycloud-table-operator-area {
  margin-bottom: 16px;
}
.dycloud-pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
