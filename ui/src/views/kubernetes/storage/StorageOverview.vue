<template>
  <PageLayout title="存储管理">
    <template #subtitle>
      管理 Kubernetes 存储资源，包括 PV、PVC、StorageClass 等
    </template>

    <div class="storage-grid">
      <!-- PV 管理卡片 -->
      <el-card class="storage-card" @click="navigateTo('/homepage/kubernetes/storage/pv')">
        <div class="card-header">
          <el-icon class="card-icon pv-icon"><HardDrive /></el-icon>
          <h3>PV 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ pvStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ pvStats.available }}</span>
              <span class="stat-label">可用</span>
            </div>
            <div class="stat-item">
              <span class="stat-value warning">{{ pvStats.bound }}</span>
              <span class="stat-label">已绑定</span>
            </div>
            <div class="stat-item">
              <span class="stat-value danger">{{ pvStats.failed }}</span>
              <span class="stat-label">失败</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- PVC 管理卡片 -->
      <el-card class="storage-card" @click="navigateTo('/homepage/kubernetes/storage/pvc')">
        <div class="card-header">
          <el-icon class="card-icon pvc-icon"><Folder /></el-icon>
          <h3>PVC 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ pvcStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ pvcStats.bound }}</span>
              <span class="stat-label">已绑定</span>
            </div>
            <div class="stat-item">
              <span class="stat-value warning">{{ pvcStats.pending }}</span>
              <span class="stat-label">等待中</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ pvcStats.capacity }}</span>
              <span class="stat-label">容量(GB)</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- StorageClass 管理卡片 -->
      <el-card class="storage-card" @click="navigateTo('/homepage/kubernetes/storage/storageclass')">
        <div class="card-header">
          <el-icon class="card-icon sc-icon"><Management /></el-icon>
          <h3>StorageClass 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ scStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ scStats.default }}</span>
              <span class="stat-label">默认</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ scStats.provisioner }}</span>
              <span class="stat-label">供应商</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ scStats.volumeBindingMode }}</span>
              <span class="stat-label">绑定模式</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 存储使用情况 -->
      <el-card class="storage-card overview-card">
        <div class="card-header">
          <el-icon class="card-icon usage-icon"><DataAnalysis /></el-icon>
          <h3>存储使用情况</h3>
        </div>
        <div class="card-content">
          <div class="usage-stats">
            <div class="usage-item">
              <div class="usage-label">已使用存储</div>
              <div class="usage-bar">
                <div class="usage-fill" :style="{ width: usagePercentage + '%' }"></div>
              </div>
              <div class="usage-text">{{ usageStats.used }}GB / {{ usageStats.total }}GB</div>
            </div>
            <div class="quick-stats">
              <div class="quick-stat">
                <span class="stat-value">{{ usageStats.pvCount }}</span>
                <span class="stat-label">PV 数量</span>
              </div>
              <div class="quick-stat">
                <span class="stat-value">{{ usageStats.pvcCount }}</span>
                <span class="stat-label">PVC 数量</span>
              </div>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  HardDrive, Folder, Management, DataAnalysis 
} from '@element-plus/icons-vue'
import PageLayout from '@/components/layout/PageLayout.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'StorageOverview'
})

const router = useRouter()

// 统计数据
const pvStats = reactive({
  total: 0,
  available: 0,
  bound: 0,
  failed: 0
})

const pvcStats = reactive({
  total: 0,
  bound: 0,
  pending: 0,
  capacity: 0
})

const scStats = reactive({
  total: 0,
  default: 0,
  provisioner: 0,
  volumeBindingMode: 0
})

const usageStats = reactive({
  total: 0,
  used: 0,
  pvCount: 0,
  pvcCount: 0
})

const recentEvents = ref([])

// 计算属性
const usagePercentage = computed(() => {
  return usageStats.total > 0 ? (usageStats.used / usageStats.total) * 100 : 0
})

// 方法
const navigateTo = (path: string) => {
  router.push(path)
}

const fetchStats = async () => {
  try {
    // TODO: 调用实际的 API 获取统计数据
    
    // 模拟数据
    Object.assign(pvStats, {
      total: 10,
      available: 4,
      bound: 5,
      failed: 1
    })
    
    Object.assign(pvcStats, {
      total: 8,
      bound: 6,
      pending: 2,
      capacity: 250
    })
    
    Object.assign(scStats, {
      total: 3,
      default: 1,
      provisioner: 3,
      volumeBindingMode: 2
    })
    
    Object.assign(usageStats, {
      total: 500,
      used: 180,
      pvCount: 10,
      pvcCount: 8
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
        object: 'PV/pv-mysql-001',
        reason: 'VolumeProvisioned',
        message: 'Successfully provisioned volume',
        time: '2024-01-01T10:00:00Z'
      },
      {
        type: 'Warning',
        object: 'PVC/data-redis-0',
        reason: 'ProvisioningFailed',
        message: 'Failed to provision volume: no storage class',
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
.storage-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.storage-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e4e7ed;
}

.storage-card:hover:not(.overview-card) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-color: #409eff;
}

.overview-card {
  cursor: default;
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

.pv-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.pvc-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.sc-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.usage-icon {
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

.usage-stats {
  padding: 0 8px;
}

.usage-item {
  margin-bottom: 20px;
}

.usage-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.usage-bar {
  height: 8px;
  background: #f0f2f5;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.usage-fill {
  height: 100%;
  background: linear-gradient(90deg, #409eff 0%, #67c23a 100%);
  transition: width 0.3s ease;
}

.usage-text {
  font-size: 12px;
  color: #909399;
  text-align: center;
}

.quick-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.quick-stat {
  text-align: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

.quick-stat .stat-value {
  font-size: 20px;
  margin-bottom: 4px;
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
  .storage-grid {
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