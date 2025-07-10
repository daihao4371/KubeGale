<template>
  <PageLayout title="RBAC权限管理">
    <template #subtitle>
      管理 Kubernetes RBAC 权限资源，包括 Role、ClusterRole、RoleBinding、ServiceAccount 等
    </template>

    <div class="rbac-grid">
      <!-- Role 管理卡片 -->
      <el-card class="rbac-card" @click="navigateTo('/homepage/kubernetes/rbac/role')">
        <div class="card-header">
          <el-icon class="card-icon role-icon"><UserFilled /></el-icon>
          <h3>Role 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ roleStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ roleStats.namespaces }}</span>
              <span class="stat-label">命名空间</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ roleStats.rules }}</span>
              <span class="stat-label">规则数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ roleStats.bindings }}</span>
              <span class="stat-label">绑定数</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- ClusterRole 管理卡片 -->
      <el-card class="rbac-card" @click="navigateTo('/homepage/kubernetes/rbac/clusterrole')">
        <div class="card-header">
          <el-icon class="card-icon clusterrole-icon"><Crown /></el-icon>
          <h3>ClusterRole 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ clusterRoleStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ clusterRoleStats.system }}</span>
              <span class="stat-label">系统角色</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ clusterRoleStats.custom }}</span>
              <span class="stat-label">自定义</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ clusterRoleStats.rules }}</span>
              <span class="stat-label">规则数</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- RoleBinding 管理卡片 -->
      <el-card class="rbac-card" @click="navigateTo('/homepage/kubernetes/rbac/rolebinding')">
        <div class="card-header">
          <el-icon class="card-icon binding-icon"><Connection /></el-icon>
          <h3>RoleBinding 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ bindingStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ bindingStats.users }}</span>
              <span class="stat-label">用户绑定</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ bindingStats.groups }}</span>
              <span class="stat-label">组绑定</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ bindingStats.serviceAccounts }}</span>
              <span class="stat-label">SA绑定</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- ServiceAccount 管理卡片 -->
      <el-card class="rbac-card" @click="navigateTo('/homepage/kubernetes/rbac/serviceaccount')">
        <div class="card-header">
          <el-icon class="card-icon sa-icon"><Avatar /></el-icon>
          <h3>ServiceAccount 管理</h3>
        </div>
        <div class="card-content">
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ saStats.total }}</span>
              <span class="stat-label">总数</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ saStats.default }}</span>
              <span class="stat-label">默认SA</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ saStats.custom }}</span>
              <span class="stat-label">自定义SA</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ saStats.tokens }}</span>
              <span class="stat-label">令牌数</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 权限分析 -->
    <div class="permission-analysis">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card class="analysis-card">
            <div class="analysis-header">
              <h4>权限分布统计</h4>
            </div>
            <div class="permission-stats">
              <div class="permission-item">
                <div class="permission-name">资源访问权限</div>
                <div class="permission-bar">
                  <div class="permission-fill read" style="width: 70%"></div>
                  <div class="permission-fill write" style="width: 30%"></div>
                </div>
                <div class="permission-labels">
                  <span class="label-item read">读取 70%</span>
                  <span class="label-item write">写入 30%</span>
                </div>
              </div>
              <div class="permission-item">
                <div class="permission-name">集群级别权限</div>
                <div class="permission-bar">
                  <div class="permission-fill admin" style="width: 20%"></div>
                  <div class="permission-fill view" style="width: 80%"></div>
                </div>
                <div class="permission-labels">
                  <span class="label-item admin">管理员 20%</span>
                  <span class="label-item view">查看 80%</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="analysis-card">
            <div class="analysis-header">
              <h4>安全风险评估</h4>
            </div>
            <div class="security-assessment">
              <div class="risk-item">
                <div class="risk-level high">高风险</div>
                <div class="risk-count">{{ securityStats.high }}</div>
                <div class="risk-desc">过度权限的角色</div>
              </div>
              <div class="risk-item">
                <div class="risk-level medium">中风险</div>
                <div class="risk-count">{{ securityStats.medium }}</div>
                <div class="risk-desc">权限集中的用户</div>
              </div>
              <div class="risk-item">
                <div class="risk-level low">低风险</div>
                <div class="risk-count">{{ securityStats.low }}</div>
                <div class="risk-desc">最小权限角色</div>
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
  UserFilled, Crown, Connection, Avatar 
} from '@element-plus/icons-vue'
import PageLayout from '@/components/layout/PageLayout.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import { formatAge } from '@/utils/date'

