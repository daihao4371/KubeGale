import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { login } from '@/api/login/login'
import { getCaptcha } from '@/api/login/captcha'

export default function useLogin() {
  const router = useRouter()
  const captchaImg = ref('')
  const captchaId = ref('')

  // 登录表单
  const loginForm = reactive({
    username: '',
    password: '',
    captcha: '',
    captchaId: ''
  })

  // 表单验证规则
  const rules = {
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
  }

  // 获取验证码
  const refreshCaptcha = async () => {
    try {
      const res = await getCaptcha()
      if (res.data.code === 0) {
        captchaImg.value = res.data.data.picPath
        captchaId.value = res.data.data.captchaId
        loginForm.captchaId = res.data.data.captchaId
      }
    } catch (error: unknown) {
      console.error('获取验证码失败:', error)
      ElMessage.error('获取验证码失败')
    }
  }

  // 登录方法
  const handleLogin = async (formEl: FormInstance) => {
    if (!formEl) return
    await formEl.validate(async (valid: boolean) => {
      if (valid) {
        try {
          const res = await login(loginForm)
          if (res.data.code === 0) {
            // 登录成功，保存token和用户信息
            localStorage.setItem('token', res.data.data.token)
            localStorage.setItem('userInfo', JSON.stringify(res.data.data.user))
            ElMessage.success('登录成功')
            // 跳转到首页
            router.push('/homepage')
          } else {
            ElMessage.error(res.data.msg || '登录失败')
            refreshCaptcha()
          }
        } catch (error) {
          ElMessage.error('登录失败，请稍后重试')
          refreshCaptcha()
        }
      }
    })
  }

  // 初始化获取验证码
  refreshCaptcha()

  return {
    loginForm,
    rules,
    captchaImg,
    refreshCaptcha,
    handleLogin
  }
}
