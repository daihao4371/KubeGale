import service from '@/utils/request';
import type { ApiResponse } from '@/api/kubernetes/cluster/k8sCluster'; // Re-use generic ApiResponse
import type { K8sNode } from '@/types/kubernetes'; // Import centralized K8sNode type

// Base URL construction for direct Kubernetes API proxying if applicable,
// though current implementation uses backend-defined paths.
// const K8sBaseUrl = (clusterId: string | number): string => {
//     return `/kubernetes/proxy/${clusterId}/api/v1/nodes`;
// };

export interface NodeListParams {
    cluster_id: string | number; // Changed from clusterId to match existing backend param
    page?: number;
    pageSize?: number;
    keyword?: string; // For searching by node name or other fields
    // labelSelector?: string; // For filtering by labels
}

export interface NodeListResponse {
    items: K8sNode[]; // Assuming the backend returns a list of K8sNode objects
    total: number;
    // page and pageSize might also be part of the response if backend supports it explicitly
    page?: number;
    pageSize?: number;
}

export const getNodesList = (params: NodeListParams): Promise<ApiResponse<NodeListResponse>> => {
    return service({
        url: '/kubernetes/nodes', // Existing backend endpoint
        method: 'get',
        params,
    });
};

export interface NodeDetailParams {
    cluster_id: string | number; // Changed from clusterId
    nodeName: string;
}

// Assuming DescribeNodeInfo returns a single K8sNode object within a data wrapper
// The original response structure in index.vue was res.data.items for a single node, which is unusual.
// Let's assume the backend API /kubernetes/nodeDetails returns { code: 0, data: { node: K8sNode }, message: '' }
// Or if it's res.data.items, then items is K8sNode. For now, assuming a more standard single object response.
export interface NodeDetailResponse {
    // items: K8sNode; // As per current index.vue consumption (res.data.items)
    node: K8sNode; // A more typical structure
}


export const describeNodeInfo = (params: NodeDetailParams): Promise<ApiResponse<{items: K8sNode}>> => { // Adjusted to current index.vue consumption
    return service({
        url: '/kubernetes/nodeDetails', // Existing backend endpoint
        method: 'get',
        params,
    });
};

// For NodeUpdate, the payload is the K8sNode object itself or relevant parts.
// The existing API takes { cluster_id, nodeName, content: NodeContent }
// NodeContent was a simplified version of K8sNode. We should use K8sNode for 'content'.
export interface NodeUpdatePayload {
    cluster_id: string | number; // Changed from clusterId
    nodeName: string; // Name of the node to update
    content: Partial<K8sNode>; // The K8sNode object with fields to update (e.g., labels, annotations, taints)
                                // Using Partial<K8sNode> as typically only a subset of fields are updatable.
}

export const updateNode = (payload: NodeUpdatePayload): Promise<ApiResponse<{items: K8sNode} | null>> => { // items for consistency with current error checks in index.vue
    return service({
        url: '/kubernetes/nodes', // Existing backend endpoint
        method: 'put',
        data: payload,
    });
};

export interface NodeDeleteParams {
    cluster_id: string | number; // Changed from clusterId
    nodeName: string;
}

// The original NodesDelete uses a dynamic URL: `/kubernetes/${cluster_id}/api/v1/nodes/${name}`
// This implies it might be hitting a proxy directly.
// For consistency with other calls, if there's a backend wrapper, it might be different.
// Let's stick to the dynamic URL for now as per original.
export const deleteNode = (params: NodeDeleteParams): Promise<ApiResponse<{items: K8sNode} | null>> => { // items for consistency
    return service({
        url: `/kubernetes/${params.cluster_id}/api/v1/nodes/${params.nodeName}`, // Original backend endpoint pattern
        method: 'delete',
    });
};

// Functions related to node operations like Cordon/Uncordon (scheduling)
// These typically involve patching the node's spec.unschedulable field.

export interface NodeSchedulePayload {
    cluster_id: string | number;
    nodeName: string;
    unschedulable: boolean;
}

// This function would likely call the same 'updateNode' endpoint,
// but the 'content' would be specifically structured to update 'spec.unschedulable'.
// Example: content: { spec: { unschedulable: true/false } }
// The current index.vue directly modifies the row object and calls NodesUpdate.
// So, a dedicated setNodeScheduling function might not be strictly needed if NodesUpdate handles partial updates to spec.

// Pod Eviction is handled by an API in pods.ts (PodsEviction) - no change needed here.
// Node Metrics are handled by an API in metrics.ts (NodeMetricsList) - no change needed here.
