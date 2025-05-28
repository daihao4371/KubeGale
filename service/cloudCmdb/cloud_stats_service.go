package cloudCmdb

import (
	"KubeGale/global"
	"KubeGale/model/cloudCmdb"
	"KubeGale/model/cloudCmdb/response" // New response DTO package
	"fmt"
	"strings"

	"go.uber.org/zap" // Added for logging
)

type CloudStatsService struct{}

// GetResourceCountsByProviderInput defines the input parameters for the service method.
type GetResourceCountsByProviderInput struct {
	ProviderIDs   []uint   // Filter by specific provider IDs (optional)
	ResourceTypes []string // Filter by specific resource types: "vm", "rds", "lb" (optional)
}

// GetResourceCountsByProvider calculates resource counts grouped by cloud provider.
func (s *CloudStatsService) GetResourceCountsByProvider(input GetResourceCountsByProviderInput) ([]response.ProviderResourceCount, error) {
	var results []response.ProviderResourceCount
	var cloudPlatforms []cloudCmdb.CloudPlatform

	db := global.KUBEGALE_DB.Model(&cloudCmdb.CloudPlatform{})

	if len(input.ProviderIDs) > 0 {
		db = db.Where("id IN ?", input.ProviderIDs)
	}

	if err := db.Find(&cloudPlatforms).Error; err != nil {
		global.KUBEGALE_LOG.Error("Failed to fetch cloud platforms", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch cloud platforms: %w", err)
	}

	// Determine which resource types to count
	countVM := true
	countRDS := true
	countLB := true
	if len(input.ResourceTypes) > 0 {
		countVM, countRDS, countLB = false, false, false
		for _, rt := range input.ResourceTypes {
			switch strings.ToLower(rt) {
			case "vm":
				countVM = true
			case "rds":
				countRDS = true
			case "lb":
				countLB = true
			}
		}
	}

	for _, cp := range cloudPlatforms {
		var vmCount, rdsCount, lbCount int64
		var err error

		rc := response.ResourceCounts{}

		if countVM {
			err = global.KUBEGALE_DB.Model(&cloudCmdb.VirtualMachine{}).Where("cloud_platform_id = ?", cp.ID).Count(&vmCount).Error
			if err != nil {
				global.KUBEGALE_LOG.Warn("Failed to count VMs for provider", zap.Uint("providerID", cp.ID), zap.Error(err))
				// Continue counting other resources
			}
			rc.VirtualMachines = &vmCount
		}

		if countRDS {
			err = global.KUBEGALE_DB.Model(&cloudCmdb.RDS{}).Where("cloud_platform_id = ?", cp.ID).Count(&rdsCount).Error
			if err != nil {
				global.KUBEGALE_LOG.Warn("Failed to count RDS instances for provider", zap.Uint("providerID", cp.ID), zap.Error(err))
			}
			rc.RDSInstances = &rdsCount
		}

		if countLB {
			err = global.KUBEGALE_DB.Model(&cloudCmdb.LoadBalancer{}).Where("cloud_platform_id = ?", cp.ID).Count(&lbCount).Error
			if err != nil {
				global.KUBEGALE_LOG.Warn("Failed to count LoadBalancers for provider", zap.Uint("providerID", cp.ID), zap.Error(err))
			}
			rc.LoadBalancers = &lbCount
		}

		results = append(results, response.ProviderResourceCount{
			ProviderID:     cp.ID,
			ProviderName:   cp.Name,
			ProviderType:   cp.Platform, // Assuming CloudPlatform.Platform stores type like "aliyun", "aws"
			ResourceCounts: rc,
		})
	}

	return results, nil
}
