<template>
  <router-view v-if="showRouterView" />
  <RbacOverview v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import RbacOverview from './RbacOverview.vue'

defineOptions({
  name: 'RbacView'
})

const route = useRoute()

// 如果路由有子路径，显示子路由，否则显示概览页面
const showRouterView = computed(() => {
  const segments = route.path.split('/').filter(Boolean)
  const rbacIndex = segments.indexOf('rbac')
  return rbacIndex >= 0 && segments.length > rbacIndex + 1
})
</script> 