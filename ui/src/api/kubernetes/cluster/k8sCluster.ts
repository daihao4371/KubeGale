import service from '@/utils/request'
import type { K8sCluster, K8sClusterUser, K8sClusterRoleDefinition, K8sApiGroup, K8sClusterCredential } from '@/types/kubernetes';

// Generic API Response structure - assuming this is standard across the app
export interface ApiResponse<T> {
  code: number;
  data: T;
  message?: string;
}

// === Cluster Management ===

export interface ClusterListParams {
  page?: number;
  pageSize?: number;
  name?: string;
  // TODO: Add other filter params like startCreatedAt, endCreatedAt if supported by backend
  [key: string]: unknown; // Allow other potential query params
}

export interface ClusterListResponse {
  list: K8sCluster[];
  total: number;
  page: number;
  pageSize: number;
}

export interface ClusterDetailResponse {
  cluster: K8sCluster;
}

export const getClustersList = (params: ClusterListParams): Promise<ApiResponse<ClusterListResponse>> => {
  return service({
    url: '/kubernetes/clusterList', // Backend endpoint for listing clusters
    method: 'get',
    params
  })
}

export const getClusterById = (clusterId: string): Promise<ApiResponse<ClusterDetailResponse>> => {
  return service({
    url: '/kubernetes/clusterById', // Backend endpoint for fetching a cluster by ID
    method: 'post', // Assuming POST as per original, though GET with /kubernetes/cluster/{id} is common
    data: { id: clusterId }
  })
}

// For CreateCluster, the payload is K8sCluster but without 'id' and potentially 'createdAt', 'updatedAt'
export type CreateClusterPayload = Omit<K8sCluster, 'id' | 'createdAt' | 'updatedAt' | 'version' | 'status' | 'provider' >;

export const createCluster = (payload: CreateClusterPayload): Promise<ApiResponse<K8sCluster>> => { // Assuming backend returns the created cluster object
  return service({
    url: '/kubernetes/cluster', // Backend endpoint for creating a cluster
    method: 'post',
    data: payload
  })
}

// For UpdateCluster, 'id' is required.
export type UpdateClusterPayload = Partial<Omit<K8sCluster, 'id' | 'createdAt' | 'updatedAt'>> & Pick<K8sCluster, 'id'>;


export const updateCluster = (payload: UpdateClusterPayload): Promise<ApiResponse<K8sCluster>> => { // Assuming backend returns the updated cluster
  return service({
    url: '/kubernetes/cluster', // Backend endpoint for updating a cluster
    method: 'put',
    data: payload
  })
}

export const deleteCluster = (clusterId: string): Promise<ApiResponse<null>> => {
  return service({
    url: '/kubernetes/cluster', // Backend endpoint for deleting a cluster
    method: 'delete',
    data: { id: clusterId } // Assuming ID is passed in the body for DELETE
  })
}

export const deleteClusterByIds = (ids: string[]): Promise<ApiResponse<null>> => {
  return service({
    url: '/kubernetes/clusterByIds', // Backend endpoint for batch deleting clusters
    method: 'delete',
    data: { IDs: ids }
  })
}

// === Cluster User Management ===
// These are application-level users associated with a cluster, not necessarily direct k8s users/serviceaccounts

export interface ClusterUserListParams {
  clusterId: string;
  page?: number;
  pageSize?: number;
  username?: string;
}

export interface ClusterUserListResponse {
  list: K8sClusterUser[];
  total: number;
}

export const getClusterUsers = (params: ClusterUserListParams): Promise<ApiResponse<ClusterUserListResponse>> => {
  return service({
    // TODO: Define actual backend endpoint for listing users of a cluster
    url: `/kubernetes/cluster/${params.clusterId}/users`,
    method: 'get',
    params: { page: params.page, pageSize: params.pageSize, username: params.username }
  })
}

export type CreateClusterUserPayload = Omit<K8sClusterUser, 'id' | 'joinedAt'>;

export const createClusterUser = (payload: CreateClusterUserPayload): Promise<ApiResponse<K8sClusterUser>> => {
  return service({
    // TODO: Define actual backend endpoint
    url: `/kubernetes/cluster/${payload.clusterId}/users`,
    method: 'post',
    data: payload
  })
}

export type UpdateClusterUserPayload = Partial<Omit<K8sClusterUser, 'clusterId' | 'userId' | 'username' | 'joinedAt'>> & Pick<K8sClusterUser, 'id' | 'clusterId'>;

export const updateClusterUser = (payload: UpdateClusterUserPayload): Promise<ApiResponse<K8sClusterUser>> => {
  return service({
    // TODO: Define actual backend endpoint
    url: `/kubernetes/cluster/${payload.clusterId}/users/${payload.id}`,
    method: 'put',
    data: payload
  })
}

export const deleteClusterUser = (clusterId: string, userId: string): Promise<ApiResponse<null>> => {
  return service({
    // TODO: Define actual backend endpoint
    url: `/kubernetes/cluster/${clusterId}/users/${userId}`,
    method: 'delete',
  })
}


// === Cluster Role Management ===
// These are roles defined within the scope of a cluster, potentially mapping to K8s Roles/ClusterRoles

export interface ClusterRoleListParams {
  clusterId: string;
  page?: number;
  pageSize?: number;
  name?: string;
}

export interface ClusterRoleListResponse {
  list: K8sClusterRoleDefinition[];
  total: number;
}

export const getClusterRoles = (params: ClusterRoleListParams): Promise<ApiResponse<ClusterRoleListResponse>> => {
  return service({
    // Original: /kubernetes/getClusterRoles (POST with IdParams)
    // TODO: Confirm if this is to list roles FOR a cluster or GET a specific role by ID. Assuming list for now.
    // url: '/kubernetes/getClusterRoles',
    // method: 'post',
    // data: { id: params.clusterId } // This seems more like getRoleById if id is roleId
    url: `/kubernetes/cluster/${params.clusterId}/roles`, // RESTful approach
    method: 'get',
    params: { page: params.page, pageSize: params.pageSize, name: params.name }
  })
}

