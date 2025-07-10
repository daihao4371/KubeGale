/**
 * 监控指标数据处理工具
 */

export interface MetricData {
  data: {
    result: Array<{
      metric: Record<string, string>
      values: Array<[number, string]>
    }>
  }
  status: string
  warnings?: string[]
}

/**
 * 标准化监控指标数据
 * @param metric 原始指标数据
 * @returns 标准化后的指标数据
 */
export function normalizeMetrics(metric: any): MetricData {
  if (!metric) {
    return {
      data: { result: [] },
      status: 'error'
    }
  }

  // 如果数据已经是标准格式，直接返回
  if (metric.data && metric.data.result) {
    return metric
  }

  // 处理不同格式的监控数据
  const result = []

  if (Array.isArray(metric)) {
    // 如果是数组格式
    result.push({
      metric: {},
      values: metric.map((item, index) => [Date.now() + index * 60000, String(item)])
    })
  } else if (typeof metric === 'object') {
    // 如果是对象格式，尝试提取值
    const values = Object.entries(metric).map(([key, value]) => [
      parseInt(key) || Date.now(),
      String(value)
    ])
    result.push({
      metric: {},
      values
    })
  }

  return {
    data: { result },
    status: 'success'
  }
}

/**
 * 格式化指标值
 * @param value 指标值
 * @param unit 单位
 * @returns 格式化后的字符串
 */
export function formatMetricValue(value: string | number, unit?: string): string {
  const numValue = typeof value === 'string' ? parseFloat(value) : value
  
  if (isNaN(numValue)) return '-'
  
  // 根据单位进行格式化
  switch (unit) {
    case 'bytes':
      return formatBytes(numValue)
    case 'cpu':
      return formatCPU(numValue)
    case 'memory':
      return formatMemory(numValue)
    default:
      return numValue.toFixed(2)
  }
}

/**
 * 格式化字节数
 */
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 格式化CPU使用率
 */
export function formatCPU(cpu: number): string {
  if (cpu < 1) {
    return `${(cpu * 1000).toFixed(0)}m`
  }
  return `${cpu.toFixed(2)}`
}

/**
 * 格式化内存使用
 */
export function formatMemory(memory: number): string {
  return formatBytes(memory)
}

/**
 * 计算使用率百分比
 */
export function calculateUsagePercentage(used: number, total: number): number {
  if (total === 0) return 0
  return Math.round((used / total) * 100)
}

/**
 * 获取指标趋势
 */
export function getMetricTrend(values: Array<[number, string]>): 'up' | 'down' | 'stable' {
  if (values.length < 2) return 'stable'
  
  const latest = parseFloat(values[values.length - 1][1])
  const previous = parseFloat(values[values.length - 2][1])
  
  if (latest > previous * 1.05) return 'up'
  if (latest < previous * 0.95) return 'down'
  return 'stable'
}