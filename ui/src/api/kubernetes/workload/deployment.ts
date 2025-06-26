import service from '@/utils/request';
import type { ApiResponse } from '@/api/kubernetes/cluster/k8sCluster'; // Re-use generic ApiResponse
import type { K8sDeployment, K8sReplicaSet } from '@/types/kubernetes'; // Import centralized types

export interface DeploymentListParams {
  clusterId: string;
  namespace: string;
  page?: number;
  pageSize?: number;
  name?: string; // For filtering by deployment name
  labelSelector?: string; // For filtering by labels (e.g., "app=nginx,tier=frontend")
}

export interface DeploymentListResponse {
  items: K8sDeployment[];
  total: number;
}

// Get list of Deployments
export const getDeployments = (params: DeploymentListParams): Promise<ApiResponse<DeploymentListResponse>> => {
  return service({
    url: `/kubernetes/proxy/${params.clusterId}/apis/apps/v1/namespaces/${params.namespace}/deployments`,
    method: 'get',
    params: {
      // Kubernetes API uses different pagination/filtering params
      limit: params.pageSize,
      continue: params.page && params.page > 1 ? `page-${params.page}` : undefined, // Example of how 'continue' might be used. Real usage is more complex.
      fieldSelector: params.name ? `metadata.name=${params.name}` : undefined,
      labelSelector: params.labelSelector,
    },
  });
};

// Get a specific Deployment by name
export const getDeploymentByName = (clusterId: string, namespace: string, name: string): Promise<ApiResponse<K8sDeployment>> => {
  return service({
    url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/deployments/${name}`,
    method: 'get',
  });
};

// Create a new Deployment
// The payload should be a valid K8sDeployment object
export const createDeployment = (clusterId: string, namespace: string, payload: K8sDeployment): Promise<ApiResponse<K8sDeployment>> => {
  return service({
    url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/deployments`,
    method: 'post',
    data: payload,
  });
};

// Update an existing Deployment
// The payload should be a valid K8sDeployment object
export const updateDeployment = (clusterId: string, namespace: string, name: string, payload: K8sDeployment): Promise<ApiResponse<K8sDeployment>> => {
  return service({
    url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/deployments/${name}`,
    method: 'put',
    data: payload,
  });
};

// Delete a Deployment
export const deleteDeployment = (clusterId: string, namespace: string, name: string): Promise<ApiResponse<any>> => { // Status object might be returned
  return service({
    url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/deployments/${name}`,
    method: 'delete',
  });
};

