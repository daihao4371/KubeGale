/**
 * 基础时间接口
 * 包含创建、更新和删除时间的通用时间字段
 */
export interface BaseTime {
  CreatedAt: string  // 记录创建时间
  UpdatedAt: string  // 记录最后更新时间
  DeletedAt: string | null  // 记录删除时间，未删除时为null
}

/**
 * 基础菜单接口
 * 定义系统菜单项的基本结构
 */
export interface BaseMenu {
  menuId: number  // 菜单唯一标识ID
  menuName: string  // 菜单名称
  path: string  // 菜单路由路径
  component: string  // 菜单对应的组件路径
  parentId: number  // 父菜单ID，用于构建菜单层级关系
}

/**
 * 基础角色接口
 * 定义系统角色的基本结构
 */
export interface BaseAuthority {
  authorityId: number  // 角色唯一标识ID
  authorityName: string  // 角色名称
  parentId: number  // 父角色ID，用于构建角色层级关系
  defaultRouter: string  // 角色默认路由路径
}

/**
 * 角色信息接口
 * 扩展基础角色接口，包含完整的角色信息
 */
export interface Authority extends BaseAuthority, BaseTime {
  dataAuthorityId: Authority[] | null  // 数据权限ID列表
  children: Authority[] | null  // 子角色列表，用于构建角色树
  menus: Record<string, unknown> | null  // 角色关联的菜单信息
}

/**
 * 基础用户信息接口
 * 定义用户基本信息结构
 */
export interface BaseUserInfo {
  userName: string  // 用户名
  nickName: string  // 用户昵称
  phone: string  // 用户手机号
  email: string  // 用户邮箱
  headerImg: string  // 用户头像图片路径
  enable: number  // 用户状态：1-启用，0-禁用
  authorityId: number  // 用户主角色ID
}

/**
 * 用户注册数据类型
 * 用于用户注册时的数据结构
 */
export interface RegisterUserData extends BaseUserInfo {
  passWord: string  // 用户密码
  authorityIds?: number[]  // 用户关联的角色ID列表（可选）
  ID?: number  // 用户ID（可选，注册时通常不需要）
}

/**
 * 用户数据类型
 * 系统中使用的用户信息结构
 */
export interface UserInfo extends BaseUserInfo {
  ID: number  // 用户唯一标识ID
  authority?: BaseTime  // 用户权限相关时间信息（可选）
  originSetting?: Record<string, unknown>  // 用户原始设置信息（可选）
  authorityIds?: number[]  // 用户关联的角色ID列表（可选）
}

/**
 * 用户详细信息类型
 * 包含完整的用户信息，用于用户详情展示
 */
export interface UserDetailInfo extends UserInfo, BaseTime {
  authorities: Array<Authority>  // 用户关联的角色列表
  uuid: string  // 用户唯一标识符
  msg: string  // 用户相关消息
}

/**
 * API接口定义
 * 描述系统中的API接口信息
 */
export interface Api {
  ID: number  // API唯一标识ID
  path: string  // API路径
  description: string  // API描述
  apiGroup: string  // API所属分组
  method: string  // API请求方法（GET、POST等）
  createdAt: string  // 创建时间
  updatedAt: string  // 更新时间
}

/**
 * API管理状态接口
 * 用于API管理页面的状态管理
 */
export interface ApiState {
  loading: boolean  // 加载状态
  apiList: Api[]  // API列表数据
  total: number  // 总记录数
  pageSize: number  // 每页显示条数
  currentPage: number  // 当前页码
  dialogVisible: boolean  // 对话框显示状态
  dialogTitle: string  // 对话框标题
  dialogType: 'create' | 'edit'  // 对话框类型：创建或编辑
  formLoading: boolean  // 表单加载状态
  apiGroups: string[]  // API分组列表
  searchForm: {  // 搜索表单
    path: string  // API路径搜索条件
    description: string  // API描述搜索条件
    apiGroup: string  // API分组搜索条件
    method: string  // API方法搜索条件
  }
  filteredApiList: Api[]  // 过滤后的API列表
}

/**
 * API表单接口
 * 用于API创建和编辑的表单数据结构
 */
export interface ApiForm {
  ID?: number  // API ID（编辑时需要）
  path: string  // API路径
  description: string  // API描述
  apiGroup: string  // API所属分组
  method: string  // API请求方法
}

/**
 * 默认API表单数据
 * 创建新API时的默认值
 */
export const defaultApiForm: ApiForm = {
  path: '',  // 默认空路径
  description: '',  // 默认空描述
  apiGroup: '',  // 默认空分组
  method: 'GET'  // 默认GET方法
}

/**
 * 角色创建表单接口
 * 用于创建新角色的表单数据结构
 */
export interface CreateRoleForm {
  authorityId: number  // 角色ID
  authorityName: string  // 角色名称
  parentId: number  // 父角色ID
  defaultRouter: string  // 默认路由
}

/**
 * 角色复制表单接口
 * 用于复制现有角色的表单数据结构
 */
