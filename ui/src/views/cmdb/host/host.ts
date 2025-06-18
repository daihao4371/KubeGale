import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import type { Host } from '@/types/cmdb'
import {
  getHostList,
  createHost,
  updateHost,
  deleteHost,
  batchDeleteHosts,
  getHostDetail,
  importHosts,
  authenticateHost
} from '@/api/cmdb/host'
import { getProjectList } from '@/api/cmdb/project'

// 定义接口
interface HostListParams {
  page: number
  pageSize: number
  name?: string
  ip?: string
  project_id?: number
}

interface ProjectTreeNode {
  id: number
  name: string
  children?: ProjectTreeNode[]
}

export default function useHost() {
  // 状态管理
  const loading = ref(false)
  const tableData = ref<Host[]>([])
  const selectedHosts = ref<Host[]>([])
  const currentPage = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const dialogVisible = ref(false)
  const dialogType = ref<'add' | 'edit'>('add')
  const formRef = ref<FormInstance>()

  // 搜索表单
  const searchForm = reactive({
    name: '',
    ip: '',
    projectId: undefined as number | undefined
  })

  // 表单数据
  const formData = reactive({
    ID: 0,
    name: '',
    serverHost: '',
    port: 22,
    username: 'root',
    password: '',
    privateKey: '',
    project: undefined as number | undefined,
    note: ''
  })

  // 表单验证规则
  const formRules = {
    name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
    serverHost: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
    port: [{ required: true, message: '请输入端口', trigger: 'blur' }],
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    project: [{ required: true, message: '请选择所属项目', trigger: 'change' }]
  }

  // 项目选项
  const projectOptions = ref<{ id: number; name: string }[]>([])

  // 项目树形数据
  const projectTreeData = ref<ProjectTreeNode[]>([])

  // 详情对话框
  const detailVisible = ref(false)
  const hostDetail = ref<Host | null>(null)
  const detailLoading = ref(false)

  // 导入对话框
  const importVisible = ref(false)
  const importForm = reactive({
    file: null as File | null,
    project: undefined as number | undefined
  })
  const importFormRef = ref<FormInstance>()
  const importRules = {
    file: [{ required: true, message: '请选择Excel文件', trigger: 'change' }],
    project: [{ required: true, message: '请选择所属项目', trigger: 'change' }]
  }

  // 新增：全量主机列表用于项目树计数
  const allHosts = ref<Host[]>([])

  // 拉取所有主机（不分页、不过滤）
  const fetchAllHosts = async () => {
    try {
      const res = await getHostList({ page: 1, pageSize: 99999 }) // 拉全量
      if (res.code === 0) {
        allHosts.value = res.data.list
      }
    } catch (error) {
      allHosts.value = []
    }
  }

  // 获取项目列表并转换为树形结构
  const fetchProjectList = async () => {
    try {
      const res = await getProjectList({
        page: 1,
        pageSize: 1000
      })
      if (res.code === 0) {
        // 转换为树形结构
        projectTreeData.value = res.data.list.map(item => ({
          id: item.ID,
          name: item.name
        }))
        // 同时更新项目选项
        projectOptions.value = res.data.list.map(item => ({
          id: item.ID,
          name: item.name
        }))
      }
    } catch (error) {
      console.error('获取项目列表失败:', error)
      ElMessage.error('获取项目列表失败')
    }
  }

  // 获取主机列表
  const fetchList = async () => {
    loading.value = true
    try {
      // 确保先获取项目列表
      if (projectOptions.value.length === 0) {
        await fetchProjectList()
      }
      const params: HostListParams = {
        page: currentPage.value,
        pageSize: pageSize.value,
        name: searchForm.name || undefined,
        ip: searchForm.ip || undefined,
        project_id: searchForm.projectId ? Number(searchForm.projectId) : undefined
      }
      const res = await getHostList(params)
      if (res.code === 0) {
        // 将项目ID映射为项目名称
        let filteredData = res.data.list
        
        // 如果选择了特定项目，确保只显示该项目的主机
        if (searchForm.projectId) {
          filteredData = filteredData.filter(host => Number(host.project) === Number(searchForm.projectId))
        }
        
        tableData.value = filteredData.map(host => ({
          ...host,
          projectName: projectOptions.value.find(p => p.id === Number(host.project))?.name || '-'
        }))
        total.value = filteredData.length // 更新总数为过滤后的数量
      }
      // 同步刷新全量主机
      fetchAllHosts()
    } catch (error) {
      console.error('获取主机列表失败:', error)
      ElMessage.error('获取主机列表失败')
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
    Object.assign(searchForm, {
      name: '',
      ip: '',
      projectId: undefined
    })
    currentPage.value = 1
    handleSearch()
  }

  // 表格选择
  const handleSelectionChange = (selection: Host[]) => {
    selectedHosts.value = selection
  }

  // 新建主机
  const handleAdd = () => {
    dialogType.value = 'add'
    formData.ID = 0
    formData.name = ''
    formData.serverHost = ''
    formData.port = 22
    formData.username = 'root'
    formData.password = ''
    formData.privateKey = ''
    formData.project = undefined
    formData.note = ''
    dialogVisible.value = true
  }

  // 编辑主机
  const handleEdit = (row: Host) => {
    dialogType.value = 'edit'
    Object.assign(formData, row)
    dialogVisible.value = true
  }

  // 删除主机
  const handleDelete = (row: Host) => {
    ElMessageBox.confirm('确认删除该主机吗？', '提示', {
      type: 'warning'
    }).then(async () => {
      try {
        const res = await deleteHost({ id: row.ID })
        if (res.code === 0) {
          ElMessage.success(res.msg || '删除成功')
          fetchList()
        }
      } catch (error) {
        console.error('删除主机失败:', error)
        ElMessage.error('删除主机失败')
      }
    })
  }

  // 批量删除
  const handleBatchDelete = () => {
    if (selectedHosts.value.length === 0) {
      ElMessage.warning('请选择要删除的主机')
      return
    }
    ElMessageBox.confirm(`确认删除选中的 ${selectedHosts.value.length} 台主机吗？`, '提示', {
      type: 'warning'
    }).then(async () => {
      try {
        const res = await batchDeleteHosts({
          ids: selectedHosts.value.map(item => item.ID)
        })
        if (res.code === 0) {
          ElMessage.success(res.msg || '批量删除成功')
          fetchList()
        }
      } catch (error) {
        console.error('批量删除主机失败:', error)
        ElMessage.error('批量删除主机失败')
      }
    })
  }

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
      if (valid) {
        try {
          const api = dialogType.value === 'add' ? createHost : updateHost
          const submitData = {
            ...formData,
            project: formData.project || 0
          }
          const res = await api(submitData)
          if (res.code === 0) {
            ElMessage.success(dialogType.value === 'add' ? '创建成功' : '更新成功')
            dialogVisible.value = false
            fetchList()
          }
        } catch (error) {
          console.error(dialogType.value === 'add' ? '创建主机失败:' : '更新主机失败:', error)
          ElMessage.error(dialogType.value === 'add' ? '创建主机失败' : '更新主机失败')
        }
      }
    })
  }

  // 分页
  const handleSizeChange = (val: number) => {
    pageSize.value = val
    fetchList()
  }

  const handleCurrentChange = (val: number) => {
    currentPage.value = val
    fetchList()
  }

  // 查看详情
  const handleDetail = async (row: Host) => {
    detailLoading.value = true
    try {
      const res = await getHostDetail(row.ID)
      if (res.code === 0) {
        hostDetail.value = res.data
        detailVisible.value = true
      }
    } catch (error) {
      console.error('获取主机详情失败:', error)
      ElMessage.error('获取主机详情失败')
    } finally {
      detailLoading.value = false
    }
  }

  // 处理文件选择
  const handleFileChange = (file: { raw: File }) => {
    importForm.file = file.raw
    return false
  }

  // 下载模板
  const downloadTemplate = () => {
    const template = [
      ['主机名称', 'IP地址', '端口', '用户名', '密码', '备注'],
      ['示例主机', '192.168.1.1', '22', 'root', 'password', '示例备注']
    ].map(row => row.join(',')).join('\n')
    
    const blob = new Blob([template], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = '主机导入模板.csv'
    link.click()
    URL.revokeObjectURL(link.href)
  }

  // 提交导入
  const handleImport = async () => {
    if (!importFormRef.value) return
    await importFormRef.value.validate(async (valid) => {
      if (valid && importForm.file && importForm.project) {
        const formData = new FormData()
        formData.append('file', importForm.file)
        formData.append('projectId', importForm.project.toString())
        
        try {
          const res = await importHosts(formData)
          if (res.code === 0) {
            ElMessage.success(res.msg || '导入成功')
            importVisible.value = false
            // 重置表单
            importForm.file = null
            importForm.project = undefined
            if (importFormRef.value) {
              importFormRef.value.resetFields()
            }
            fetchList()
          }
        } catch (error) {
          console.error('导入主机失败:', error)
          ElMessage.error('导入主机失败')
        }
      }
    })
  }

  // 处理项目选择
  const handleProjectSelect = (data: ProjectTreeNode) => {
    if (data && data.id) {
      searchForm.projectId = data.id
      currentPage.value = 1 // 切换项目时回到第一页
      handleSearch()
    }
  }

  // 初始化
  onMounted(() => {
    fetchList()
    fetchProjectList()
  })

  return {
    loading,
    tableData,
    selectedHosts,
    currentPage,
    pageSize,
    total,
    dialogVisible,
    dialogType,
    formRef,
    searchForm,
    formData,
    formRules,
    projectOptions,
    projectTreeData,
    detailVisible,
    hostDetail,
    detailLoading,
    handleSearch,
    handleReset,
    handleSelectionChange,
    handleAdd,
    handleEdit,
    handleDelete,
    handleBatchDelete,
    handleSubmit,
    handleSizeChange,
    handleCurrentChange,
    handleDetail,
    importVisible,
    importForm,
    importFormRef,
    importRules,
    handleFileChange,
    downloadTemplate,
    handleImport,
    handleProjectSelect,
    allHosts
  }
}
