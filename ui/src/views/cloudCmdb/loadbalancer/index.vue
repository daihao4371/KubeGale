<template>
  <LoadBalancerManager />
</template>

<script setup lang="ts">
import LoadBalancerManager from './LoadBalancerView.vue'

// 定义组件名称为多词名称
defineOptions({
  name: 'LoadBalancerManagerIndex'
})
</script>

<style>
.el-form-item__label {
  white-space: nowrap;
}
</style> 