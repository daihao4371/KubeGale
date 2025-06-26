export interface Props {
  form: {
    id: number
  }
}

export interface Rule {
  apiGroups: string[]
  resources: string[]
  verbs: string[]
  resourceOptions?: string[]
  verbOptions?: string[]
}

export interface RoleForm {
  name: string
  description: string
  rules: Rule[]
}

export interface GroupResource {
  group: string
  resources: {
    resource: string
    verbs: string[]
  }[]
}

export interface ClusterRole {
  metadata: {
    name: string
    annotations: {
      description: string
      [key: string]: string
    }
    creationTimestamp: string
  }
  rules: Rule[]
}

export interface ClusterFormData {
  id?: string;
  name: string;
  kube_type: number;
  kube_config: string;
  api_address: string;
  prometheus_url?: string;
  prometheus_auth_type: number;
  prometheus_user?: string;
  prometheus_pwd?: string;
  createdAt?: string;
}

// --- Enhanced Kubernetes Types ---

// Represents a Kubernetes Cluster managed by the system
export interface K8sCluster {
  id?: string; // System DB ID
  name: string;
  kube_config: string; // Raw KubeConfig content
  api_address: string; // API server endpoint
  // kube_type: 1 for KubeConfig based, 2 for Agent based (as per current API)
  // Consider using an enum or string literal type if more types are added
  kube_type: number;
  createdAt?: string; // ISO date string
  updatedAt?: string; // ISO date string

  // Prometheus Integration (as per current API)
  prometheus_url?: string;
  // prometheus_auth_type: 0 for None, 1 for Basic Auth
  prometheus_auth_type: number;
  prometheus_user?: string;
  prometheus_pwd?: string; // Consider security implications if storing this long-term

  // Additional metadata from the cluster itself (example, might be fetched on demand)
  version?: string; // Kubernetes version
  status?: 'active' | 'inactive' | 'error'; // System-level status
  provider?: string; // e.g., GKE, EKS, AKS, OnPrem
}

// User associated with a K8sCluster (specific to this application's user management)
export interface K8sClusterUser {
  id?: string; // System DB ID for the user-cluster link or user ID itself
  clusterId: string; // ID of the K8sCluster
  userId: string; // ID of the user in the main user system
  username: string; // Username
  k8sUsername?: string; // Actual username within Kubernetes (if different, e.g. from a cert)
  roles: string[]; // List of K8sClusterRole names or IDs assigned to this user in this cluster
  joinedAt?: string; // ISO date string
}

// Represents a Role within a K8sCluster (could be Kubernetes ClusterRole or a system-defined role)
// This should align with Kubernetes RBAC's Role or ClusterRole structure if managing them directly
export interface K8sClusterRoleDefinition {
  id?: string; // System DB ID or Kubernetes role name
  clusterId: string; // ID of the K8sCluster this role is defined in (if not a global template)
  name: string; // Role name (e.g., 'view-only', 'cluster-admin')
  description?: string;
  // Based on Kubernetes RBAC PolicyRule
  rules: K8sPolicyRule[];
  createdAt?: string; // ISO date string
  updatedAt?: string; // ISO date string
}

// Kubernetes PolicyRule for RBAC
// Replaces the existing 'Rule' type for more clarity
export interface K8sPolicyRule {
  apiGroups: string[]; // e.g., ["", "apps", "batch"]
  resources: string[]; // e.g., ["pods", "deployments", "services"]
  verbs: string[]; // e.g., ["get", "list", "watch", "create", "update", "patch", "delete"]
  resourceNames?: string[]; // Specific resource instances by name
  nonResourceURLs?: string[]; // For non-resource endpoints like /healthz
}


// API Group information for a cluster
export interface K8sApiGroup {
  name: string; // e.g., "apps"
  versions: string[]; // e.g., ["v1"]
  preferredVersion: string; // e.g., "v1"
}

// Credentials for accessing a cluster (could be more abstract)
// The current `k8sCluster.ts` uses ClusterData for CreateCredential, which needs refinement.
// This is a placeholder, as "Credential" can mean many things (kubeconfig, token, certs).
// For now, KubeConfig is the primary credential via K8sCluster.kube_config.
export interface K8sClusterCredential {
  id?: string;
  clusterId: string;
  type: 'kubeconfig' | 'token' | 'serviceAccount'; // Example types
  value: string; // The credential itself (e.g., token string, or reference to SA)
  description?: string;
}

// --- Existing types to be reviewed/merged/removed ---
// 'Props' seems generic and might not be needed here or should be more specific.
// 'RoleForm' can be replaced by K8sClusterRoleDefinition for form data.
// 'GroupResource' seems like a custom structure, review if K8sPolicyRule covers its intent.
// 'ClusterRole' (existing one) is similar to K8sClusterRoleDefinition but less detailed.
// 'ClusterFormData' is now covered by K8sCluster.


// --- Workload Types ---

