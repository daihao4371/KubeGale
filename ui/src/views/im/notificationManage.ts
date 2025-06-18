import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { createFeiShu, getNotificationList, deleteNotification, updateFeiShu } from '@/api/im/notification'
import { getUserList } from '@/api/system/user'
import type { CreateFeiShuParams } from '@/types/im'
import type { ApiResponse } from '@/types/cmdb'

interface UserListResponse {
  list: Array<{
    userName: string
    nickName: string
  }>
  total: number
  page: number
  pageSize: number
}

interface NotificationItem {
  id: number
  name: string
  type: string
  notification_policy: string
  robot_url: string
  created_at: string
}

interface NotificationListResponse {
  list: NotificationItem[]
  total: number
  page: number
  pageSize: number
}

// 状态定义
const searchName = ref('')
const tableData = ref<NotificationItem[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const feishuDialogVisible = ref(false)
const formLoading = ref(false)
const userOptions = ref<{ label: string; value: string }[]>([])

// 表单数据
const feishuForm = ref<CreateFeiShuParams>({
  name: '',
  type: 'feishu',
  enabled: true,
  webhook_url: '',
  description: '',
  tags: [],
  notify_events: ['alert', 'warning'],
  receivers: [],
  send_daily_stats: true
})

// 表单验证规则
const feishuFormRules: FormRules = {
  name: [
    { required: true, message: '请输入名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9\u4e00-\u9fa5_-]+$/, message: '名称只能包含中文、英文、数字、下划线和横线', trigger: 'blur' }
  ],
  webhook_url: [
    { required: true, message: '请输入机器人 Webhook 地址', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' }
  ],
  notify_events: [
    { required: true, message: '请选择通知事件', trigger: 'change' }
  ],
  receivers: [
    { required: true, message: '请选择接收者', trigger: 'change' }
  ]
}

// 编辑状态
const isEdit = ref(false)
const currentEditId = ref<number>()

// 获取通知列表
const handleSearch = async () => {
  loading.value = true
  try {
    const response = await getNotificationList({
      page: currentPage.value,
      page_size: pageSize.value
    })
    if (response.code === 0 && response.data) {
      tableData.value = response.data.list
      total.value = response.data.total
      currentPage.value = response.data.page
      pageSize.value = response.data.pageSize
    }
  } catch (error) {
    console.error('获取通知列表失败:', error)
    ElMessage.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

// 分页处理
const handleSizeChange = (val: number) => {
  pageSize.value = val
  handleSearch()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  handleSearch()
}

// 重置搜索
const handleReset = () => {
  searchName.value = ''
  currentPage.value = 1
  handleSearch()
}

// 获取用户列表
const fetchUserList = async () => {
  try {
    const response = await getUserList({ page: 1, pageSize: 1000 })
    if (response.data?.code === 0 && response.data.data) {
      userOptions.value = response.data.data.list.map((user: { userName: string; nickName: string }) => ({
        label: user.nickName,
        value: user.userName
      }))
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  }
}

// 编辑通知
const handleEdit = async (row: NotificationItem) => {
  await fetchUserList()
  isEdit.value = true
  currentEditId.value = row.id
  feishuForm.value = {
    name: row.name,
    type: 'feishu',
    enabled: true,
    webhook_url: row.robot_url,
    description: '',
    tags: [],
    notify_events: row.notification_policy.split(','),
    receivers: [],
    send_daily_stats: true
  }
  feishuDialogVisible.value = true
}

// 提交表单
const submitFeishuForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate(async (valid) => {
    if (valid) {
      formLoading.value = true
      try {
        let response
        if (isEdit.value && currentEditId.value) {
          response = await updateFeiShu({
            id: currentEditId.value,
            name: feishuForm.value.name,
            notification_policy: feishuForm.value.notify_events.join(','),
            robot_url: feishuForm.value.webhook_url,
            send_daily_stats: feishuForm.value.send_daily_stats,
            card_content: {
              alert_level: 'warning',
              alert_name: '测试告警',
              notification_policy: 'immediate',
              alert_content: '这是一个测试告警内容',
              alert_time: new Date().toISOString(),
              notified_users: feishuForm.value.receivers,
              last_similar_alert: new Date().toISOString(),
              alert_handler: '系统管理员',
              claim_alert: false,
              resolve_alert: false,
              mute_alert: false,
              unresolved_alert: true
            }
          })
        } else {
          response = await createFeiShu(feishuForm.value)
        }

        if (response.code === 0) {
          ElMessage.success(isEdit.value ? '更新成功' : '添加成功')
          feishuDialogVisible.value = false
          handleSearch()
          formEl.resetFields()
          isEdit.value = false
          currentEditId.value = undefined
        } else {
          const errorMsg = response.msg || ''
          if (errorMsg.includes('Duplicate entry') && errorMsg.includes('idx_name_type_unique')) {
            ElMessage.error('该名称的飞书通知已存在，请使用其他名称')
          } else if (errorMsg.includes('飞书通知配置不存在')) {
            ElMessage.error('该通知配置已被删除或不存在')
            feishuDialogVisible.value = false
            handleSearch()
          } else {
            ElMessage.error(errorMsg || (isEdit.value ? '更新失败' : '添加失败'))
          }
        }
      } catch (error) {
        console.error(isEdit.value ? '更新失败:' : '添加失败:', error)
        ElMessage.error(isEdit.value ? '更新失败' : '添加失败')
      } finally {
        formLoading.value = false
      }
    }
  })
}

// 添加飞书通知
const handleAddFeishu = async () => {
  isEdit.value = false
  currentEditId.value = undefined
  await fetchUserList()
  feishuForm.value = {
    name: '',
    type: 'feishu',
    enabled: true,
    webhook_url: '',
    description: '',
    tags: [],
    notify_events: ['alert', 'warning'],
    receivers: [],
    send_daily_stats: true
  }
  feishuDialogVisible.value = true
}

// 删除通知
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该通知配置吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const response = await deleteNotification({ id, type: 'feishu' })
    if (response.code === 0) {
      ElMessage.success('删除成功')
      handleSearch()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 初始加载
handleSearch()

export {
  searchName,
  handleSearch,
  handleReset,
  tableData,
  handleAddFeishu,
  feishuDialogVisible,
  feishuForm,
  feishuFormRules,
  submitFeishuForm,
  loading,
  formLoading,
  currentPage,
  pageSize,
  total,
  handleSizeChange,
  handleCurrentChange,
  userOptions,
  handleDelete,
  handleEdit,
  isEdit
}
