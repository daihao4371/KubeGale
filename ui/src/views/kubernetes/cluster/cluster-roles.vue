<template>
  <div class="cluster-roles-page">
    <el-card>
      <template #header>
        <div>
          <span>角色管理 (集群: {{ clusterId }})</span>
           <el-button style="float: right; padding: 3px 0" type="primary" link @click="goBack">返回集群详情</el-button>
        </div>
      </template>

      <div class="dycloud-table-operator-area">
        <el-button @click="handleAddRole" type="primary" icon="plus">添加角色</el-button>
      </div>

      <el-table :data="roles" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="角色名称"></el-table-column>
        <el-table-column prop="description" label="描述"></el-table-column>
        <el-table-column prop="createdAt" label="创建时间">
           <template #default="scope">
            <span>{{ formatDateTime(scope.row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button @click="handleEditRole(scope.row)" type="primary" link icon="edit">编辑</el-button>
            <el-button @click="handleRemoveRole(scope.row)" type="danger" link icon="delete">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

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

      <!-- Dialog for Add/Edit Role -->
      <el-dialog :title="dialogTitle" v-model="dialogFormVisible" width="50%">
        <el-form :model="roleForm" ref="roleFormRef" label-width="100px">
          <el-form-item label="角色名称" prop="name" :rules="[{ required: true, message: '请输入角色名称', trigger: 'blur' }]">
            <el-input v-model="roleForm.name" placeholder="输入角色名称"></el-input>
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input type="textarea" v-model="roleForm.description" placeholder="输入角色描述"></el-input>
          </el-form-item>
          <!-- Further fields for role permissions (rules) will be complex and added later -->
          <el-form-item label="权限规则">
            <p>权限规则配置界面将在此处实现。</p>
            <!-- This would involve a more complex UI for defining API groups, resources, verbs etc. -->
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitRoleForm">提交</el-button>
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
// TODO: Import actual API functions for cluster role management when available
// import { getClusterRoles, addClusterRole, updateClusterRole, removeClusterRole } from '@/api/kubernetes/cluster/k8sCluster';

// Placeholder types - these should be defined in your types/kubernetes.ts or API files
interface ClusterRole {
  id: string;
  name: string;
  description?: string;
  rules?: any[]; // This will be a complex structure based on Kubernetes RBAC rules
  createdAt?: string;
}

const route = useRoute();
const router = useRouter();

const clusterId = ref<string | null>(null);
const roles = ref<ClusterRole[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const dialogFormVisible = ref(false);
const dialogTitle = ref('');
const roleFormRef = ref<FormInstance>();
const roleForm = reactive<ClusterRole>({
  id: '',
  name: '',
  description: '',
  rules: [],
});
const isEditMode = ref(false);

const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '';
  return new Date(dateTime).toLocaleString();
};

const fetchRoles = async () => {
  if (!clusterId.value) return;
  loading.value = true;
  // Replace with actual API call
  console.log(`Fetching roles for cluster ${clusterId.value}, page: ${page.value}, pageSize: ${pageSize.value}`);
  // Simulating API call
  await new Promise(resolve => setTimeout(resolve, 500));
  const mockRoles: ClusterRole[] = [
    { id: 'role1', name: 'cluster-admin', description: 'Full control over the cluster', createdAt: new Date().toISOString() },
    { id: 'role2', name: 'view-only', description: 'Can view resources', createdAt: new Date().toISOString() },
  ];
  roles.value = mockRoles;
  total.value = mockRoles.length;
  loading.value = false;
  // try {
  //   const response = await getClusterRoles({ clusterId: clusterId.value, page: page.value, pageSize: pageSize.value });
  //   if (response.code === 0) {
  //     roles.value = response.data.list;
  //     total.value = response.data.total;
  //   } else {
  //     ElMessage.error(response.message || '获取角色列表失败');
  //   }
  // } catch (error) {
  //   ElMessage.error('获取角色列表失败');
  // } finally {
  //   loading.value = false;
  // }
};

const handleAddRole = () => {
  isEditMode.value = false;
  dialogTitle.value = '添加角色';
  roleForm.id = '';
  roleForm.name = '';
  roleForm.description = '';
  roleForm.rules = [];
  if(roleFormRef.value) roleFormRef.value.resetFields();
  dialogFormVisible.value = true;
};

const handleEditRole = (role: ClusterRole) => {
  isEditMode.value = true;
  dialogTitle.value = '编辑角色';
  // Deep copy for editing complex objects like rules
  Object.assign(roleForm, JSON.parse(JSON.stringify(role)));
  dialogFormVisible.value = true;
};

const submitRoleForm = async () => {
  if (!roleFormRef.value) return;
  await roleFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;
      // Replace with actual API call
      console.log('Submitting role form:', roleForm);
      // Simulating API call
      await new Promise(resolve => setTimeout(resolve, 500));
      ElMessage.success(isEditMode.value ? '角色更新成功' : '角色添加成功');
      dialogFormVisible.value = false;
      fetchRoles(); // Refresh
      // try {
      //   let response;
      //   if (isEditMode.value) {
      //     response = await updateClusterRole(clusterId.value, roleForm);
      //   } else {
      //     response = await addClusterRole(clusterId.value, { name: roleForm.name, description: roleForm.description, rules: roleForm.rules });
      //   }
      //   if (response.code === 0) {
      //     ElMessage.success(isEditMode.value ? '角色更新成功' : '角色添加成功');
      //     dialogFormVisible.value = false;
      //     fetchRoles();
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

const handleRemoveRole = (role: ClusterRole) => {
  ElMessageBox.confirm(`确定要删除角色 "${role.name}" 吗?`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    loading.value = true;
    // Replace with actual API call
    console.log('Removing role:', role);
    // Simulating API call
    await new Promise(resolve => setTimeout(resolve, 500));
    ElMessage.success('角色删除成功');
    fetchRoles(); // Refresh
    // try {
    //   const response = await removeClusterRole(clusterId.value, role.id);
    //   if (response.code === 0) {
    //     ElMessage.success('角色删除成功');
    //     fetchRoles();
    //   } else {
    //     ElMessage.error(response.message || '删除失败');
    //   }
    // } catch (error) {
    //   ElMessage.error('删除失败');
    // } finally {
    //   loading.value = false;
    // }
  }).catch(() => {
    ElMessage.info('已取消删除');
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
  fetchRoles();
};

const handleCurrentChange = (val: number) => {
  page.value = val;
  fetchRoles();
};

onMounted(() => {
  clusterId.value = route.params.clusterId as string;
  if (clusterId.value) {
    fetchRoles();
  } else {
    ElMessage.error('未提供集群ID');
    // Optionally redirect
    // router.push({ name: 'ClusterList' });
  }
});
</script>

<style scoped>
.cluster-roles-page {
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
