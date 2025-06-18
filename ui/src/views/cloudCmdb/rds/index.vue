<template>
  <RdsManager />
</template>

<script setup lang="ts">
import RdsManager from './CloudRds.vue'

// 定义组件名称为多词名称
defineOptions({
  name: 'RdsManagerIndex'
})
</script>

<style>
.el-form-item__label {
  white-space: nowrap;
}
</style>
