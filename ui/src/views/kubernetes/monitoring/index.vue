<template>
  <div class="monitoring-container">
    <div class="page-header">
      <h2>监控</h2>
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>Kubernetes</el-breadcrumb-item>
        <el-breadcrumb-item>监控</el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    
    <el-card>
      <div class="coming-soon">
        <el-empty description="监控功能正在开发中...">
          <el-button type="primary">敬请期待</el-button>
        </el-empty>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'MonitoringIndex'
})
</script>

<style scoped>
.monitoring-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 10px 0;
  color: var(--el-text-color-primary);
}

.coming-soon {
  text-align: center;
  padding: 40px 0;
}
</style> 