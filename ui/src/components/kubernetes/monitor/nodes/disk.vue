<template>
  <div class="monitor-disk">
    <div class="chart-container">
      <div ref="chartRef" class="chart" style="width: 100%; height: 400px;"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

interface Props {
  data: any
}

const props = defineProps<Props>()
const chartRef = ref<HTMLElement>()

onMounted(() => {
  if (chartRef.value) {
    chartRef.value.innerHTML = `
      <div style="display: flex; align-items: center; justify-content: center; height: 100%; background: #f5f7fa; border: 1px dashed #ddd; color: #666;">
        <div style="text-align: center;">
          <h3>磁盘监控图表</h3>
          <p>请集成 ECharts 或其他图表库来显示监控数据</p>
          <p>数据点数量: ${props.data ? Object.keys(props.data).length : 0}</p>
        </div>
      </div>
    `
  }
})

watch(() => props.data, () => {
  console.log('磁盘监控数据更新:', props.data)
})
</script>

<style scoped>
.monitor-disk {
  width: 100%;
  height: 100%;
}

.chart-container {
  padding: 20px;
}

.chart {
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>