// Scale a Deployment (update replicas)
export interface ScalePayload {
  spec: {
    replicas: number;
  };
}
export const scaleDeployment = (clusterId: string, namespace: string, name: string, replicas: number): Promise<ApiResponse<K8sDeployment>> => { // Or a Scale object
  // Kubernetes API uses PATCH with strategic merge or JSON patch for scaling typically.
  // Or it can be a PUT on the /scale subresource.
  // Using PATCH for this example.
  const payload: ScalePayload = { spec: { replicas } };
  return service({
    url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/deployments/${name}/scale`, // Or just .../deployments/${name} with PATCH
    method: 'put', // Or 'patch' with 'Content-Type: application/merge-patch+json' or 'application/json-patch+json'
    data: payload, // For PATCH, the payload structure would differ based on patch type. For PUT on /scale, it's a Scale object.
                  // Let's assume the backend proxy handles this abstraction and accepts a simple replica count or a Scale object.
                  // For direct K8s API, a Scale object is: { apiVersion: 'autoscaling/v1', kind: 'Scale', metadata: { name, namespace }, spec: { replicas }}
  });
};


// Rollback a Deployment (to a specific revision or undo)
// This is a more complex operation involving ReplicaSets.
// Placeholder for API structure.
export interface RollbackConfig {
  name: string; // Deployment name
  rollbackTo: {
    revision: number; // 0 means to undo the last rollout
  };
}
export interface CreateDeploymentRollbackPayload {
    apiVersion: "apps/v1"; // or "extensions/v1beta1" for older clusters
    kind: "DeploymentRollback";
    name: string; // Name of the Deployment
    updatedAnnotations?: Record<string, string>;
    rollbackTo: {
      revision: number; // Specify the revision to roll back to. 0 means to roll back to the previous revision.
    };
}
export const rollbackDeployment = (clusterId: string, namespace: string, deploymentName: string, revision: number = 0): Promise<ApiResponse<any>> => {
  const payload: CreateDeploymentRollbackPayload = {
    apiVersion: "apps/v1",
    kind: "DeploymentRollback",
    name: deploymentName,
    rollbackTo: {
      revision: revision,
    },
  };
  return service({
    // The DeploymentRollback kind was removed in k8s 1.16. Rollbacks are now done by `kubectl rollout undo deployment/...`
    // which translates to updating the deployment with a specific replicaset template or using `PATCH` to set `spec.template`.
    // Or managing ReplicaSets directly.
    // A common approach is to use `PATCH` to set `spec.paused=true`, then update `spec.template` from a previous ReplicaSet, then `spec.paused=false`.
    // For simplicity, let's assume a backend endpoint that handles this logic.
    url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/deployments/${deploymentName}/rollback`, // This is a custom endpoint
    method: 'post',
    data: { revision }, // Backend would interpret this
  });
};


// Get Deployment's revision history (ReplicaSets)
// Placeholder - this would list ReplicaSets with a specific label selector matching the deployment.
export const getDeploymentRevisionHistory = (clusterId: string, namespace: string, deploymentName: string): Promise<ApiResponse<any>> => {
    return service({
        url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/replicasets`,
        method: 'get',
        params: {
            // This requires knowing the label selector used by the deployment to find its ReplicaSets.
            // Usually, it's something like `app=<deployment-name>` and potentially a unique `pod-template-hash`.
            // This needs to be fetched from the deployment object first.
            labelSelector: `app=${deploymentName}`, // This is a simplified example
        }
    });
};

// Update Deployment with raw YAML/JSON
export const updateDeploymentWithYaml = (clusterId: string, namespace: string, name: string, payload: K8sDeployment): Promise<ApiResponse<K8sDeployment>> => {
  return service({
    url: `/kubernetes/proxy/${clusterId}/apis/apps/v1/namespaces/${namespace}/deployments/${name}`,
    method: 'put',
    data: payload,
    headers: {
        'Content-Type': 'application/yaml' // Or let backend auto-detect from payload if it's an object
    }
  });
};

// Get Deployment Events
// Placeholder - this would list Events related to the deployment.
export const getDeploymentEvents = (clusterId: string, namespace: string, deploymentUID: string): Promise<ApiResponse<any>> => {
    return service({
        url: `/kubernetes/proxy/${clusterId}/api/v1/namespaces/${namespace}/events`,
        method: 'get',
        params: {
            fieldSelector: `involvedObject.uid=${deploymentUID},involvedObject.kind=Deployment`
        }
    });
};

// Get Pods for a Deployment
// Placeholder - this would list Pods with a specific label selector matching the deployment's selector.
export const getPodsForDeployment = (clusterId: string, namespace: string, labelSelector: string): Promise<ApiResponse<any>> => {
    return service({
        url: `/kubernetes/proxy/${clusterId}/api/v1/namespaces/${namespace}/pods`,
        method: 'get',
        params: {
            labelSelector: labelSelector,
        }
    });
};

// It's often better to use the official Kubernetes client types if possible, e.g.,
// import { V1Deployment, V1DeploymentList } from '@kubernetes/client-node';
// Then K8sDeployment would be V1Deployment, DeploymentListResponse items would be V1Deployment[].
// This ensures type safety and compatibility with the Kubernetes API.
// The current K8sDeployment interface is a simplified version.
