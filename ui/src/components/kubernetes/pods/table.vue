<template>
  <div class="pods-table">
    <K8sTable
      :data="pods"
      :loading="loading"
      :show-selection="false"
      :show-pagination="false"
    >
      <el-table-column prop="metadata.name" label="名称" min-width="150" />
      <el-table-column prop="metadata.namespace" label="命名空间" width="120" />
      <el-table-column prop="status.phase" label="状态" width="100">
        <template #default="{ row }">
          <K8sStatusBadge :status="row.status?.phase || 'Unknown'" />
        </template>
      </el-table-column>
      <el-table-column label="就绪状态" width="100">
        <template #default="{ row }">
          {{ getReadyStatus(row) }}
        </template>
      </el-table-column>
      <el-table-column label="重启次数" width="100">
        <template #default="{ row }">
          {{ getRestartCount(row) }}
        </template>
      </el-table-column>
      <el-table-column label="运行时间" width="120">
        <template #default="{ row }">
          {{ formatAge(row.metadata?.creationTimestamp) }}
        </template>
      </el-table-column>
    </K8sTable>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import K8sTable from '@/components/K8sTable.vue'
import K8sStatusBadge from '@/components/K8sStatusBadge.vue'
import { formatAge } from '@/utils/date'

interface Props {
  pods?: any[]
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  pods: () => [],
  loading: false,
})

const getReadyStatus = (pod: any): string => {
  const conditions = pod.status?.conditions || []
  const readyCondition = conditions.find((c: any) => c.type === 'Ready')
  return readyCondition?.status === 'True' ? '1/1' : '0/1'
}

const getRestartCount = (pod: any): number => {
  const containerStatuses = pod.status?.containerStatuses || []
  return containerStatuses.reduce((total: number, status: any) => total + (status.restartCount || 0), 0)
}
</script>

<style scoped>
.pods-table {
  width: 100%;
}
</style>