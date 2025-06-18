import { usePermissionStore } from '@/pinia/modules/permission'

interface PermissionBinding {
  value: {
    path: string
    method: string
  }
}

export default {
  mounted(el: HTMLElement, binding: PermissionBinding) {
    const { value } = binding
    const permissionStore = usePermissionStore()
    if (!permissionStore.hasPermission(value.path, value.method)) {
      el.style.display = 'none'
    }
  }
} 