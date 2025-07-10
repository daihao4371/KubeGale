# Kubernetes 模块 API 接口规范

## 概述

本文档定义了 KubeGale 项目中 Kubernetes 模块的完整 API 接口规范，用于前后端联调开发。

## 通用规范

### 请求格式
- 基础URL: `/api/v1/kubernetes`
- 内容类型: `application/json`
- 认证: Bearer Token (JWT)

### 响应格式
```json
{
  "code": 0,          // 状态码：0-成功，非0-失败
  "data": {},         // 响应数据
  "msg": "success"    // 响应消息
}
```

### 分页格式
```json
{
  "items": [],        // 数据列表
  "total": 100,       // 总条数
  "page": 1,          // 当前页
  "pageSize": 20      // 每页条数
}
```

## 1. 集群管理 API

### 1.1 获取集群列表
- **接口**: `GET /clusters`
- **参数**: 
  ```json
  {
    "page": 1,
    "pageSize": 20,
    "keyword": "cluster-name"
  }
  ```
- **响应**:
  ```json
  {
    "code": 0,
    "data": {
      "list": [
        {
          "id": 1,
          "name": "production-cluster",
          "description": "生产环境集群",
          "kubeconfig": "...",
          "status": "active",
          "version": "v1.28.0",
          "created_at": "2024-01-01T00:00:00Z"
        }
      ],
      "total": 1,
      "page": 1,
      "pageSize": 20
    },
    "msg": "success"
  }
  ```

### 1.2 获取集群详情
- **接口**: `GET /clusters/{id}`
- **响应**: 包含完整集群信息

### 1.3 创建集群
- **接口**: `POST /clusters`
- **参数**:
  ```json
  {
    "name": "new-cluster",
    "description": "新集群",
    "kubeconfig": "..."
  }
  ```

### 1.4 更新集群
- **接口**: `PUT /clusters/{id}`

### 1.5 删除集群
- **接口**: `DELETE /clusters/{id}`

### 1.6 获取集群命名空间列表
- **接口**: `GET /clusters/{id}/namespaces`
- **响应**:
  ```json
  {
    "code": 0,
    "data": {
      "namespaces": [
        {
          "metadata": {
            "name": "default",
            "uid": "...",
            "creationTimestamp": "2024-01-01T00:00:00Z"
          },
          "status": {
            "phase": "Active"
          }
        }
      ]
    }
  }
  ```

## 2. 工作负载管理 API

### 2.1 Pod 管理

#### 2.1.1 获取 Pod 列表
- **接口**: `GET /pods`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "namespace": "default",
    "page": 1,
    "pageSize": 20,
    "keyword": "pod-name",
    "labelSelector": "app=nginx",
    "fieldSelector": "status.phase=Running"
  }
  ```
- **响应**:
  ```json
  {
    "code": 0,
    "data": {
      "items": [
        {
          "metadata": {
            "name": "nginx-pod",
            "namespace": "default",
            "uid": "...",
            "creationTimestamp": "2024-01-01T00:00:00Z",
            "labels": {
              "app": "nginx"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "nginx",
                "image": "nginx:1.20",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ]
              }
            ],
            "nodeName": "worker-node-1"
          },
          "status": {
            "phase": "Running",
            "podIP": "10.244.1.1",
            "hostIP": "192.168.1.10",
            "containerStatuses": [
              {
                "name": "nginx",
                "ready": true,
                "restartCount": 0,
                "state": {
                  "running": {
                    "startedAt": "2024-01-01T00:00:00Z"
                  }
                }
              }
            ]
          }
        }
      ],
      "total": 1,
      "page": 1,
      "pageSize": 20
    }
  }
  ```

#### 2.1.2 获取 Pod 详情
- **接口**: `GET /pods/{name}`
- **参数**: `cluster_id`, `namespace`

#### 2.1.3 创建 Pod
- **接口**: `POST /pods`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "namespace": "default",
    "content": {
      "apiVersion": "v1",
      "kind": "Pod",
      "metadata": {
        "name": "new-pod"
      },
      "spec": {
        "containers": [...]
      }
    }
  }
  ```

#### 2.1.4 删除 Pod
- **接口**: `DELETE /pods/{name}`
- **参数**: `cluster_id`, `namespace`

