import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createProject, getProjectList, updateProject, deleteProject, batchDeleteProjects, getProjectDetail } from '@/api/cmdb/project'

// 项目类型定义
export interface Project {
  ID: number
  name: string
  description: string
  manager: string
  CreatedAt: string
  UpdatedAt: string
  CreatedBy: number
  UpdatedBy: number
  DeletedBy: number
  DeletedAt: string | null
}

export default function useProject() {
  // 状态定义
  const tableData = ref<Project[]>([])
  const loading = ref(false)
  const searchForm = ref({
    name: '',
    description: ''
  })
  const pageSize = ref(10)
  const currentPage = ref(1)
  const total = ref(0)
  const dialogVisible = ref(false)
  const dialogTitle = ref('')
  const dialogType = ref<'add' | 'edit'>('add')
  const selectedRows = ref<Project[]>([])
  const formData = ref<Partial<Project>>({
    name: '',
    manager: '',
    description: ''
  })
  const detailVisible = ref(false)
  const detailData = ref<Project>({} as Project)

  // 表单校验规则
  const rules = {
    name: [
      { required: true, message: '请输入项目名称', trigger: 'blur' },
      { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
    ],
    manager: [
      { required: true, message: '请输入项目负责人', trigger: 'blur' }
    ]
  }

  // 获取项目列表
  const fetchList = async () => {
    loading.value = true
    try {
      const res = await getProjectList({
        page: currentPage.value,
        pageSize: pageSize.value,
        name: searchForm.value.name,
        description: searchForm.value.description
      })
      if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
      } else {
        ElMessage.error(res.msg || '获取项目列表失败')
      }
    } catch (error) {
      ElMessage.error('获取项目列表失败')
    } finally {
      loading.value = false
    }
  }

  // 搜索
  const handleSearch = () => {
    currentPage.value = 1
    fetchList()
  }

  // 重置搜索
  const handleReset = () => {
    searchForm.value = {
      name: '',
      description: ''
    }
    handleSearch()
  }

  // 新建项目
  const handleAdd = () => {
    dialogType.value = 'add'
    dialogTitle.value = '新建项目'
    formData.value = {
      name: '',
      manager: '',
      description: ''
    }
    dialogVisible.value = true
  }

  // 编辑项目
  const handleEdit = (row: Project) => {
    dialogType.value = 'edit'
    dialogTitle.value = '编辑项目'
    formData.value = { ...row }
    dialogVisible.value = true
  }

  // 删除项目
  const handleDelete = (row: Project) => {
    ElMessageBox.confirm(
      '确定要删除该项目吗？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(async () => {
      try {
        const res = await deleteProject({ id: row.ID })
        if (res.code === 0) {
          ElMessage.success('删除成功')
          fetchList()
        } else {
          ElMessage.error(res.msg || '删除失败')
        }
      } catch (error) {
        ElMessage.error('删除失败')
      }
    }).catch(() => {})
  }

  // 批量删除
  const handleBatchDelete = () => {
    if (selectedRows.value.length === 0) {
      ElMessage.warning('请选择要删除的项目')
      return
    }
    ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个项目吗？`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(async () => {
      try {
        const ids = selectedRows.value.map(row => row.ID)
        const res = await batchDeleteProjects({ ids })
        if (res.code === 0) {
          ElMessage.success('删除成功')
          fetchList()
        } else {
          ElMessage.error(res.msg || '删除失败')
        }
      } catch (error) {
        ElMessage.error('删除失败')
      }
    }).catch(() => {})
  }

  // 提交表单
  const handleSubmit = async () => {
    try {
      if (dialogType.value === 'add') {
        const res = await createProject({
          name: formData.value.name!,
          manager: formData.value.manager!,
          description: formData.value.description || ''
        })
        if (res.code === 0) {
          ElMessage.success('创建成功')
          dialogVisible.value = false
          fetchList()
        } else {
          ElMessage.error(res.msg || '创建失败')
        }
      } else {
        const res = await updateProject({
          id: formData.value.ID!,
          name: formData.value.name!,
          manager: formData.value.manager!,
          description: formData.value.description || ''
        })
        if (res.code === 0) {
          ElMessage.success('更新成功')
          dialogVisible.value = false
          fetchList()
        } else {
          ElMessage.error(res.msg || '更新失败')
        }
      }
    } catch (error) {
      ElMessage.error(dialogType.value === 'add' ? '创建失败' : '更新失败')
    }
  }

  // 分页
  const handleCurrentChange = (page: number) => {
    currentPage.value = page
    fetchList()
  }

  const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchList()
  }

  // 查看详情
  const handleDetail = async (row: Project) => {
    try {
      const res = await getProjectDetail(row.ID)
      if (res.code === 0) {
        detailData.value = res.data
        detailVisible.value = true
      } else {
        ElMessage.error(res.msg || '获取详情失败')
      }
    } catch (error) {
      ElMessage.error('获取详情失败')
    }
  }

  // 初始化加载数据
  fetchList()

  return {
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
  }
}
