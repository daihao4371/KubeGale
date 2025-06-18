<template>
  <router-view></router-view>
</template>

<script setup lang="ts">
defineOptions({
  name: 'KubernetesView'
})
</script>

<style scoped>
.kubernetes-container {
  padding: var(--spacing-lg); /* Use theme variable */
}
</style>
