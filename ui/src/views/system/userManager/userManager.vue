<template>
  <div class="user-manager-container">
    <el-card class="box-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="header-title">用户管理</span>
            <span class="header-subtitle">管理系统用户信息和权限</span>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="fetchUserInfo">
              <el-icon class="el-icon--left"><User /></el-icon>个人信息
            </el-button>
            <el-button type="primary" @click="handleChangePassword">
              <el-icon class="el-icon--left"><Lock /></el-icon>修改密码
            </el-button>
            <el-button type="primary" @click="handleAddUser">
              <el-icon class="el-icon--left"><Plus /></el-icon>添加用户
            </el-button>
          </div>
        </div>
      </template>

      <!-- 用户统计卡片 -->
      <div class="stats-cards">
        <el-card class="stats-card" shadow="hover">
          <template #header>
            <div class="stats-card-header">
              <el-icon><User /></el-icon>
              <span>总用户数</span>
            </div>
          </template>
          <div class="stats-card-content">
            <span class="stats-number">{{ userStats.total }}</span>
            <span class="stats-label">人</span>
          </div>
        </el-card>
        <el-card class="stats-card" shadow="hover">
          <template #header>
            <div class="stats-card-header">
              <el-icon><CircleCheck /></el-icon>
              <span>启用用户</span>
            </div>
          </template>
          <div class="stats-card-content">
            <span class="stats-number">{{ userStats.enabled }}</span>
            <span class="stats-label">人</span>
          </div>
        </el-card>
        <el-card class="stats-card" shadow="hover">
          <template #header>
            <div class="stats-card-header">
              <el-icon><CircleClose /></el-icon>
              <span>禁用用户</span>
            </div>
          </template>
          <div class="stats-card-content">
            <span class="stats-number">{{ userStats.disabled }}</span>
            <span class="stats-label">人</span>
          </div>
        </el-card>
      </div>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="searchForm.userName" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchForm.nickName" placeholder="请输入昵称" clearable />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="searchForm.phone" placeholder="请输入手机号" clearable />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="searchForm.email" placeholder="请输入邮箱" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchUserList">
            <el-icon class="el-icon--left"><Search /></el-icon>搜索
          </el-button>
          <el-button @click="resetSearchForm">
            <el-icon class="el-icon--left"><Refresh /></el-icon>重置
          </el-button>
        </el-form-item>
      </el-form>

      <el-table
        :data="userList"
        style="width: 100%"
        v-loading="loading"
        border
        stripe
        highlight-current-row
        :header-cell-style="{backgroundColor: '#f5f7fa', color: '#606266'}"
      >
        <el-table-column prop="ID" label="ID" min-width="200" show-overflow-tooltip />
        <el-table-column prop="userName" label="用户名" min-width="120" />
        <el-table-column prop="nickName" label="昵称" min-width="120" />
        <el-table-column prop="phone" label="手机号" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
        <el-table-column prop="authorityId" label="用户角色" min-width="200">
          <template #default="scope">
            <el-cascader
              :model-value="getUserRoleIds(scope.row)"
              :options="roleOptions"
              :props="{
                checkStrictly: true,
                value: 'value',
                label: 'label',
                children: 'children',
                multiple: true,
                emitPath: false
              }"
              placeholder="请选择角色"
              clearable
              collapse-tags
              collapse-tags-tooltip
              @change="(value: string[]) => handleRoleChange(value.map(v => Number(v)), scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="enable" label="启用" min-width="80">
          <template #default="scope">
            <el-switch
              v-model="scope.row.enable"
              :active-value="1"
              :inactive-value="2"
              :active-color="'#13ce66'"
              :inactive-color="'#ff4949'"
              @change="handleEnableChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="280">
          <template #default="scope">
            <div class="operation-buttons">
              <el-button size="small" type="primary" text @click="handleEdit(scope.row)">
                <el-icon class="el-icon--left"><Edit /></el-icon>编辑
              </el-button>
              <el-button size="small" type="primary" text @click="handleResetPassword(scope.row)">
                <el-icon class="el-icon--left"><Key /></el-icon>重设密码
              </el-button>
              <el-button size="small" type="danger" text @click="handleDelete(scope.row)">
                <el-icon class="el-icon--left"><Delete /></el-icon>删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 30, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>

      <!-- 添加/编辑用户对话框 -->
      <el-dialog
        v-model="dialogVisible"
        :title="dialogTitle"
        width="500px"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
      >
        <el-form
          ref="formRef"
          :model="userForm"
          :rules="userFormRules"
          label-width="100px"
          v-loading="dialogLoading"
        >
          <el-form-item label="用户名" prop="userName">
            <el-input v-model="userForm.userName" placeholder="请输入用户名" :disabled="dialogTitle === '编辑用户'" />
          </el-form-item>
          <el-form-item label="密码" prop="passWord" v-if="dialogTitle === '添加用户'">
            <el-input v-model="userForm.passWord" type="password" placeholder="请输入密码" />
          </el-form-item>
          <el-form-item label="昵称" prop="nickName">
            <el-input v-model="userForm.nickName" placeholder="请输入昵称" />
          </el-form-item>
          <el-form-item label="手机号" prop="phone">
            <el-input v-model="userForm.phone" placeholder="请输入手机号" />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="userForm.email" placeholder="请输入邮箱" />
          </el-form-item>
          <el-form-item label="角色" prop="authorityId">
            <el-cascader
              v-model="userForm.authorityId"
              :options="roleOptions"
              :props="{
                checkStrictly: true,
                value: 'value',
                label: 'label',
                children: 'children',
                multiple: true,
                emitPath: false
              }"
              placeholder="请选择角色"
              clearable
              collapse-tags
              collapse-tags-tooltip
            />
          </el-form-item>
          <el-form-item label="状态" prop="enable">
            <el-switch
              v-model="userForm.enable"
              :active-value="1"
              :inactive-value="0"
              :active-color="'#13ce66'"
              :inactive-color="'#ff4949'"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="dialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitForm" :loading="dialogLoading">确定</el-button>
          </span>
        </template>
      </el-dialog>

      <!-- 用户详细信息对话框 -->
      <el-dialog
        v-model="userInfoDialogVisible"
        title="用户详细信息"
        width="600px"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
      >
        <div v-loading="userInfoLoading">
          <el-form
            ref="userInfoFormRef"
            :model="userInfoForm"
            :rules="userInfoFormRules"
            label-width="100px"
          >
            <el-form-item label="用户ID">
              <span>{{ userDetailInfo?.ID }}</span>
            </el-form-item>
            <el-form-item label="用户名">
              <span>{{ userDetailInfo?.userName }}</span>
            </el-form-item>
            <el-form-item label="昵称">
              <span>{{ userDetailInfo?.nickName }}</span>
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="userInfoForm.phone" placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="userInfoForm.email" placeholder="请输入邮箱" />
            </el-form-item>
            <el-form-item label="角色">
              <el-tag v-for="(auth, index) in userDetailInfo?.authorities" :key="index" type="info" effect="plain" class="auth-tag">
                {{ auth.authorityName }}
              </el-tag>
            </el-form-item>
            <el-form-item label="状态">
              <el-tag :type="userDetailInfo?.enable === 1 ? 'success' : 'danger'">
                {{ userDetailInfo?.enable === 1 ? '启用' : '禁用' }}
              </el-tag>
            </el-form-item>
            <el-form-item label="创建时间">
              <span>{{ userDetailInfo?.CreatedAt }}</span>
            </el-form-item>
            <el-form-item label="更新时间">
              <span>{{ userDetailInfo?.UpdatedAt }}</span>
            </el-form-item>
          </el-form>
        </div>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="userInfoDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitUserInfoForm" :loading="userInfoLoading">保存</el-button>
          </span>
        </template>
      </el-dialog>

      <!-- 修改密码对话框 -->
      <el-dialog
        v-model="changePasswordDialogVisible"
        title="修改密码"
        width="500px"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
      >
        <el-form
          ref="changePasswordFormRef"
          :model="changePasswordForm"
          :rules="changePasswordRules"
          label-width="100px"
          v-loading="changePasswordLoading"
        >
          <el-form-item label="当前密码" prop="Password">
            <el-input v-model="changePasswordForm.Password" type="password" placeholder="请输入当前密码" show-password />
          </el-form-item>
          <el-form-item label="新密码" prop="NewPassword">
            <el-input v-model="changePasswordForm.NewPassword" type="password" placeholder="请输入新密码" show-password />
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="changePasswordDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="submitChangePasswordForm" :loading="changePasswordLoading">确定</el-button>
          </span>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { Plus, Edit, Delete, Key, User, Lock, Search, Refresh, CircleCheck, CircleClose } from '@element-plus/icons-vue'
