<template>
  <router-view v-if="showRouterView" />
  <NetworkOverview v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import NetworkOverview from './NetworkOverview.vue'

defineOptions({
  name: 'NetworkView'
})

const route = useRoute()

// 如果路由有子路径，显示子路由，否则显示概览页面
const showRouterView = computed(() => {
  const segments = route.path.split('/').filter(Boolean)
  const networkIndex = segments.indexOf('network')
  return networkIndex >= 0 && segments.length > networkIndex + 1
})
</script> 