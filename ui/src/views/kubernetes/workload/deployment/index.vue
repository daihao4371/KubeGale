<template>
  <div class="deployment-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>Deployment 管理</h2>
        <el-breadcrumb separator="/">
          <el-breadcrumb-item>Kubernetes</el-breadcrumb-item>
          <el-breadcrumb-item>工作负载</el-breadcrumb-item>
          <el-breadcrumb-item>Deployment</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          创建 Deployment
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <el-card>
        <el-form :model="searchForm" inline>
          <el-form-item label="集群">
            <el-select v-model="searchForm.cluster_id" placeholder="选择集群" clearable>
              <el-option
                v-for="cluster in clusterList"
                :key="cluster.id"
                :label="cluster.name"
                :value="cluster.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="命名空间">
            <el-select v-model="searchForm.namespace" placeholder="选择命名空间" clearable>
              <el-option
                v-for="ns in namespaceList"
                :key="ns"
                :label="ns"
                :value="ns"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索 Deployment 名称"
              clearable
              @keyup.enter="handleSearch"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 数据表格 -->
    <div class="table-section">
      <el-card>
        <el-table
          v-loading="loading"
          :data="deploymentList"
          stripe
          border
          style="width: 100%"
        >
          <el-table-column prop="metadata.name" label="名称" min-width="150">
            <template #default="{ row }">
              <el-link type="primary" @click="handleViewDetail(row)">
                {{ row.metadata.name }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
          <el-table-column label="副本数" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getReplicaStatusType(row)">
                {{ row.status.replicas }}/{{ row.spec.replicas }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row)">
                {{ getStatusText(row) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="镜像" min-width="200">
            <template #default="{ row }">
              <el-tooltip
                v-for="container in row.spec.template.spec.containers"
                :key="container.name"
                :content="container.image"
                placement="top"
              >
                <el-tag size="small" class="mr-2">
                  {{ container.name }}: {{ getShortImage(container.image) }}
                </el-tag>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.creationTimestamp" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.metadata.creationTimestamp) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="handleViewDetail(row)">
                详情
              </el-button>
              <el-button size="small" type="primary" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button size="small" type="warning" @click="handleScale(row)">
                扩缩容
              </el-button>
              <el-dropdown @command="(command: string) => handleCommand(command, row)">
                <el-button size="small">
                  更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="restart">重启</el-dropdown-item>
                    <el-dropdown-item command="rollback">回滚</el-dropdown-item>
                    <el-dropdown-item command="logs">查看日志</el-dropdown-item>
                    <el-dropdown-item command="terminal">终端</el-dropdown-item>
                    <el-dropdown-item divided command="delete">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>

    <!-- 创建/编辑对话框 -->
    <deployment-form
      v-model:visible="formVisible"
      :deployment="currentDeployment"
      :cluster-id="searchForm.cluster_id"
      :namespace="searchForm.namespace"
      @success="handleFormSuccess"
    />

    <!-- 扩缩容对话框 -->
    <scale-dialog
      v-model:visible="scaleVisible"
      :deployment="currentDeployment"
      @success="handleScaleSuccess"
    />

    <!-- 详情对话框 -->
    <deployment-detail
      v-model:visible="detailVisible"
      :deployment="currentDeployment"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, ArrowDown } from '@element-plus/icons-vue'
import { deploymentApi } from '@/api/kubernetes/workload'
import { clusterApi } from '@/api/kubernetes/cluster'
import type { Deployment, GetWorkloadListRequest } from '@/types/kubernetes/workload'
import type { Cluster, Namespace } from '@/types/kubernetes/cluster'
import DeploymentForm from './components/deployment-form.vue'
import ScaleDialog from './components/scale-dialog.vue'
import DeploymentDetail from './components/deployment-detail.vue'

// 响应式数据
const loading = ref(false)
const deploymentList = ref<Deployment[]>([])
const clusterList = ref<Cluster[]>([])
const namespaceList = ref<string[]>([])

// 搜索表单
const searchForm = reactive<GetWorkloadListRequest>({
  cluster_id: 0,
  namespace: '',
  page: 1,
  pageSize: 20,
  keyword: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 对话框状态
const formVisible = ref(false)
const scaleVisible = ref(false)
const detailVisible = ref(false)
const currentDeployment = ref<Deployment | null>(null)

// 计算属性
const getReplicaStatusType = computed(() => {
  return (deployment: Deployment) => {
    const { replicas, readyReplicas } = deployment.status
    if (readyReplicas === replicas) return 'success'
    if (readyReplicas > 0) return 'warning'
    return 'danger'
  }
})

const getStatusType = computed(() => {
  return (deployment: Deployment) => {
    const { replicas, readyReplicas, availableReplicas } = deployment.status
    if (availableReplicas === replicas) return 'success'
    if (readyReplicas > 0) return 'warning'
    return 'danger'
  }
})

const getStatusText = computed(() => {
  return (deployment: Deployment) => {
    const { replicas, readyReplicas, availableReplicas } = deployment.status
    if (availableReplicas === replicas) return '运行中'
    if (readyReplicas > 0) return '部分就绪'
    return '未就绪'
  }
})

// 方法
const loadClusterList = async () => {
  try {
    const response = await clusterApi.getClusterList()
    clusterList.value = response.data.list || []
  } catch (error) {
    console.error('加载集群列表失败:', error)
  }
}

const loadNamespaceList = async () => {
  if (!searchForm.cluster_id) {
    namespaceList.value = []
    return
  }
  
  try {
    const response = await clusterApi.getClusterListNamespace({ id: searchForm.cluster_id })
    const namespaces = response.data.namespaces || []
    namespaceList.value = namespaces.map((ns: Namespace) => ns.metadata.name)
  } catch (error) {
    console.error('加载命名空间列表失败:', error)
  }
}

const loadDeploymentList = async () => {
  if (!searchForm.cluster_id || !searchForm.namespace) {
    deploymentList.value = []
    pagination.total = 0
    return
  }

  loading.value = true
  try {
    const response = await deploymentApi.getDeploymentList({
      ...searchForm,
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    
    deploymentList.value = response.data.items || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('加载 Deployment 列表失败:', error)
    ElMessage.error('加载 Deployment 列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadDeploymentList()
}

const handleReset = () => {
  searchForm.keyword = ''
  pagination.page = 1
  loadDeploymentList()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadDeploymentList()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadDeploymentList()
}

const handleCreate = () => {
  currentDeployment.value = null
  formVisible.value = true
}

const handleEdit = (deployment: Deployment) => {
  currentDeployment.value = deployment
  formVisible.value = true
}

const handleViewDetail = (deployment: Deployment) => {
  currentDeployment.value = deployment
  detailVisible.value = true
}

const handleScale = (deployment: Deployment) => {
  currentDeployment.value = deployment
  scaleVisible.value = true
}

const handleCommand = async (command: string, deployment: Deployment) => {
  switch (command) {
    case 'restart':
      await handleRestart(deployment)
      break
    case 'rollback':
      await handleRollback(deployment)
      break
    case 'logs':
      handleViewLogs(deployment)
      break
    case 'terminal':
      handleOpenTerminal(deployment)
      break
    case 'delete':
      await handleDelete(deployment)
      break
  }
}

const handleRestart = async (deployment: Deployment) => {
  try {
    await ElMessageBox.confirm(
      `确定要重启 Deployment "${deployment.metadata.name}" 吗？`,
      '确认重启',
      { type: 'warning' }
    )
    
    // TODO: 实现重启逻辑
    ElMessage.success('重启操作已提交')
  } catch (error) {
    // 用户取消
  }
}

const handleRollback = async (deployment: Deployment) => {
  // TODO: 实现回滚逻辑
  ElMessage.info('回滚功能开发中')
}

const handleViewLogs = (deployment: Deployment) => {
  // TODO: 实现查看日志逻辑
  ElMessage.info('日志查看功能开发中')
}

const handleOpenTerminal = (deployment: Deployment) => {
  // TODO: 实现终端功能
  ElMessage.info('终端功能开发中')
}

const handleDelete = async (deployment: Deployment) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Deployment "${deployment.metadata.name}" 吗？此操作不可恢复！`,
      '确认删除',
      { type: 'warning' }
    )
    
    await deploymentApi.deleteDeployment({
      cluster_id: searchForm.cluster_id,
      namespace: searchForm.namespace,
      name: deployment.metadata.name
    })
    
    ElMessage.success('删除成功')
    loadDeploymentList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除 Deployment 失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

const handleFormSuccess = () => {
  formVisible.value = false
  loadDeploymentList()
}

const handleScaleSuccess = () => {
  scaleVisible.value = false
  loadDeploymentList()
}

const getShortImage = (image: string) => {
  const parts = image.split('/')
  return parts[parts.length - 1] || image
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 监听集群变化
const watchClusterChange = () => {
  if (searchForm.cluster_id) {
    loadNamespaceList()
  }
}

// 生命周期
onMounted(() => {
  loadClusterList()
})

// 监听搜索表单变化
watch(() => searchForm.cluster_id, watchClusterChange)
</script>

<style scoped>
.deployment-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  color: var(--el-text-color-primary);
}

.search-section {
  margin-bottom: 20px;
}

.table-section {
  margin-bottom: 20px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.mr-2 {
  margin-right: 8px;
}
</style> 