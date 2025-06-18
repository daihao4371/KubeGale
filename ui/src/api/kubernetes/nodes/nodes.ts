import service from '@/utils/request'

const BaseUrl = (cluster_id: string | number): string => {
    return `/kubernetes/${cluster_id}/api/v1/nodes`
}

interface NodeParams {
    cluster_id: string | number
    page?: number
    pageSize?: number
    keyword?: string
    nodeName?: string
}

interface NodeContent {
    apiVersion: string
    kind: string
    metadata: {
        name: string
        annotations?: Record<string, string>
        labels?: Record<string, string>
    }
    spec?: {
        unschedulable?: boolean
    }
    status?: {
        allocatable?: {
            cpu: string
            memory: string
        }
    }
}

interface NodeUpdateData {
    cluster_id: string | number
    nodeName: string
    content: NodeContent
}

// export const NodesList = (cluster_id, page, pageSize, keywords) => {
//     return service({
//         url: `${BaseUrl(cluster_id)}?search=true&keywords=${keywords}&page=${page}&pageSize=${pageSize}`,
//         method: 'get',
//     })
// }
export const GetNodesList = (params: NodeParams) => {
    return service({
        url: '/kubernetes/nodes',
        method: 'get',
        params
    })
}

export const DescribeNodeInfo = (params: NodeParams) => {
    return service({
        url: '/kubernetes/nodeDetails',
        method: 'get',
        params
    })
}
// export const NodesUpdate = (cluster_id, name, data) => {
//     return service({
//         url: `${BaseUrl(cluster_id)}/${name}`,
//         method: 'put',
//         data
//     })
// }
export const NodesUpdate = (data: NodeUpdateData) => {
    return service({
        url: '/kubernetes/nodes',
        method: 'put',
        data
    })
}
export const NodesDelete = (cluster_id: string | number, name: string) => {
    return service({
        url: `${BaseUrl(cluster_id)}/${name}`,
        method: 'delete',
    })
}
