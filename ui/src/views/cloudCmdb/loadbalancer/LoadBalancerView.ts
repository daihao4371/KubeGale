import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { syncLoadBalancer, getLoadBalancerList } from '@/api/cloudCmdb/cloud_load_balancer'
import type { LoadBalancerItem } from '@/api/cloudCmdb/cloud_load_balancer'
import type { PlatformRegionTreeItem, RegionNode } from '@/api/cloudCmdb/cloud_rds'

// 类型定义
export interface SearchInfo {
  keyword: string
}

export default function useLoadBalancer() {
  // 状态定义
  const tableData = ref<LoadBalancerItem[]>([])
  const loading = ref(false)
  const currentProviderId = ref(0)
  const searchInfo = ref<SearchInfo>({ keyword: '' })
  const pageSize = ref(10)
  const currentPage = ref(1)
  const total = ref(0)
  const detailDialogVisible = ref(false)
  const currentDetail = ref<LoadBalancerItem | null>(null)

  // 获取列表数据
  const fetchList = async () => {
    loading.value = true
    try {
      const params = {
        loadBalancer: {
          name: searchInfo.value.keyword || '',
          instanceId: searchInfo.value.keyword || '',
          region: ''
        },
        page: currentPage.value,
        pageSize: pageSize.value,
        keyword: searchInfo.value.keyword || '',
        field: '',
        orderKey: 'id' as const,
        desc: false
      }
      const res = await getLoadBalancerList(params)
      if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
      } else {
        ElMessage.error(res.msg || '获取列表失败')
      }
    } catch (error) {
      ElMessage.error('获取列表失败')
    } finally {
      loading.value = false
    }
  }

  // 初始化
  const init = () => {
    fetchList()
  }

  // 处理平台选择
  const handlePlatformSelect = (payload: { platform: PlatformRegionTreeItem, region: RegionNode | null }) => {
    currentProviderId.value = payload.platform.id
    currentPage.value = 1
    fetchList()
  }

  // 搜索
  const onSearch = () => {
    currentPage.value = 1
    fetchList()
  }

  const onReset = () => {
    searchInfo.value.keyword = ''
    currentPage.value = 1
    fetchList()
  }

  // 分页
  const handleCurrentChange = (page: number) => {
    currentPage.value = page
    fetchList()
  }
  const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchList()
  }

  // 同步
  const onSync = async () => {
    if (!currentProviderId.value) {
      ElMessage.warning('请先选择云厂商')
      return
    }
    loading.value = true
    try {
      const res = await syncLoadBalancer(currentProviderId.value)
      if (res.code === 0) {
        ElMessage.success(res.msg || '同步操作成功')
        fetchList()
      } else {
        ElMessage.error(res.msg || '同步操作失败')
      }
    } catch (error) {
      ElMessage.error('同步操作失败')
    } finally {
      loading.value = false
    }
  }

  // 详情
  const handleDetail = (row: LoadBalancerItem) => {
    currentDetail.value = row
    detailDialogVisible.value = true
  }

  return {
    tableData,
    loading,
    currentProviderId,
    searchInfo,
    pageSize,
    currentPage,
    total,
    detailDialogVisible,
    currentDetail,
    init,
    handlePlatformSelect,
    onSearch,
    onReset,
    onSync,
    handleCurrentChange,
    handleSizeChange,
    handleDetail
  }
}