export interface CopyRoleForm {
  authorityId: number  // 新角色ID
  authorityName: string  // 新角色名称
  parentId: number  // 父角色ID
  oldAuthorityId: number  // 被复制的角色ID
  oldAuthorityName: string  // 被复制的角色名称
}

/**
 * 复制角色请求接口
 * 用于发送复制角色的API请求
 */
export interface CopyAuthorityRequest {
  authority: {  // 新角色信息
    authorityId: number  // 角色ID
    authorityName: string  // 角色名称
    parentId: number  // 父角色ID
    defaultRouter: string  // 默认路由
  }
  oldAuthorityId: number  // 被复制的角色ID
}

/**
 * 删除角色参数接口
 * 用于删除角色的API请求参数
 */
export interface DeleteAuthorityParams {
  AuthorityId: number  // 要删除的角色ID
}

/**
 * 角色管理状态接口
 * 用于角色管理页面的状态管理
 */
export interface RoleState {
  roleList: Authority[]  // 角色列表
  loading: boolean  // 加载状态
  createRoleDialogVisible: boolean  // 创建角色对话框显示状态
  createRoleLoading: boolean  // 创建角色加载状态
  createRoleForm: CreateRoleForm  // 创建角色表单数据
  editRoleDialogVisible: boolean  // 编辑角色对话框显示状态
  editRoleLoading: boolean  // 编辑角色加载状态
  editRoleForm: CreateRoleForm  // 编辑角色表单数据
  copyRoleDialogVisible: boolean  // 复制角色对话框显示状态
  copyRoleLoading: boolean  // 复制角色加载状态
  copyRoleForm: CopyRoleForm  // 复制角色表单数据
  permissionDialogVisible: boolean  // 权限设置对话框显示状态
  currentRole: Authority | null  // 当前选中的角色
}

/**
 * 角色管理表单验证规则
 * 用于验证角色表单输入的合法性
 */
export const createRoleRules = {
  authorityId: [
    { required: true, message: '请输入角色ID', trigger: 'blur' },  // 角色ID必填验证
    { type: 'number', message: '角色ID必须为数字', trigger: 'blur' }  // 角色ID类型验证
  ],
  authorityName: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },  // 角色名称必填验证
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }  // 角色名称长度验证
  ]
}

/**
 * 默认角色创建表单数据
 * 创建新角色时的默认值
 */
export const defaultCreateRoleForm: CreateRoleForm = {
  authorityId: 0,  // 默认角色ID为0
  authorityName: '',  // 默认空角色名称
  parentId: 0,  // 默认父角色ID为0
  defaultRouter: 'dashboard'  // 默认路由为dashboard
}

/**
 * 默认角色复制表单数据
 * 复制角色时的默认值
 */
export const defaultCopyRoleForm: CopyRoleForm = {
  authorityId: 0,  // 默认新角色ID为0
  authorityName: '',  // 默认空角色名称
  parentId: 0,  // 默认父角色ID为0
  oldAuthorityId: 0,  // 默认被复制的角色ID为0
  oldAuthorityName: ''  // 默认空被复制的角色名称
}

/**
 * 用户接口
 * 定义系统中用户的基本信息结构
 */
export interface User {
  id: number  // 用户ID
  username: string  // 用户名
  nickname?: string  // 用户昵称（可选）
  realName?: string  // 用户真实姓名（可选）
}

/**
 * 系统操作记录接口
 * 记录用户在系统中的操作信息
 */
export interface SysOperationRecord {
  id: number  // 记录ID
  created_at: string  // 创建时间（下划线命名）
  updated_at: string  // 更新时间（下划线命名）
  ip: string  // 操作IP地址
  method: string  // 请求方法
  path: string  // 请求路径
  status: number  // 响应状态码
  latency: string | number  // 请求耗时
  agent: string  // 用户代理信息
  error_message: string  // 错误信息
  body: string  // 请求体
  resp: string  // 响应内容
  user_id: number  // 用户ID
  user: User  // 用户信息
  operator_name?: string  // 操作者名称（可选）
  operator_real_name?: string  // 操作者真实姓名（可选）
  ID: number  // 记录ID（大写命名）
  CreatedAt?: string  // 创建时间（可选，大写命名）
  UpdatedAt?: string  // 更新时间（可选，大写命名）
}

/**
 * 操作记录查询参数接口
 * 用于查询操作记录的参数结构
 */
export interface OperationRecordParams {
  page?: number  // 页码（可选）
  pageSize?: number  // 每页条数（可选）
  path?: string  // 请求路径过滤（可选）
  method?: string  // 请求方法过滤（可选）
  status?: number  // 状态码过滤（可选）
  ip?: string  // IP地址过滤（可选）
  userId?: number  // 用户ID过滤（可选）
  startTime?: string  // 开始时间过滤（可选）
  endTime?: string  // 结束时间过滤（可选）
}

/**
 * 操作记录状态接口
 * 用于操作记录页面的状态管理
 */
export interface OperationRecordState {
  loading: boolean  // 加载状态
  operationList: SysOperationRecord[]  // 操作记录列表
  total: number  // 总记录数
  pageSize: number  // 每页显示条数
  currentPage: number  // 当前页码
  searchForm: OperationRecordParams  // 搜索表单参数
}
