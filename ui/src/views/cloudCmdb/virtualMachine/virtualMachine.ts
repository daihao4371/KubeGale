import type { VirtualMachine } from '@/api/cloudCmdb/cloud_virtual_machine'
import type { PlatformRegionTreeItem, RegionNode } from '@/api/cloudCmdb/cloud_rds'
import { ref } from 'vue'

export interface VirtualMachineState {
  tableData: VirtualMachine[]
  loading: boolean
  syncing: boolean
  page: number
  pageSize: number
  total: number
  selectedPlatformId: number | null
}

export interface VirtualMachineMethods {
  fetchData: () => Promise<void>
  handleSync: () => Promise<void>
  handlePlatformSelect: (payload: { platform: PlatformRegionTreeItem, region: RegionNode | null }) => void
  handleSizeChange: (val: number) => void
  handlePageChange: (val: number) => void
}

// 新增：搜索表单
export const searchForm = ref({
  name: '',
  instanceId: ''
})

// 新增：处理搜索
export const handleSearch = async (fetchData: () => Promise<void>) => {
  // 重置页码为1
  // 这里fetchData会自动带上searchForm参数
  await fetchData()
}

// 新增：重置搜索
export const handleResetSearch = async (fetchData: () => Promise<void>) => {
  searchForm.value.name = ''
  searchForm.value.instanceId = ''
  await fetchData()
}
