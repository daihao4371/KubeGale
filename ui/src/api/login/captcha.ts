import request from '@/api/request'

export interface CaptchaResponse {
  code: number
  data: {
    captchaId: string
    picPath: string
    captchaLength: number
    openCaptcha: boolean
  }
  msg: string
}

export function getCaptcha() {
  return request.post<CaptchaResponse>('/base/captcha')
}