defineOptions({
  name: 'RbacOverview'
})

const router = useRouter()

// 统计数据
const roleStats = reactive({
  total: 0,
  namespaces: 0,
  rules: 0,
  bindings: 0
})

const clusterRoleStats = reactive({
  total: 0,
  system: 0,
  custom: 0,
  rules: 0
})

const bindingStats = reactive({
  total: 0,
  users: 0,
  groups: 0,
  serviceAccounts: 0
})

const saStats = reactive({
  total: 0,
  default: 0,
  custom: 0,
  tokens: 0
})

const securityStats = reactive({
  high: 0,
  medium: 0,
  low: 0
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
    Object.assign(roleStats, {
      total: 12,
      namespaces: 5,
      rules: 35,
      bindings: 18
    })
    
    Object.assign(clusterRoleStats, {
      total: 25,
      system: 20,
      custom: 5,
      rules: 180
    })
    
    Object.assign(bindingStats, {
      total: 22,
      users: 8,
      groups: 6,
      serviceAccounts: 8
    })
    
    Object.assign(saStats, {
      total: 15,
      default: 5,
      custom: 10,
      tokens: 15
    })
    
    Object.assign(securityStats, {
      high: 2,
      medium: 5,
      low: 15
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
        object: 'Role/app-reader',
        reason: 'Created',
        message: 'Role created successfully',
        time: '2024-01-01T10:00:00Z'
      },
      {
        type: 'Warning',
        object: 'RoleBinding/admin-binding',
        reason: 'PermissionDenied',
        message: 'User does not have permission to create RoleBinding',
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
.rbac-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.rbac-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e4e7ed;
}

.rbac-card:hover {
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

.role-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.clusterrole-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.binding-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.sa-icon {
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

.stat-label {
  font-size: 12px;
  color: #6c757d;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.permission-analysis {
  margin-bottom: 30px;
}

.analysis-card {
  height: 100%;
}

.analysis-header {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f2f5;
}

.analysis-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

.permission-stats {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.permission-item {
  padding: 16px 0;
}

.permission-name {
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 8px;
}

.permission-bar {
  height: 8px;
  background: #f0f2f5;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  margin-bottom: 8px;
}

.permission-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.permission-fill.read {
  background: #67c23a;
}

.permission-fill.write {
  background: #e6a23c;
}

.permission-fill.admin {
  background: #f56c6c;
}

.permission-fill.view {
  background: #409eff;
}

.permission-labels {
  display: flex;
  gap: 16px;
}

.label-item {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  color: #fff;
}

.label-item.read {
  background: #67c23a;
}

.label-item.write {
  background: #e6a23c;
}

.label-item.admin {
  background: #f56c6c;
}

.label-item.view {
  background: #409eff;
}

.security-assessment {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.risk-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.risk-item:last-child {
  border-bottom: none;
}

.risk-level {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  color: #fff;
  min-width: 60px;
  text-align: center;
}

.risk-level.high {
  background: #f56c6c;
}

.risk-level.medium {
  background: #e6a23c;
}

.risk-level.low {
  background: #67c23a;
}

.risk-count {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  min-width: 40px;
}

.risk-desc {
  flex: 1;
  color: #606266;
  font-size: 14px;
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
  .rbac-grid {
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
  
  .risk-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>