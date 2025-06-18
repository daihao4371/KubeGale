<template>
  <div>
    <button v-if="!visible" class="ai-entry-btn" @click="visible = true">
      <img src="@/assets/gpt.png" alt="AI助手" />
    </button>
    <AIAssistantDrawer v-model="visible" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AIAssistantDrawer from './AIAssistantDrawer.vue'
const visible = ref(false)
</script>

<style scoped>
.ai-entry-btn {
  position: fixed;
  right: 32px;
  bottom: 32px;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 4px 24px rgba(0,0,0,0.18);
  border: none;
  outline: none;
  cursor: pointer;
  z-index: 10000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  transition: box-shadow 0.2s, transform 0.2s;
}
.ai-entry-btn:hover {
  box-shadow: 0 8px 32px rgba(64,158,255,0.18);
  transform: scale(1.08);
}
.ai-entry-btn img {
  width: 40px;
  height: 40px;
  object-fit: contain;
  pointer-events: none;
}
</style> 