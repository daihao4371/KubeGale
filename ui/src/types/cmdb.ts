export interface Host {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  name: string
  port: number
  project: number
  username: string
  serverHost: string
  privateKey: string
  password: string
  labels: string
  note: string
  cpuModel: string
  cpuCount: string
  memory: string
  diskTotal: string
  diskInfo: string
  os: string
  osVersion: string
  osArch: string
  status: string
  publicIP: string
  privateIP: string
  CreatedBy: number
  UpdatedBy: number
  DeletedBy: number
}

export interface ApiResponse<T = Record<string, unknown>> {
  code: number
  data: T
  msg: string
}

export interface AddUserReq {
  username: string
  password: string
  nickname?: string
  email?: string
  phone?: string
  authorityId?: number
}

export interface UpdateUserReq {
  ID: number
  nickname?: string
  email?: string
  phone?: string
  authorityId?: number
}

export interface ChangePasswordReq {
  oldPassword: string
  newPassword: string
} 