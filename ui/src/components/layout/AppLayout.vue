<template>
  <el-container class="layout-container">
    <el-aside width="200px" class="sidebar">
      <div class="logo">
        <h2>KubeGale</h2>
      </div>
      <el-menu
        router
        :default-active="activeMenu"
        background-color="var(--bg-color-sidebar)" /* Assuming --bg-color-sidebar will be defined, e.g., #242424 or #304156 */
      text-color="var(--vt-c-text-dark-2)" /* Light text for dark sidebar */
      active-text-color="var(--color-primary)"
      >
      <el-menu-item index="/dashboard">
        <el-icon><DataLine /></el-icon>
        <span>仪表盘</span>
      </el-menu-item>
      <el-menu-item index="/kubernetes">
        <el-icon><Connection /></el-icon>
        <span>Kubernetes管理</span>
      </el-menu-item>
      <el-menu-item index="/cicd">
        <el-icon><SetUp /></el-icon>
        <span>CI/CD管理</span>
      </el-menu-item>
      <el-menu-item index="/config">
        <el-icon><Files /></el-icon>
        <span>配置中心</span>
      </el-menu-item>
      <el-menu-item index="/cmdb">
        <el-icon><DataAnalysis /></el-icon>
        <span>CMDB管理</span>
      </el-menu-item>
      <el-menu-item index="/docker">
        <el-icon><Box /></el-icon>
        <span>Docker管理</span>
      </el-menu-item>
      <el-menu-item index="/prometheus">
        <el-icon><Monitor /></el-icon>
        <span>Prometheus监控</span>
      </el-menu-item>
      <el-menu-item index="/system">
        <el-icon><Setting /></el-icon>
        <span>系统管理</span>
      </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-icon class="toggle-sidebar"><Fold /></el-icon>
          <breadcrumb />
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="30" />
              <span>管理员</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人中心</el-dropdown-item>
                <el-dropdown-item>修改密码</el-dropdown-item>
                <el-dropdown-item divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
      <el-footer class="footer">
        <p>KubeGale © {{ new Date().getFullYear() }} - 云原生Kubernetes管理平台</p>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  DataLine,
  Connection,
  SetUp,
  Files,
  DataAnalysis,
  Box,
  Monitor,
  Setting,
  Fold
} from '@element-plus/icons-vue'
import Breadcrumb from './Breadcrumb.vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: var(--bg-color-sidebar); /* Assuming --bg-color-sidebar will be defined, e.g., #242424 or #304156 */
  color: var(--vt-c-text-dark-2); /* Default text color for items directly in sidebar not in el-menu */
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--vt-c-text-dark-1); /* Brighter text for logo */
}

.header {
  background-color: var(--bg-color-header);
  border-bottom: 1px solid var(--border-color-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.toggle-sidebar {
  font-size: 20px;
  margin-right: 20px;
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.user-info span {
  margin-left: 10px;
}

.footer {
  background-color: var(--bg-color-card);
  border-top: 1px solid var(--border-color-light);
  text-align: center;
  font-size: 12px;
  color: var(--text-color-secondary);
  padding: 10px 0;
}

/* Ensure el-main has appropriate padding using our global variable */
.el-main { /* Targeting the el-main component directly. Scoped styles will handle specificity. */
  padding: var(--spacing-lg);
}
</style>
