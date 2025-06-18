<template>
  <div class="virtual-machine-container">
    <div class="platform-tree">
      <CloudPlatformTree @select="handlePlatformSelect" />
    </div>
    <div class="virtual-machine-list">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>云服务器管理</span>
            <div class="header-buttons">
              <el-button 
                type="primary" 
                :loading="syncing" 
                :disabled="!selectedPlatformId"
                @click="handleSync"
              >
                <el-icon><Refresh /></el-icon>同步数据
              </el-button>
            </div>
          </div>
        </template>

        <!-- 搜索栏 -->
        <el-form :inline="true" :model="searchForm" class="search-form" style="margin-bottom: 20px;">
          <el-form-item label="实例名称">
            <el-input v-model="searchForm.name" placeholder="请输入实例名称" clearable />
          </el-form-item>
          <el-form-item label="实例ID">
            <el-input v-model="searchForm.instanceId" placeholder="请输入实例ID" clearable />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSearch">
              搜索
            </el-button>
            <el-button @click="onResetSearch">
              重置
            </el-button>
          </el-form-item>
        </el-form>

        <div class="table-container">
          <el-table :data="tableData" v-loading="loading" height="100%" style="width: 100%">
            <el-table-column prop="instance_id" label="实例ID/名称" min-width="180">
              <template #default="scope">
                <div>
                  <div>{{ scope.row.instance_id }}</div>
                  <div class="vm-name">{{ scope.row.name }}</div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" min-width="100" />
            <el-table-column prop="os" label="操作系统" min-width="120" />
            <el-table-column prop="region_name" label="可用区" min-width="120" />
            <el-table-column label="配置" min-width="120">
              <template #default="scope">
                <span>{{ scope.row.cpu }}核/{{ (scope.row.memory / 1024).toFixed(2) }}GB</span>
              </template>
            </el-table-column>
            <el-table-column label="IP地址" min-width="160">
              <template #default="scope">
                <div>内网：{{ scope.row.private_addr }}</div>
                <div>公网：{{ scope.row.public_addr }}</div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="scope">
                <el-button
                  type="primary"
                  link
                  @click="handleViewDetails(scope.row)"
                >
                  <el-icon><View /></el-icon>
                  详情
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </el-card>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="虚拟机详情"
      width="60%"
      destroy-on-close
    >
      <template v-if="currentVM">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="实例ID">{{ currentVM.instance_id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ currentVM.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ currentVM.status }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{ currentVM.os }}</el-descriptions-item>
          <el-descriptions-item label="可用区">{{ currentVM.region_name }}</el-descriptions-item>
          <el-descriptions-item label="配置">
            {{ currentVM.cpu }}核/{{ (currentVM.memory / 1024).toFixed(2) }}GB
          </el-descriptions-item>
          <el-descriptions-item label="内网IP">{{ currentVM.private_addr }}</el-descriptions-item>
          <el-descriptions-item label="公网IP">{{ currentVM.public_addr }}</el-descriptions-item>
        </el-descriptions>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { virtualMachineList, syncMachine } from '@/api/cloudCmdb/cloud_virtual_machine'
import type { VirtualMachine } from '@/api/cloudCmdb/cloud_virtual_machine'
import CloudPlatformTree from '@/views/cloudCmdb/rds/components/CloudPlatformTree.vue'
import type { PlatformRegionTreeItem, RegionNode } from '@/api/cloudCmdb/cloud_rds'
import { ElMessage } from 'element-plus'
import { Refresh, View } from '@element-plus/icons-vue'
import { searchForm, handleSearch, handleResetSearch } from './virtualMachine'

const tableData = ref<VirtualMachine[]>([])
const loading = ref(false)
const syncing = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedPlatformId = ref<number | null>(null)
const dialogVisible = ref(false)
const currentVM = ref<VirtualMachine | null>(null)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await virtualMachineList({ 
      page: page.value, 
      pageSize: pageSize.value,
      name: searchForm.value.name,
      instance_id: searchForm.value.instanceId
    })
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

const handleSync = async () => {
  if (!selectedPlatformId.value) {
    ElMessage.warning('请先选择云平台')
    return
  }

  syncing.value = true
  try {
    const res = await syncMachine({ id: selectedPlatformId.value })
    if (res.code === 0) {
      ElMessage.success('同步成功')
      fetchData() // 同步成功后刷新数据
    } else {
      ElMessage.error(res.msg || '同步失败')
    }
  } catch {
    ElMessage.error('同步失败')
  } finally {
    syncing.value = false
  }
}

const handlePlatformSelect = (payload: { platform: PlatformRegionTreeItem, region: RegionNode | null }) => {
  selectedPlatformId.value = payload.platform.id
  // 如果选择了区域，可以在这里添加区域过滤逻辑
}

onMounted(fetchData)
watch([page, pageSize], fetchData)

const handleSizeChange = (val: number) => {
  pageSize.value = val
  page.value = 1
}

const handlePageChange = (val: number) => {
  page.value = val
}

const handleViewDetails = (row: VirtualMachine) => {
  currentVM.value = row
  dialogVisible.value = true
}

const onSearch = () => {
  page.value = 1
  handleSearch(fetchData)
}
const onResetSearch = () => {
  page.value = 1
  handleResetSearch(fetchData)
}
</script>

<style scoped>
@import './virtualMachine.css';
</style>
