<template>
  <div class="kubernetes-overview">
    <PageLayout>
      <template #header>
        <h1>Kubernetes 管理概览</h1>
      </template>
      
      <template #actions>
        <el-button type="primary" @click="refreshOverview">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </template>
      
      <template #default>
        <div class="overview-content">
          <!-- 统计卡片 -->
          <div class="stats-grid">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon cluster">
                  <el-icon><Box /></el-icon>
                </div>
                <div class="stat-info">
                  <h3>{{ clusterCount }}</h3>
                  <p>集群数量</p>
                </div>
              </div>
            </el-card>
            
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon nodes">
                  <el-icon><Monitor /></el-icon>
                </div>
                <div class="stat-info">
                  <h3>{{ nodeCount }}</h3>
                  <p>节点数量</p>
                </div>
              </div>
            </el-card>
            
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon pods">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="stat-info">
                  <h3>{{ podCount }}</h3>
                  <p>Pod数量</p>
                </div>
              </div>
            </el-card>
            
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon namespaces">
                  <el-icon><Files /></el-icon>
                </div>
                <div class="stat-info">
                  <h3>{{ namespaceCount }}</h3>
                  <p>命名空间</p>
                </div>
              </div>
            </el-card>
          </div>
          
          <!-- 快速导航 -->
          <div class="quick-nav">
            <h2>快速导航</h2>
            <div class="nav-grid">
              <div class="nav-section">
                <h3>工作负载</h3>
                <div class="nav-items">
                  <el-button 
                    v-for="item in workloadItems" 
                    :key="item.path"
                    @click="navigateTo(item.path)"
                    class="nav-item"
                  >
                    <el-icon><component :is="item.icon" /></el-icon>
                    {{ item.title }}
                  </el-button>
                </div>
              </div>
              
              <div class="nav-section">
                <h3>网络管理</h3>
                <div class="nav-items">
                  <el-button 
                    v-for="item in networkItems" 
                    :key="item.path"
                    @click="navigateTo(item.path)"
                    class="nav-item"
                  >
                    <el-icon><component :is="item.icon" /></el-icon>
                    {{ item.title }}
                  </el-button>
                </div>
              </div>
              
              <div class="nav-section">
                <h3>存储配置</h3>
                <div class="nav-items">
                  <el-button 
                    v-for="item in storageItems" 
                    :key="item.path"
                    @click="navigateTo(item.path)"
                    class="nav-item"
                  >
                    <el-icon><component :is="item.icon" /></el-icon>
                    {{ item.title }}
                  </el-button>
                </div>
              </div>
              
              <div class="nav-section">
                <h3>系统管理</h3>
                <div class="nav-items">
                  <el-button 
                    v-for="item in systemItems" 
                    :key="item.path"
                    @click="navigateTo(item.path)"
                    class="nav-item"
                  >
                    <el-icon><component :is="item.icon" /></el-icon>
                    {{ item.title }}
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </PageLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Box, Monitor, Document, Files, Refresh, Connection, Folder, Setting, User } from '@element-plus/icons-vue'
import PageLayout from '@/components/layout/PageLayout.vue'

defineOptions({
  name: 'KubernetesOverview'
})

const router = useRouter()

// 统计数据
const clusterCount = ref(0)
const nodeCount = ref(0)
const podCount = ref(0)
const namespaceCount = ref(0)

// 快速导航配置
const workloadItems = [
  { title: 'Pod管理', path: '/homepage/kubernetes/workload/pod', icon: 'Document' },
  { title: 'Deployment', path: '/homepage/kubernetes/workload/deployment', icon: 'Document' },
  { title: 'StatefulSet', path: '/homepage/kubernetes/workload/statefulset', icon: 'Document' },
  { title: 'DaemonSet', path: '/homepage/kubernetes/workload/daemonset', icon: 'Document' }
]

const networkItems = [
  { title: 'Service', path: '/homepage/kubernetes/network/service', icon: 'Connection' },
  { title: 'Ingress', path: '/homepage/kubernetes/network/ingress', icon: 'Connection' },
  { title: 'Endpoint', path: '/homepage/kubernetes/network/endpoint', icon: 'Connection' },
  { title: 'NetworkPolicy', path: '/homepage/kubernetes/network/networkpolicy', icon: 'Connection' }
]

const storageItems = [
  { title: 'PersistentVolume', path: '/homepage/kubernetes/storage/pv', icon: 'Folder' },
  { title: 'PVC', path: '/homepage/kubernetes/storage/pvc', icon: 'Folder' },
  { title: 'StorageClass', path: '/homepage/kubernetes/storage/storageclass', icon: 'Folder' },
  { title: 'ConfigMap', path: '/homepage/kubernetes/config/configmap', icon: 'Setting' }
]

const systemItems = [
  { title: '集群管理', path: '/homepage/kubernetes/cluster', icon: 'Box' },
  { title: '节点管理', path: '/homepage/kubernetes/nodes', icon: 'Monitor' },
  { title: '命名空间', path: '/homepage/kubernetes/namespace', icon: 'Files' },
  { title: 'RBAC权限', path: '/homepage/kubernetes/rbac', icon: 'User' }
]

// 方法
const fetchOverviewData = async () => {
  try {
    // TODO: 调用API获取概览数据
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 模拟数据
    clusterCount.value = 3
    nodeCount.value = 12
    podCount.value = 145
    namespaceCount.value = 8
    
  } catch (error) {
    console.error('获取概览数据失败:', error)
    ElMessage.error('获取概览数据失败')
  }
}

const refreshOverview = () => {
  fetchOverviewData()
  ElMessage.success('数据已刷新')
}

const navigateTo = (path: string) => {
  router.push(path)
}

// 生命周期
onMounted(() => {
  fetchOverviewData()
})
</script>

<style scoped>
.kubernetes-overview {
  min-height: 100%;
}

.overview-content {
  padding: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 24px;
  color: white;
}

.stat-icon.cluster {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.nodes {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.pods {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.namespaces {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info h3 {
  margin: 0 0 4px 0;
  font-size: 28px;
  font-weight: 700;
  color: #2c3e50;
}

.stat-info p {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
}

.quick-nav h2 {
  margin: 0 0 24px 0;
  color: #2c3e50;
  font-size: 20px;
  font-weight: 600;
}

.nav-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

.nav-section h3 {
  margin: 0 0 16px 0;
  color: #34495e;
  font-size: 16px;
  font-weight: 600;
  padding-bottom: 8px;
  border-bottom: 2px solid #ecf0f1;
}

.nav-items {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}

.nav-item {
  height: 60px;
  padding: 16px;
  border: 1px solid #e4e7ed;
  background: #fafbfc;
  border-radius: 6px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
}

.nav-item:hover {
  background: #409eff;
  color: white;
  border-color: #409eff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.nav-item .el-icon {
  font-size: 18px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
  }
  
  .nav-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .nav-items {
    grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
    gap: 10px;
  }
  
  .overview-content {
    padding: 16px;
  }
}
</style>