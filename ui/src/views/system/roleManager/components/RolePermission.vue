<template>
  <div class="role-permission-container">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="角色API" name="api">
        <div class="api-filter">
          <el-input
            v-model="filterTextName"
            placeholder="请输入API组名称"
            clearable
            @clear="handleFilterClear"
          />
          <el-input
            v-model="filterTextPath"
            placeholder="请输入API路径"
            clearable
            @clear="handleFilterClear"
          />
        </div>
        <el-tree
          ref="apiTree"
          :data="apiTreeData"
          :props="apiProps"
          show-checkbox
          node-key="onlyId"
          :default-checked-keys="apiTreeIds"
          @check="handleApiCheck"
          :filter-node-method="filterNode"
        >
          <template #default="{ data }">
            <span class="custom-tree-node">
              <span>{{ data.description }}</span>
              <span class="node-path">{{ data.path }}</span>
              <el-tag size="small" :type="getMethodType(data.method)">{{ data.method }}</el-tag>
            </span>
          </template>
        </el-tree>
      </el-tab-pane>
    </el-tabs>
    <div class="permission-footer">
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getAllApis, getPolicyPathByAuthorityId, updateCasbin } from '@/api/system/roles'
import { freshCasbin } from '@/api/system/apis'

import type { Authority } from '@/types/system'

// 定义接口
interface Api {
  ID: number
  path: string
  description: string
  apiGroup: string
  method: string
  onlyId?: string
}

interface ApiGroup {
  ID: string
  description: string
  children: Api[]
}

const props = defineProps<{
  role: Authority | null
}>()

const emit = defineEmits(['close'])

// 状态定义
const activeTab = ref('api')
const saving = ref(false)
const filterTextName = ref('')
const filterTextPath = ref('')
const apiTreeIds = ref<string[]>([])

// 树形数据
const apiTreeData = ref<ApiGroup[]>([])

// 树形控件引用
const apiTree = ref()

// 树形配置
const apiProps = {
  children: 'children',
  label: 'description'
}

// 获取方法类型对应的标签样式
const getMethodType = (method: string) => {
  const types: Record<string, string> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger'
  }
  return types[method] || 'info'
}

// 创建API树
const buildApiTree = (apis: Api[]): ApiGroup[] => {
  const apiObj: Record<string, Api[]> = {}
  
  apis.forEach(item => {
    item.onlyId = `p:${item.path}m:${item.method}`
    if (apiObj[item.apiGroup]) {
      apiObj[item.apiGroup].push(item)
    } else {
      apiObj[item.apiGroup] = [item]
    }
  })

  return Object.entries(apiObj).map(([key, value]) => ({
    ID: key,
    description: `${key}组`,
    children: value
  }))
}

// 加载API树
const loadApiTree = async () => {
  try {
    // 获取所有API
    const res = await getAllApis()
    if (res.data?.code === 0 && Array.isArray(res.data.data.apis)) {
      apiTreeData.value = buildApiTree(res.data.data.apis)
      
      // 获取当前角色已授权的API
      if (props.role) {
        const policyRes = await getPolicyPathByAuthorityId({
          authorityId: props.role.authorityId
        })
        if (policyRes.data?.code === 0 && Array.isArray(policyRes.data.data.paths)) {
          // 构建已授权API的唯一标识
          apiTreeIds.value = policyRes.data.data.paths.map((item: { path: string; method: string }) => 
            `p:${item.path}m:${item.method}`
          )
          // 设置树形控件的选中状态
          if (apiTree.value) {
            apiTree.value.setCheckedKeys(apiTreeIds.value)
          }
        } else {
          console.error('获取角色API权限失败:', policyRes)
          apiTreeIds.value = []
        }
      }
    } else {
      console.error('获取所有API失败:', res)
      apiTreeData.value = []
    }
  } catch (error) {
    console.error('加载API树失败:', error)
    ElMessage.error('加载API树失败')
  }
}

// 处理API选中
const handleApiCheck = () => {
  // 不需要在这里处理，保存时会获取选中的节点
}

// 过滤节点
const filterNode = (value: string, data: Api | ApiGroup) => {
  if (!filterTextName.value && !filterTextPath.value) return true
  
  let matchesName = true
  let matchesPath = true
  
  if (filterTextName.value) {
    matchesName = Boolean(data.description && data.description.includes(filterTextName.value))
  }
  
  if (filterTextPath.value) {
    matchesPath = Boolean('path' in data && data.path && data.path.includes(filterTextPath.value))
  }
  
  return matchesName && matchesPath
}

// 处理过滤清除
const handleFilterClear = () => {
  apiTree.value.filter('')
}

// 监听过滤条件变化
watch([filterTextName, filterTextPath], () => {
  apiTree.value.filter('')
})

// 保存权限设置
const handleSave = async () => {
  saving.value = true
  try {
    await saveApiPermissions()
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 保存API权限
const saveApiPermissions = async () => {
  if (!props.role) {
    ElMessage.error('角色信息不存在')
    return
  }

  // 获取选中的API节点
  const checkArr = apiTree.value.getCheckedNodes(true)
  const casbinInfos = checkArr.map((item: Api) => ({
    path: item.path,
    method: item.method
  }))

  // 调用接口保存API权限
  const res = await updateCasbin({
    authorityId: props.role.authorityId,
    casbinInfos
  })

  if (res.data?.code === 0) {
    // 刷新 Casbin 规则
    try {
      await freshCasbin()
      ElMessage.success('保存API权限成功')
      emit('close')
    } catch (error) {
      console.error('刷新 Casbin 规则失败:', error)
      ElMessage.error('刷新 Casbin 规则失败')
    }
  } else {
    ElMessage.error(res.data?.msg || '保存API权限失败')
  }
}

// 取消
const handleCancel = () => {
  emit('close')
}

// 初始化
onMounted(async () => {
  if (props.role) {
    await loadApiTree()
  }
})

// 监听角色变化
watch(() => props.role, async (newVal) => {
  if (newVal) {
    // 重置状态
    apiTreeIds.value = []
    // 重新加载API权限
    await loadApiTree()
  }
}, { immediate: true })
</script>

<style scoped>
.role-permission-container {
  padding: 20px;
}

.api-filter {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  padding-right: 8px;
}

.node-path {
  color: #909399;
  margin: 0 8px;
}

.permission-footer {
  margin-top: 20px;
  text-align: right;
}
</style>