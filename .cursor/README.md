# KubeGale API 测试数据

本文档提供了 KubeGale 系统各个模块的 Postman 测试数据。

## 通用说明

1. 所有请求需要在 Header 中添加：
   ```
   Content-Type: application/json
   Authorization: Bearer {token}
   ```

2. 分页参数说明：
   - page: 从 1 开始
   - pageSize: 每页数量

## 云平台管理 (Cloud Platform)

### 创建云平台
```json
POST /cloud_platform/create
{
    "name": "阿里云",
    "type": "aliyun",
    "accessKey": "your_access_key",
    "secretKey": "your_secret_key",
    "region": "cn-hangzhou",
    "description": "阿里云测试环境"
}
```

### 获取云平台列表
```json
POST /cloud_platform/list
{
    "page": 1,
    "pageSize": 10,
    "name": "",
    "type": ""
}
```

### 更新云平台
```json
PUT /cloud_platform/update
{
    "id": 1,
    "name": "阿里云更新",
    "type": "aliyun",
    "accessKey": "new_access_key",
    "secretKey": "new_secret_key",
    "region": "cn-shanghai",
    "description": "更新后的描述"
}
```

### 获取单个云平台
```json
POST /cloud_platform/getById
{
    "id": 1
}
```

### 删除云平台
```json
DELETE /cloud_platform/delete
{
    "id": 1
}
```

### 批量删除云平台
```json
DELETE /cloud_platform/deleteByIds
{
    "ids": [1, 2, 3]
}
```

## 云主机管理 (Virtual Machine)

### 同步云主机
```json
POST /virtualMachine/sync
{
    "platformId": 1,
    "regionId": "cn-hangzhou"
}
```

### 获取云主机目录树
```json
POST /virtualMachine/tree
{
    "platformId": 1
}
```

### 获取云主机列表
```json
POST /virtualMachine/list
{
    "page": 1,
    "pageSize": 10,
    "platformId": 1,
    "regionId": "",
    "instanceId": "",
    "instanceName": ""
}
```

## 负载均衡管理 (Load Balancer)

### 同步负载均衡
```json
POST /loadBalancer/sync
{
    "platformId": 1,
    "regionId": "cn-hangzhou"
}
```

### 获取负载均衡目录树
```json
POST /loadBalancer/tree
{
    "platformId": 1
}
```

### 获取负载均衡列表
```json
POST /loadBalancer/list
{
    "page": 1,
    "pageSize": 10,
    "platformId": 1,
    "regionId": "",
    "loadBalancerId": "",
    "loadBalancerName": ""
}
```

## RDS管理

### 同步RDS
```json
POST /rds/sync
{
    "platformId": 1,
    "regionId": "cn-hangzhou"
}
```

### 获取RDS目录树
```json
POST /rds/tree
{
    "platformId": 1
}
```

### 获取RDS列表
```json
POST /rds/list
{
    "page": 1,
    "pageSize": 10,
    "platformId": 1,
    "regionId": "",
    "dbInstanceId": "",
    "dbInstanceName": ""
}
```

## 云区域管理

### 同步区域信息
```json
POST /cloud_region/syncRegion
{
    "platformId": 1
}
```

## CMDB项目管理

### 创建项目
```json
POST /cmdb/projects
{
    "name": "测试项目",
    "description": "这是一个测试项目",
    "owner": "admin",
    "status": "active"
}
```

### 获取项目列表
```json
GET /cmdb/projects
{
    "page": 1,
    "pageSize": 10,
    "name": "",
    "status": ""
}
```

### 更新项目
```json
PUT /cmdb/projects
{
    "id": 1,
    "name": "更新后的项目",
    "description": "更新后的描述",
    "owner": "admin",
    "status": "inactive"
}
```

### 删除项目
```json
DELETE /cmdb/projects
{
    "id": 1
}
```

### 批量删除项目
```json
DELETE /cmdb/projectsByIds
{
    "ids": [1, 2, 3]
}
```

## CMDB主机管理

### 创建主机
```json
POST /cmdb/hosts
{
    "name": "测试主机",
    "ip": "192.168.1.100",
    "port": 22,
    "username": "root",
    "password": "password",
    "projectId": 1,
    "description": "测试主机描述"
}
```

### SSH认证主机
```json
POST /cmdb/hosts/authentication
{
    "id": 1,
    "password": "new_password"
}
```

### 批量导入主机
```json
POST /cmdb/hosts/import
{
    "projectId": 1,
    "hosts": [
        {
            "name": "主机1",
            "ip": "192.168.1.101",
            "port": 22,
            "username": "root",
            "password": "password1"
        },
        {
            "name": "主机2",
            "ip": "192.168.1.102",
            "port": 22,
            "username": "root",
            "password": "password2"
        }
    ]
}
```

### 获取主机列表
```json
GET /cmdb/hosts
{
    "page": 1,
    "pageSize": 10,
    "name": "",
    "ip": "",
    "projectId": 1
}
```

## 批量操作

### 执行批量命令
```json
POST /cmdb/batchOperations/execute
{
    "projectId": 1,
    "hostIds": [1, 2, 3],
    "command": "ls -l",
    "timeout": 30
}
```

### 获取执行记录
```json
GET /cmdb/batchOperations/execLogs
{
    "page": 1,
    "pageSize": 10,
    "projectId": 1,
    "hostId": 1
}
```

## 注意事项

1. 所有请求都需要先进行登录认证，获取 token
2. 密码等敏感信息建议使用加密传输
3. 时间戳使用毫秒级时间戳
4. 实际使用时需要根据实际情况修改参数值
5. 分页接口的 page 从 1 开始 