<template>
  <div class="operation-record-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="header-title">操作记录</span>
            <span class="header-subtitle">查看系统操作日志</span>
          </div>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :model="operationState.searchForm" inline class="search-form">
        <el-form-item label="请求方法">
          <el-select v-model="operationState.searchForm.method" placeholder="请选择请求方法" clearable style="width: 220px">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="请求路径">
          <el-input v-model="operationState.searchForm.path" placeholder="请输入请求路径" clearable style="width: 220px" />
        </el-form-item>
        <el-form-item label="状态码">
          <el-input v-model="operationState.searchForm.status" placeholder="请输入状态码" clearable style="width: 220px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 操作按钮 -->
      <div class="operation-buttons">
        <el-button type="danger" :disabled="!operationState.selectedRows.length" @click="handleBatchDelete">
          <el-icon><Delete /></el-icon> 批量删除
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table
        v-loading="operationState.loading"
        :data="operationState.tableData"
        border
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="操作人" width="150">
          <template #default="scope">
            <div :class="{ 'system-user': !scope.row.user || (!scope.row.user.username && !scope.row.user.nickname) }">
              <template v-if="scope.row.user && (scope.row.user.username || scope.row.user.userName) && (scope.row.user.nickname || scope.row.user.nickName) && (scope.row.user.username || scope.row.user.userName) !== (scope.row.user.nickname || scope.row.user.nickName)">
                <div class="user-info">
                  <span class="user-name">{{ (scope.row.user.username || scope.row.user.userName) }}</span>
                  <span class="user-real-name">{{ (scope.row.user.nickname || scope.row.user.nickName) }}</span>
                </div>
              </template>
              <template v-else>
                {{ formatUser(scope.row.user, scope.row) }}
              </template>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="CreatedAt" label="操作时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态码" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="请求IP" width="150" />
        <el-table-column prop="method" label="请求方法" width="100">
          <template #default="{ row }">
            <el-tag :type="getMethodType(row.method)">
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="请求路径" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">
              <el-icon><View /></el-icon> 详情
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              <el-icon><Delete /></el-icon> 删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="operationState.currentPage"
          v-model:page-size="operationState.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="operationState.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="operationState.detailDialogVisible"
      title="操作详情"
      width="800px"
      :close-on-click-modal="false"
      class="operation-detail-dialog"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="操作人">
          <span :class="{ 'system-user': !operationState.currentDetail?.user || (!operationState.currentDetail?.user.username && !operationState.currentDetail?.user.nickname) }">
            {{ operationState.currentDetail ? formatUser(operationState.currentDetail.user, operationState.currentDetail) : '-' }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="操作时间">
          {{ formatDate(operationState.currentDetail?.CreatedAt ?? '') }}
        </el-descriptions-item>
        <el-descriptions-item label="状态码">{{ operationState.currentDetail?.status }}</el-descriptions-item>
        <el-descriptions-item label="请求IP">{{ operationState.currentDetail?.ip }}</el-descriptions-item>
        <el-descriptions-item label="请求方法">{{ operationState.currentDetail?.method }}</el-descriptions-item>
        <el-descriptions-item label="请求路径">{{ operationState.currentDetail?.path }}</el-descriptions-item>
        <el-descriptions-item label="请求体" :span="2">
          <pre class="json-content">{{ formatJSON(operationState.currentDetail?.body) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应体" :span="2">
          <pre class="json-content">{{ formatJSON(operationState.currentDetail?.resp) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2">
          <pre>{{ operationState.currentDetail?.error_message }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { Delete, View } from '@element-plus/icons-vue'
import {
  operationState,
  fetchData,
  handleSearch,
  resetSearch,
  handleSelectionChange,
  handleDetail,
  handleDelete,
  handleBatchDelete,
  handleSizeChange,
  handleCurrentChange,
  getStatusType,
  getMethodType,
  formatJSON
} from './operationRecord'

// 修改 formatUser 函数的类型定义
function formatUser(user: Record<string, unknown> | undefined, row: Record<string, unknown>): string {
  let username = '';
  let realName = '';
  
  if (user) {
    // 获取用户名
    username = (user.username || user.userName || user.operator_name || '-') as string;
    // 获取真实姓名
    realName = (user.nickname || user.nickName || user.operator_real_name || '-') as string;
  } else {
    // 如果没有user对象，则尝试从row中获取
    username = (row.username || row.userName || row.operator_name || '-') as string;
    realName = (row.nickname || row.nickName || row.operator_real_name || '-') as string;
  }
  
  // 如果用户名和真实姓名都存在且不同，则同时显示
  if (username !== '-' && realName !== '-' && username !== realName) {
    return `${username} (${realName})`;
  } else if (username !== '-') {
    return username;
  } else {
    return realName;
  }
}

function formatDate(date: string) {
  if (!date) return '-'
  return date.replace('T', ' ').replace(/\..+/, '')
}

// 页面加载时获取数据
onMounted(() => {
  fetchData()
})
</script>

<style src="./operationRecord.css" scoped></style>
