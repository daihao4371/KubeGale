<template>
  <el-tag
    :type="statusType"
    :effect="effect"
    :size="size"
    class="k8s-status-badge"
  >
    <el-icon v-if="showIcon && statusIcon" class="status-icon">
      <component :is="statusIcon" />
    </el-icon>
    <span :class="{ 'with-icon': showIcon && statusIcon }">{{ displayText }}</span>
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { 
  CircleCheck, Warning, CircleClose, InfoFilled, 
  QuestionFilled 
} from '@element-plus/icons-vue'

interface Props {
  status: string
  effect?: 'dark' | 'light' | 'plain'
  size?: 'large' | 'default' | 'small'
  showIcon?: boolean
  customText?: string
}

const props = withDefaults(defineProps<Props>(), {
  effect: 'light',
  size: 'default',
  showIcon: false,
  customText: '',
})

// 显示文本
const displayText = computed(() => {
  return props.customText || props.status
})

// 状态类型映射
const statusType = computed(() => {
  const status = props.status.toLowerCase()
  
  // Pod状态
  if (status === 'running') return 'success'
  if (status === 'pending') return 'warning'
  if (status === 'failed' || status === 'error') return 'danger'
  if (status === 'succeeded') return 'success'
  if (status === 'terminating') return 'info'
  if (status === 'crashloopbackoff') return 'danger'
  if (status === 'imagepullbackoff') return 'danger'
  if (status === 'errimagepull') return 'danger'
  if (status === 'createcontainerconfigerror') return 'danger'
  if (status === 'invalidimagenameerror') return 'danger'
  if (status === 'completed') return 'success'
  
  // 部署状态
  if (status === 'available') return 'success'
  if (status === 'progressing') return 'warning'
  if (status === 'replicafailure') return 'danger'
  
  // 服务状态
  if (status === 'active') return 'success'
  if (status === 'inactive') return 'info'
  
  // 节点状态
  if (status === 'ready') return 'success'
  if (status === 'notready') return 'danger'
  if (status === 'unknown') return 'info'
  if (status === 'schedulingdisabled') return 'warning'
  
  // 存储状态
  if (status === 'bound') return 'success'
  if (status === 'available') return 'success'
  if (status === 'released') return 'warning'
  if (status === 'failed') return 'danger'
  
  // 网络状态
  if (status === 'loadbalancer') return 'success'
  if (status === 'clusterip') return 'primary'
  if (status === 'nodeport') return 'warning'
  if (status === 'externalname') return 'info'
  
  // 任务状态
  if (status === 'complete') return 'success'
  if (status === 'backofflimitexceeded') return 'danger'
  if (status === 'deadlineexceeded') return 'danger'
  
  // HPA状态
  if (status === 'scaledup') return 'success'
  if (status === 'scaleddown') return 'info'
  if (status === 'desiredreplicascurrentreplicas') return 'warning'
  
  // 通用状态
  if (status === 'true') return 'success'
  if (status === 'false') return 'danger'
  if (status === 'enabled') return 'success'
  if (status === 'disabled') return 'info'
  if (status === 'healthy') return 'success'
  if (status === 'unhealthy') return 'danger'
  
  // 默认
  return 'info'
})

// 状态图标映射
const statusIcon = computed(() => {
  if (!props.showIcon) return null
  
  const status = props.status.toLowerCase()
  
  // 成功状态
  if (['running', 'succeeded', 'available', 'active', 'ready', 'bound', 'complete', 'true', 'enabled', 'healthy'].includes(status)) {
    return 'CircleCheck'
  }
  
  // 警告状态
  if (['pending', 'progressing', 'schedulingdisabled', 'released', 'nodeport', 'desiredreplicascurrentreplicas'].includes(status)) {
    return 'Warning'
  }
  
  // 错误状态
  if (['failed', 'error', 'crashloopbackoff', 'imagepullbackoff', 'errimagepull', 'replicafailure', 'notready', 'backofflimitexceeded', 'deadlineexceeded', 'false', 'unhealthy'].includes(status)) {
    return 'CircleClose'
  }
  
  // 信息状态
  if (['terminating', 'inactive', 'unknown', 'scaleddown', 'disabled', 'clusterip', 'externalname'].includes(status)) {
    return 'InfoFilled'
  }
  
  return 'QuestionFilled'
})
</script>

<style scoped>
.k8s-status-badge {
  text-transform: capitalize;
  display: inline-flex;
  align-items: center;
}

.status-icon {
  font-size: 12px;
  margin-right: 4px;
}

.with-icon {
  display: inline-flex;
  align-items: center;
}

/* 不同尺寸的图标大小 */
.k8s-status-badge.el-tag--small .status-icon {
  font-size: 10px;
  margin-right: 2px;
}

.k8s-status-badge.el-tag--large .status-icon {
  font-size: 14px;
  margin-right: 6px;
}

/* 状态动画 */
.k8s-status-badge :deep(.el-icon) {
  transition: all 0.3s ease;
}

/* 脉冲动画用于等待状态 */
@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}

.k8s-status-badge.el-tag--warning :deep(.el-icon) {
  animation: pulse 2s infinite;
}
</style> 