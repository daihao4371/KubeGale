import { defineStore } from 'pinia'
import { getApisByAuthority } from '@/api/system/authorityApi'

interface ApiItem {
  path: string
  method: string
  [key: string]: unknown
}

interface PermissionState {
  apiList: ApiItem[]
}

export const usePermissionStore = defineStore('permission', {
  state: (): PermissionState => ({
    apiList: []
  }),
  actions: {
    async fetchApiList(authorityId: number) {
      const res = await getApisByAuthority(authorityId)
      if (res?.data?.code === 0) {
        this.apiList = res.data.data as ApiItem[]
      }
    },
    hasPermission(path: string, method: string): boolean {
      return this.apiList.some(api => api.path === path && api.method === method)
    }
  }
}) 