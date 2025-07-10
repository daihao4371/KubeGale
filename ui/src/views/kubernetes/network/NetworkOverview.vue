<template>
  <PageLayout title="网络管理">
    <template #subtitle>
      管理 Kubernetes 网络资源，包括 Service、Ingress、NetworkPolicy 等
    </template>

    <div class="network-grid">
      <!-- Service 管理卡片 -->
      <el-card class="network-card" @click="navigateTo('/homepage/kubernetes/network/service')">
        <div class="card-header">
          <el-icon class="card-icon service-icon"><Connection /></el-icon>
          <h3>Service 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ serviceStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ serviceStats.clusterIP }}</span>
              <span class="stat-label">ClusterIP</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ serviceStats.nodePort }}</span>
              <span class="stat-label">NodePort</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ serviceStats.loadBalancer }}</span>
              <span class="stat-label">LoadBalancer</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Ingress 管理卡片 -->
      <el-card class="network-card" @click="navigateTo('/homepage/kubernetes/network/ingress')">
        <div class="card-header">
          <el-icon class="card-icon ingress-icon"><Link /></el-icon>
          <h3>Ingress 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ ingressStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ ingressStats.active }}</span>
              <span class="stat-label">活跃</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ ingressStats.hosts }}</span>
              <span class="stat-label">主机数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ ingressStats.tls }}</span>
              <span class="stat-label">TLS</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Endpoint 管理卡片 -->
      <el-card class="network-card" @click="navigateTo('/homepage/kubernetes/network/endpoint')">
        <div class="card-header">
          <el-icon class="card-icon endpoint-icon"><Position /></el-icon>
          <h3>Endpoint 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ endpointStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ endpointStats.ready }}</span>
              <span class="stat-label">就绪</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ endpointStats.addresses }}</span>
              <span class="stat-label">地址数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ endpointStats.ports }}</span>
              <span class="stat-label">端口数</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- NetworkPolicy 管理卡片 -->
      <el-card class="network-card" @click="navigateTo('/homepage/kubernetes/network/networkpolicy')">
        <div class="card-header">
          <el-icon class="card-icon policy-icon"><Lock /></el-icon>
          <h3>NetworkPolicy 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ policyStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ policyStats.ingress }}</span>
              <span class="stat-label">入站规则</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ policyStats.egress }}</span>
              <span class="stat-label">出站规则</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ policyStats.namespaces }}</span>
              <span class="stat-label">命名空间</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 最近活动 -->
    <div class="recent-activity">
      <h3>最近活动</h3>
      <el-table :data="recentEvents" size="small">
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <K8sStatusBadge :status="row.type" size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="object" label="对象" width="150" />
        <el-table-column prop="reason" label="原因" width="120" />
        <el-table-column prop="message" label="消息" show-overflow-tooltip />
        <el-table-column prop="time" label="时间" width="150">
          <template #default="{ row }">
            {{ formatAge(row.time) }}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Connection, Link, Position, Lock 
} from '@element-plus/icons-vue'
import PageLayout from '@/components/layout/PageLayout.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'NetworkOverview'
})

const router = useRouter()

// 统计数据
const serviceStats = reactive({
  total: 0,
  clusterIP: 0,
  nodePort: 0,
  loadBalancer: 0
})

const ingressStats = reactive({
  total: 0,
  active: 0,
  hosts: 0,
  tls: 0
})

const endpointStats = reactive({
  total: 0,
  ready: 0,
  addresses: 0,
  ports: 0
})

const policyStats = reactive({
  total: 0,
  ingress: 0,
  egress: 0,
  namespaces: 0
})

const recentEvents = ref([])

// 方法
const navigateTo = (path: string) => {
  router.push(path)
}

const fetchStats = async () => {
  try {
    // TODO: 调用实际的 API 获取统计数据
    
    // 模拟数据
    Object.assign(serviceStats, {
      total: 12,
      clusterIP: 8,
      nodePort: 3,
      loadBalancer: 1
    })
    
    Object.assign(ingressStats, {
      total: 5,
      active: 4,
      hosts: 8,
      tls: 3
    })
    
    Object.assign(endpointStats, {
      total: 15,
      ready: 13,
      addresses: 28,
      ports: 45
    })
    
    Object.assign(policyStats, {
      total: 6,
      ingress: 8,
      egress: 12,
      namespaces: 4
    })
    
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const fetchRecentEvents = async () => {
  try {
    // TODO: 调用实际的 API 获取最近事件
    
    // 模拟数据
    recentEvents.value = [
      {
        type: 'Normal',
        object: 'Service/nginx-service',
        reason: 'Created',
        message: 'Service created successfully',
        time: '2024-01-01T10:00:00Z'
      },
      {
        type: 'Warning',
        object: 'Ingress/app-ingress',
        reason: 'CertificateError',
        message: 'Failed to obtain SSL certificate',
        time: '2024-01-01T09:55:00Z'
      }
    ]
    
  } catch (error) {
    console.error('获取最近事件失败:', error)
  }
}

// 生命周期
onMounted(async () => {
  await Promise.all([
    fetchStats(),
    fetchRecentEvents()
  ])
})
</script>

<style scoped>
.network-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.network-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e4e7ed;
}

.network-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-color: #409eff;
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f2f5;
}

.card-icon {
  font-size: 24px;
  margin-right: 12px;
  padding: 8px;
  border-radius: 6px;
  color: #fff;
}

.service-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.ingress-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.endpoint-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.policy-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-item {
  text-align: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

.stat-value {
  display: block;
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 4px;
}

.stat-value.running {
  color: #67c23a;
}

.stat-value.warning {
  color: #e6a23c;
}

.stat-value.danger {
  color: #f56c6c;
}

.stat-label {
  font-size: 12px;
  color: #6c757d;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.recent-activity {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.recent-activity h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .network-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .stats {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }
  
  .stat-item {
    padding: 8px;
  }
  
  .stat-value {
    font-size: 20px;
  }
}
</style>