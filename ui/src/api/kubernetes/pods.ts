import service from '@/utils/request'

interface PodParams {
    cluster_id: string | number
    namespace?: string
    podName?: string
    fieldSelector?: string
    selector?: string
}

export const GetPodsList = (params: PodParams) => {
    return service({
        url: '/kubernetes/pods',
        method: 'get',
        params
    })
}

export const PodsDelete = (params: PodParams) => {
    return service({
        url: '/kubernetes/pods',
        method: 'delete',
        params
    })
}

export const PodsEviction = (params: PodParams) => {
    return service({
        url: '/kubernetes/pods/eviction',
        method: 'post',
        data: params
    })
} 