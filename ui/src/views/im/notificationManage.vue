<template>
  <div class="notification-manage">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>通知管理</span>
          <div class="header-buttons">
            <el-button type="primary" :icon="Plus" @click="handleAddFeishu">添加飞书机器人</el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-form">
        <el-input
          v-model="searchName"
          placeholder="请输入通知名称搜索"
          class="search-input"
          :prefix-icon="Search"
          clearable
          @clear="handleReset"
        />
        <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
        <el-button :icon="RefreshRight" @click="handleReset">重置</el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        border
      >
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag type="warning">
              {{ row.type === 'feishu' ? '飞书' : row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="notification_policy" label="通知策略" width="200" />
        <el-table-column prop="robot_url" label="Webhook地址" width="300" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              @click="handleDelete(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 飞书机器人对话框 -->
    <el-dialog
      v-model="feishuDialogVisible"
      :title="isEdit ? '编辑飞书机器人' : '添加飞书机器人'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="feishuFormRef"
        :model="feishuForm"
        :rules="feishuFormRules"
        label-width="120px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="feishuForm.name" placeholder="请输入机器人名称" />
        </el-form-item>
        <el-form-item label="Webhook地址" prop="webhook_url">
          <el-input v-model="feishuForm.webhook_url" placeholder="请输入飞书机器人的Webhook地址" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="feishuForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述信息"
          />
        </el-form-item>
        <el-form-item label="标签" prop="tags">
          <el-select
            v-model="feishuForm.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或输入标签"
          >
            <el-option label="告警" value="告警" />
            <el-option label="监控" value="监控" />
            <el-option label="通知" value="通知" />
          </el-select>
        </el-form-item>
        <el-form-item label="通知事件" prop="notify_events">
          <el-select
            v-model="feishuForm.notify_events"
            multiple
            placeholder="请选择通知事件"
          >
            <el-option label="告警" value="alert" />
            <el-option label="警告" value="warning" />
          </el-select>
        </el-form-item>
        <el-form-item label="接收者" prop="receivers">
          <el-select
            v-model="feishuForm.receivers"
            multiple
            filterable
            placeholder="请选择接收者"
          >
            <el-option
              v-for="user in userOptions"
              :key="user.value"
              :label="user.label"
              :value="user.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="每日统计" prop="send_daily_stats">
          <el-switch v-model="feishuForm.send_daily_stats" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="feishuDialogVisible = false">
            <el-icon><Close /></el-icon>取消
          </el-button>
          <el-button type="primary" :loading="formLoading" @click="submitFeishuForm(feishuFormRef)">
            <el-icon><Check /></el-icon>确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Plus, Search, RefreshRight, Close, Check } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import dayjs from 'dayjs'
import {
  searchName,
  handleSearch,
  handleReset,
  tableData,
  handleAddFeishu,
  feishuDialogVisible,
  feishuForm,
  feishuFormRules,
  submitFeishuForm,
  loading,
  formLoading,
  currentPage,
  pageSize,
  total,
  handleSizeChange,
  handleCurrentChange,
  userOptions,
  handleDelete,
  handleEdit,
  isEdit
} from './notificationManage'

const feishuFormRef = ref<FormInstance>()

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}
</script>

<style scoped>
@import './notificationManage.css';
</style>
