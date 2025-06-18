<template>
  <div class="project-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="项目名称">
          <el-input v-model="searchForm.name" placeholder="请输入项目名称" clearable />
        </el-form-item>
        <el-form-item label="项目描述">
          <el-input v-model="searchForm.description" placeholder="请输入项目描述" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 操作按钮 -->
    <div class="action-bar">
      <div class="action-right">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>新建项目
        </el-button>
        <el-button type="danger" :disabled="selectedRows.length === 0" @click="handleBatchDelete">
          <el-icon><Delete /></el-icon>批量删除
        </el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      v-loading="loading"
      :data="tableData"
      border
      style="width: 100%"
      @selection-change="selectedRows = $event"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="项目名称" min-width="100" show-overflow-tooltip />
      <el-table-column prop="manager" label="项目负责人" width="120" />
      <el-table-column prop="description" label="项目描述" min-width="200" show-overflow-tooltip />
      <el-table-column prop="CreatedAt" label="创建时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.CreatedAt).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column prop="UpdatedAt" label="更新时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.UpdatedAt).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleDetail(row)">
            <el-icon><View /></el-icon>详情
          </el-button>
          <el-button type="primary" link @click="handleEdit(row)">
            <el-icon><Edit /></el-icon>编辑
          </el-button>
          <el-button type="danger" link @click="handleDelete(row)">
            <el-icon><Delete /></el-icon>删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination">
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

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="项目详情"
      width="600px"
      destroy-on-close
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="项目名称">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item label="项目负责人">{{ detailData.manager }}</el-descriptions-item>
        <el-descriptions-item label="项目描述">{{ detailData.description }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ new Date(detailData.CreatedAt).toLocaleString() }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ new Date(detailData.UpdatedAt).toLocaleString() }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 新建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="项目名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入项目名称" />
        </el-form-item>
        <el-form-item label="项目负责人" prop="manager">
          <el-input v-model="formData.manager" placeholder="请输入项目负责人" />
        </el-form-item>
        <el-form-item label="项目描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入项目描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import useProject from './project'
import { Plus, Delete, Edit, View } from '@element-plus/icons-vue'

defineOptions({
  name: 'ProjectView'
})

const {
  tableData,
  loading,
  searchForm,
  pageSize,
  currentPage,
  total,
  dialogVisible,
  dialogTitle,
  formData,
  rules,
  selectedRows,
  detailVisible,
  detailData,
  handleSearch,
  handleReset,
  handleAdd,
  handleEdit,
  handleDelete,
  handleBatchDelete,
  handleSubmit,
  handleCurrentChange,
  handleSizeChange,
  handleDetail
} = useProject()
</script>

<style src="./project.css" scoped></style>
