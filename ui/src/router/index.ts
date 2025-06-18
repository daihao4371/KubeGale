import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'  // 修改为重定向到仪表盘
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/login/Login.vue')
    },
    {
      path: '/terminal',
      name: 'Terminal',
      component: () => import('@/views/terminal/index.vue'),
      meta: {
        title: '终端',
        requiresAuth: true
      }
    },
    {
      path: '/homepage',
      name: 'homepage',
      component: () => import('../views/homepage/Homepage.vue'),
      meta: {
        requiresAuth: true
      },
      children: [
        // 添加仪表盘路由
        {
          path: '/dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/index.vue'),
          meta: {
            requiresAuth: true,
            title: '仪表盘'
          }
        },
        // 系统管理相关路由
        {
          path: '/system',
          component: () => import('@/views/system/index.vue'),
          children: [
            {
              path: '',
              redirect: '/system/userManager'
            },
            {
              path: 'userManager',
              component: () => import('@/views/system/userManager/index.vue')
            },
            // 添加角色管理路由
            {
              path: 'roleManager',
              component: () => import('@/views/system/roleManager/index.vue'),
              meta: {
                requiresAuth: true,
                title: '角色管理'
              }
            },
            // 添加API管理路由
            {
              path: 'apiManager',
              component: () => import('@/views/system/apiManager/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'API管理'
              }
            },
            // 添加操作记录路由
            {
              path: 'menuManager',
              component: () => import('@/views/system/operationRecord/index.vue'),
              meta: {
                requiresAuth: true,
                title: '操作记录'
              }
            }
          ]
        },
        // 添加IM通知管理路由
        {
          path: 'im',
          name: 'im',
          component: () => import('@/views/im/index.vue'),
          meta: {
            requiresAuth: true,
            title: 'IM通知管理'
          }
        },
        // 添加云资产管理路由
        {
          path: 'cloud-assets',
          name: 'cloud-assets',
          component: () => import('@/views/cloudCmdb/index.vue'),
          meta: {
            requiresAuth: true,
            title: '云资产管理'
          },
          children: [
            {
              path: '',
              name: 'cloud-assets-default',
              redirect: '/homepage/cloud-assets/provider'
            },
            {
              path: 'provider',
              name: 'cloud-provider',
              component: () => import('@/views/cloudCmdb/provider/index.vue'),
              meta: {
                requiresAuth: true,
                title: '云厂商管理'
              }
            },
            {
              path: 'virtualMachine',
              name: 'cloud-virtual-machine',
              component: () => import('@/views/cloudCmdb/virtualMachine/inde.vue'),
              meta: {
                requiresAuth: true,
                title: '云服务器'
              }
            },
            {
              path: 'loadbalancer',
              name: 'cloud-loadbalancer',
              component: () => import('@/views/cloudCmdb/loadbalancer/index.vue'),
              meta: {
                requiresAuth: true,
                title: '负载均衡'
              }
            },
            {
              path: 'rds',
              name: 'cloud-rds',
              component: () => import('@/views/cloudCmdb/rds/index.vue'),
              meta: {
                requiresAuth: true,
                title: '云数据库'
              }
            }
          ]
        },
        {
          path: 'cmdb',
          name: 'cmdb',
          component: () => import('../views/cmdb/index.vue'),
          meta: {
            requiresAuth: true,
            title: 'CMDB资产管理'
          },
          children: [
            {
              path: '',
              name: 'cmdb-default',
              redirect: '/homepage/cmdb/project'
            },
            {
              path: 'project',
              name: 'cmdb-project',
              component: () => import('../views/cmdb/project/index.vue'),
              meta: {
                requiresAuth: true,
                title: '项目管理'
              }
            },
            {
              path: 'host',
              name: 'cmdb-host',
              component: () => import('../views/cmdb/host/index.vue'),
              meta: {
                requiresAuth: true,
                title: '主机管理'
              }
            },
            {
              path: 'batchtask',
              name: 'cmdb-batchtask',
              component: () => import('../views/cmdb/batchtask/index.vue'),
              meta: {
                requiresAuth: true,
                title: '批量任务'
              }
            }
          ]
        },
        {
          path: 'kubernetes',
          name: 'kubernetes',
          component: () => import('../views/kubernetes/index.vue'),
          meta: {
            requiresAuth: true
          },
          children: [
            {
              path: '',
              redirect: '/homepage/kubernetes/cluster'
            },
            {
              path: 'cluster',
              name: 'kubernetes-cluster',
              component: () => import('../views/kubernetes/k8sCluster/k8sCluster.vue'),
              meta: {
                requiresAuth: true,
                title: '集群管理'
              }
            }
          ]
        },
        {
          path: 'prometheus',
          name: 'prometheus',
          component: () => import('../views/prometheus/index.vue'),
          meta: {
            requiresAuth: true
          }
        },
        {
          path: 'config',
          name: 'config',
          component: () => import('../views/config/index.vue'),
          meta: {
            requiresAuth: true
          }
        },
        {
          path: 'docker',
          name: 'docker',
          component: () => import('../views/docker/index.vue'),
          meta: {
            requiresAuth: true
          }
        },
        {
          path: 'cicd',
          name: 'cicd',
          component: () => import('../views/cicd/index.vue'),
          meta: {
            requiresAuth: true
          }
        }
      ]
    }
  ]
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - KubeGale` : 'KubeGale'
  
  // 检查路由是否需要认证
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // 检查用户是否已登录
    const token = localStorage.getItem('token')
    if (!token) {
      ElMessage.warning('请先登录')
      // 重定向到登录页
      next({ path: '/login' })
    } else {
      // 已登录，允许访问
      next()
    }
  } else {
    // 不需要认证的路由，直接放行
    next()
  }
})

export default router
