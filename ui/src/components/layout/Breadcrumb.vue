<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index" :to="item.path">
      {{ item.meta.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>

  <!-- 历史记录下拉菜单 -->
  <el-dropdown class="history-dropdown" trigger="click">
    <el-button type="primary" size="small" plain>
      历史记录 <el-icon class="el-icon--right"><arrow-down /></el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item v-if="visitedPages.length === 0" disabled>暂无历史记录</el-dropdown-item>
        <el-dropdown-item
          v-for="(page, index) in visitedPages"
          :key="index"
          @click="navigateTo(page.path)"
        >
          {{ page.meta.title }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'

defineOptions({
  name: 'BreadcrumbNav'
})

interface BreadcrumbItem {
  path: string
  meta: {
    title: string
    [key: string]: string | number | boolean | undefined
  }
}

const route = useRoute()
const router = useRouter()
const breadcrumbs = ref<BreadcrumbItem[]>([])
const visitedPages = ref<BreadcrumbItem[]>([])
const MAX_HISTORY = 10 // 最多保存10条历史记录

// 从本地存储加载历史记录
const loadVisitedPages = () => {
  const savedPages = localStorage.getItem('visitedPages')
  if (savedPages) {
    try {
      visitedPages.value = JSON.parse(savedPages)
    } catch (e) {
      console.error('解析历史记录失败:', e)
      visitedPages.value = []
    }
  }
}

// 保存历史记录到本地存储
const saveVisitedPages = () => {
  localStorage.setItem('visitedPages', JSON.stringify(visitedPages.value))
}

// 添加页面到历史记录
const addToHistory = (item: BreadcrumbItem) => {
  // 如果页面没有标题，不添加到历史记录
  if (!item.meta?.title) return

  // 如果已经存在相同路径的页面，先移除它
  const existingIndex = visitedPages.value.findIndex(page => page.path === item.path)
  if (existingIndex !== -1) {
    visitedPages.value.splice(existingIndex, 1)
  }

  // 添加到历史记录开头
  visitedPages.value.unshift(item)

  // 限制历史记录数量
  if (visitedPages.value.length > MAX_HISTORY) {
    visitedPages.value = visitedPages.value.slice(0, MAX_HISTORY)
  }

  // 保存到本地存储
  saveVisitedPages()
}

// 导航到指定路径
const navigateTo = (path: string) => {
  router.push(path)
}

const getBreadcrumbs = () => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)

  const result = matched.map(item => {
    return {
      path: item.path,
      meta: {
        title: item.meta.title as string,
        ...item.meta
      }
    } as BreadcrumbItem
  })

  breadcrumbs.value = result

  // 如果当前页面有标题，添加到历史记录
  if (result.length > 0) {
    addToHistory(result[result.length - 1])
  }
}

// 组件挂载时加载历史记录
onMounted(() => {
  loadVisitedPages()
})

watch(
  () => route.path,
  () => getBreadcrumbs(),
  { immediate: true }
)
</script>

<style scoped>
.el-breadcrumb {
  margin-left: 8px;
  line-height: 60px;
  display: inline-block;
}

.history-dropdown {
  margin-left: 16px;
  display: inline-block;
}
</style>
