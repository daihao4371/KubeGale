<template>
  <ProviderManager />
</template>

<script setup lang="ts">
import ProviderManager from './provider.vue'

// 定义组件名称为多词名称
defineOptions({
  name: 'ProviderManagerIndex'
})
</script>

<style>
.el-form-item__label {
  white-space: nowrap;
}
</style> 