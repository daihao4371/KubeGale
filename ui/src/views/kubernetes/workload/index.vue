<template>
  <router-view v-if="showRouterView" />
  <WorkloadOverview v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import WorkloadOverview from './WorkloadOverview.vue'

defineOptions({
  name: 'WorkloadView'
})

const route = useRoute()

// 如果路由有子路径，显示子路由，否则显示概览页面
const showRouterView = computed(() => {
  const segments = route.path.split('/').filter(Boolean)
  const workloadIndex = segments.indexOf('workload')
  return workloadIndex >= 0 && segments.length > workloadIndex + 1
})
</script> 