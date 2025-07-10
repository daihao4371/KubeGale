import { ref, reactive, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { logout } from '@/api/login/login'
import { 
  Setting, 
  Document, 
  Ship, 
  Monitor, 
  Tools, 
  Box, 
  Connection,
  SwitchButton,
  ArrowDown,
  ArrowRight,
  User,
  UserFilled,
  Menu,
  HomeFilled,
  DataAnalysis,
  ChatDotRound,
  Cloudy,
  Operation,
  Folder,
  Files,
  Grid,
  List,
  Histogram,
  Bell,
  Platform,
  Odometer,
  SetUp,
  Goods,
  Timer
} from '@element-plus/icons-vue'

interface MenuItem {
  id: string
  title: string
  icon: string
  path: string
  children?: MenuItem[]
}

interface BreadcrumbItem {
  title: string
  path: string
  timestamp: number
}

export default function useHomepage() {
  const router = useRouter()
  const getInitialActiveMenu = () => {
    return localStorage.getItem('activeMenu') || 'dashboard'
  }
  const activeMenu = ref(getInitialActiveMenu())
  const recentPages = ref<BreadcrumbItem[]>([])
  const maxRecentPages = 10

  const menuItems: MenuItem[] = reactive([
    { id: 'dashboard', title: '仪表盘', icon: 'Odometer', path: '/dashboard' },

    { 
      id: 'system', 
      title: '系统管理', 
      icon: 'Setting', 
      path: '/system',
      children: [
        { id: 'system-user', title: '用户管理', icon: 'User', path: '/system/userManager' },
        { id: 'system-role', title: '角色管理', icon: 'UserFilled', path: '/system/roleManager' },
        { id: 'system-api', title: 'API管理', icon: 'Connection', path: '/system/apiManager' },
        { id: 'system-menu', title: '操作记录', icon: 'Timer', path: '/system/menuManager' },
      ]
    },
    { 
      id: 'cloud-assets', 
      title: '云资产管理', 
      icon: 'Cloudy', 
      path: '/homepage/cloud-assets',
      children: [
        { id: 'cloud-provider', title: '云厂商管理', icon: 'Platform', path: '/homepage/cloud-assets/provider' },
        { id: 'cloud-virtual-machine', title: '云服务器', icon: 'Monitor', path: '/homepage/cloud-assets/virtualMachine' },
        { id: 'cloud-loadbalancer', title: '负载均衡', icon: 'Operation', path: '/homepage/cloud-assets/loadbalancer' },
        { id: 'cloud-rds', title: '云数据库', icon: 'DataAnalysis', path: '/homepage/cloud-assets/rds' },
      ]
    },
    { 
      id: 'cmdb', 
      title: 'CMDB资产管理', 
      icon: 'Grid', 
      path: '/homepage/cmdb',
      children: [
        { id: 'cmdb-project', title: '项目管理', icon: 'List', path: '/homepage/cmdb/project' },
        { id: 'cmdb-host', title: '主机管理', icon: 'Monitor', path: '/homepage/cmdb/host' },
        { id: 'cmdb-batchtask', title: '批量任务管理', icon: 'Document', path: '/homepage/cmdb/batchtask' }
      ]
    },
    { id: 'im', title: 'IM通知管理', icon: 'Bell', path: '/homepage/im' },

    { 
      id: 'kubernetes', 
      title: 'k8s管理', 
      icon: 'Ship', 
      path: '/homepage/kubernetes',
      children: [
        { id: 'kubernetes-cluster', title: '集群管理', icon: 'Box', path: '/homepage/kubernetes/cluster' },
        { 
          id: 'kubernetes-workload', 
          title: '工作负载', 
          icon: 'Grid', 
          path: '/homepage/kubernetes/workload',
          children: [
            { id: 'kubernetes-pod', title: 'Pod管理', icon: 'Document', path: '/homepage/kubernetes/workload/pod' },
            { id: 'kubernetes-deployment', title: 'Deployment', icon: 'Document', path: '/homepage/kubernetes/workload/deployment' },
            { id: 'kubernetes-statefulset', title: 'StatefulSet', icon: 'Document', path: '/homepage/kubernetes/workload/statefulset' },
            { id: 'kubernetes-daemonset', title: 'DaemonSet', icon: 'Document', path: '/homepage/kubernetes/workload/daemonset' },
            { id: 'kubernetes-job', title: 'Job', icon: 'Document', path: '/homepage/kubernetes/workload/job' },
            { id: 'kubernetes-cronjob', title: 'CronJob', icon: 'Document', path: '/homepage/kubernetes/workload/cronjob' }
          ]
        },
        { 
          id: 'kubernetes-network', 
          title: '网络管理', 
          icon: 'Connection', 
          path: '/homepage/kubernetes/network',
          children: [
            { id: 'kubernetes-service', title: 'Service', icon: 'Document', path: '/homepage/kubernetes/network/service' },
            { id: 'kubernetes-ingress', title: 'Ingress', icon: 'Document', path: '/homepage/kubernetes/network/ingress' },
            { id: 'kubernetes-endpoint', title: 'Endpoint', icon: 'Document', path: '/homepage/kubernetes/network/endpoint' },
            { id: 'kubernetes-networkpolicy', title: 'NetworkPolicy', icon: 'Document', path: '/homepage/kubernetes/network/networkpolicy' }
          ]
        },
        { 
          id: 'kubernetes-storage', 
          title: '存储管理', 
          icon: 'Folder', 
          path: '/homepage/kubernetes/storage',
          children: [
            { id: 'kubernetes-pv', title: 'PersistentVolume', icon: 'Document', path: '/homepage/kubernetes/storage/pv' },
            { id: 'kubernetes-pvc', title: 'PersistentVolumeClaim', icon: 'Document', path: '/homepage/kubernetes/storage/pvc' },
            { id: 'kubernetes-storageclass', title: 'StorageClass', icon: 'Document', path: '/homepage/kubernetes/storage/storageclass' }
          ]
        },
        { 
          id: 'kubernetes-config', 
          title: '配置管理', 
          icon: 'Setting', 
          path: '/homepage/kubernetes/config',
          children: [
            { id: 'kubernetes-configmap', title: 'ConfigMap', icon: 'Document', path: '/homepage/kubernetes/config/configmap' },
            { id: 'kubernetes-secret', title: 'Secret', icon: 'Document', path: '/homepage/kubernetes/config/secret' },
            { id: 'kubernetes-resourcequota', title: 'ResourceQuota', icon: 'Document', path: '/homepage/kubernetes/config/resourcequota' },
            { id: 'kubernetes-hpa', title: 'HPA', icon: 'Document', path: '/homepage/kubernetes/config/hpa' }
          ]
        },
        { 
          id: 'kubernetes-rbac', 
          title: 'RBAC权限', 
          icon: 'User', 
          path: '/homepage/kubernetes/rbac',
          children: [
            { id: 'kubernetes-role', title: 'Role', icon: 'Document', path: '/homepage/kubernetes/rbac/role' },
            { id: 'kubernetes-clusterrole', title: 'ClusterRole', icon: 'Document', path: '/homepage/kubernetes/rbac/clusterrole' },
            { id: 'kubernetes-rolebinding', title: 'RoleBinding', icon: 'Document', path: '/homepage/kubernetes/rbac/rolebinding' },
            { id: 'kubernetes-serviceaccount', title: 'ServiceAccount', icon: 'Document', path: '/homepage/kubernetes/rbac/serviceaccount' }
          ]
        },
        { id: 'kubernetes-nodes', title: '节点管理', icon: 'Monitor', path: '/homepage/kubernetes/nodes' },
        { id: 'kubernetes-namespace', title: '命名空间', icon: 'Files', path: '/homepage/kubernetes/namespace' },
        { id: 'kubernetes-monitoring', title: '监控管理', icon: 'DataAnalysis', path: '/homepage/kubernetes/monitoring' },
        { id: 'kubernetes-terminal', title: 'Web终端', icon: 'Document', path: '/homepage/kubernetes/terminal' }
      ]
    },
    { id: 'prometheus', title: 'Prometheus监控管理', icon: 'Histogram', path: '/homepage/prometheus' },
    { id: 'config', title: '配置中心', icon: 'SetUp', path: '/homepage/config' },
    { id: 'docker', title: 'docker管理', icon: 'Goods', path: '/homepage/docker' },
    { id: 'cicd', title: 'CICD', icon: 'Timer', path: '/homepage/cicd' }
  ])

  // ----------- 菜单展开记忆 start -----------
  const getAllParentMenuIds = (items: MenuItem[]): string[] => {
    const ids: string[] = []
    items.forEach(item => {
      if (item.children && item.children.length) {
        ids.push(item.id)
      }
    })
    return ids
  }

  const getInitialExpandedMenus = () => {
    const saved = localStorage.getItem('expandedMenus')
    if (saved) {
      try {
        return JSON.parse(saved)
      } catch {
        return getAllParentMenuIds(menuItems)
      }
    }
    return getAllParentMenuIds(menuItems)
  }

  const expandedMenus = ref<string[]>(getInitialExpandedMenus())

  watch(expandedMenus, (val) => {
    localStorage.setItem('expandedMenus', JSON.stringify(val))
  }, { deep: true })
  // ----------- 菜单展开记忆 end -----------
  
  // ----------- 菜单高亮和页面记忆 start -----------
  watch(activeMenu, (val) => {
    localStorage.setItem('activeMenu', val)
  })

  // 递归查找菜单项
  const findMenuItem = (id: string, items: MenuItem[]): MenuItem | undefined => {
    for (const item of items) {
      if (item.id === id) {
        return item
      }
      if (item.children && item.children.length > 0) {
        const found = findMenuItem(id, item.children)
        if (found) {
          return found
        }
      }
    }
    return undefined
  }

  onMounted(() => {
    if (activeMenu.value) {
      const menu = findMenuItem(activeMenu.value, menuItems)
      if (menu && menu.path) {
        router.replace(menu.path)
      }
    }
  })
  // ----------- 菜单高亮和页面记忆 end -----------
  
  // 初始化时导航到仪表盘
  router.push('/dashboard')
  
  // 添加最近访问页面
  const addRecentPage = (menu: MenuItem) => {
    const newPage = {
      title: menu.title,
      path: menu.path,
      timestamp: Date.now()
    }
    
    // 移除已存在的相同路径的页面
    recentPages.value = recentPages.value.filter(page => page.path !== menu.path)
    
    // 添加到最近访问列表的开头
    recentPages.value.unshift(newPage)
    
    // 保持最多10个记录
    if (recentPages.value.length > maxRecentPages) {
      recentPages.value = recentPages.value.slice(0, maxRecentPages)
    }
  }
  
  // 更新面包屑
  const updateBreadcrumbs = (menuId: string) => {
    const selectedMenu = findMenuItem(menuId, menuItems)
    if (selectedMenu) {
      if (selectedMenu.children) {
        // 如果是父菜单，只显示父菜单
        addRecentPage(selectedMenu)
      } else {
        // 如果是子菜单，显示父菜单和子菜单
        const parentMenu = findParentMenu(menuId, menuItems)
        if (parentMenu) {
          addRecentPage(selectedMenu)
        } else {
          addRecentPage(selectedMenu)
        }
      }
    }
  }
  
  // 查找父菜单
  const findParentMenu = (childId: string, items: MenuItem[]): MenuItem | undefined => {
    for (const item of items) {
      if (item.children) {
        if (item.children.some(child => child.id === childId)) {
          return item
        }
        const found = findParentMenu(childId, item.children)
        if (found) {
          return found
        }
      }
    }
    return undefined
  }
  
  const selectMenu = (menuId: string) => {
    const selectedMenu = findMenuItem(menuId, menuItems)
    
    if (selectedMenu) {
      if (selectedMenu.children && selectedMenu.children.length > 0) {
        if (expandedMenus.value.includes(menuId)) {
          expandedMenus.value = expandedMenus.value.filter(id => id !== menuId)
        } else {
          expandedMenus.value.push(menuId)
        }
        activeMenu.value = menuId
      } else {
        activeMenu.value = menuId
        router.push(selectedMenu.path)
      }
      // 更新面包屑
      updateBreadcrumbs(menuId)
    }
  }
  
  // 检查菜单项是否有激活的子菜单
  const hasActiveChild = (item: MenuItem): boolean => {
    if (!item.children) return false
    
    return item.children.some(child => child.id === activeMenu.value) || 
           (expandedMenus.value.includes(item.id) && item.id === activeMenu.value)
  }
  
  const username = ref('Admin')
  const currentTime = ref(new Date().toLocaleString())

  // 更新时间
  setInterval(() => {
    currentTime.value = new Date().toLocaleString()
  }, 1000)

  // 退出登录方法
  const handleLogout = async () => {
    try {
      const res = await logout()
      if (res.data.code === 0) {
        ElMessage.success('退出登录成功')
        // 清除本地存储的token和用户信息
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        // 跳转到登录页
        router.push('/login')
      } else {
        ElMessage.error(res.data.msg || '退出登录失败')
      }
    } catch (error) {
      ElMessage.error('退出登录失败，请稍后重试')
    }
  }
  
  return {
    menuItems,
    activeMenu,
    expandedMenus,
    selectMenu,
    username,
    currentTime,
    handleLogout,
    hasActiveChild,
    recentPages
  }
}