# KubeGale 项目文档

## 项目概述

KubeGale 是一个基于 Vue 3 + Go 的 Kubernetes 集群管理平台，提供完整的 K8s 资源管理和监控功能。

## 文档结构

```
docs/
├── README.md                           # 本文档
├── kubernetes-frontend-implementation.md  # K8s模块前端实现文档
├── api-documentation.md               # API接口文档 (待创建)
├── deployment-guide.md                # 部署指南 (待创建)
└── user-manual.md                     # 用户手册 (待创建)
```

## 当前状态

### 后端开发状态 ✅
- [x] 集群管理 API 完成
- [x] 工作负载管理 API 完成
- [x] 节点管理 API 完成
- [x] 命名空间管理 API 完成
- [x] 网络管理 API 完成
- [x] 存储管理 API 完成
- [x] 配置管理 API 完成
- [x] 角色管理 API 完成
- [x] 监控指标 API 完成
- [x] 终端管理 API 完成

### 前端开发状态 🚧
- [x] 项目基础架构搭建
- [x] 集群管理基础页面 (部分完成)
- [x] 节点管理页面 (部分完成)
- [x] **Deployment 管理页面** ✅ **已完成**
  - [x] Deployment 列表页面
  - [x] Deployment 创建/编辑表单
  - [x] Deployment 详情查看
  - [x] Deployment 扩缩容功能
  - [x] Deployment 删除功能
  - [x] 相关 API 接口和类型定义
- [ ] Pod 管理页面 (待开发)
- [ ] 命名空间管理页面 (待开发)
- [ ] 网络管理页面 (待开发)
- [ ] 存储管理页面 (待开发)
- [ ] 配置管理页面 (待开发)
- [ ] 角色管理页面 (待开发)
- [ ] 监控功能页面 (待开发)
- [ ] 终端功能页面 (待开发)

## 开发计划

### 第一阶段：核心功能实现 (Week 1-2) ✅
- [x] 完善集群管理功能
- [x] 实现 Deployment 管理
- [ ] 实现 Pod 管理
- [ ] 完善节点管理功能

### 第二阶段：基础资源管理 (Week 3-4)
- [ ] 实现命名空间管理
- [ ] 实现 Service 管理
- [ ] 实现 ConfigMap/Secret 管理
- [ ] 实现基础监控功能

### 第三阶段：高级功能实现 (Week 5-6)
- [ ] 实现其他工作负载类型
- [ ] 实现存储管理
- [ ] 实现角色管理
- [ ] 实现终端功能

### 第四阶段：优化和测试 (Week 7-8)
- [ ] 性能优化
- [ ] 用户体验优化
- [ ] 功能测试
- [ ] 文档完善

## 技术栈

### 前端技术栈
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **代码规范**: ESLint + Prettier

### 后端技术栈
- **语言**: Go 1.24.3
- **Web 框架**: Gin
- **数据库**: MySQL + GORM
- **认证**: JWT
- **日志**: Zap
- **配置管理**: Viper

## 项目结构

```
KubeGale/
├── api/                    # 后端 API 层
│   └── v1/
│       └── kubernetes/     # K8s 相关 API
├── model/                  # 数据模型层
│   └── kubernetes/         # K8s 相关模型
├── service/                # 业务逻辑层
│   └── kubernetes/         # K8s 相关服务
├── router/                 # 路由层
│   └── kubernetes/         # K8s 相关路由
├── ui/                     # 前端代码
│   └── src/
│       ├── views/          # 页面组件
│       │   └── kubernetes/ # K8s 相关页面
│       │       └── workload/
│       │           └── deployment/  # Deployment 管理
│       ├── api/            # API 接口
│       │   └── kubernetes/ # K8s 相关接口
│       │       └── workload/        # 工作负载 API
│       └── types/          # TypeScript 类型定义
│           └── kubernetes/ # K8s 相关类型
└── docs/                   # 项目文档
```

## 已实现功能详情

### Deployment 管理功能 ✅

#### 1. 列表页面 (`ui/src/views/kubernetes/workload/deployment/index.vue`)
- **功能特性**:
  - 集群和命名空间筛选
  - 关键词搜索
  - 分页显示
  - 状态标签显示
  - 副本数状态显示
  - 镜像信息展示
  - 操作按钮（详情、编辑、扩缩容、删除等）

#### 2. 创建/编辑表单 (`ui/src/views/kubernetes/workload/deployment/components/deployment-form.vue`)
- **功能特性**:
  - 基本信息配置（名称、副本数）
  - 容器配置（镜像、端口、环境变量、资源限制）
  - 更新策略配置
  - 表单验证
  - 动态添加/删除容器、端口、环境变量

