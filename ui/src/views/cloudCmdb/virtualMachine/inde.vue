<template>
  <VirtualMachineManager />
</template>

<script setup lang="ts">
import VirtualMachineManager from './virtualMachine.vue'

defineOptions({
  name: 'VirtualMachineManagerIndex'
})
</script>

<style>
.el-form-item__label {
  white-space: nowrap;
}
</style>
