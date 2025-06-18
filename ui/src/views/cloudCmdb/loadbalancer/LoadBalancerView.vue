<template>
  <div class="loadbalancer-layout">
    <!-- 左侧树形目录 -->
    <div class="sidebar">
      <el-card>
        <div class="sidebar-title">资源分布</div>
        <CloudPlatformTree @select="handlePlatformSelect" />
      </el-card>
    </div>
    <!-- 右侧表格内容 -->
    <div class="main-content">
      <el-card class="main-card">
        <!-- 搜索栏 -->
        <el-form :inline="true" :model="searchInfo" class="search-form">
          <el-form-item label="识别搜索">
            <el-input v-model="searchInfo.keyword" placeholder="请输入搜索关键词" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSearch">查询</el-button>
            <el-button @click="onReset">重置</el-button>
            <el-button type="primary" @click="onSync">同步</el-button>
          </el-form-item>
        </el-form>
        <!-- 表格内容 -->
        <div class="table-wrapper">
          <el-table
            :data="tableData"
            v-loading="loading"
            class="table-container"
            :max-height="500"
            style="width: 100%;"
            border
          >
            <el-table-column type="selection" width="55" />
            <el-table-column prop="name" label="实例名称" min-width="180" show-overflow-tooltip />
            <el-table-column prop="instance_id" label="实例ID" min-width="280" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === '运行中' ? 'success' : 'danger'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="region_name" label="区域" width="120" />
            <el-table-column label="IP地址" width="200">
              <template #default="{ row }">
                <div v-if="row.private_addr">私网: {{ row.private_addr }}</div>
                <div v-if="row.public_addr">公网: {{ row.public_addr }}</div>
              </template>
            </el-table-column>
            <el-table-column prop="creation_time" label="创建时间" width="180">
              <template #default="{ row }">
                {{ new Date(row.creation_time).toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleDetail(row)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <!-- 分页条 -->
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
    v-model="detailDialogVisible"
    title="负载均衡器详情"
    width="800px"
    destroy-on-close
  >
    <el-descriptions :column="2" border>
      <el-descriptions-item label="实例名称" :span="2">{{ currentDetail?.name }}</el-descriptions-item>
      <el-descriptions-item label="实例ID" :span="2">{{ currentDetail?.instance_id }}</el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="currentDetail?.status === '运行中' ? 'success' : 'danger'">
          {{ currentDetail?.status }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="区域">{{ currentDetail?.region_name }}</el-descriptions-item>
      <el-descriptions-item label="私网IP">{{ currentDetail?.private_addr }}</el-descriptions-item>
      <el-descriptions-item label="公网IP">{{ currentDetail?.public_addr }}</el-descriptions-item>
      <el-descriptions-item label="创建时间" :span="2">
        {{ currentDetail?.creation_time ? new Date(currentDetail.creation_time).toLocaleString() : '-' }}
      </el-descriptions-item>
    </el-descriptions>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import './LoadBalancerView.css'
import useLoadBalancer from './LoadBalancerView'
import CloudPlatformTree from '../rds/components/CloudPlatformTree.vue'

const {
  loading,
  searchInfo,
  pageSize,
  currentPage,
  total,
  tableData,
  detailDialogVisible,
  currentDetail,
  init,
  handlePlatformSelect,
  onSearch,
  onReset,
  onSync,
  handleCurrentChange,
  handleSizeChange,
  handleDetail
} = useLoadBalancer()

// 初始化
onMounted(() => {
  init()
})
</script> 