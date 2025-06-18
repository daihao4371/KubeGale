# 项目结构索引（细化版）

## 1. 入口与主文件
- `src/main.ts`：Vue 应用入口，初始化根组件和全局配置
- `src/App.vue`：根组件，应用布局和全局容器
- `index.html`：HTML 模板，挂载点和基础 meta 配置

## 2. 路由
- `src/router/index.ts`：全局路由配置，定义页面路由表和路由守卫

## 3. 视图页面（Views）
- `src/views/`：所有业务页面
  - `im/`：IM 通知与消息相关模块
    - `notificationManage.vue`：IM 通知管理页面 UI
    - `notificationManage.ts`：IM 通知管理页面的业务逻辑（组合式 API）
    - `notificationManage.css`：IM 通知管理页面样式
    - `index.vue`：IM 模块入口页面
  - `cmdb/`：CMDB 资产管理模块（主机、项目、批量任务等）
  - `system/`：系统管理模块（用户、角色、权限、操作日志等）
  - `homepage/`：首页及仪表盘模块（云资产、CMDB、快捷入口等）
  - `prometheus/`：Prometheus 监控管理模块（监控大盘、告警等）
  - `kubernetes/`：K8s 容器管理模块（集群、命名空间、工作负载等）
  - `docker/`：Docker 容器管理模块（镜像、容器、仓库等）
  - `config/`：配置中心模块（配置项、环境变量等）
  - `cicd/`：CI/CD 持续集成与部署模块（流水线、构建、发布等）
  - `terminal/`：Web终端/命令行模块
  - `login/`：登录与认证模块
  - `dashboard/`：通用仪表盘模块
  - `cloudCmdb/`：云 CMDB 资产管理模块

## 4. 组件（Components）
- `src/components/`：通用组件
  - `common/Pagination.vue`：分页组件，统一分页交互
  - `layout/AppLayout.vue`：主布局组件，包含侧边栏、头部等
  - `layout/Breadcrumb.vue`：面包屑导航组件

## 5. API 接口
- `src/api/`：所有后端接口请求
  - `im/notification.ts`：IM 通知相关接口（增删改查、测试、卡片内容等）
  - `system/user.ts`：用户相关接口（用户列表、增删改查等）
  - `request.ts`：全局 axios 封装，统一请求拦截、响应处理
  - 其他如 `cmdb/`、`cloudCmdb/`、`login/`、`homepage/`、`prometheus/`、`cicd/`、`config/`、`dashboard/`、`docker/`、`kubernetes/` 等：各业务模块接口

## 6. 类型定义（Types）
- `src/types/`：全局和业务类型定义
  - `im.ts`：IM 通知相关类型（通知、卡片内容、API 参数等）
  - `system.ts`：系统用户相关类型（用户、角色、权限等）
  - `cmdb.ts`、`cloudCmdb.ts`：CMDB 相关类型

## 7. 工具函数
- `src/utils/`：通用工具函数
  - `format.ts`：格式化相关（如日期、字符串等）
  - `request.ts`：请求相关辅助函数

## 8. 静态资源
- `src/assets/`：静态资源与全局样式
  - `main.css`：全局主样式
  - `base.css`：基础样式
  - `utilities.css`：工具类样式
  - `images/`：图片资源

## 9. 配置文件
- `vite.config.ts`：Vite 构建工具配置
- `package.json`：依赖与脚本配置
- `tsconfig.json`、`tsconfig.app.json`、`tsconfig.node.json`、`tsconfig.vitest.json`：TypeScript 配置
- `.editorconfig`、`.prettierrc.json`、`.gitignore` 等：开发规范与工具配置

---

> 本索引文件细化到主要文件，并为各模块和关键文件添加了功能注释，便于开发者快速理解和定位项目结构。 