package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"KubeGale/model/im/request"
	// "KubeGale/model/im/response" // May be needed for TestNotification
	// "KubeGale/utils" // For Verify if used directly in service, or for initializing mocks
	"errors"
	"regexp" // For go-sqlmock if used
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock" // Example mocking library
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Helper function to initialize a mock GORM DB
func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when initializing gorm", err)
	}
	return gormDB, mock
}

func TestNotificationService_CreateFeiShu(t *testing.T) {
	// Setup mock DB
	mockDB, mock := newMockDB(t)
	global.KUBEGALE_DB = mockDB // Replace global DB with mock

	service := NotificationService{}

	req := request.CreateFeiShuRequest{
		Name:           "Test FeiShu",
		NotifyEvents:   []string{"alert"},
		SendDailyStats: true,
		WebhookURL:     "http://fakehook.com/feishu",
	}

	// Expected NotificationConfig creation
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `im_notification_configs`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), req.Name, im.NotificationTypeFeiShu, strings.Join(req.NotifyEvents, ","), req.SendDailyStats).
		WillReturnResult(sqlmock.NewResult(1, 1)) // ID 1, 1 row affected

	// Expected FeiShuConfig creation
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `im_fei_shu_configs`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), uint(1), req.WebhookURL). // Assumes NotificationConfigID is 1
		WillReturnResult(sqlmock.NewResult(1, 1)) // ID 1, 1 row affected
	mock.ExpectCommit()

	// Expected query after creation to preload NotificationConfig
	// This part is tricky with sqlmock directly for Preload.
	// Often, you'd mock the rows returned by the preload query.
	// For simplicity, we might skip testing the preload directly with sqlmock's exec expectations
	// and trust GORM if the creation queries are correct.
	// Or, expect a SELECT query that GORM would make for preloading.
	// Example:
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `im_fei_shu_configs` WHERE `im_fei_shu_configs`.`id` = ? ORDER BY `im_fei_shu_configs`.`id` LIMIT 1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "notification_config_id", "robot_url"}).AddRow(1, 1, req.WebhookURL))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `im_notification_configs` WHERE `im_notification_configs`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type"}).AddRow(1, req.Name, im.NotificationTypeFeiShu))


	createdConfig, err := service.CreateFeiShu(req)

	assert.NoError(t, err)
	assert.NotNil(t, createdConfig)
	assert.Equal(t, req.WebhookURL, createdConfig.RobotURL)
	assert.NotNil(t, createdConfig.NotificationConfig)
	assert.Equal(t, req.Name, createdConfig.NotificationConfig.Name)
	assert.Equal(t, im.NotificationTypeFeiShu, createdConfig.NotificationConfig.Type)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNotificationService_CreateDingTalk(t *testing.T) {
	// Setup mock DB
	mockDB, mock := newMockDB(t)
	global.KUBEGALE_DB = mockDB

	service := NotificationService{}

	req := request.CreateDingTalkRequest{
		Name:           "Test DingTalk",
		NotifyEvents:   []string{"event"},
		SendDailyStats: false,
		WebhookURL:     "http://fakehook.com/dingtalk",
		Secret:         "sec123",
	}

	mock.ExpectBegin()
	// NotificationConfig
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `im_notification_configs`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), req.Name, im.NotificationTypeDingTalk, strings.Join(req.NotifyEvents, ","), req.SendDailyStats).
		WillReturnResult(sqlmock.NewResult(2, 1)) // ID 2
	// DingTalkConfig
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `im_ding_talk_configs`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), uint(2), req.WebhookURL, req.Secret).
		WillReturnResult(sqlmock.NewResult(1, 1)) // ID 1 (specific table ID)
	mock.ExpectCommit()

	// Mocking the preload part
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `im_ding_talk_configs` WHERE `im_ding_talk_configs`.`id` = ? ORDER BY `im_ding_talk_configs`.`id` LIMIT 1")).
		WithArgs(1). // Assuming DingTalkConfig ID is 1
		WillReturnRows(sqlmock.NewRows([]string{"id", "notification_config_id", "webhook_url", "secret"}).AddRow(1, 2, req.WebhookURL, req.Secret))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `im_notification_configs` WHERE `im_notification_configs`.`id` = ?")).
		WithArgs(2). // NotificationConfig ID is 2
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type"}).AddRow(2, req.Name, im.NotificationTypeDingTalk))


	createdConfig, err := service.CreateDingTalk(req)

	assert.NoError(t, err)
	assert.NotNil(t, createdConfig)
	assert.Equal(t, req.WebhookURL, createdConfig.WebhookURL)
	assert.Equal(t, req.Secret, createdConfig.Secret)
	assert.NotNil(t, createdConfig.NotificationConfig)
	assert.Equal(t, req.Name, createdConfig.NotificationConfig.Name)
	assert.Equal(t, im.NotificationTypeDingTalk, createdConfig.NotificationConfig.Type)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TODO: Add more tests for Update, Delete, GetById, GetList, TestNotification for both types.
// - Test cases for error conditions (e.g., name conflict, DB errors).
// - Test cases for Update operations (partial updates, updating card content).
// - Test GetNotificationList with various filters and pagination.
// - Test TestNotification by mocking the message sender (e.g. messageIm.MessageServiceApp).

// Example for mocking MessageServiceApp for TestNotification (conceptual)
/*
type MockMessageService struct {
	SendFeiShuFunc func(config response.NotificationDetailConfig, cardContent response.CardContentDetail, message string) error
	SendDingTalkFunc func(config response.NotificationDetailConfig, cardContent response.CardContentDetail, message string) error
}
func (m *MockMessageService) SendFeiShuMessage(config response.NotificationDetailConfig, cardContent response.CardContentDetail, message string) error {
	if m.SendFeiShuFunc != nil {
		return m.SendFeiShuFunc(config, cardContent, message)
	}
	return nil
}
func (m *MockMessageService) SendDingTalkMessage(config response.NotificationDetailConfig, cardContent response.CardContentDetail, message string) error {
	if m.SendDingTalkFunc != nil {
		return m.SendDingTalkFunc(config, cardContent, message)
	}
	return nil
}

func TestNotificationService_TestNotification_DingTalk(t *testing.T) {
    // ... setup DB mock for fetching configs ...
    
    // Setup mock message sender
    originalMessageService := messageIm.MessageServiceApp
    mockSender := &MockMessageService{}
    messageIm.MessageServiceApp = mockSender // Replace global instance
    defer func() { messageIm.MessageServiceApp = originalMessageService }() // Restore

    mockSender.SendDingTalkFunc = func(config response.NotificationDetailConfig, cardContent response.CardContentDetail, message string) error {
        assert.Equal(t, "expected_webhook_url", config.RobotURL)
        assert.Equal(t, "test message", message)
        return nil // Simulate success
    }

    // ... call service.TestNotification ...
    // ... assertions ...
}
*/
