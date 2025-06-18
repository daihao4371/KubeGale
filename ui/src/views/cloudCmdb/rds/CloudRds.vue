<template>
  <div class="cloud-rds-container">
    <div class="platform-tree">
      <CloudPlatformTree @select="handlePlatformSelect" />
    </div>
    <div class="cloud-rds-list">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>云数据库管理</span>
            <div class="header-buttons">
              <el-button 
                type="primary" 
                :loading="loading" 
                :disabled="!selectedPlatformId"
                @click="handleSync"
              >
                <el-icon><Refresh /></el-icon>同步数据
              </el-button>
            </div>
          </div>
        </template>

        <!-- 搜索栏 -->
        <el-form :inline="true" :model="searchForm" class="search-form">
            <el-form-item label="实例名称">
              <el-input v-model="searchForm.name" placeholder="请输入实例名称" clearable />
            </el-form-item>
            <el-form-item label="实例ID">
              <el-input v-model="searchForm.instanceId" placeholder="请输入实例ID" clearable />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch">
                <el-icon><Search /></el-icon>搜索
              </el-button>
              <el-button @click="handleResetSearch">
                <el-icon><RefreshRight /></el-icon>重置
              </el-button>
            </el-form-item>
          </el-form>
          
          <div class="table-container">
            <el-table
              :data="tableData"
              style="width: 100%"
              v-loading="loading"
              border
            >
              <el-table-column prop="name" label="实例名称" min-width="150" />
              <el-table-column prop="instance_id" label="实例ID" min-width="150" />
              <el-table-column prop="private_addr" label="私有地址" min-width="200" />
              <el-table-column prop="public_addr" label="公网地址" min-width="200">
                <template #default="{ row }">
                  {{ row.public_addr || '-' }}
                </template>
              </el-table-column>
              <el-table-column prop="region_name" label="区域" width="120" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="getStatusType(row.status)">{{ row.status }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="creation_time" label="创建时间" width="180">
                <template #default="{ row }">
                  {{ formatDateTime(row.creation_time) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="handleView(row)">
                    <el-icon><View /></el-icon>详情
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
          
          <!-- 分页组件 -->
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
      </div>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="实例详情"
      width="800px"
    >
      <el-descriptions
        v-if="currentInstance"
        :column="2"
        border
      >
        <el-descriptions-item label="实例名称" :span="2">
          {{ currentInstance.name }}
        </el-descriptions-item>
        <el-descriptions-item label="实例ID" :span="2">
          {{ currentInstance.instance_id }}
        </el-descriptions-item>
        <el-descriptions-item label="私有地址" :span="2">
          {{ currentInstance.private_addr }}
        </el-descriptions-item>
        <el-descriptions-item label="公网地址" :span="2">
          {{ currentInstance.public_addr || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="区域">
          {{ currentInstance.region_name }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentInstance.status)">
            {{ currentInstance.status }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">
          {{ formatDateTime(currentInstance.creation_time) }}
        </el-descriptions-item>
        <el-descriptions-item label="云平台" :span="2">
          {{ currentInstance.cloud_platform?.name || '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { 
  View,
  Search, 
  Refresh, 
  RefreshRight
} from '@element-plus/icons-vue'
import CloudPlatformTree from './components/CloudPlatformTree.vue'
import './CloudRds.css'
import {
  loading,
  currentPage,
  total,
  pageSize,
  tableData,
  detailVisible,
  currentInstance,
  searchForm,
  formatDateTime,
  getStatusType,
  loadData,
  handleSearch,
  handleResetSearch,
  handleSync,
  handlePlatformSelect,
  handleView,
  handleSizeChange,
  handleCurrentChange,
  selectedPlatformId
} from './CloudRds'

defineOptions({
  name: 'RdsManagerView'
})

// 页面加载时获取数据
onMounted(() => {
  loadData()
})
</script>
