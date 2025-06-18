export interface Region {
  id: number
  name: string
}

export interface CloudProvider {
  id: number
  name: string
  platform: string
  access_key_id: string
  access_key_secret: string
  created_at?: string
  updated_at?: string
}

export interface CloudProviderState {
  loading: boolean
  providerList: CloudProvider[]
  total: number
  pageSize: number
  currentPage: number
  dialogVisible: boolean
  dialogTitle: string
  dialogType: 'create' | 'edit'
  formLoading: boolean
  searchForm: {
    name: string
    type: string
    status: string
  }
  filteredProviderList: CloudProvider[]
}

export interface CloudProviderForm {
  id?: number
  name: string
  platform: string
  access_key_id: string
  access_key_secret: string
}

export const defaultCloudProviderForm: CloudProviderForm = {
  name: '',
  platform: '',
  access_key_id: '',
  access_key_secret: ''
}

export interface ApiResponse<T> {
  code: number
  msg: string
  data: T
}

export interface CloudProviderListResponse {
  list: CloudProvider[]
  total: number
  page: number
  pageSize: number
}

export interface CloudProviderDetailResponse {
  cloud_platform: CloudProvider
  regions: Region[]
}

export interface CloudProviderCreateResponse {
  code: number
  data: Record<string, never>
  msg: string
}

export interface RdsInstance {
  id: number
  name: string
  instance_id: string
  private_addr: string
  public_addr: string
  region: string
  region_name: string
  status: string
  creation_time: string
  expired_time: string
  cloud_platform_id: number
  cloud_platform: CloudProvider
  created_at: string
  updated_at: string
}

export interface RdsState {
  loading: boolean
  rdsList: RdsInstance[]
  total: number
  pageSize: number
  currentPage: number
  dialogVisible: boolean
  dialogTitle: string
  dialogType: 'create' | 'edit'
  formLoading: boolean
  searchForm: {
    name: string
    instance_id: string
    region: string
    status: string
  }
  filteredRdsList: RdsInstance[]
}

export interface RdsForm {
  id?: number
  name: string
  engine: string
  engine_version: string
  instance_type: string
  instance_class: string
  storage_type: string
  storage_size: number
  region: string
  vpc_id: string
  subnet_id: string
  security_group_ids: string[]
  master_username: string
  master_password: string
  backup_retention_period: number
  maintenance_window: string
  cloud_platform_id: number
}

export const defaultRdsForm: RdsForm = {
  name: '',
  engine: '',
  engine_version: '',
  instance_type: '',
  instance_class: '',
  storage_type: '',
  storage_size: 20,
  region: '',
  vpc_id: '',
  subnet_id: '',
  security_group_ids: [],
  master_username: '',
  master_password: '',
  backup_retention_period: 7,
  maintenance_window: '',
  cloud_platform_id: 0
}

export interface RdsListResponse {
  list: RdsInstance[]
  total: number
  page: number
  pageSize: number
} 