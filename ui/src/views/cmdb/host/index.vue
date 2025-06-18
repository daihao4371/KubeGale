<template>
  <HostView />
</template>

<script setup lang="ts">
import HostView from './host.vue'

defineOptions({
  name: 'HostIndexView'
})
</script>

<style scoped>
.host-index {
  height: 100%;
}
</style>