#### 2.1.5 重启 Pod
- **接口**: `POST /pods/{name}/restart`
- **参数**: `cluster_id`, `namespace`

#### 2.1.6 获取 Pod 日志
- **接口**: `GET /pods/{name}/logs`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "namespace": "default",
    "container": "nginx",
    "tail_lines": 100,
    "follow": false
  }
  ```
- **响应**:
  ```json
  {
    "code": 0,
    "data": {
      "logs": "2024-01-01 00:00:00 [info] nginx started\n..."
    }
  }
  ```

#### 2.1.7 获取 Pod 终端
- **接口**: `GET /pods/{name}/terminal`
- **参数**: `cluster_id`, `namespace`, `container`
- **响应**:
  ```json
  {
    "code": 0,
    "data": {
      "terminal_url": "ws://localhost:8080/ws/terminal/...",
      "token": "..."
    }
  }
  ```

#### 2.1.8 获取 Pod 事件
- **接口**: `GET /pods/{name}/events`
- **参数**: `cluster_id`, `namespace`

#### 2.1.9 获取 Pod 指标
- **接口**: `GET /pods/metrics`
- **参数**: `cluster_id`, `namespace`

### 2.2 Deployment 管理

#### 2.2.1 获取 Deployment 列表
- **接口**: `GET /deployments`
- **参数**: 同 Pod 列表参数

#### 2.2.2 获取 Deployment 详情
- **接口**: `GET /deployments/{name}`

#### 2.2.3 创建 Deployment
- **接口**: `POST /deployments`

#### 2.2.4 更新 Deployment
- **接口**: `PUT /deployments/{name}`

#### 2.2.5 删除 Deployment
- **接口**: `DELETE /deployments/{name}`

#### 2.2.6 扩缩容 Deployment
- **接口**: `POST /deployments/{name}/scale`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "namespace": "default",
    "replicas": 3
  }
  ```

#### 2.2.7 回滚 Deployment
- **接口**: `POST /deployments/{name}/rollback`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "namespace": "default",
    "revision": 1
  }
  ```

#### 2.2.8 获取 Deployment 历史版本
- **接口**: `GET /deployments/{name}/revisions`

### 2.3 StatefulSet 管理
- 接口规范同 Deployment，路径替换为 `/statefulsets`

### 2.4 DaemonSet 管理
- 接口规范同 Deployment，路径替换为 `/daemonsets`

### 2.5 Job 管理
- 接口规范同 Deployment，路径替换为 `/jobs`

### 2.6 CronJob 管理
- 接口规范同 Deployment，路径替换为 `/cronjobs`

## 3. 服务与网络 API

### 3.1 Service 管理

#### 3.1.1 获取 Service 列表
- **接口**: `GET /services`

#### 3.1.2 获取 Service 详情
- **接口**: `GET /services/{name}`

#### 3.1.3 创建 Service
- **接口**: `POST /services`

#### 3.1.4 更新 Service
- **接口**: `PUT /services/{name}`

#### 3.1.5 删除 Service
- **接口**: `DELETE /services/{name}`

### 3.2 Ingress 管理
- 接口规范同 Service，路径替换为 `/ingresses`

### 3.3 Endpoint 管理
- 接口规范同 Service，路径替换为 `/endpoints`

## 4. 存储管理 API

### 4.1 PersistentVolume 管理
- **接口**: `GET /persistentvolumes`
- **接口**: `GET /persistentvolumes/{name}`
- **接口**: `POST /persistentvolumes`
- **接口**: `PUT /persistentvolumes/{name}`
- **接口**: `DELETE /persistentvolumes/{name}`

### 4.2 PersistentVolumeClaim 管理
- **接口**: `GET /persistentvolumeclaims`
- 其他接口同 PV

### 4.3 StorageClass 管理
- **接口**: `GET /storageclasses`
- 其他接口同 PV

## 5. 配置管理 API

### 5.1 ConfigMap 管理
- **接口**: `GET /configmaps`
- **接口**: `GET /configmaps/{name}`
- **接口**: `POST /configmaps`
- **接口**: `PUT /configmaps/{name}`
- **接口**: `DELETE /configmaps/{name}`

### 5.2 Secret 管理
- **接口**: `GET /secrets`
- 其他接口同 ConfigMap

## 6. 权限管理 API

### 6.1 Role 管理
- **接口**: `GET /roles`
- **接口**: `GET /roles/{name}`
- **接口**: `POST /roles`
- **接口**: `PUT /roles/{name}`
- **接口**: `DELETE /roles/{name}`

### 6.2 ClusterRole 管理
- **接口**: `GET /clusterroles`
- 其他接口同 Role

### 6.3 RoleBinding 管理
- **接口**: `GET /rolebindings`
- 其他接口同 Role

### 6.4 ClusterRoleBinding 管理
- **接口**: `GET /clusterrolebindings`
- 其他接口同 Role

### 6.5 ServiceAccount 管理
- **接口**: `GET /serviceaccounts`
- 其他接口同 Role

## 7. 节点管理 API

### 7.1 获取节点列表
- **接口**: `GET /nodes`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "page": 1,
    "pageSize": 20,
    "keyword": "node-name"
  }
  ```

