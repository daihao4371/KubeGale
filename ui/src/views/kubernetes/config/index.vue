<template>
  <router-view v-if="showRouterView" />
  <ConfigOverview v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import ConfigOverview from './ConfigOverview.vue'

defineOptions({
  name: 'ConfigView'
})

const route = useRoute()

// 如果路由有子路径，显示子路由，否则显示概览页面
const showRouterView = computed(() => {
  const segments = route.path.split('/').filter(Boolean)
  const configIndex = segments.indexOf('config')
  return configIndex >= 0 && segments.length > configIndex + 1
})
</script> 