import { onMounted } from 'vue'
import {
  currentPage,
  pageSize,
  total,
  loading,
  userList,
  fetchUserList,
  handleCurrentChange,
  handleSizeChange,
  handleAddUser,
  handleEdit,
  handleDelete,
  handleEnableChange,
  // 新增导入
  dialogVisible,
  dialogTitle,
  dialogLoading,
  userForm,
  userFormRules,
  formRef,
  submitForm,
  // 用户详情相关
  userDetailInfo,
  userInfoDialogVisible,
  userInfoLoading,
  fetchUserInfo,
  userInfoFormRef,
  userInfoForm,
  userInfoFormRules,
  submitUserInfoForm,
  // 角色相关
  fetchRoleList,
  roleOptions,
  handleRoleChange,
  getUserRoleIds,
  // 搜索相关
  searchForm,
  resetSearchForm,
  // 统计相关
  userStats
} from './userManager'

import {
  changePasswordDialogVisible,
  changePasswordLoading,
  changePasswordFormRef,
  changePasswordForm,
  changePasswordRules,
  handleChangePassword,
  submitChangePasswordForm,
  handleResetPassword
} from './passwordManager'

// 页面加载时获取用户列表和角色列表
onMounted(() => {
  fetchUserList()
  fetchRoleList()
})
</script>

<style src="./userManager.css" scoped></style>