// Kubernetes Deployment (apps/v1.Deployment)
// This is a more detailed representation.
// Consider using official types from @kubernetes/client-node if possible for full accuracy.
export interface K8sDeployment {
  apiVersion?: 'apps/v1';
  kind?: 'Deployment';
  metadata: {
    name: string;
    namespace: string;
    labels?: Record<string, string>;
    annotations?: Record<string, string>;
    uid?: string;
    creationTimestamp?: string;
    generation?: number;
    [key: string]: any; // Allow other metadata fields
  };
  spec: {
    replicas?: number;
    selector: { // V1LabelSelector
      matchLabels?: Record<string, string>;
      matchExpressions?: Array<{ key: string; operator: string; values?: string[] }>;
    };
    template: { // V1PodTemplateSpec
      metadata?: {
        labels?: Record<string, string>;
        annotations?: Record<string, string>;
        [key: string]: any;
      };
      spec: K8sPodSpec; // V1PodSpec
    };
    strategy?: { // V1DeploymentStrategy
      type?: 'RollingUpdate' | 'Recreate';
      rollingUpdate?: { // V1RollingUpdateDeployment
        maxUnavailable?: number | string;
        maxSurge?: number | string;
      };
    };
    minReadySeconds?: number;
    revisionHistoryLimit?: number;
    paused?: boolean;
    progressDeadlineSeconds?: number;
  };
  status?: { // V1DeploymentStatus
    observedGeneration?: number;
    replicas?: number;
    updatedReplicas?: number;
    readyReplicas?: number;
    availableReplicas?: number;
    unavailableReplicas?: number;
    conditions?: Array<K8sDeploymentCondition>; // V1DeploymentCondition
    collisionCount?: number;
  };
}

// Kubernetes DeploymentCondition (apps/v1.DeploymentCondition)
export interface K8sDeploymentCondition {
  type: string; // e.g., Available, Progressing, ReplicaFailure
  status: 'True' | 'False' | 'Unknown';
  lastUpdateTime?: string;
  lastTransitionTime?: string;
  reason?: string;
  message?: string;
}

// Kubernetes PodSpec (core/v1.PodSpec) - Simplified
export interface K8sPodSpec {
  containers: Array<K8sContainer>;
  volumes?: Array<K8sVolume>; // V1Volume
  restartPolicy?: 'Always' | 'OnFailure' | 'Never';
  terminationGracePeriodSeconds?: number;
  dnsPolicy?: 'ClusterFirst' | 'Default' | 'None' | 'ClusterFirstWithHostNet';
  serviceAccountName?: string;
  nodeSelector?: Record<string, string>;
  affinity?: any; // V1Affinity
  tolerations?: any[]; // V1Toleration
  imagePullSecrets?: Array<{ name: string }>; // V1LocalObjectReference
  // ... other PodSpec fields
  [key: string]: any;
}

// Kubernetes Container (core/v1.Container) - Simplified
export interface K8sContainer {
  name: string;
  image: string;
  imagePullPolicy?: 'Always' | 'IfNotPresent' | 'Never';
  command?: string[];
  args?: string[];
  ports?: Array<K8sContainerPort>; // V1ContainerPort
  env?: Array<K8sEnvVar>; // V1EnvVar
  resources?: K8sResourceRequirements; // V1ResourceRequirements
  volumeMounts?: Array<K8sVolumeMount>; // V1VolumeMount
  livenessProbe?: K8sProbe; // V1Probe
  readinessProbe?: K8sProbe; // V1Probe
  startupProbe?: K8sProbe; // V1Probe
  // ... other Container fields
  [key: string]: any;
}

// Kubernetes ContainerPort (core/v1.ContainerPort)
export interface K8sContainerPort {
  name?: string;
  hostPort?: number;
  containerPort: number;
  protocol?: 'TCP' | 'UDP' | 'SCTP';
  hostIP?: string;
}

// Kubernetes EnvVar (core/v1.EnvVar)
export interface K8sEnvVar {
  name: string;
  value?: string;
  valueFrom?: any; // V1EnvVarSource
}

// Kubernetes ResourceRequirements (core/v1.ResourceRequirements)
export interface K8sResourceRequirements {
  limits?: Record<string, string>; // e.g., { cpu: "500m", memory: "128Mi" }
  requests?: Record<string, string>;
}

// Kubernetes Volume (core/v1.Volume) - Highly simplified, actual is a union type
export interface K8sVolume {
  name: string;
  // Examples of volume sources, actual type is complex
  configMap?: { name?: string; items?: Array<{key: string; path: string}>; defaultMode?: number; optional?: boolean };
  secret?: { secretName?: string; items?: Array<{key: string; path: string}>; defaultMode?: number; optional?: boolean };
  emptyDir?: any; // V1EmptyDirVolumeSource
  persistentVolumeClaim?: { claimName: string; readOnly?: boolean; };
  // ... other volume types
  [key: string]: any;
}

// Kubernetes VolumeMount (core/v1.VolumeMount)
export interface K8sVolumeMount {
  name: string;
  mountPath: string;
  readOnly?: boolean;
  subPath?: string;
  mountPropagation?: string;
  subPathExpr?: string;
}

