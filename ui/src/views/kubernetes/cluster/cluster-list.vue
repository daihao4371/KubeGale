<template>
  <div>
    <!-- Search bar -->
    <div class="dycloud-search-box">
      <el-form :inline="true" :model="searchInfo" ref="searchForm">
        <el-form-item label="名称">
          <el-input v-model="searchInfo.name" placeholder="输入集群名称查询"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button @click="onSubmit" type="primary" icon="search">查询</el-button>
          <el-button @click="onReset" icon="refresh">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- Action buttons -->
    <div class="dycloud-table-operator-area">
      <el-button @click="handleCreate" type="primary" icon="plus">新建集群</el-button>
      <!-- Add other batch actions if needed -->
    </div>

    <!-- Table to display clusters -->
    <div class="dycloud-table-box">
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="集群名称"></el-table-column>
        <el-table-column prop="id" label="集群ID"></el-table-column>
        <el-table-column prop="api_address" label="API地址"></el-table-column>
        <el-table-column prop="kube_type" label="类型">
          <template #default="scope">
            <span>{{ formatKubeType(scope.row.kube_type) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间">
          <template #default="scope">
            <span>{{ formatDateTime(scope.row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="scope">
            <el-button @click="handleDetail(scope.row)" type="primary" link icon="view">详情</el-button>
            <el-button @click="handleEdit(scope.row)" type="primary" link icon="edit">编辑</el-button>
            <el-button @click="handleDelete(scope.row)" type="danger" link icon="delete">删除</el-button>
            <el-button @click="handleManageUsers(scope.row)" type="primary" link>用户管理</el-button>
            <el-button @click="handleManageRoles(scope.row)" type="primary" link>角色管理</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Pagination -->
    <div class="dycloud-pagination">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="page"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total">
      </el-pagination>
    </div>

    <!-- Dialog for Create/Edit Cluster -->
    <el-dialog :title="dialogTitle" v-model="dialogFormVisible" width="50%">
      <ClusterForm v-if="dialogFormVisible" :cluster="selectedCluster" @submit="onFormSubmit" @cancel="dialogFormVisible = false" />
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getClustersList, DeleteCluster } from '@/api/kubernetes/cluster/k8sCluster';
// Assuming ClusterForm will be created later
import ClusterForm from './cluster-form.vue';
import { useRouter } from 'vue-router';
import type { ClusterData } from '@/api/kubernetes/cluster/k8sCluster'; // Assuming types will be more detailed later

const router = useRouter();

const searchInfo = reactive({
  name: '',
});

const tableData = ref<ClusterData[]>([]);
const loading = ref(true);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const dialogFormVisible = ref(false);
const dialogTitle = ref('');
const selectedCluster = ref<ClusterData | null>(null);

const formatKubeType = (type: number) => {
  // This should be mapped to actual meanings if available
  return type === 1 ? 'KubeConfig' : type === 2 ? 'Agent' : '未知';
};

const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '';
  return new Date(dateTime).toLocaleString();
};

const fetchClusters = async () => {
  loading.value = true;
  try {
    const response = await getClustersList({
      page: page.value,
      pageSize: pageSize.value,
      name: searchInfo.name || undefined,
    });
    if (response.code === 0) {
      tableData.value = response.data.list;
      total.value = response.data.total;
    } else {
      ElMessage.error(response.message || '获取集群列表失败');
    }
  } catch (error) {
    console.error('Error fetching clusters:', error);
    ElMessage.error('获取集群列表失败');
  } finally {
    loading.value = false;
  }
};

const onSubmit = () => {
  page.value = 1; // Reset to first page for new search
  fetchClusters();
};

const onReset = () => {
  searchInfo.name = '';
  page.value = 1;
  fetchClusters();
};

const handleCreate = () => {
  selectedCluster.value = null; // Important for create mode
  dialogTitle.value = '新建集群';
  dialogFormVisible.value = true;
};

const handleEdit = (row: ClusterData) => {
  selectedCluster.value = JSON.parse(JSON.stringify(row)); // Deep copy for editing
  dialogTitle.value = '编辑集群';
  dialogFormVisible.value = true;
};

const handleDetail = (row: ClusterData) => {
  // Navigate to cluster detail page
  // This will be implemented once cluster-detail.vue and routing are set up
  router.push({ name: 'ClusterDetail', params: { id: row.id } });
  console.log('Navigate to detail for:', row);
};

const handleDelete = (row: ClusterData) => {
  ElMessageBox.confirm(`确定要删除集群 "${row.name}"吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    if (!row.id) {
      ElMessage.error('集群ID不存在');
      return;
    }
    try {
      const response = await DeleteCluster({ id: row.id });
      if (response.code === 0) {
        ElMessage.success('删除成功');
        fetchClusters(); // Refresh list
      } else {
        ElMessage.error(response.message || '删除失败');
      }
    } catch (error) {
      console.error('Error deleting cluster:', error);
      ElMessage.error('删除失败');
    }
  }).catch(() => {
    ElMessage.info('已取消删除');
  });
};

const handleManageUsers = (row: ClusterData) => {
  // Navigate to user management page for this cluster
  // This will be implemented once cluster-users.vue and routing are set up
  router.push({ name: 'ClusterUsers', params: { clusterId: row.id } });
  console.log('Manage users for:', row);
};

const handleManageRoles = (row: ClusterData) => {
  // Navigate to role management page for this cluster
  // This will be implemented once cluster-roles.vue and routing are set up
  router.push({ name: 'ClusterRoles', params: { clusterId: row.id } });
  console.log('Manage roles for:', row);
};

const onFormSubmit = () => {
  dialogFormVisible.value = false;
  fetchClusters(); // Refresh list after create/edit
};

const handleSizeChange = (val: number) => {
  pageSize.value = val;
  fetchClusters();
};

const handleCurrentChange = (val: number) => {
  page.value = val;
  fetchClusters();
};

onMounted(() => {
  fetchClusters();
});

</script>

<style scoped>
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
