<template>
  <div class="kubernetes-container">
    <router-view v-if="$route.path !== '/homepage/kubernetes'" />
    <KubernetesOverview v-else />
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import KubernetesOverview from './overview/index.vue'

defineOptions({
  name: 'KubernetesView'
})

const $route = useRoute()
</script>

<style scoped>
.kubernetes-container {
  width: 100%;
  height: 100%;
}
</style>
