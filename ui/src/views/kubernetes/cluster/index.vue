<template>
  <div class="cluster-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1>集群管理</h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建集群
        </el-button>
        <el-button @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <el-form :model="searchForm" inline>
        <el-form-item label="集群名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入集群名称"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="认证类型">
          <el-select v-model="searchForm.kube_type" placeholder="请选择认证类型" clearable>
            <el-option label="KubeConfig" :value="1" />
            <el-option label="Token" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 数据表格 -->
    <div class="table-section">
      <K8sTable
        :data="clusterStore.clusters"
        :loading="clusterStore.loading"
        :show-selection="true"
        :show-pagination="false"
        @selection-change="handleSelectionChange"
      >
        <el-table-column prop="name" label="集群名称" min-width="150">
          <template #default="{ row }">
            <el-link type="primary" @click="handleView(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        
        <el-table-column prop="kube_type" label="认证类型" width="120">
          <template #default="{ row }">
            <el-tag :type="row.kube_type === 1 ? 'success' : 'warning'">
              {{ row.kube_type === 1 ? 'KubeConfig' : 'Token' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="api_address" label="API地址" min-width="200" show-overflow-tooltip />
        
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button size="small" type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </K8sTable>
    </div>

    <!-- 创建/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="集群名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入集群名称" />
        </el-form-item>
        
        <el-form-item label="认证类型" prop="kube_type">
          <el-radio-group v-model="formData.kube_type">
            <el-radio :label="1">KubeConfig</el-radio>
            <el-radio :label="2">Token</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item
          v-if="formData.kube_type === 1"
          label="KubeConfig"
          prop="kube_config"
        >
          <el-input
            v-model="formData.kube_config"
            type="textarea"
            :rows="8"
            placeholder="请输入KubeConfig内容"
          />
        </el-form-item>
        
        <el-form-item
          v-if="formData.kube_type === 2"
          label="Token"
          prop="kube_config"
        >
          <el-input
            v-model="formData.kube_config"
            type="textarea"
            :rows="4"
            placeholder="请输入Token"
          />
        </el-form-item>
        
        <el-form-item label="API地址" prop="api_address">
          <el-input v-model="formData.api_address" placeholder="请输入API地址" />
        </el-form-item>
        
        <el-form-item label="CA证书" prop="kube_ca_crt">
          <el-input
            v-model="formData.kube_ca_crt"
            type="textarea"
            :rows="4"
            placeholder="请输入CA证书内容（可选）"
          />
        </el-form-item>
        
        <el-form-item label="Prometheus地址" prop="prometheus_url">
          <el-input v-model="formData.prometheus_url" placeholder="请输入Prometheus地址（可选）" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { useClusterStore } from '@/stores/kubernetes/cluster'
import K8sTable from '@/components/K8sTable.vue'
import type { Cluster, CreateClusterRequest, UpdateClusterRequest } from '@/types/kubernetes/cluster'
import { useRouter } from 'vue-router'

// 状态管理
const clusterStore = useClusterStore()
const router = useRouter()

// 响应式数据
const searchForm = reactive({
  name: '',
  kube_type: undefined as number | undefined,
})

const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const selectedClusters = ref<Cluster[]>([])

// 表单数据
const formData = reactive<CreateClusterRequest & { id?: number }>({
  name: '',
  kube_type: 1,
  kube_config: '',
  kube_ca_crt: '',
  api_address: '',
  prometheus_url: '',
  prometheus_auth_type: 0,
  prometheus_user: '',
  prometheus_pwd: '',
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入集群名称', trigger: 'blur' },
  ],
  kube_type: [
    { required: true, message: '请选择认证类型', trigger: 'change' },
  ],
  kube_config: [
    { required: true, message: '请输入认证信息', trigger: 'blur' },
  ],
  api_address: [
    { required: true, message: '请输入API地址', trigger: 'blur' },
  ],
}

// 计算属性
const dialogTitle = computed(() => {
  return dialogType.value === 'create' ? '创建集群' : '编辑集群'
})

// 方法
const handleSearch = () => {
  // TODO: 实现搜索功能
  console.log('搜索参数:', searchForm)
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.kube_type = undefined
  handleSearch()
}

const handleRefresh = async () => {
  try {
    await clusterStore.fetchClusters()
    ElMessage.success('刷新成功')
  } catch (error) {
    console.error('刷新失败:', error)
  }
}

const handleCreate = () => {
  dialogType.value = 'create'
  dialogVisible.value = true
  resetForm()
}

const handleEdit = (row: Cluster) => {
  dialogType.value = 'edit'
  dialogVisible.value = true
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    kube_type: row.kube_type,
    kube_config: row.kube_config,
    kube_ca_crt: row.kube_ca_crt,
    api_address: row.api_address,
    prometheus_url: row.prometheus_url,
    prometheus_auth_type: row.prometheus_auth_type,
    prometheus_user: row.prometheus_user,
    prometheus_pwd: row.prometheus_pwd,
  })
}

const handleView = (row: Cluster) => {
  viewDetail(row.id)
}

const handleDelete = async (row: Cluster) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除集群 "${row.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await clusterStore.deleteCluster(row.id)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const handleSelectionChange = (selection: Cluster[]) => {
  selectedClusters.value = selection
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    if (dialogType.value === 'create') {
      await clusterStore.createCluster(formData as CreateClusterRequest)
      dialogVisible.value = false
      resetForm()
    } else {
      await clusterStore.updateCluster(formData as UpdateClusterRequest)
      dialogVisible.value = false
      resetForm()
    }
  } catch (error) {
    console.error('提交失败:', error)
    // 错误消息已经在store中处理，这里不需要重复显示
  } finally {
    submitting.value = false
  }
}

const handleDialogClose = () => {
  resetForm()
}

const resetForm = () => {
  Object.assign(formData, {
    id: undefined,
    name: '',
    kube_type: 1,
    kube_config: '',
    kube_ca_crt: '',
    api_address: '',
    prometheus_url: '',
    prometheus_auth_type: 0,
    prometheus_user: '',
    prometheus_pwd: '',
  })
  formRef.value?.clearValidate()
}

const formatDate = (dateString: string) => {
  if (!dateString) return '未知时间'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) {
      return '无效日期'
    }
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    console.error('日期格式化错误:', error, dateString)
    return '日期错误'
  }
}

const viewDetail = (id: number) => {
  router.push(`/kubernetes/cluster/${id}`)
}

// 生命周期
onMounted(async () => {
  try {
    await clusterStore.fetchClusters()
  } catch (error) {
    console.error('获取集群列表失败:', error)
  }
})
</script>

<style scoped>
.cluster-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.search-section {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.table-section {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style> 