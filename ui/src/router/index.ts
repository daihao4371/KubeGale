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
            // 集群管理 - 已存在
            {
              path: 'cluster',
              name: 'kubernetes-cluster',
              component: () => import('../views/kubernetes/cluster/index.vue'),
              meta: {
                requiresAuth: true,
                title: '集群管理'
              }
            },
            {
              path: 'cluster/:id',
              name: 'kubernetes-cluster-detail',
              component: () => import('../views/kubernetes/cluster/detail.vue'),
              meta: {
                requiresAuth: true,
                title: '集群详情'
              }
            },
            // 工作负载管理 - 已存在
            {
              path: 'workload',
              name: 'kubernetes-workload',
              component: () => import('../views/kubernetes/workload/index.vue'),
              meta: {
                requiresAuth: true,
                title: '工作负载'
              }
            },
            {
              path: 'workload/pod',
              name: 'kubernetes-pod',
              component: () => import('../views/kubernetes/workload/pod/PodManagement.vue'),
              meta: {
                requiresAuth: true,
                title: 'Pod管理'
              }
            },
            {
              path: 'workload/deployment',
              name: 'kubernetes-deployment',
              component: () => import('../views/kubernetes/workload/deployment/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Deployment管理'
              }
            },
            {
              path: 'workload/statefulset',
              name: 'kubernetes-statefulset',
              component: () => import('../views/kubernetes/workload/statefulset/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'StatefulSet管理'
              }
            },
            {
              path: 'workload/daemonset',
              name: 'kubernetes-daemonset',
              component: () => import('../views/kubernetes/workload/daemonset/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'DaemonSet管理'
              }
            },
            {
              path: 'workload/job',
              name: 'kubernetes-job',
              component: () => import('../views/kubernetes/workload/job/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Job管理'
              }
            },
            {
              path: 'workload/cronjob',
              name: 'kubernetes-cronjob',
              component: () => import('../views/kubernetes/workload/cronjob/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'CronJob管理'
              }
            },
            // 网络管理 - 已存在
            {
              path: 'network',
              name: 'kubernetes-network',
              component: () => import('../views/kubernetes/network/index.vue'),
              meta: {
                requiresAuth: true,
                title: '网络管理'
              }
            },
            {
              path: 'network/service',
              name: 'kubernetes-service',
              component: () => import('../views/kubernetes/network/service/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Service管理'
              }
            },
            {
              path: 'network/ingress',
              name: 'kubernetes-ingress',
              component: () => import('../views/kubernetes/network/ingress/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Ingress管理'
              }
            },
            {
              path: 'network/endpoint',
              name: 'kubernetes-endpoint',
              component: () => import('../views/kubernetes/network/endpoint/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Endpoint管理'
              }
            },
            {
              path: 'network/networkpolicy',
              name: 'kubernetes-networkpolicy',
              component: () => import('../views/kubernetes/network/networkpolicy/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'NetworkPolicy管理'
              }
            },
            // 存储管理 - 已存在
            {
              path: 'storage',
              name: 'kubernetes-storage',
              component: () => import('../views/kubernetes/storage/index.vue'),
              meta: {
                requiresAuth: true,
                title: '存储管理'
              }
            },
            {
              path: 'storage/pv',
              name: 'kubernetes-pv',
              component: () => import('../views/kubernetes/storage/pv/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'PersistentVolume管理'
              }
            },
            {
              path: 'storage/pvc',
              name: 'kubernetes-pvc',
              component: () => import('../views/kubernetes/storage/pvc/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'PersistentVolumeClaim管理'
              }
            },
            {
              path: 'storage/storageclass',
              name: 'kubernetes-storageclass',
              component: () => import('../views/kubernetes/storage/storageclass/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'StorageClass管理'
              }
            },
            // 配置管理 - 已存在
            {
              path: 'config',
              name: 'kubernetes-config',
              component: () => import('../views/kubernetes/config/index.vue'),
              meta: {
                requiresAuth: true,
                title: '配置管理'
              }
            },
            {
              path: 'config/configmap',
              name: 'kubernetes-configmap',
              component: () => import('../views/kubernetes/config/configmap/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'ConfigMap管理'
              }
            },
            {
              path: 'config/secret',
              name: 'kubernetes-secret',
              component: () => import('../views/kubernetes/config/secret/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Secret管理'
              }
            },
            {
              path: 'config/resourcequota',
              name: 'kubernetes-resourcequota',
              component: () => import('../views/kubernetes/config/resourcequota/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'ResourceQuota管理'
              }
            },
            {
              path: 'config/hpa',
              name: 'kubernetes-hpa',
              component: () => import('../views/kubernetes/config/hpa/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'HPA管理'
              }
            },
            // RBAC权限管理 - 已存在
            {
              path: 'rbac',
              name: 'kubernetes-rbac',
              component: () => import('../views/kubernetes/rbac/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'RBAC权限管理'
              }
            },
            {
              path: 'rbac/role',
              name: 'kubernetes-role',
              component: () => import('../views/kubernetes/rbac/role/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Role管理'
              }
            },
            {
              path: 'rbac/clusterrole',
              name: 'kubernetes-clusterrole',
              component: () => import('../views/kubernetes/rbac/clusterrole/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'ClusterRole管理'
              }
            },
            {
              path: 'rbac/rolebinding',
              name: 'kubernetes-rolebinding',
              component: () => import('../views/kubernetes/rbac/rolebinding/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'RoleBinding管理'
              }
            },
            {
              path: 'rbac/serviceaccount',
              name: 'kubernetes-serviceaccount',
              component: () => import('../views/kubernetes/rbac/serviceaccount/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'ServiceAccount管理'
              }
            },
            // 节点管理 - 已存在
            {
              path: 'nodes',
              name: 'kubernetes-nodes',
              component: () => import('../views/kubernetes/nodes/index.vue'),
              meta: {
                requiresAuth: true,
                title: '节点管理'
              }
            },
            {
              path: 'nodes/:name',
              name: 'kubernetes-node-detail',
              component: () => import('../views/kubernetes/nodes/detail.vue'),
              meta: {
                requiresAuth: true,
                title: '节点详情'
              }
            },
            {
              path: 'nodes/:name/monitor',
              name: 'kubernetes-node-monitor',
              component: () => import('../views/kubernetes/nodes/monitor.vue'),
              meta: {
                requiresAuth: true,
                title: '节点监控'
              }
            },
            // 命名空间管理 - 已存在
            {
              path: 'namespace',
              name: 'kubernetes-namespace',
              component: () => import('../views/kubernetes/namespace/index.vue'),
              meta: {
                requiresAuth: true,
                title: '命名空间管理'
              }
            },
            // 监控管理 - 已存在
            {
              path: 'monitoring',
              name: 'kubernetes-monitoring',
              component: () => import('../views/kubernetes/monitoring/index.vue'),
              meta: {
                requiresAuth: true,
                title: '监控管理'
              }
            },
            // 终端 - 已存在
            {
              path: 'terminal',
              name: 'kubernetes-terminal',
              component: () => import('../views/kubernetes/terminal/index.vue'),
              meta: {
                requiresAuth: true,
                title: 'Web终端'
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