### 7.2 获取节点详情
- **接口**: `GET /nodes/{name}`

### 7.3 更新节点
- **接口**: `PUT /nodes/{name}`

### 7.4 删除节点
- **接口**: `DELETE /nodes/{name}`

### 7.5 节点调度控制
- **接口**: `POST /nodes/{name}/schedule`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "schedulable": true
  }
  ```

### 7.6 节点 Pod 驱逐
- **接口**: `POST /nodes/{name}/drain`

### 7.7 获取节点指标
- **接口**: `GET /nodes/metrics`

### 7.8 节点终端
- **接口**: `GET /nodes/{name}/terminal`

## 8. 命名空间管理 API

### 8.1 获取命名空间列表
- **接口**: `GET /namespaces`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "page": 1,
    "pageSize": 20,
    "keyword": "namespace-name"
  }
  ```

### 8.2 获取命名空间详情
- **接口**: `GET /namespaces/{name}`

### 8.3 创建命名空间
- **接口**: `POST /namespaces`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "content": {
      "apiVersion": "v1",
      "kind": "Namespace",
      "metadata": {
        "name": "new-namespace"
      }
    }
  }
  ```

### 8.4 删除命名空间
- **接口**: `DELETE /namespaces/{name}`

### 8.5 获取命名空间资源配额
- **接口**: `GET /namespaces/{name}/resourcequotas`

### 8.6 获取命名空间限制范围
- **接口**: `GET /namespaces/{name}/limitranges`

## 9. 监控与指标 API

### 9.1 获取集群指标
- **接口**: `GET /metrics/cluster`
- **参数**: `cluster_id`

### 9.2 获取节点指标
- **接口**: `GET /metrics/nodes`
- **参数**: `cluster_id`, `node_name` (可选)

### 9.3 获取 Pod 指标
- **接口**: `GET /metrics/pods`
- **参数**: `cluster_id`, `namespace`, `pod_name` (可选)

### 9.4 获取资源使用统计
- **接口**: `GET /metrics/stats`
- **参数**: `cluster_id`, `time_range`

## 10. 事件管理 API

### 10.1 获取事件列表
- **接口**: `GET /events`
- **参数**:
  ```json
  {
    "cluster_id": 1,
    "namespace": "default",
    "resource_type": "Pod",
    "resource_name": "nginx-pod",
    "page": 1,
    "pageSize": 20
  }
  ```

## 错误代码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 集群不存在 |
| 1003 | 命名空间不存在 |
| 1004 | 资源不存在 |
| 1005 | 权限不足 |
| 1006 | 集群连接失败 |
| 1007 | 资源创建失败 |
| 1008 | 资源更新失败 |
| 1009 | 资源删除失败 |
| 2001 | 认证失败 |
| 2002 | Token 过期 |
| 5001 | 内部服务器错误 |

## 待实现功能

1. WebSocket 实时日志流
2. 文件上传下载到 Pod
3. 资源模板管理
4. YAML 编辑器集成
5. 多集群统一管理
6. 资源拓扑图
7. 告警规则管理
8. 自动化运维任务

## 注意事项

1. 所有涉及集群操作的接口都需要验证集群连接状态
2. 敏感操作（删除、更新）需要增加确认机制
3. 大量数据返回时要考虑分页和性能优化
4. WebSocket 连接需要处理断线重连
5. 文件操作需要考虑安全性验证