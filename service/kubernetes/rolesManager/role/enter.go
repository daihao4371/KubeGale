package role

import (
	"KubeGale/service/kubernetes/rolesManager/podSecurityPolicies"
)

type ServiceGroup struct {
	K8sRoleService
	PodSecurityPoliciesServiceGroup podSecurityPolicies.ServiceGroup
}
