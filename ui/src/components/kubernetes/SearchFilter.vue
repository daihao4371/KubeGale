<template>
  <div class="search-filter">
    <el-form :model="searchForm" :inline="true" class="search-form">
      <!-- 名称搜索 -->
      <el-form-item label="名称" v-if="showNameFilter">
        <el-input
          v-model="searchForm.name"
          placeholder="请输入名称"
          clearable
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        />
      </el-form-item>

      <!-- 命名空间筛选 -->
      <el-form-item label="命名空间" v-if="showNamespaceFilter">
        <el-select 
          v-model="searchForm.namespace" 
          placeholder="请选择命名空间" 
          clearable
          @change="handleSearch"
        >
          <el-option 
            v-for="ns in namespaces" 
            :key="ns" 
            :label="ns" 
            :value="ns"
          />
        </el-select>
      </el-form-item>

      <!-- 状态筛选 -->
      <el-form-item label="状态" v-if="showStatusFilter">
        <el-select 
          v-model="searchForm.status" 
          placeholder="请选择状态" 
          clearable
          @change="handleSearch"
        >
          <el-option 
            v-for="status in statusOptions" 
            :key="status.value" 
            :label="status.label" 
            :value="status.value"
          />
        </el-select>
      </el-form-item>

      <!-- 标签筛选 -->
      <el-form-item label="标签" v-if="showLabelFilter">
        <el-input
          v-model="searchForm.labels"
          placeholder="key=value"
          clearable
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        />
      </el-form-item>

      <!-- 自定义筛选器 -->
      <slot name="filters" :search-form="searchForm" :handle-search="handleSearch" />

      <!-- 操作按钮 -->
      <el-form-item>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'

interface StatusOption {
  label: string
  value: string
}

interface SearchForm {
  name: string
  namespace: string
  status: string
  labels: string
  [key: string]: any
}

interface Props {
  modelValue?: SearchForm
  showNameFilter?: boolean
  showNamespaceFilter?: boolean
  showStatusFilter?: boolean
  showLabelFilter?: boolean
  namespaces?: string[]
  statusOptions?: StatusOption[]
  defaultSearch?: Partial<SearchForm>
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => ({
    name: '',
    namespace: '',
    status: '',
    labels: ''
  }),
  showNameFilter: true,
  showNamespaceFilter: true,
  showStatusFilter: true,
  showLabelFilter: true,
  namespaces: () => [],
  statusOptions: () => [
    { label: '运行中', value: 'Running' },
    { label: '等待中', value: 'Pending' },
    { label: '失败', value: 'Failed' },
    { label: '成功', value: 'Succeeded' },
    { label: '终止中', value: 'Terminating' }
  ],
  defaultSearch: () => ({})
})

const emit = defineEmits<{
  'update:modelValue': [value: SearchForm]
  search: [params: SearchForm]
  reset: []
}>()

const searchForm = reactive<SearchForm>({
  name: '',
  namespace: '',
  status: '',
  labels: '',
  ...props.defaultSearch,
  ...props.modelValue
})

// 监听搜索表单变化
watch(
  () => searchForm,
  (newValue) => {
    emit('update:modelValue', { ...newValue })
  },
  { deep: true }
)

// 监听外部值变化
watch(
  () => props.modelValue,
  (newValue) => {
    Object.assign(searchForm, newValue)
  },
  { deep: true }
)

const handleSearch = () => {
  emit('search', { ...searchForm })
}

const handleReset = () => {
  Object.assign(searchForm, {
    name: '',
    namespace: '',
    status: '',
    labels: '',
    ...props.defaultSearch
  })
  emit('reset')
  handleSearch()
}
</script>

<style scoped>
.search-filter {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.search-form {
  margin: 0;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 12px;
}

.search-form :deep(.el-form-item__label) {
  color: #606266;
  font-weight: 500;
}

.search-form :deep(.el-input) {
  width: 200px;
}

.search-form :deep(.el-select) {
  width: 200px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .search-form :deep(.el-form-item) {
    margin-right: 10px;
  }
}

@media (max-width: 768px) {
  .search-filter {
    padding: 15px;
  }
  
  .search-form :deep(.el-form-item) {
    width: 100%;
    margin-right: 0;
  }
  
  .search-form :deep(.el-input),
  .search-form :deep(.el-select) {
    width: 100%;
  }
}
</style>