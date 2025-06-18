import { defineStore } from 'pinia'

interface UserInfo {
  authorityId: number
  userName?: string
  nickName?: string
  [key: string]: string | number | undefined
}

interface UserState {
  userInfo: UserInfo | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    userInfo: {
      authorityId: 888 // 默认角色ID，可根据实际登录逻辑动态赋值
    }
  }),
  actions: {
    setUserInfo(info: UserInfo) {
      this.userInfo = info
    }
  }
}) 