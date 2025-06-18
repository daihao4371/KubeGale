<script setup lang="ts">
import { RouterView } from 'vue-router'
import AIAssistantEntry from '@/components/ai/AIAssistantEntry.vue'
import { StagewiseToolbar } from '@stagewise/toolbar-vue'
import { VuePlugin } from '@stagewise-plugins/vue'
import { onMounted } from 'vue'
import { usePermissionStore } from '@/pinia/modules/permission'
import { useUserStore } from '@/pinia/modules/user'

const isDev = import.meta.env.DEV
const permissionStore = usePermissionStore()
const userStore = useUserStore()

onMounted(() => {
  if (userStore.userInfo?.authorityId) {
    permissionStore.fetchApiList(userStore.userInfo.authorityId)
  }
})
</script>

<template>
  <RouterView />
  <AIAssistantEntry />
  <StagewiseToolbar v-if="isDev" :config="{ plugins: [VuePlugin] }" />
</template>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

#app {
  width: 100%;
  height: 100vh;
}
</style>
