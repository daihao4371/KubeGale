import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { resetPassword, changePassword } from '@/api/system/user'

// 修改密码请求接口
export interface ChangePasswordReq {
  Password: string
  NewPassword: string
}

// 修改密码对话框可见性
export const changePasswordDialogVisible = ref(false)
export const changePasswordLoading = ref(false)
export const changePasswordFormRef = ref()

// 修改密码表单数据
export const changePasswordForm = reactive<ChangePasswordReq>({
  Password: '',
  NewPassword: ''
})

// 修改密码表单验证规则
export const changePasswordRules = {
  Password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  NewPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

// 打开修改密码对话框
export const handleChangePassword = () => {
  changePasswordDialogVisible.value = true
  // 重置表单
  if (changePasswordFormRef.value) {
    changePasswordFormRef.value.resetFields()
  }
  changePasswordForm.Password = ''
  changePasswordForm.NewPassword = ''
}

// 提交修改密码表单
export const submitChangePasswordForm = async () => {
  if (!changePasswordFormRef.value) return
  
  await changePasswordFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      changePasswordLoading.value = true
      try {
        const res = await changePassword(changePasswordForm)
        
        if (res.data && res.data.code === 0) {
          ElMessage.success('密码修改成功')
          changePasswordDialogVisible.value = false
        } else {
          ElMessage.error(res.data?.msg || '密码修改失败')
        }
      } catch (error) {
        console.error('修改密码失败:', error)
        ElMessage.error('修改密码失败，请稍后重试')
      } finally {
        changePasswordLoading.value = false
      }
    }
  })
}

// 重设密码
export const handleResetPassword = (row: { ID: number }) => {
  ElMessageBox.confirm(
    '是否将此用户密码重置为123456?',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
      draggable: true,
    }
  )
    .then(async () => {
      try {
        const response = await resetPassword({ 
          ID: row.ID, 
          newPassword: '123456' 
        });
        
        if (response.data && response.data.code === 0) {
          ElMessage({
            message: '密码重置成功',
            type: 'success'
          });
        } else {
          ElMessage({
            message: response.data?.msg || '密码重置失败',
            type: 'error'
          });
        }
      } catch (error) {
        console.error('密码重置失败:', error);
        ElMessage({
          message: '密码重置失败，请重试',
          type: 'error'
        });
      }
    })
    .catch(() => {
      // 用户取消操作，不做任何处理
    });
}; 