// Kubernetes Probe (core/v1.Probe) - Simplified
export interface K8sProbe {
  exec?: { command?: string[] }; // V1ExecAction
  httpGet?: { path?: string; port: number | string; scheme?: string; host?: string; httpHeaders?: Array<{name: string; value: string}> }; // V1HTTPGetAction
  tcpSocket?: { port: number | string; host?: string }; // V1TCPSocketAction
  grpc?: { port: number; service?: string }; // V1GRPCAction
  initialDelaySeconds?: number;
  timeoutSeconds?: number;
  periodSeconds?: number;
  successThreshold?: number;
  failureThreshold?: number;
  terminationGracePeriodSeconds?: number;
}


// Kubernetes ReplicaSet (apps/v1.ReplicaSet) - Simplified
export interface K8sReplicaSet {
  apiVersion?: 'apps/v1';
  kind?: 'ReplicaSet';
  metadata: {
    name: string;
    namespace: string;
    labels?: Record<string, string>;
    annotations?: Record<string, string>;
    ownerReferences?: Array<{ apiVersion: string; kind: string; name: string; uid: string; controller?: boolean; blockOwnerDeletion?:boolean }>;
    uid?: string;
    creationTimestamp?: string;
  };
  spec: {
    replicas?: number;
    selector: { // V1LabelSelector
      matchLabels?: Record<string, string>;
      matchExpressions?: Array<{ key: string; operator: string; values?: string[] }>;
    };
    template?: { // V1PodTemplateSpec, same as Deployment's template
      metadata?: {
        labels?: Record<string, string>;
      };
      spec: K8sPodSpec;
    };
    minReadySeconds?: number;
  };
  status?: { // V1ReplicaSetStatus
    replicas: number;
    fullyLabeledReplicas?: number;
    readyReplicas?: number;
    availableReplicas?: number;
    observedGeneration?: number;
    conditions?: Array<any>; // V1ReplicaSetCondition
  };
}

// --- Node Types ---

// Kubernetes Node (core/v1.Node)
export interface K8sNode {
  apiVersion?: 'v1';
  kind?: 'Node';
  metadata: {
    name: string;
    uid?: string;
    labels?: Record<string, string>;
    annotations?: Record<string, string>;
    creationTimestamp?: string;
    [key: string]: any;
  };
  spec?: K8sNodeSpec;
  status?: K8sNodeStatus;
}

// Kubernetes NodeSpec (core/v1.NodeSpec)
export interface K8sNodeSpec {
  podCIDR?: string;
  podCIDRs?: string[];
  providerID?: string;
  unschedulable?: boolean;
  taints?: K8sTaint[];
  // configSource?: V1NodeConfigSource; // More complex type
  [key: string]: any;
}

// Kubernetes Taint (core/v1.Taint)
export interface K8sTaint {
  key: string;
  value?: string;
  effect: 'NoSchedule' | 'PreferNoSchedule' | 'NoExecute';
  timeAdded?: string;
}

// Kubernetes NodeStatus (core/v1.NodeStatus)
export interface K8sNodeStatus {
  capacity?: Record<string, string>; // ResourceList (e.g., cpu: "2", memory: "4Gi")
  allocatable?: Record<string, string>; // ResourceList
  conditions?: K8sNodeCondition[];
  addresses?: K8sNodeAddress[];
  daemonEndpoints?: { // V1NodeDaemonEndpoints
    kubeletEndpoint?: {
      Port: number;
    };
  };
  nodeInfo?: K8sNodeSystemInfo; // V1NodeSystemInfo
  images?: any[]; // V1ContainerImage
  volumesInUse?: string[];
  volumesAttached?: any[]; // V1AttachedVolume
  // config?: V1NodeConfigStatus; // More complex type
  [key: string]: any;
}

// Kubernetes NodeCondition (core/v1.NodeCondition)
export interface K8sNodeCondition {
  type: string; // e.g., Ready, MemoryPressure, DiskPressure, PIDPressure, NetworkUnavailable
  status: 'True' | 'False' | 'Unknown';
  lastHeartbeatTime?: string;
  lastTransitionTime?: string;
  reason?: string;
  message?: string;
}

// Kubernetes NodeAddress (core/v1.NodeAddress)
export interface K8sNodeAddress {
  type: 'Hostname' | 'InternalIP' | 'ExternalIP'; // More types exist
  address: string;
}

// Kubernetes NodeSystemInfo (core/v1.NodeSystemInfo)
export interface K8sNodeSystemInfo {
  machineID: string;
  systemUUID: string;
  bootID: string;
  kernelVersion: string;
  osImage: string;
  containerRuntimeVersion: string;
  kubeletVersion: string;
  kubeProxyVersion: string;
  operatingSystem: string;
  architecture: string;
}

// Node Metrics (custom structure, typically from metrics-server or Prometheus)
// This is a generic structure; specific metrics might vary.
export interface K8sNodeMetrics {
  metadata: {
    name: string; // Node name
    [key: string]: any;
  };
  timestamp?: string;
  window?: string; // e.g., "1m"
  usage?: {
    cpu?: string; // e.g., "500m"
    memory?: string; // e.g., "1Gi"
    pods?: string; // Number of pods
    // Potentially disk usage, network I/O etc.
    [key: string]: any;
  };
}