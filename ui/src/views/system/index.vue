<template>
  <div class="system-container">
    <router-view />
  </div>
</template>

<script setup lang="ts">
// 使用 defineOptions 定义组件名称
defineOptions({
  name: 'SystemLayout'
})
</script>

<style scoped>
.system-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
  border: none; /* 确保没有边框 */
}
</style>
