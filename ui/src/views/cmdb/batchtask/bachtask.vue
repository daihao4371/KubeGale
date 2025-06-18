<template>
  <div class="batch-task-root">
    <el-card class="main-card">
      <div class="main-title">批量任务</div>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="命令执行" name="execute">
          <el-card class="section-card">
            <div class="section-title">批量命令执行</div>
            <!-- 1. 选择目标主机 -->
            <div class="sub-section">
              <div class="sub-title">1. 选择目标主机</div>
              <el-table 
                v-loading="loading"
                :data="hosts" 
                style="width: 100%" 
                border
                @selection-change="handleSelectionChange"
              >
                <el-table-column type="selection" width="55" />
                <el-table-column prop="name" label="主机名称" />
                <el-table-column prop="publicIP" label="公网IP" />
                <el-table-column prop="privateIP" label="私网IP" />
                <el-table-column prop="os" label="操作系统" />
                <el-table-column prop="osVersion" label="系统版本" />
                <el-table-column prop="status" label="状态" />
              </el-table>
              <div class="table-footer">已选择 {{ selectedHosts.length }} 台主机</div>
              <el-pagination
                background
                layout="total, sizes, prev, pager, next, jumper"
                :total="total"
                :page-size="pageSize"
                :current-page="currentPage"
                :page-sizes="[10, 20, 50]"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                style="margin-top: 10px;"
              />
            </div>
            <!-- 2. 输入执行命令 -->
            <div class="sub-section">
              <div class="sub-title">2. 输入执行命令</div>
              <el-radio-group v-model="execType" class="exec-type-radio">
                <el-radio label="Shell">Shell</el-radio>
                <el-radio label="Python">Python</el-radio>
              </el-radio-group>
              <el-input
                type="textarea"
                :rows="4"
                v-model="command"
                placeholder="请输入要执行的命令，例如：df -h|free -m"
                class="command-input"
              />
            </div>
            <!-- 3. 执行与结果 -->
            <div class="sub-section">
              <div class="sub-title">3. 执行与结果</div>
              <el-button type="primary" :loading="executing" @click="executeCommand">执行命令</el-button>
              
              <!-- 执行结果 -->
              <div v-if="executionResult" class="execution-result">
                <div class="result-header">
                  <div class="result-item">
                    <span class="label">执行状态：</span>
                    <el-tag :type="executionResult.status === 'success' ? 'success' : 'danger'">
                      {{ executionResult.status === 'success' ? '成功' : '失败' }}
                    </el-tag>
                  </div>
                  <div class="result-item">
                    <span class="label">总主机数：</span>
                    <span>{{ executionResult.allHosts.length }}</span>
                  </div>
                  <div class="result-item">
                    <span class="label">成功主机数：</span>
                    <span>{{ executionResult.successHosts?.length || 0 }}</span>
                  </div>
                  <div class="result-item">
                    <span class="label">失败主机数：</span>
                    <span>{{ executionResult.failureHosts?.length || 0 }}</span>
                  </div>
                </div>
                
                <!-- 执行日志 -->
                <div class="execution-logs">
                  <div v-for="log in executionResult.executionLogs" :key="log.host" class="log-item">
                    <div class="log-header">
                      <span class="host">{{ log.host }}</span>
                      <el-tag size="small" type="success" v-if="executionResult.successHosts?.includes(log.host)">成功</el-tag>
                      <el-tag size="small" type="danger" v-else-if="executionResult.failureHosts?.includes(log.host)">失败</el-tag>
                    </div>
                    <pre class="log-content">{{ log.output || '无输出' }}</pre>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-tab-pane>
        <el-tab-pane label="执行日志" name="log">
          <el-card class="section-card">
            <div class="section-title">批量操作记录</div>
            <div class="log-filter-bar">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                style="margin-right: 8px;"
              />
              <el-input v-model="filterUser" placeholder="执行用户" style="width: 150px; margin-right: 8px;" />
              <el-input v-model="filterUserId" placeholder="用户ID或名称" style="width: 180px; margin-right: 8px;" />
              <el-select v-model="filterStatus" placeholder="执行状态" style="width: 120px; margin-right: 8px;">
                <el-option label="全部" value="" />
                <el-option label="成功" value="success" />
                <el-option label="失败" value="failed" />
              </el-select>
              <el-button type="primary" @click="fetchExecutionLogs">查询</el-button>
              <el-button @click="resetFilters">重置</el-button>
            </div>
            <el-table 
              v-loading="logsLoading"
              :data="executionLogs" 
              style="width: 100%; margin-top: 10px;" 
              border
            >
              <el-table-column prop="ID" label="ID" width="60" />
              <el-table-column label="执行用户" width="120">
                <template #default="{ row }">
                  {{ userMap[row.UserId] || row.UserId }}
                </template>
              </el-table-column>
              <el-table-column prop="Command" label="命令" show-overflow-tooltip />
              <el-table-column label="执行时间" width="180">
                <template #default="{ row }">
                  {{ formatTime(row.CreatedAt) }}
                </template>
              </el-table-column>
              <el-table-column label="执行状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.Status === 'success' ? 'success' : 'danger'">
                    {{ row.Status === 'success' ? '成功' : '失败' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120">
                <template #default="{ row }">
                  <el-button 
                    type="primary" 
                    link
                    @click="showExecutionDetail(row)"
                  >
                    查看详情
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="logsTotal"
              :page-size="logsPageSize"
              :current-page="logsCurrentPage"
              :page-sizes="[10, 20, 50]"
              @size-change="handleLogsSizeChange"
              @current-change="handleLogsCurrentChange"
              style="margin-top: 10px;"
            />
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 执行详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="执行详情"
      width="800px"
      destroy-on-close
    >
      <div v-if="currentDetail" class="execution-detail">
        <div class="detail-header">
          <div class="detail-item">
            <span class="label">命令：</span>
            <span class="value">{{ currentDetail.Command }}</span>
          </div>
          <div class="detail-item">
            <span class="label">执行时间：</span>
            <span class="value">{{ formatTime(currentDetail.CreatedAt) }}</span>
          </div>
          <div class="detail-item">
            <span class="label">执行状态：</span>
            <el-tag :type="currentDetail.Status === 'success' ? 'success' : 'danger'">
              {{ currentDetail.Status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </div>
        </div>
        <div class="detail-logs">
          <div v-for="log in parseExecutionLogs(currentDetail.ExecutionLogs)" :key="log.host" class="log-item">
            <div class="log-header">
              <span class="host">{{ log.host }}</span>
              <el-tag 
                size="small" 
                :type="currentDetail.SuccessHosts.includes(log.host) ? 'success' : 'danger'"
              >
                {{ currentDetail.SuccessHosts.includes(log.host) ? '成功' : '失败' }}
              </el-tag>
            </div>
            <pre class="log-content">{{ log.output || '无输出' }}</pre>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getHostList } from '@/api/cmdb/host'
import { batchExecuteCommand, getExecutionLogs } from '@/api/cmdb/batchtask'
import { getUserList } from '@/api/system/user'
import type { Host } from '@/types/cmdb'
import type { BatchExecuteCommandRes, ExecutionLog } from '@/api/cmdb/batchtask'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { UserInfo } from '@/types/system'

const activeTab = ref('execute')
const execType = ref<'shell' | 'python'>('shell')
const command = ref('')
const hosts = ref<Host[]>([])
const selectedHosts = ref<Host[]>([])
const dateRange = ref<[Date, Date] | null>(null)
const filterUser = ref('')
const filterUserId = ref('')
const filterStatus = ref('')
const loading = ref(false)
const executing = ref(false)
const total = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)
const executionResult = ref<BatchExecuteCommandRes | null>(null)
const executionLogs = ref<ExecutionLog[]>([])
const logsLoading = ref(false)
const logsTotal = ref(0)
const logsPageSize = ref(10)
const logsCurrentPage = ref(1)
const detailDialogVisible = ref(false)
const currentDetail = ref<ExecutionLog | null>(null)
const userMap = ref<Record<number, string>>({})

// 获取主机列表
const fetchHostList = async () => {
  try {
    loading.value = true
    const res = await getHostList({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    hosts.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取主机列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取执行日志列表
const fetchExecutionLogs = async () => {
  try {
    logsLoading.value = true
    const res = await getExecutionLogs({
      page: logsCurrentPage.value,
      pageSize: logsPageSize.value
    })
    executionLogs.value = res.data
    logsTotal.value = res.data.length
  } catch (error) {
    console.error('获取执行日志失败:', error)
    ElMessage.error('获取执行日志失败')
  } finally {
    logsLoading.value = false
  }
}

// 解析执行日志
const parseExecutionLogs = (logs: string) => {
  try {
    return JSON.parse(logs)
  } catch (error) {
    console.error('解析执行日志失败:', error)
    return []
  }
}

// 格式化时间
const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

// 执行命令
const executeCommand = async () => {
  if (!command.value.trim()) {
    ElMessage.warning('请输入要执行的命令')
    return
  }
  if (selectedHosts.value.length === 0) {
    ElMessage.warning('请选择要执行命令的主机')
    return
  }

  try {
    executing.value = true
    const res = await batchExecuteCommand({
      hosts: selectedHosts.value.map(host => host.publicIP),
      users: selectedHosts.value.map(() => 'root'),
      ports: selectedHosts.value.map(() => 22),
      command: command.value,
      language: execType.value
    })
    executionResult.value = res.data
    ElMessage.success('命令执行成功')
  } catch (error) {
    console.error('执行命令失败:', error)
    ElMessage.error('执行命令失败')
  } finally {
    executing.value = false
  }
}

const handleSelectionChange = (selected: Host[]) => {
  selectedHosts.value = selected
}

const handleSizeChange = (newSize: number) => {
  pageSize.value = newSize
  fetchHostList()
}

const handleCurrentChange = (newPage: number) => {
  currentPage.value = newPage
  fetchHostList()
}

// 处理日志分页
const handleLogsSizeChange = (newSize: number) => {
  logsPageSize.value = newSize
  fetchExecutionLogs()
}

const handleLogsCurrentChange = (newPage: number) => {
  logsCurrentPage.value = newPage
  fetchExecutionLogs()
}

// 显示执行详情
const showExecutionDetail = (log: ExecutionLog) => {
  currentDetail.value = log
  detailDialogVisible.value = true
}

// 重置过滤器
const resetFilters = () => {
  dateRange.value = null
  filterUser.value = ''
  filterUserId.value = ''
  filterStatus.value = ''
  fetchExecutionLogs()
}

const fetchAllUsers = async () => {
  let page = 1
  const pageSize = 100
  let hasMore = true
  const map: Record<number, string> = {}

  while (hasMore) {
    const res = await getUserList({ page, pageSize })
    if (res.data.code === 0 && res.data.data?.list?.length) {
      res.data.data.list.forEach((u: UserInfo) => {
        map[u.ID] = u.nickName || u.userName || String(u.ID)
      })
      hasMore = res.data.data.list.length === pageSize
      page++
    } else {
      hasMore = false
    }
  }
  userMap.value = map
}

onMounted(() => {
  fetchHostList()
  fetchExecutionLogs()
  fetchAllUsers()
})

defineOptions({ name: 'BatchTaskPage' })
</script>

<style scoped>
.batch-task-root {
  padding: 24px;
  background: #f7f8fa;
  min-height: 100vh;
}
.main-card {
  max-width: 1200px;
  margin: 0 auto;
  background: #fff;
}
.main-title {
  font-size: 22px;
  font-weight: bold;
  margin-bottom: 18px;
}
.section-card {
  margin-bottom: 18px;
  background: #fafbfc;
}
.section-title {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 12px;
}
.sub-section {
  margin-bottom: 24px;
  background: #fff;
  border-radius: 6px;
  padding: 18px 18px 8px 18px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.03);
}
.sub-title {
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 10px;
}
.table-footer {
  margin-top: 8px;
  color: #888;
  font-size: 13px;
}
.exec-type-radio {
  margin-bottom: 10px;
}
.command-input {
  width: 100%;
  margin-bottom: 10px;
}
.log-filter-bar {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  flex-wrap: wrap;
}
.execution-result {
  margin-top: 20px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 16px;
}
.result-header {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}
.result-item {
  display: flex;
  align-items: center;
  gap: 8px;
}
.label {
  color: #606266;
}
.execution-logs {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.log-item {
  border: 1px solid #ebeef5;
  border-radius: 4px;
}
.log-header {
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.host {
  font-weight: 500;
  color: #303133;
  font-family: monospace;
}
.log-content {
  margin: 0;
  padding: 12px;
  background-color: #1e1e1e;
  color: #fff;
  font-family: monospace;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 300px;
  overflow-y: auto;
  font-size: 13px;
  line-height: 1.5;
}
.execution-detail {
  padding: 16px;
}
.detail-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}
.detail-item {
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}
.detail-item .label {
  width: 80px;
  color: #606266;
}
.detail-item .value {
  flex: 1;
  word-break: break-all;
}
.detail-logs {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
