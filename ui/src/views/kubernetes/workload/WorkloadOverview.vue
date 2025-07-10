<template>
  <PageLayout title="工作负载管理">
    <template #subtitle>
      管理 Kubernetes 工作负载资源，包括 Pod、Deployment、StatefulSet 等
    </template>

    <div class="workload-grid">
      <!-- Pod 管理卡片 -->
      <el-card class="workload-card" @click="navigateTo('/homepage/kubernetes/workload/pod')">
        <div class="card-header">
          <el-icon class="card-icon pod-icon"><Box /></el-icon>
          <h3>Pod 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ podStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ podStats.running }}</span>
              <span class="stat-label">运行中</span>
            </div>
            <div class="stat-item">
              <span class="stat-value warning">{{ podStats.pending }}</span>
              <span class="stat-label">等待中</span>
            </div>
            <div class="stat-item">
              <span class="stat-value danger">{{ podStats.failed }}</span>
              <span class="stat-label">失败</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Deployment 管理卡片 -->
      <el-card class="workload-card" @click="navigateTo('/homepage/kubernetes/workload/deployment')">
        <div class="card-header">
          <el-icon class="card-icon deployment-icon"><Grid /></el-icon>
          <h3>Deployment 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ deploymentStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ deploymentStats.available }}</span>
              <span class="stat-label">可用</span>
            </div>
            <div class="stat-item">
              <span class="stat-value warning">{{ deploymentStats.progressing }}</span>
              <span class="stat-label">进行中</span>
            </div>
            <div class="stat-item">
              <span class="stat-value danger">{{ deploymentStats.failed }}</span>
              <span class="stat-label">失败</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- StatefulSet 管理卡片 -->
      <el-card class="workload-card" @click="navigateTo('/homepage/kubernetes/workload/statefulset')">
        <div class="card-header">
          <el-icon class="card-icon statefulset-icon"><Files /></el-icon>
          <h3>StatefulSet 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ statefulSetStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ statefulSetStats.ready }}</span>
              <span class="stat-label">就绪</span>
            </div>
            <div class="stat-item">
              <span class="stat-value warning">{{ statefulSetStats.updating }}</span>
              <span class="stat-label">更新中</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ statefulSetStats.replicas }}</span>
              <span class="stat-label">副本数</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- DaemonSet 管理卡片 -->
      <el-card class="workload-card" @click="navigateTo('/homepage/kubernetes/workload/daemonset')">
        <div class="card-header">
          <el-icon class="card-icon daemonset-icon"><Monitor /></el-icon>
          <h3>DaemonSet 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ daemonSetStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ daemonSetStats.desired }}</span>
              <span class="stat-label">期望</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ daemonSetStats.current }}</span>
              <span class="stat-label">当前</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ daemonSetStats.ready }}</span>
              <span class="stat-label">就绪</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Job 管理卡片 -->
      <el-card class="workload-card" @click="navigateTo('/homepage/kubernetes/workload/job')">
        <div class="card-header">
          <el-icon class="card-icon job-icon"><Timer /></el-icon>
          <h3>Job 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ jobStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ jobStats.completed }}</span>
              <span class="stat-label">完成</span>
            </div>
            <div class="stat-item">
              <span class="stat-value warning">{{ jobStats.active }}</span>
              <span class="stat-label">活跃</span>
            </div>
            <div class="stat-item">
              <span class="stat-value danger">{{ jobStats.failed }}</span>
              <span class="stat-label">失败</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- CronJob 管理卡片 -->
      <el-card class="workload-card" @click="navigateTo('/homepage/kubernetes/workload/cronjob')">
        <div class="card-header">
          <el-icon class="card-icon cronjob-icon"><Clock /></el-icon>
          <h3>CronJob 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ cronJobStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value running">{{ cronJobStats.active }}</span>
              <span class="stat-label">活跃</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ cronJobStats.suspended }}</span>
              <span class="stat-label">挂起</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ cronJobStats.lastSchedule }}</span>
              <span class="stat-label">最近调度</span>
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
  Box, Grid, Files, Monitor, Timer, Clock 
} from '@element-plus/icons-vue'
import PageLayout from '@/components/layout/PageLayout.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'WorkloadManagement'
})

const router = useRouter()

// 统计数据
const podStats = reactive({
  total: 0,
  running: 0,
  pending: 0,
  failed: 0
})

const deploymentStats = reactive({
  total: 0,
  available: 0,
  progressing: 0,
  failed: 0
})

const statefulSetStats = reactive({
  total: 0,
  ready: 0,
  updating: 0,
  replicas: 0
})

const daemonSetStats = reactive({
  total: 0,
  desired: 0,
  current: 0,
  ready: 0
})

const jobStats = reactive({
  total: 0,
  completed: 0,
  active: 0,
  failed: 0
})

const cronJobStats = reactive({
  total: 0,
  active: 0,
  suspended: 0,
  lastSchedule: '2h ago'
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
    Object.assign(podStats, {
      total: 15,
      running: 12,
      pending: 2,
      failed: 1
    })
    
    Object.assign(deploymentStats, {
      total: 8,
      available: 7,
      progressing: 1,
      failed: 0
    })
    
    Object.assign(statefulSetStats, {
      total: 3,
      ready: 3,
      updating: 0,
      replicas: 9
    })
    
    Object.assign(daemonSetStats, {
      total: 2,
      desired: 6,
      current: 6,
      ready: 6
    })
    
    Object.assign(jobStats, {
      total: 5,
      completed: 4,
      active: 1,
      failed: 0
    })
    
    Object.assign(cronJobStats, {
      total: 3,
      active: 2,
      suspended: 1,
      lastSchedule: '2h ago'
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
        object: 'Pod/nginx-pod-1',
        reason: 'Started',
        message: 'Started container nginx',
        time: '2024-01-01T10:00:00Z'
      },
      {
        type: 'Warning',
        object: 'Pod/app-pod-2',
        reason: 'Failed',
        message: 'Error: ImagePullBackOff',
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
.workload-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.workload-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e4e7ed;
}

.workload-card:hover {
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

.pod-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.deployment-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.statefulset-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.daemonset-icon {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.job-icon {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.cronjob-icon {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  color: #333 !important;
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
  .workload-grid {
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