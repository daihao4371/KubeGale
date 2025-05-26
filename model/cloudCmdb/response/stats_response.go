package response

// ResourceCounts holds the counts for different resource types.
// Pointers are used to allow omitting a field if the count is zero or the type wasn't requested.
type ResourceCounts struct {
	VirtualMachines *int64 `json:"virtual_machines,omitempty"`
	RDSInstances    *int64 `json:"rds_instances,omitempty"`
	LoadBalancers   *int64 `json:"load_balancers,omitempty"`
}

// ProviderResourceCount represents the resource counts for a single cloud provider.
type ProviderResourceCount struct {
	ProviderID     uint           `json:"provider_id"`
	ProviderName   string         `json:"provider_name"`   // From CloudPlatform.Name
	ProviderType   string         `json:"provider_type"`   // From CloudPlatform.Platform
	ResourceCounts ResourceCounts `json:"resource_counts"`
}
