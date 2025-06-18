<template>
  <div class="api-manager-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>API管理</span>
          <div class="header-buttons">
            <el-button type="success" @click="apiManager.handleFreshCasbin">
              <el-icon><Refresh /></el-icon>刷新缓存
            </el-button>
            <el-button type="primary" @click="apiManager.handleAddApi">
              <el-icon><Plus /></el-icon>添加API
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 搜索栏 -->
      <el-form :inline="true" :model="apiManager.state.searchForm" class="search-form">
        <el-form-item label="API路径">
          <el-input v-model="apiManager.state.searchForm.path" placeholder="请输入API路径" clearable />
        </el-form-item>
        <el-form-item label="API简介">
          <el-input v-model="apiManager.state.searchForm.description" placeholder="请输入API简介" clearable />
        </el-form-item>
        <el-form-item label="API分组">
          <el-select v-model="apiManager.state.searchForm.apiGroup" placeholder="请选择API分组" clearable style="width: 160px">
            <el-option
              v-for="group in apiManager.state.apiGroups"
              :key="group"
              :label="group"
              :value="group"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="请求方式">
          <el-select v-model="apiManager.state.searchForm.method" placeholder="请选择请求方式" clearable style="width: 160px">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="apiManager.handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="apiManager.handleResetSearch">
            <el-icon><RefreshRight /></el-icon>重置
          </el-button>
        </el-form-item>
      </el-form>
      
      <el-table
        :data="apiManager.state.filteredApiList.slice((apiManager.state.currentPage - 1) * apiManager.state.pageSize, apiManager.state.currentPage * apiManager.state.pageSize)"
        style="width: 100%"
        v-loading="apiManager.state.loading"
        border
      >
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="path" label="API路径" min-width="200" />
        <el-table-column prop="apiGroup" label="API分组" min-width="120" />
        <el-table-column prop="description" label="API简介" min-width="200" />
        <el-table-column prop="method" label="请求方式" width="100">
          <template #default="{ row }">
            <el-tag :type="apiManager.getMethodType(row.method)">{{ row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="apiManager.handleEdit(row)">
              <el-icon><Edit /></el-icon>编辑
            </el-button>
            <el-button type="danger" link @click="apiManager.handleDelete(row)">
              <el-icon><Delete /></el-icon>删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="apiManager.state.currentPage"
          v-model:page-size="apiManager.state.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="apiManager.state.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="apiManager.handleSizeChange"
          @current-change="apiManager.handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑API对话框 -->
    <el-dialog
      v-model="apiManager.state.dialogVisible"
      :title="apiManager.state.dialogTitle"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="apiManager.form.value"
        :rules="rules"
        label-width="100px"
        v-loading="apiManager.state.formLoading"
      >
        <el-form-item label="API路径" prop="path">
          <el-input v-model="apiManager.form.value.path" placeholder="请输入API路径" />
        </el-form-item>
        <el-form-item label="API分组" prop="apiGroup">
          <el-select
            v-model="apiManager.form.value.apiGroup"
            placeholder="请选择或输入API分组"
            style="width: 100%"
            filterable
            allow-create
            default-first-option
          >
            <el-option
              v-for="group in apiManager.state.apiGroups"
              :key="group"
              :label="group"
              :value="group"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="请求方式" prop="method">
          <el-select v-model="apiManager.form.value.method" placeholder="请选择请求方式" style="width: 100%">
            <el-option label="创建(POST请求)" value="POST" />
            <el-option label="查看(GET)" value="GET" />
            <el-option label="更新(PUT)" value="PUT" />
            <el-option label="删除(DELETE)" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="API简介" prop="description">
          <el-input
            v-model="apiManager.form.value.description"
            type="textarea"
            :rows="3"
            placeholder="请输入API简介"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="apiManager.state.dialogVisible = false">
            <el-icon><Close /></el-icon>取消
          </el-button>
          <el-button type="primary" @click="apiManager.handleSubmit">
            <el-icon><Check /></el-icon>确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import useApiManager from './apiManager'
import { 
  Plus, 
  Edit, 
  Delete, 
  Search, 
  Refresh, 
  RefreshRight, 
  Check, 
  Close 
} from '@element-plus/icons-vue'

const formRef = ref<FormInstance>()
const rules: FormRules = {
  path: [
    { required: true, message: '请输入API路径', trigger: 'blur' },
    { pattern: /^\/[a-zA-Z0-9-_/]+$/, message: 'API路径格式不正确', trigger: 'blur' }
  ],
  apiGroup: [
    { required: true, message: '请选择API分组', trigger: 'change' }
  ],
  method: [
    { required: true, message: '请选择请求方式', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入API简介', trigger: 'blur' }
  ]
}

const apiManager = useApiManager()

onMounted(() => {
  apiManager.loadApiList()
})
</script>

<style src="./apiManager.css" scoped></style>