import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOperationRecordList, deleteOperationRecord, batchDeleteOperationRecord } from '@/api/system/operationRecord'
import type { OperationRecordParams, SysOperationRecord } from '@/types/system'

// 操作记录状态
export const operationState = reactive({
  loading: false,
  tableData: [] as SysOperationRecord[],
  total: 0,
  currentPage: 1,
  pageSize: 10,
  selectedRows: [] as SysOperationRecord[],
  searchForm: {
    page: 1,
    pageSize: 10
  } as OperationRecordParams,
  detailDialogVisible: false,
  currentDetail: null as SysOperationRecord | null
})

// 获取列表数据
export const fetchData = async () => {
  operationState.loading = true
  try {
    const res = await getOperationRecordList(operationState.searchForm)
    if (res.data?.code === 0) {
      operationState.tableData = res.data.data.list
      operationState.total = res.data.data.total
    } else {
      ElMessage.error(res.data?.msg || '获取数据失败')
    }
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  } finally {
    operationState.loading = false
  }
}

// 搜索
export const handleSearch = () => {
  operationState.currentPage = 1
  operationState.searchForm.page = 1
  fetchData()
}

// 重置搜索
export const resetSearch = () => {
  Object.assign(operationState.searchForm, {
    page: 1,
    pageSize: operationState.pageSize,
    path: undefined,
    method: undefined,
    status: undefined
  })
  fetchData()
}

// 处理选择变化
export const handleSelectionChange = (rows: SysOperationRecord[]) => {
  operationState.selectedRows = rows
}

// 处理详情
export const handleDetail = (row: SysOperationRecord) => {
  operationState.currentDetail = row
  operationState.detailDialogVisible = true
}

// 处理删除
export const handleDelete = async (row: SysOperationRecord) => {
  try {
    await ElMessageBox.confirm('确定要删除该条记录吗？', '提示', {
      type: 'warning'
    })
    // 使用ID字段而不是id字段
    const ID = row.ID || row.id
    const res = await deleteOperationRecord(ID)
    if (res.data?.code === 0) {
      ElMessage.success('删除成功')
      fetchData()
    } else {
      ElMessage.error(res.data?.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 处理批量删除
export const handleBatchDelete = async () => {
  if (!operationState.selectedRows.length) return
  try {
    await ElMessageBox.confirm('确定要删除选中的记录吗？', '提示', {
      type: 'warning'
    })
    // 使用ID字段而不是id字段，确保兼容性
    const ids = operationState.selectedRows.map(row => row.ID || row.id)
    const res = await batchDeleteOperationRecord(ids)
    if (res.data?.code === 0) {
      ElMessage.success('批量删除成功')
      fetchData()
    } else {
      ElMessage.error(res.data?.msg || '批量删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

// 处理分页大小变化
export const handleSizeChange = (val: number) => {
  operationState.pageSize = val
  operationState.searchForm.pageSize = val
  fetchData()
}

// 处理页码变化
export const handleCurrentChange = (val: number) => {
  operationState.currentPage = val
  operationState.searchForm.page = val
  fetchData()
}

// 获取状态码类型
export const getStatusType = (status: number) => {
  if (status >= 200 && status < 300) return 'success'
  if (status >= 300 && status < 400) return 'warning'
  if (status >= 400) return 'danger'
  return 'info'
}

// 获取请求方法类型
export const getMethodType = (method: string) => {
  const types: Record<string, string> = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger'
  }
  return types[method] || 'info'
}

// 格式化JSON字符串
export const formatJSON = (jsonString: string | undefined) => {
  if (!jsonString) return '-'
  try {
    // 尝试解析JSON字符串
    const parsed = JSON.parse(jsonString)
    // 格式化为带缩进的JSON字符串
    return JSON.stringify(parsed, null, 2)
  } catch (e) {
    // 如果解析失败，直接返回原始字符串
    return jsonString
  }
}
