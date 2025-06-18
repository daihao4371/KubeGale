import request from '@/utils/request'
import type { ApiResponse } from '@/types/cloudCmdb'

// 同步Region
export const syncRegion = (data: { id: number }): Promise<ApiResponse<null>> => {
  return request<ApiResponse<null>>({
    url: '/cloud_region/syncRegion',
    method: 'post',
    data
  })
}
