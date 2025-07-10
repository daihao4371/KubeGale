<template>
  <div class="page-layout">
    <!-- 页面头部 -->
    <header class="page-header">
      <div class="header-left">
        <div v-if="$slots.header" class="page-title-slot">
          <slot name="header" />
        </div>
        <h1 v-else-if="title" class="page-title">{{ title }}</h1>
        <div v-if="$slots.subtitle" class="page-subtitle">
          <slot name="subtitle" />
        </div>
      </div>
      <div class="header-right">
        <slot name="actions" />
      </div>
    </header>

    <!-- 搜索区域 -->
    <div v-if="$slots.search" class="search-section">
      <slot name="search" />
    </div>

    <!-- 主要内容 -->
    <main class="page-content">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
}

defineProps<Props>()
</script>

<style scoped>
.page-layout {
  padding: 20px;
  min-height: 100vh;
  background: #f5f7fa;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left {
  flex: 1;
}

.page-title-slot h1,
.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  line-height: 1.2;
}

.page-subtitle {
  margin-top: 8px;
  color: #6c757d;
  font-size: 14px;
  line-height: 1.4;
}

.header-right {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-shrink: 0;
}

.search-section {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.page-content {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-layout {
    padding: 10px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 15px;
  }
  
  .header-right {
    width: 100%;
    justify-content: flex-end;
  }
  
  .page-title-slot h1,
  .page-title {
    font-size: 20px;
  }
  
  .search-section {
    padding: 15px;
  }
}
</style>