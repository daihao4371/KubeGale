package im

import (
	"KubeGale/global"
	"KubeGale/model/im/response" // Assuming response DTOs are here
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// DingTalkMessageRequest defines the structure for a DingTalk text message.
type DingTalkMessageRequest struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		IsAtAll bool `json:"isAtAll"`
		// Add AtMobiles []string `json:"atMobiles,omitempty"` if needed
	} `json:"at"`
}

// SendDingTalkMessage sends a message to a DingTalk webhook.
// For now, it sends the 'message' parameter as a simple text message.
// 'cardContent' is available for future enhancements (e.g., richer messages).
func SendDingTalkMessage(config response.NotificationDetailConfig, cardContent response.CardContentDetail, message string) error {
	if config.RobotURL == "" {
		return fmt.Errorf("dingtalk webhook URL is empty")
	}

	webhookURL := config.RobotURL
	var err error

	// If a secret is provided, sign the request
	if config.Secret != nil && *config.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		stringToSign := fmt.Sprintf("%s\n%s", timestamp, *config.Secret)
		
		mac := hmac.New(sha256.New, []byte(*config.Secret))
		mac.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		// Add timestamp and sign to webhook URL
        parsedURL, parseErr := url.Parse(config.RobotURL)
        if parseErr != nil {
            return fmt.Errorf("error parsing dingtalk webhook_url: %w", parseErr)
        }
        query := parsedURL.Query()
        query.Set("timestamp", timestamp)
        query.Set("sign", signature)
        parsedURL.RawQuery = query.Encode()
        webhookURL = parsedURL.String()
	}

	// For now, use the direct 'message' for simplicity, as in TestNotification
	// In a real scenario, you might format content from cardContent:
	// content := fmt.Sprintf("Alert: %s\nDetails: %s\nPolicy: %s", cardContent.AlertName, cardContent.AlertContent, cardContent.NotificationPolicy)
	content := message // Use the direct message for TestNotification

	msgPayload := DingTalkMessageRequest{
		MsgType: "text",
	}
	msgPayload.Text.Content = content
	msgPayload.At.IsAtAll = false // Configure as needed, e.g., based on cardContent.NotifiedUsers

	payloadBytes, err := json.Marshal(msgPayload)
	if err != nil {
		global.KUBEGALE_LOG.Error("Failed to marshal DingTalk message payload", zap.Error(err))
		return fmt.Errorf("failed to marshal dingtalk payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		global.KUBEGALE_LOG.Error("Failed to create DingTalk HTTP request", zap.Error(err))
		return fmt.Errorf("failed to create dingtalk request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.KUBEGALE_LOG.Error("Failed to send message to DingTalk", zap.String("url", webhookURL), zap.Error(err))
		return fmt.Errorf("failed to send to dingtalk: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		global.KUBEGALE_LOG.Error("DingTalk API request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_body", string(body)))
		return fmt.Errorf("dingtalk api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// DingTalk success response usually contains {"errcode":0,"errmsg":"ok"}
	var dingTalkResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &dingTalkResp); err == nil {
		if dingTalkResp.ErrCode != 0 {
			global.KUBEGALE_LOG.Error("DingTalk API returned an error",
				zap.Int("errcode", dingTalkResp.ErrCode),
				zap.String("errmsg", dingTalkResp.ErrMsg))
			return fmt.Errorf("dingtalk api returned error: code %d, msg: %s", dingTalkResp.ErrCode, dingTalkResp.ErrMsg)
		}
	} else {
        global.KUBEGALE_LOG.Warn("Could not parse DingTalk response body, but status was 200", zap.String("body", string(body)))
    }


	global.KUBEGALE_LOG.Info("Successfully sent message to DingTalk", zap.String("url", config.RobotURL))
	return nil
}
