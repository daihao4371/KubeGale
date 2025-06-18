import { reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAllApis,
  getApiGroups,
  createApi,
  getApiById,
  updateApi,
  deleteApi,
  freshCasbin
} from '@/api/system/apis'
import type { Api, ApiState, ApiForm } from '@/types/system'
import { defaultApiForm } from '@/types/system'

export default function useApiManager() {
  const state = reactive<ApiState>({
    loading: false,
    apiList: [],
    total: 0,
    pageSize: 10,
    currentPage: 1,
    dialogVisible: false,
    dialogTitle: '',
    dialogType: 'create',
    formLoading: false,
    apiGroups: [],
    searchForm: {
      path: '',
      description: '',
      apiGroup: '',
      method: ''
    },
    filteredApiList: []
  })

  const form = ref<ApiForm>({ ...defaultApiForm })

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

  // 加载API列表
  const loadApiList = async () => {
    state.loading = true
    try {
      const res = await getAllApis()
      if (res.data?.code === 0) {
        state.apiList = res.data.data.apis
        state.total = state.apiList.length
        state.filteredApiList = state.apiList
        
        // 加载API分组，确保搜索栏可以使用
        await loadApiGroups()
      } else {
        ElMessage.error(res.data?.msg || '获取API列表失败')
      }
    } catch (error) {
      console.error('获取API列表失败:', error)
      ElMessage.error('获取API列表失败')
    } finally {
      state.loading = false
    }
  }

  // 加载API分组
  const loadApiGroups = async () => {
    try {
      const res = await getApiGroups()
      if (res.data?.code === 0) {
        state.apiGroups = res.data.data.groups
      }
    } catch (error) {
      console.error('获取API分组失败:', error)
    }
  }

  // 处理页码变化
  const handleCurrentChange = (page: number) => {
    state.currentPage = page
  }

  // 处理每页条数变化
  const handleSizeChange = (size: number) => {
    state.pageSize = size
    state.currentPage = 1
  }

  // 打开创建对话框
  const handleAddApi = async () => {
    state.dialogType = 'create'
    state.dialogTitle = '创建API'
    state.dialogVisible = true
    form.value = { ...defaultApiForm }
    await loadApiGroups()
  }

  // 打开编辑对话框
  const handleEdit = async (row: Api) => {
    state.dialogType = 'edit'
    state.dialogTitle = '编辑API'
    state.dialogVisible = true

    // 先填充当前行的数据，提升用户体验
    form.value = {
      ID: row.ID,
      path: row.path,
      description: row.description,
      apiGroup: row.apiGroup,
      method: row.method
    }

    // 加载API分组
    await loadApiGroups()

  }

  // 处理删除
  const handleDelete = async (row: Api) => {
    try {
      await ElMessageBox.confirm(
        `确定要删除API "${row.description}" 吗？`,
        '删除确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      )

      const res = await deleteApi({ ID: row.ID })
      if (res.data?.code === 0) {
        ElMessage.success('删除成功')
        await freshCasbin()
        loadApiList()
      } else {
        ElMessage.error(res.data?.msg || '删除失败')
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('删除API失败:', error)
        ElMessage.error('删除API失败')
      }
    }
  }

  // 提交表单
  const handleSubmit = async () => {
    state.formLoading = true
    try {
      if (state.dialogType === 'create') {
        const res = await createApi(form.value)
        if (res.data?.code === 0) {
          ElMessage.success('创建成功')
          state.dialogVisible = false
          await freshCasbin()
          loadApiList()
        } else {
          ElMessage.error(res.data?.msg || '创建失败')
        }
      } else {
        if (!form.value.ID) {
          ElMessage.error('更新失败：缺少ID')
          return
        }
        // 只传递后端需要的字段，避免 createdAt/updatedAt 导致报错
        const { ID, path, description, apiGroup, method } = form.value
        const updateData = { ID, path, description, apiGroup, method }
        const res = await updateApi(updateData as Required<ApiForm>)
        if (res.data?.code === 0) {
          ElMessage.success('更新成功')
          state.dialogVisible = false
          await freshCasbin()
          loadApiList()
        } else {
          ElMessage.error(res.data?.msg || '更新失败')
        }
      }
    } catch (error) {
      console.error('提交失败:', error)
      ElMessage.error('提交失败')
    } finally {
      state.formLoading = false
    }
  }

  // 搜索方法
  const handleSearch = () => {
    let filtered = state.apiList
    if (state.searchForm.path) {
      filtered = filtered.filter(api => api.path.includes(state.searchForm.path))
    }
    if (state.searchForm.description) {
      filtered = filtered.filter(api => api.description.includes(state.searchForm.description))
    }
    if (state.searchForm.apiGroup) {
      filtered = filtered.filter(api => api.apiGroup === state.searchForm.apiGroup)
    }
    if (state.searchForm.method) {
      filtered = filtered.filter(api => api.method === state.searchForm.method)
    }
    state.filteredApiList = filtered
    state.total = filtered.length
    state.currentPage = 1
  }

  // 重置搜索
  const handleResetSearch = () => {
    state.searchForm = { path: '', description: '', apiGroup: '', method: '' }
    state.filteredApiList = state.apiList
    state.total = state.apiList.length
    state.currentPage = 1
  }

  // 刷新Casbin缓存
  const handleFreshCasbin = async () => {
    try {
      state.loading = true
      const res = await freshCasbin()
      if (res.data?.code === 0) {
        ElMessage.success('刷新缓存成功')
      } else {
        ElMessage.error(res.data?.msg || '刷新缓存失败')
      }
    } catch (error) {
      console.error('刷新缓存失败:', error)
      ElMessage.error('刷新缓存失败')
    } finally {
      state.loading = false
    }
  }

  return {
    state,
    form,
    getMethodType,
    loadApiList,
    handleCurrentChange,
    handleSizeChange,
    handleAddApi,
    handleEdit,
    handleDelete,
    handleSubmit,
    handleSearch,
    handleResetSearch,
    handleFreshCasbin  // 添加到返回对象中
  }
}
