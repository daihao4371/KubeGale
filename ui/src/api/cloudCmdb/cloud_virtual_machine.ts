import service from '@/utils/request'

interface ListRequest {
  page: number
  pageSize: number
  name?: string
  instance_id?: string
}

export interface CloudPlatform {
  id: number
  name: string
  access_key_id: string
  access_key_secret: string
  platform: string
  created_at: string
  updated_at: string
}

export interface VirtualMachine {
  id: number
  name: string
  instance_id: string
  cpu: number
  memory: number
  os: string
  os_type: string
  private_addr: string
  public_addr: string
  region: string
  region_name: string
  status: string
  creation_time: string
  expired_time: string
  cloud_platform_id: number
  cloud_platform: CloudPlatform
  created_at: string
  updated_at: string
}

interface ListResponse {
  code: number
  data: {
    list: VirtualMachine[]
    total: number
    page: number
    pageSize: number
  }
  msg: string
}

interface SyncResponse {
  code: number
  msg: string
}

export const virtualMachineList = (data: ListRequest) => {
  return service<ListResponse>({
    url: '/virtualMachine/list',
    method: 'post',
    data
  })
}

export const syncMachine = (data: { id: number }) => {
  return service<SyncResponse>({
    url: '/virtualMachine/sync',
    method: 'post',
    data
  })
}

export const virtualMachineTree = (data: { platform_id: number }) => {
  return service({
    url: '/virtualMachine/tree',
    method: 'post',
    data
  })
}

