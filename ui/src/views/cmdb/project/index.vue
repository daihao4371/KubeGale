<template>
  <ProjectView />
</template>

<script setup lang="ts">
import ProjectView from './project.vue'

defineOptions({
  name: 'ProjectIndexView'
})
</script>

<style scoped>
.project-container {
  padding: var(--spacing-lg); /* Use theme variable */
}
</style>
