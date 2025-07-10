<template>
  <router-view v-if="showRouterView" />
  <StorageOverview v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import StorageOverview from './StorageOverview.vue'

defineOptions({
  name: 'StorageView'
})

const route = useRoute()

// 如果路由有子路径，显示子路由，否则显示概览页面
const showRouterView = computed(() => {
  const segments = route.path.split('/').filter(Boolean)
  const storageIndex = segments.indexOf('storage')
  return storageIndex >= 0 && segments.length > storageIndex + 1
})
</script> 