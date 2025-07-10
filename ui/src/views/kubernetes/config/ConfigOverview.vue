<template>
  <PageLayout title="配置管理">
    <template #subtitle>
      管理 Kubernetes 配置资源，包括 ConfigMap、Secret、ResourceQuota、HPA 等
    </template>

    <div class="config-grid">
      <!-- ConfigMap 管理卡片 -->
      <el-card class="config-card" @click="navigateTo('/homepage/kubernetes/config/configmap')">
        <div class="card-header">
          <el-icon class="card-icon configmap-icon"><Document /></el-icon>
          <h3>ConfigMap 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ configMapStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ configMapStats.namespaces }}</span>
              <span class="stat-label">命名空间</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ configMapStats.keys }}</span>
              <span class="stat-label">配置项</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ configMapStats.size }}</span>
              <span class="stat-label">大小(KB)</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Secret 管理卡片 -->
      <el-card class="config-card" @click="navigateTo('/homepage/kubernetes/config/secret')">
        <div class="card-header">
          <el-icon class="card-icon secret-icon"><Key /></el-icon>
          <h3>Secret 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ secretStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ secretStats.opaque }}</span>
              <span class="stat-label">Opaque</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ secretStats.tls }}</span>
              <span class="stat-label">TLS</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ secretStats.dockerRegistry }}</span>
              <span class="stat-label">Registry</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- ResourceQuota 管理卡片 -->
      <el-card class="config-card" @click="navigateTo('/homepage/kubernetes/config/resourcequota')">
        <div class="card-header">
          <el-icon class="card-icon quota-icon"><Odometer /></el-icon>
          <h3>ResourceQuota 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ quotaStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ quotaStats.cpu }}</span>
              <span class="stat-label">CPU限制</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ quotaStats.memory }}</span>
              <span class="stat-label">内存限制</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ quotaStats.pods }}</span>
              <span class="stat-label">Pod限制</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- HPA 管理卡片 -->
      <el-card class="config-card" @click="navigateTo('/homepage/kubernetes/config/hpa')">
        <div class="card-header">
          <el-icon class="card-icon hpa-icon"><TrendCharts /></el-icon>
          <h3>HPA 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ hpaStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ hpaStats.active }}</span>
              <span class="stat-label">活跃</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ hpaStats.targets }}</span>
              <span class="stat-label">目标</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ hpaStats.replicas }}</span>
              <span class="stat-label">副本数</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 配置使用情况 -->
    <div class="usage-overview">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="usage-card">
            <div class="usage-header">
              <h4>命名空间配置分布</h4>
            </div>
            <div class="namespace-stats">
              <div v-for="ns in namespaceStats" :key="ns.name" class="namespace-item">
                <div class="namespace-name">{{ ns.name }}</div>
                <div class="namespace-counts">
                  <span class="count-item">CM: {{ ns.configMaps }}</span>
                  <span class="count-item">Secret: {{ ns.secrets }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="usage-card">
            <div class="usage-header">
              <h4>配置类型统计</h4>
            </div>
            <div class="type-stats">
              <div class="type-item">
                <div class="type-icon configmap-small"></div>
                <div class="type-info">
                  <div class="type-name">ConfigMap</div>
                  <div class="type-count">{{ configMapStats.total }} 个</div>
                </div>
              </div>
              <div class="type-item">
                <div class="type-icon secret-small"></div>
                <div class="type-info">
                  <div class="type-name">Secret</div>
                  <div class="type-count">{{ secretStats.total }} 个</div>
                </div>
              </div>
              <div class="type-item">
                <div class="type-icon quota-small"></div>
                <div class="type-info">
                  <div class="type-name">ResourceQuota</div>
                  <div class="type-count">{{ quotaStats.total }} 个</div>
                </div>
              </div>
              <div class="type-item">
                <div class="type-icon hpa-small"></div>
                <div class="type-info">
                  <div class="type-name">HPA</div>
                  <div class="type-count">{{ hpaStats.total }} 个</div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
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
  Document, Key, Odometer, TrendCharts 
} from '@element-plus/icons-vue'
import PageLayout from '@/components/layout/PageLayout.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'ConfigOverview'
})

const router = useRouter()

// 统计数据
const configMapStats = reactive({
  total: 0,
  namespaces: 0,
  keys: 0,
  size: 0
})

const secretStats = reactive({
  total: 0,
  opaque: 0,
  tls: 0,
  dockerRegistry: 0
})

const quotaStats = reactive({
  total: 0,
  cpu: 0,
  memory: 0,
  pods: 0
})

const hpaStats = reactive({
  total: 0,
  active: 0,
  targets: 0,
  replicas: 0
})

const namespaceStats = ref([])
const recentEvents = ref([])

// 方法
const navigateTo = (path: string) => {
  router.push(path)
}

const fetchStats = async () => {
  try {
    // TODO: 调用实际的 API 获取统计数据
    
    // 模拟数据
    Object.assign(configMapStats, {
      total: 15,
      namespaces: 5,
      keys: 42,
      size: 128
    })
    
    Object.assign(secretStats, {
      total: 8,
      opaque: 4,
      tls: 2,
      dockerRegistry: 2
    })
    
    Object.assign(quotaStats, {
      total: 3,
      cpu: 12,
      memory: 24,
      pods: 50
    })
    
    Object.assign(hpaStats, {
      total: 4,
      active: 3,
      targets: 6,
      replicas: 18
    })
    
    namespaceStats.value = [
      { name: 'default', configMaps: 5, secrets: 3 },
      { name: 'kube-system', configMaps: 8, secrets: 4 },
      { name: 'production', configMaps: 2, secrets: 1 }
    ]
    
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
        object: 'ConfigMap/app-config',
        reason: 'Created',
        message: 'ConfigMap created successfully',
        time: '2024-01-01T10:00:00Z'
      },
      {
        type: 'Warning',
        object: 'Secret/tls-secret',
        reason: 'CertificateExpired',
        message: 'TLS certificate will expire soon',
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
.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.config-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e4e7ed;
}

.config-card:hover {
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

.configmap-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.secret-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.quota-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.hpa-icon {
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

.usage-overview {
  margin-bottom: 30px;
}

.usage-card {
  height: 100%;
}

.usage-header {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f2f5;
}

.usage-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.namespace-stats {
  max-height: 200px;
  overflow-y: auto;
}

.namespace-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.namespace-item:last-child {
  border-bottom: none;
}

.namespace-name {
  font-weight: 500;
  color: #2c3e50;
}

.namespace-counts {
  display: flex;
  gap: 12px;
}

.count-item {
  font-size: 12px;
  color: #606266;
  background: #f0f2f5;
  padding: 2px 8px;
  border-radius: 4px;
}

.type-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.type-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.type-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
}

.configmap-small {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.secret-small {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.quota-small {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.hpa-small {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.type-info {
  flex: 1;
}

.type-name {
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 4px;
}

.type-count {
  font-size: 12px;
  color: #606266;
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
  .config-grid {
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