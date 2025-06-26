<template>
  <div class="cluster-users-page">
    <el-card>
      <template #header>
        <div>
          <span>用户管理 (集群: {{ clusterId }})</span>
          <el-button style="float: right; padding: 3px 0" type="primary" link @click="goBack">返回集群详情</el-button>
        </div>
      </template>

      <div class="dycloud-table-operator-area">
        <el-button @click="handleAddUser" type="primary" icon="plus">添加用户</el-button>
      </div>

      <el-table :data="users" v-loading="loading" style="width: 100%">
        <el-table-column prop="username" label="用户名"></el-table-column>
        <el-table-column prop="role" label="角色"></el-table-column>
        <el-table-column prop="joinedAt" label="加入时间">
           <template #default="scope">
            <span>{{ formatDateTime(scope.row.joinedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button @click="handleEditUser(scope.row)" type="primary" link icon="edit">编辑</el-button>
            <el-button @click="handleRemoveUser(scope.row)" type="danger" link icon="delete">移除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination can be added here if user list is long -->
       <div class="dycloud-pagination" v-if="total > pageSize">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="page"
          :page-sizes="[10, 20, 50]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total">
        </el-pagination>
      </div>

      <!-- Dialog for Add/Edit User -->
      <el-dialog :title="dialogTitle" v-model="dialogFormVisible" width="40%">
        <el-form :model="userForm" ref="userFormRef" label-width="100px">
          <el-form-item label="用户名" prop="username" :rules="[{ required: true, message: '请输入用户名', trigger: 'blur' }]">
            <el-input v-model="userForm.username" placeholder="输入用户名"></el-input>
          </el-form-item>
          <el-form-item label="角色" prop="role" :rules="[{ required: true, message: '请选择角色', trigger: 'change' }]">
            <el-select v-model="userForm.role" placeholder="选择角色" style="width:100%">
              <!-- Replace with actual roles fetched from API -->
              <el-option label="Admin" value="admin"></el-option>
              <el-option label="Developer" value="developer"></el-option>
              <el-option label="Viewer" value="viewer"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitUserForm">提交</el-button>
            <el-button @click="dialogFormVisible = false">取消</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>

    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus';
// TODO: Import actual API functions for cluster user management when available
// import { getClusterUsers, addClusterUser, updateClusterUser, removeClusterUser } from '@/api/kubernetes/cluster/k8sCluster';

// Placeholder types - these should be defined in your types/kubernetes.ts or API files
interface ClusterUser {
  id: string;
  username: string;
  role: string; // This might be a more complex object
  joinedAt?: string;
}

const route = useRoute();
const router = useRouter();

const clusterId = ref<string | null>(null);
const users = ref<ClusterUser[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);


const dialogFormVisible = ref(false);
const dialogTitle = ref('');
const userFormRef = ref<FormInstance>();
const userForm = reactive({
  id: '',
  username: '',
  role: '',
});
const isEditMode = ref(false);

const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '';
  return new Date(dateTime).toLocaleString();
};

const fetchUsers = async () => {
  if (!clusterId.value) return;
  loading.value = true;
  // Replace with actual API call
  console.log(`Fetching users for cluster ${clusterId.value}, page: ${page.value}, pageSize: ${pageSize.value}`);
  // Simulating API call
  await new Promise(resolve => setTimeout(resolve, 500));
  const mockUsers: ClusterUser[] = [
    { id: 'user1', username: 'alice', role: 'Admin', joinedAt: new Date().toISOString() },
    { id: 'user2', username: 'bob', role: 'Developer', joinedAt: new Date().toISOString() },
  ];
  users.value = mockUsers;
  total.value = mockUsers.length;
  loading.value = false;
  // try {
  //   const response = await getClusterUsers({ clusterId: clusterId.value, page: page.value, pageSize: pageSize.value });
  //   if (response.code === 0) {
  //     users.value = response.data.list;
  //     total.value = response.data.total;
  //   } else {
  //     ElMessage.error(response.message || '获取用户列表失败');
  //   }
  // } catch (error) {
  //   ElMessage.error('获取用户列表失败');
  // } finally {
  //   loading.value = false;
  // }
};

const handleAddUser = () => {
  isEditMode.value = false;
  dialogTitle.value = '添加用户';
  userForm.id = '';
  userForm.username = '';
  userForm.role = '';
  if(userFormRef.value) userFormRef.value.resetFields();
  dialogFormVisible.value = true;
};

const handleEditUser = (user: ClusterUser) => {
  isEditMode.value = true;
  dialogTitle.value = '编辑用户';
  userForm.id = user.id;
  userForm.username = user.username;
  userForm.role = user.role;
  dialogFormVisible.value = true;
};

const submitUserForm = async () => {
  if (!userFormRef.value) return;
  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;
      // Replace with actual API call
      console.log('Submitting user form:', userForm);
      // Simulating API call
      await new Promise(resolve => setTimeout(resolve, 500));
      ElMessage.success(isEditMode.value ? '用户更新成功' : '用户添加成功');
      dialogFormVisible.value = false;
      fetchUsers(); // Refresh
      // try {
      //   let response;
      //   if (isEditMode.value) {
      //     response = await updateClusterUser(clusterId.value, userForm);
      //   } else {
      //     response = await addClusterUser(clusterId.value, { username: userForm.username, role: userForm.role });
      //   }
      //   if (response.code === 0) {
      //     ElMessage.success(isEditMode.value ? '用户更新成功' : '用户添加成功');
      //     dialogFormVisible.value = false;
      //     fetchUsers();
      //   } else {
      //     ElMessage.error(response.message || '操作失败');
      //   }
      // } catch (error) {
      //   ElMessage.error('操作失败');
      // } finally {
      //   loading.value = false;
      // }
    }
  });
};

const handleRemoveUser = (user: ClusterUser) => {
  ElMessageBox.confirm(`确定要从集群中移除用户 "${user.username}" 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    loading.value = true;
    // Replace with actual API call
    console.log('Removing user:', user);
    // Simulating API call
    await new Promise(resolve => setTimeout(resolve, 500));
    ElMessage.success('用户移除成功');
    fetchUsers(); // Refresh
    // try {
    //   const response = await removeClusterUser(clusterId.value, user.id);
    //   if (response.code === 0) {
    //     ElMessage.success('用户移除成功');
    //     fetchUsers();
    //   } else {
    //     ElMessage.error(response.message || '移除失败');
    //   }
    // } catch (error) {
    //   ElMessage.error('移除失败');
    // } finally {
    //   loading.value = false;
    // }
  }).catch(() => {
    ElMessage.info('已取消移除');
  });
};

const goBack = () => {
    if (clusterId.value) {
        router.push({ name: 'ClusterDetail', params: { id: clusterId.value } });
    } else {
        router.push({ name: 'ClusterList' }); // Fallback
    }
};

const handleSizeChange = (val: number) => {
  pageSize.value = val;
  fetchUsers();
};

const handleCurrentChange = (val: number) => {
  page.value = val;
  fetchUsers();
};

onMounted(() => {
  clusterId.value = route.params.clusterId as string;
  if (clusterId.value) {
    fetchUsers();
  } else {
    ElMessage.error('未提供集群ID');
    // Optionally redirect
    // router.push({ name: 'ClusterList' });
  }
});
</script>

<style scoped>
.cluster-users-page {
  padding: 20px;
}
.dycloud-table-operator-area {
  margin-bottom: 16px;
}
.dycloud-pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
