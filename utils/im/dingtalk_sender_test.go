package im

import (
	"KubeGale/model/im/response"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1" // Example HTTP mocking library
)

// Mock a logger for tests if needed, or ensure global.KUBEGALE_LOG is test-safe
func init() {
	// Basic zap logger for tests if not already initialized
	// global.KUBEGALE_LOG = zap.NewNop()
	// Or initialize your actual logger in a test-friendly way
}

func TestSendDingTalkMessage_Success(t *testing.T) {
	defer gock.Off() // Clean up all gocks after this test

	config := response.NotificationDetailConfig{
		RobotURL: "http://dingtalk.example.com/webhook",
	}
	message := "Hello DingTalk"

	gock.New("http://dingtalk.example.com").
		Post("/webhook").
		JSON(map[string]interface{}{ // Use map for flexibility with JSON matching
			"msgtype": "text",
			"text": map[string]string{
				"content": message,
			},
			"at": map[string]bool{
				"isAtAll": false,
			},
		}).
		Reply(200).
		JSON(map[string]interface{}{"errcode": 0, "errmsg": "ok"})

	err := SendDingTalkMessage(config, response.CardContentDetail{}, message)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone()) // Verify all mocks were matched
}

func TestSendDingTalkMessage_WithSecret(t *testing.T) {
	defer gock.Off()

	secret := "sec12345"
	config := response.NotificationDetailConfig{
		RobotURL: "http://dingtalk.example.com/webhook",
		Secret:   &secret,
	}
	message := "Signed Message"

	// Gock doesn't easily support matching query params that are dynamic (like timestamp & sign)
	// So, we match the path and method, and then can inspect the request in the handler if needed
	// Or use a more flexible matcher like MatchFunc
	gock.New("http://dingtalk.example.com").
		Post("/webhook").
		// MatchFunc to verify signature if possible, or just check that it's called
		AddMatcher(func(req *http.Request, ereq *gock.Request) (bool, error) {
			// Check for timestamp and sign query parameters
			q := req.URL.Query()
			assert.NotEmpty(t, q.Get("timestamp"))
			assert.NotEmpty(t, q.Get("sign"))
			// More detailed signature verification is complex here,
			// as it requires knowing the exact timestamp generated.
			// Trust the signing logic if the params are present.
			return true, nil
		}).
		Reply(200).
		JSON(map[string]interface{}{"errcode": 0, "errmsg": "ok"})

	err := SendDingTalkMessage(config, response.CardContentDetail{}, message)
	assert.NoError(t, err)
	assert.True(t, gock.IsDone())
}

func TestSendDingTalkMessage_HttpError(t *testing.T) {
	defer gock.Off()
	config := response.NotificationDetailConfig{RobotURL: "http://dingtalk.example.com/webhook"}

	gock.New("http://dingtalk.example.com").
		Post("/webhook").
		Reply(500).
		BodyString("Internal Server Error")

	err := SendDingTalkMessage(config, response.CardContentDetail{}, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dingtalk api error: status 500")
}

func TestSendDingTalkMessage_DingTalkError(t *testing.T) {
	defer gock.Off()
	config := response.NotificationDetailConfig{RobotURL: "http://dingtalk.example.com/webhook"}

	gock.New("http://dingtalk.example.com").
		Post("/webhook").
		Reply(200).
		JSON(map[string]interface{}{"errcode": 310000, "errmsg": "sign not match"})

	err := SendDingTalkMessage(config, response.CardContentDetail{}, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dingtalk api returned error: code 310000")
}

// TODO: Add test for webhook URL parsing error
// TODO: Add test for JSON marshal error for payload (harder to trigger unless struct is changed)
