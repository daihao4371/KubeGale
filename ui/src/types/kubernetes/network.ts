// 网络资源基础类型
export interface NetworkBase {
  metadata: {
    name: string
    namespace: string
    uid: string
    creationTimestamp: string
    labels?: Record<string, string>
    annotations?: Record<string, string>
  }
}

// Service 相关类型
export interface Service extends NetworkBase {
  spec: {
    type: string // ClusterIP, NodePort, LoadBalancer, ExternalName
    selector?: Record<string, string>
    ports?: Array<{
      name?: string
      port: number
      targetPort: number | string
      protocol: string
      nodePort?: number
    }>
    clusterIP?: string
    clusterIPs?: string[]
    externalIPs?: string[]
    sessionAffinity?: string
    loadBalancerIP?: string
    loadBalancerSourceRanges?: string[]
    externalName?: string
  }
  status: {
    loadBalancer?: {
      ingress?: Array<{
        ip?: string
        hostname?: string
      }>
    }
  }
}

// Ingress 相关类型
export interface Ingress extends NetworkBase {
  spec: {
    ingressClassName?: string
    defaultBackend?: IngressBackend
    tls?: Array<{
      hosts?: string[]
      secretName?: string
    }>
    rules?: Array<{
      host?: string
      http?: {
        paths: Array<{
          path?: string
          pathType: string
          backend: IngressBackend
        }>
      }
    }>
  }
  status: {
    loadBalancer?: {
      ingress?: Array<{
        ip?: string
        hostname?: string
      }>
    }
  }
}

export interface IngressBackend {
  service?: {
    name: string
    port: {
      number?: number
      name?: string
    }
  }
  resource?: {
    apiGroup?: string
    kind: string
    name: string
  }
}

// Endpoint 相关类型
export interface Endpoint extends NetworkBase {
  subsets?: Array<{
    addresses?: Array<{
      ip: string
      hostname?: string
      nodeName?: string
      targetRef?: {
        kind: string
        name: string
        namespace: string
        uid: string
      }
    }>
    notReadyAddresses?: Array<{
      ip: string
      hostname?: string
      nodeName?: string
      targetRef?: {
        kind: string
        name: string
        namespace: string
        uid: string
      }
    }>
    ports?: Array<{
      name?: string
      port: number
      protocol: string
    }>
  }>
}

// NetworkPolicy 相关类型
export interface NetworkPolicy extends NetworkBase {
  spec: {
    podSelector: {
      matchLabels?: Record<string, string>
      matchExpressions?: Array<{
        key: string
        operator: string
        values?: string[]
      }>
    }
    policyTypes: string[]
    ingress?: Array<{
      from?: Array<{
        podSelector?: {
          matchLabels?: Record<string, string>
        }
        namespaceSelector?: {
          matchLabels?: Record<string, string>
        }
        ipBlock?: {
          cidr: string
          except?: string[]
        }
      }>
      ports?: Array<{
        protocol?: string
        port?: number | string
        endPort?: number
      }>
    }>
    egress?: Array<{
      to?: Array<{
        podSelector?: {
          matchLabels?: Record<string, string>
        }
        namespaceSelector?: {
          matchLabels?: Record<string, string>
        }
        ipBlock?: {
          cidr: string
          except?: string[]
        }
      }>
      ports?: Array<{
        protocol?: string
        port?: number | string
        endPort?: number
      }>
    }>
  }
  status: {}
}

// API 请求类型
export interface GetNetworkListRequest {
  cluster_id: number
  namespace?: string
  page?: number
  pageSize?: number
  keyword?: string
  labelSelector?: string
  fieldSelector?: string
}

export interface GetNetworkDetailRequest {
  cluster_id: number
  namespace: string
  name: string
}

export interface CreateNetworkRequest {
  cluster_id: number
  namespace: string
  content: any
}

export interface UpdateNetworkRequest {
  cluster_id: number
  namespace: string
  name: string
  content: any
}

export interface DeleteNetworkRequest {
  cluster_id: number
  namespace: string
  name: string
}

// API 响应类型
export interface NetworkListResponse<T = any> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface NetworkDetailResponse<T = any> {
  items: T
}