package cloudCmdb

import (
	"KubeGale/global"
	"KubeGale/model/common/response"
	"KubeGale/service" // Assuming CloudStatsService is accessible via service.ServiceGroupApp
	"KubeGale/service/cloudCmdb" // For GetResourceCountsByProviderInput
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudStatsApi struct{}

// GetResourceCountsByProviderHandler
// @Tags CloudCMDB Statistics
// @Summary Get resource counts grouped by cloud provider
// @Produce application/json
// @Param resource_types query string false "Comma-separated list of resource types (e.g., vm,rds,lb). If empty, all types are counted."
// @Param provider_ids query string false "Comma-separated list of cloud platform IDs. If empty, counts for all providers."
// @Success 200 {object} response.Response{data=[]response.ProviderResourceCount,msg=string} "Successfully retrieved resource counts"
// @Router /cloudCmdb/stats/countsByProvider [get]
func (s *CloudStatsApi) GetResourceCountsByProviderHandler(c *gin.Context) {
	var input cloudCmdb.GetResourceCountsByProviderInput // Use the input struct from the service package

	// Parse resource_types
	resourceTypesQuery := c.Query("resource_types")
	if resourceTypesQuery != "" {
		input.ResourceTypes = strings.Split(resourceTypesQuery, ",")
		for i, rt := range input.ResourceTypes {
			input.ResourceTypes[i] = strings.TrimSpace(strings.ToLower(rt))
		}
	}

	// Parse provider_ids
	providerIDsQuery := c.Query("provider_ids")
	if providerIDsQuery != "" {
		idStrings := strings.Split(providerIDsQuery, ",")
		for _, idStr := range idStrings {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
			if err != nil {
				global.KUBEGALE_LOG.Error("Invalid provider_id format", zap.String("idStr", idStr), zap.Error(err))
				response.FailWithMessage("Invalid provider_id format: "+idStr, c)
				return
			}
			input.ProviderIDs = append(input.ProviderIDs, uint(id))
		}
	}
	
	// Access the service - assuming CloudStatsService is part of CloudCmdbServiceGroup
	cloudStatsService := service.ServiceGroupApp.CloudCmdbServiceGroup.CloudStatsService 
	counts, err := cloudStatsService.GetResourceCountsByProvider(input)
	if err != nil {
		global.KUBEGALE_LOG.Error("Failed to get resource counts by provider", zap.Error(err))
		response.FailWithMessage("Failed to get resource counts: "+err.Error(), c)
		return
	}

	response.OkWithData(counts, c)
}
