package cloudCmdb

import (
	"KubeGale/global"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Helper function to initialize a mock GORM DB (can be moved to a common test helper)
func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) // Use regexp matcher for flexibility
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing gorm", err)
	}
	return gormDB, mock
}

func TestCloudStatsService_GetResourceCountsByProvider_AllProvidersAllTypes(t *testing.T) {
	mockDB, mock := newMockDB(t)
	global.KUBEGALE_DB = mockDB

	service := CloudStatsService{}
	input := GetResourceCountsByProviderInput{} // No filters

	// Mock fetching CloudPlatforms
	platformRows := sqlmock.NewRows([]string{"id", "name", "platform"}).
		AddRow(1, "Aliyun East", "aliyun").
		AddRow(2, "AWS North", "aws")
	mock.ExpectQuery("^SELECT \\* FROM `cloud_platform`").WillReturnRows(platformRows)

	// Mock counts for Aliyun (ID 1)
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_virtual_machine` WHERE cloud_platform_id = \\?$").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_rds` WHERE cloud_platform_id = \\?$").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_load_balancer` WHERE cloud_platform_id = \\?$").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock counts for AWS (ID 2)
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_virtual_machine` WHERE cloud_platform_id = \\?$").
		WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_rds` WHERE cloud_platform_id = \\?$").
		WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0)) // Example of zero count
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_load_balancer` WHERE cloud_platform_id = \\?$").
		WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	results, err := service.GetResourceCountsByProvider(input)

	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Check Aliyun results
	assert.Equal(t, uint(1), results[0].ProviderID)
	assert.Equal(t, "Aliyun East", results[0].ProviderName)
	assert.Equal(t, "aliyun", results[0].ProviderType)
	assert.Equal(t, int64(10), *results[0].ResourceCounts.VirtualMachines)
	assert.Equal(t, int64(2), *results[0].ResourceCounts.RDSInstances)
	assert.Equal(t, int64(1), *results[0].ResourceCounts.LoadBalancers)

	// Check AWS results
	assert.Equal(t, uint(2), results[1].ProviderID)
	assert.Equal(t, "AWS North", results[1].ProviderName)
	assert.Equal(t, "aws", results[1].ProviderType)
	assert.Equal(t, int64(20), *results[1].ResourceCounts.VirtualMachines)
	assert.Equal(t, int64(0), *results[1].ResourceCounts.RDSInstances)
	assert.Equal(t, int64(3), *results[1].ResourceCounts.LoadBalancers)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCloudStatsService_GetResourceCountsByProvider_FilterByType(t *testing.T) {
	mockDB, mock := newMockDB(t)
	global.KUBEGALE_DB = mockDB

	service := CloudStatsService{}
	input := GetResourceCountsByProviderInput{
		ResourceTypes: []string{"vm", "lb"}, // Only VM and LB
	}

	platformRows := sqlmock.NewRows([]string{"id", "name", "platform"}).
		AddRow(1, "Test Provider", "test")
	mock.ExpectQuery("^SELECT \\* FROM `cloud_platform`").WillReturnRows(platformRows)

	// Expect VM and LB counts, but not RDS
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_virtual_machine` WHERE cloud_platform_id = \\?$").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectQuery("^SELECT count\\(\\*\\) FROM `cloud_load_balancer` WHERE cloud_platform_id = \\?$").
		WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	results, err := service.GetResourceCountsByProvider(input)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, uint(1), results[0].ProviderID)
	assert.NotNil(t, results[0].ResourceCounts.VirtualMachines)
	assert.Equal(t, int64(5), *results[0].ResourceCounts.VirtualMachines)
	assert.Nil(t, results[0].ResourceCounts.RDSInstances) // RDS should not be counted
	assert.NotNil(t, results[0].ResourceCounts.LoadBalancers)
	assert.Equal(t, int64(2), *results[0].ResourceCounts.LoadBalancers)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
