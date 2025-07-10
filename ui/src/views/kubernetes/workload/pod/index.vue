<template>
  <router-view></router-view>
</template>

<script setup lang="ts">
defineOptions({
  name: 'PodView'
})
</script>

<style scoped>
.pod-container {
  padding: var(--spacing-lg);
}
</style> 