export type CreateClusterRolePayload = Omit<K8sClusterRoleDefinition, 'id' | 'createdAt' | 'updatedAt'>;

export const createClusterRole = (payload: CreateClusterRolePayload): Promise<ApiResponse<K8sClusterRoleDefinition>> => {
  return service({
    // Original: /kubernetes/createClusterRole (POST with ClusterData) -> Needs dedicated payload
    url: `/kubernetes/cluster/${payload.clusterId}/roles`,
    method: 'post',
    data: payload
  })
}

export type UpdateClusterRolePayload = Partial<Omit<K8sClusterRoleDefinition, 'clusterId' | 'createdAt' | 'updatedAt'>> & Pick<K8sClusterRoleDefinition, 'id' | 'clusterId'>;

export const updateClusterRole = (payload: UpdateClusterRolePayload): Promise<ApiResponse<K8sClusterRoleDefinition>> => {
  return service({
    // Original: /kubernetes/updateClusterRole (PUT with ClusterData) -> Needs dedicated payload
    url: `/kubernetes/cluster/${payload.clusterId}/roles/${payload.id}`,
    method: 'put',
    data: payload
  })
}

export const deleteClusterRole = (clusterId: string, roleId: string): Promise<ApiResponse<null>> => {
  return service({
    // Original: /kubernetes/deleteClusterRole (DELETE with IdParams for roleId)
    url: `/kubernetes/cluster/${clusterId}/roles/${roleId}`,
    method: 'delete'
  })
}

// === Cluster API Group Management ===

export interface ClusterApiGroupListParams {
  clusterId: string;
}

export interface ClusterApiGroupListResponse {
  apiGroups: K8sApiGroup[];
}

export const getClusterApiGroups = (params: ClusterApiGroupListParams): Promise<ApiResponse<ClusterApiGroupListResponse>> => {
  return service({
    // Original: /kubernetes/getClusterApiGroups (POST with IdParams for clusterId)
    // url: '/kubernetes/getClusterApiGroups',
    // method: 'post',
    // data: { id: params.clusterId }
    url: `/kubernetes/cluster/${params.clusterId}/apigroups`, // RESTful approach
    method: 'get',
  })
}

// === Cluster Credential Management ===
// The document implies managing credentials beyond the initial KubeConfig.
// The existing `CreateCredential` takes `ClusterData`, which is not suitable.
// This section needs more clarity on what kind of credentials (e.g., service account tokens for specific tasks).
// Placeholder for now.

export interface CreateCredentialPayload { // This needs proper definition based on backend
    clusterId: string;
    type: K8sClusterCredential['type'];
    value: string; // Example: token string
    description?: string;
}

export const createClusterCredential = (payload: CreateCredentialPayload): Promise<ApiResponse<K8sClusterCredential>> => {
  return service({
    url: `/kubernetes/cluster/${payload.clusterId}/credentials`, // Example RESTful endpoint
    method: 'post',
    data: payload
  })
}

// Get User by ID (Original: getUserById) - Unclear if this is a general system user or cluster-specific. Assuming general for now.
// This might belong to a more general user API module if not strictly for k8s cluster context.
// For now, keeping it as per original file structure, but flagged for review.
export interface UserDetailResponse { // Define a proper User type if not already available globally
    id: string;
    username: string;
    email?: string; // etc.
}
export const getUserById = (userId: string): Promise<ApiResponse<UserDetailResponse>> => {
  return service({
    url: '/kubernetes/getUserById', // This endpoint seems kubernetes-specific by its path
    method: 'post', // Assuming POST as per original
    data: { id: userId }
  })
}


// getClusterUserNamespace & getClusterListNamespace - These seem specific to a user's access within a cluster.
// Their purpose needs to be clarified:
// - getClusterUserNamespace: Namespaces a specific user (from your system) has access to in a cluster?
// - getClusterListNamespace: All namespaces in a cluster? (This is usually `GET /api/v1/namespaces` to k8s API)

export interface NamespaceListResponse {
    // Assuming a simple list of namespace names for now
    // This should be `V1NamespaceList` or similar if directly from Kubernetes
    items: Array<{ name: string; [key: string]: any }>; // Replace with proper Namespace type
    total?: number;
}

// Assuming this gets namespaces accessible by the *currently authenticated user* for that cluster,
// or for a specific system user if `userId` is provided.
export const getAccessibleNamespaces = (clusterId: string, userId?: string): Promise<ApiResponse<NamespaceListResponse>> => {
    let url = `/kubernetes/cluster/${clusterId}/namespaces`;
    const params: any = {};
    if (userId) {
        params.userId = userId; // Or however the backend identifies the user
    }
    return service({
        // url: '/kubernetes/getClusterUserNamespace', // Original endpoint
        // method: 'post',
        // data: { id: userId, cluster_id: clusterId } // Original payload was just 'id'
        url: url,
        method: 'get', // More RESTful for fetching a list
        params: params
    });
}

// Assuming this gets all namespaces in a cluster, typically an admin/privileged operation.
export const getAllClusterNamespaces = (clusterId: string): Promise<ApiResponse<NamespaceListResponse>> => {
    return service({
        // url: `/kubernetes/getClusterListNamespace`, // Original endpoint
        // method: 'post',
        // data: { id: clusterId } // Original payload was just 'id'
        url: `/kubernetes/cluster/${clusterId}/all-namespaces`, // Example RESTful
        method: 'get'
    });
}