#### 3. 详情查看 (`ui/src/views/kubernetes/workload/deployment/components/deployment-detail.vue`)
- **功能特性**:
  - 基本信息展示
  - 状态信息展示
  - 容器详细信息
  - 更新策略信息
  - 标签和注解展示
  - 条件状态表格

#### 4. 扩缩容对话框 (`ui/src/views/kubernetes/workload/deployment/components/scale-dialog.vue`)
- **功能特性**:
  - 当前副本数显示
  - 目标副本数设置
  - 状态预览
  - 操作确认

#### 5. API 接口 (`ui/src/api/kubernetes/workload/`)
- **已实现接口**:
  - `deploymentApi.getDeploymentList()` - 获取列表
  - `deploymentApi.getDeploymentDetail()` - 获取详情
  - `deploymentApi.createDeployment()` - 创建
  - `deploymentApi.updateDeployment()` - 更新
  - `deploymentApi.deleteDeployment()` - 删除
  - `deploymentApi.scaleDeployment()` - 扩缩容
  - `deploymentApi.rollbackDeployment()` - 回滚

#### 6. 类型定义 (`ui/src/types/kubernetes/workload.ts`)
- **已定义类型**:
  - `Deployment` - Deployment 资源类型
  - `Pod` - Pod 资源类型
  - `Container` - 容器类型
  - `PodSpec` - Pod 规格类型
  - 各种 API 请求和响应类型

## 开发规范

### 代码规范
1. **命名规范**
   - 组件名使用 PascalCase
   - 文件名使用 kebab-case
   - 变量和函数使用 camelCase
   - 常量使用 UPPER_CASE

2. **文件组织**
   - 按功能模块组织文件
   - 相关文件放在同一目录下
   - 使用 index 文件进行统一导出

3. **组件设计**
   - 使用 Composition API
   - 保持组件的单一职责
   - 合理使用 props 和 emits
   - 实现适当的错误处理

### Git 工作流
1. **分支管理**
   - `main`: 主分支，用于生产环境
   - `develop`: 开发分支，用于集成测试
   - `feature/*`: 功能分支，用于新功能开发
   - `hotfix/*`: 热修复分支，用于紧急修复

2. **提交规范**
   - `feat`: 新功能
   - `fix`: 修复 bug
   - `docs`: 文档更新
   - `style`: 代码格式调整
   - `refactor`: 代码重构
   - `test`: 测试相关
   - `chore`: 构建过程或辅助工具的变动

## 部署说明

### 开发环境
1. **后端启动**
   ```bash
   cd KubeGale
   go mod tidy
   go run resource/cmd/main.go
   ```

2. **前端启动**
   ```bash
   cd ui
   npm install
   npm run dev
   ```

### 生产环境
1. **后端部署**
   - 使用 Docker 容器化部署
   - 配置环境变量
   - 设置数据库连接

2. **前端部署**
   - 构建生产版本
   - 部署到 Nginx 或 CDN
   - 配置反向代理

## 测试策略

### 单元测试
- 使用 Vitest 进行前端单元测试
- 使用 Go 标准测试包进行后端单元测试
- 测试覆盖率要求 > 80%

### 集成测试
- API 接口测试
- 前后端集成测试
- 数据库集成测试

### 端到端测试
- 使用 Playwright 进行 E2E 测试
- 覆盖主要用户流程
- 自动化测试部署

## 监控和日志

### 前端监控
- 错误监控和上报
- 性能监控
- 用户行为分析

### 后端监控
- 应用性能监控
- 数据库性能监控
- 系统资源监控

### 日志管理
- 结构化日志记录
- 日志级别管理
- 日志聚合和分析

## 安全考虑

### 前端安全
- XSS 防护
- CSRF 防护
- 敏感信息保护
- 输入验证

### 后端安全
- 身份认证和授权
- API 访问控制
- 数据加密
- SQL 注入防护

## 性能优化

### 前端优化
- 代码分割和懒加载
- 图片优化
- 缓存策略
- 打包优化

### 后端优化
- 数据库查询优化
- 缓存使用
- 并发处理
- 资源限制

## 贡献指南

1. **Fork 项目**
2. **创建功能分支**
3. **提交代码**
4. **创建 Pull Request**
5. **代码审查**
6. **合并代码**

## 联系方式

- **项目维护者**: [维护者信息]
- **技术支持**: [支持邮箱]
- **问题反馈**: [GitHub Issues]

## 许可证

本项目采用 [许可证类型] 许可证，详见 LICENSE 文件。

---

**最后更新**: 2024-01-XX
**版本**: v1.0.0 