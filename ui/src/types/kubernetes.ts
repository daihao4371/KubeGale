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