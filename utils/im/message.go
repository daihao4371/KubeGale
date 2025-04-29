package im

import (
	"KubeGale/model/im"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MessageService struct{}

var MessageServiceApp = new(MessageService)

// 钉钉消息结构
type DingTalkMessage struct {
	MsgType  string                  `json:"msgtype"`
	Markdown DingTalkMarkdownContent `json:"markdown,omitempty"`
}

// 钉钉 Markdown 消息内容
type DingTalkMarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// 飞书消息结构
type FeiShuMessage struct {
	MsgType string               `json:"msg_type"`
	Content FeiShuMessageContent `json:"content"`
}

// 飞书消息内容
type FeiShuMessageContent struct {
	Text string `json:"text"`
}

// 飞书卡片消息
type FeiShuCardMessage struct {
	MsgType string            `json:"msg_type"`
	Card    FeiShuCardContent `json:"card"`
}

// 飞书卡片内容
type FeiShuCardContent struct {
	Header   FeiShuCardHeader    `json:"header"`
	Elements []FeiShuCardElement `json:"elements"`
}

// 飞书卡片标题
type FeiShuCardHeader struct {
	Title    FeiShuCardText `json:"title"`
	Template string         `json:"template"` // 卡片标题背景色，取值范围：blue、wathet、turquoise、green、yellow、orange、red、carmine、violet、purple、indigo、grey
}

// 飞书卡片文本
type FeiShuCardText struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

// 飞书卡片元素
type FeiShuCardElement struct {
	Tag  string         `json:"tag"`
	Text FeiShuCardText `json:"text,omitempty"`
}

// @function: SendDingTalkMessage
// @description: 发送钉钉消息
func (messageService *MessageService) SendDingTalkMessage(dingTalkConfig im.DingTalkConfig, cardContent im.CardContentConfig, message string) error {
	// 构建签名
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	stringToSign := timestamp + "\n" + dingTalkConfig.SignatureKey

	mac := hmac.New(sha256.New, []byte(dingTalkConfig.SignatureKey))
	mac.Write([]byte(stringToSign))
	sign := url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	// 构建请求URL
	requestURL := fmt.Sprintf("%s&timestamp=%s&sign=%s", dingTalkConfig.RobotURL, timestamp, sign)

	// 构建消息内容
	title := "系统通知"
	content := message

	// 如果有卡片内容配置，则使用卡片内容
	if cardContent.ID != 0 {
		title = cardContent.AlertName
		content = fmt.Sprintf("### %s\n\n**告警等级**: %s\n\n**告警内容**: %s\n\n**告警时间**: %s\n\n**通知人**: %s\n\n**上次相似告警**: %s\n\n**告警处理人**: %s",
			cardContent.AlertName,
			cardContent.AlertLevel,
			cardContent.AlertContent,
			cardContent.AlertTime.Format("2006-01-02 15:04:05"),
			cardContent.NotifiedUsers,
			cardContent.LastSimilarAlert,
			cardContent.AlertHandler)
	}

	// 构建消息体
	msg := DingTalkMessage{
		MsgType: "markdown",
		Markdown: DingTalkMarkdownContent{
			Title: title,
			Text:  content,
		},
	}

	// 发送请求
	return sendRequest(requestURL, msg)
}

// @function: SendFeiShuMessage
// @description: 发送飞书消息
func (messageService *MessageService) SendFeiShuMessage(feiShuConfig im.FeiShuConfig, cardContent im.CardContentConfig, message string) error {
	// 构建请求URL
	requestURL := feiShuConfig.RobotURL

	// 如果有卡片内容配置，则使用卡片消息
	if cardContent.ID != 0 {
		// 根据告警等级设置卡片颜色
		template := "blue"
		switch cardContent.AlertLevel {
		case "Critical":
			template = "red"
		case "Warning":
			template = "orange"
		case "Info":
			template = "blue"
		}

		// 构建卡片消息
		cardMsg := FeiShuCardMessage{
			MsgType: "interactive",
			Card: FeiShuCardContent{
				Header: FeiShuCardHeader{
					Title: FeiShuCardText{
						Tag:     "plain_text",
						Content: cardContent.AlertName,
					},
					Template: template,
				},
				Elements: []FeiShuCardElement{
					{
						Tag: "div",
						Text: FeiShuCardText{
							Tag:     "lark_md",
							Content: fmt.Sprintf("**告警等级**: %s", cardContent.AlertLevel),
						},
					},
					{
						Tag: "div",
						Text: FeiShuCardText{
							Tag:     "lark_md",
							Content: fmt.Sprintf("**告警内容**: %s", cardContent.AlertContent),
						},
					},
					{
						Tag: "div",
						Text: FeiShuCardText{
							Tag:     "lark_md",
							Content: fmt.Sprintf("**告警时间**: %s", cardContent.AlertTime.Format("2006-01-02 15:04:05")),
						},
					},
					{
						Tag: "div",
						Text: FeiShuCardText{
							Tag:     "lark_md",
							Content: fmt.Sprintf("**通知人**: %s", cardContent.NotifiedUsers),
						},
					},
					{
						Tag: "div",
						Text: FeiShuCardText{
							Tag:     "lark_md",
							Content: fmt.Sprintf("**上次相似告警**: %s", cardContent.LastSimilarAlert),
						},
					},
					{
						Tag: "div",
						Text: FeiShuCardText{
							Tag:     "lark_md",
							Content: fmt.Sprintf("**告警处理人**: %s", cardContent.AlertHandler),
						},
					},
				},
			},
		}

		// 发送请求
		return sendRequest(requestURL, cardMsg)
	}

	// 构建普通文本消息
	msg := FeiShuMessage{
		MsgType: "text",
		Content: FeiShuMessageContent{
			Text: message,
		},
	}

	// 发送请求
	return sendRequest(requestURL, msg)
}

// @function: sendRequest
// @description: 发送HTTP请求
func sendRequest(url string, data interface{}) error {
	// 将消息转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 发送POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
