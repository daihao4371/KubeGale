package im

import (
	"KubeGale/global"
	"KubeGale/model/im"
	"fmt"
	"time"
)

type CronService struct{}

var CronServiceApp = new(CronService)

// @function: SendDailyStats
// @description: 发送每日统计信息
func (cronService *CronService) SendDailyStats() error {
	// 获取所有启用了每日统计的通知配置
	var dingTalkConfigs []im.DingTalkConfig
	var feiShuConfigs []im.FeiShuConfig
	
	// 查询钉钉配置
	if err := global.KUBEGALE_DB.Where("send_daily_stats = ?", true).Find(&dingTalkConfigs).Error; err != nil {
		return err
	}
	
	// 查询飞书配置
	if err := global.KUBEGALE_DB.Where("send_daily_stats = ?", true).Find(&feiShuConfigs).Error; err != nil {
		return err
	}
	
	// 生成统计信息
	stats := generateDailyStats()
	
	// 发送钉钉通知
	for _, config := range dingTalkConfigs {
		var cardContent im.CardContentConfig
		global.KUBEGALE_DB.Where("notification_id = ?", config.ID).First(&cardContent)
		
		err := MessageServiceApp.SendDingTalkMessage(config, cardContent, stats)
		if err != nil {
			global.KUBEGALE_LOG.Error(fmt.Sprintf("发送钉钉每日统计失败: %s", err.Error()))
		}
	}
	
	// 发送飞书通知
	for _, config := range feiShuConfigs {
		var cardContent im.CardContentConfig
		global.KUBEGALE_DB.Where("notification_id = ?", config.ID).First(&cardContent)
		
		err := MessageServiceApp.SendFeiShuMessage(config, cardContent, stats)
		if err != nil {
			global.KUBEGALE_LOG.Error(fmt.Sprintf("发送飞书每日统计失败: %s", err.Error()))
		}
	}
	
	return nil
}

// @function: generateDailyStats
// @description: 生成每日统计信息
func generateDailyStats() string {
	// 这里应该实现实际的统计逻辑，例如查询数据库获取今日告警数量、处理情况等
	// 这里只是一个示例
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	
	// 查询昨日告警数量（示例）
	var alertCount int64
	global.KUBEGALE_DB.Model(&im.CardContentConfig{}).Where("alert_time BETWEEN ? AND ?", yesterday.Format("2006-01-02"), now.Format("2006-01-02")).Count(&alertCount)
	
	// 查询已解决告警数量（示例）
	var resolvedCount int64
	global.KUBEGALE_DB.Model(&im.CardContentConfig{}).Where("alert_time BETWEEN ? AND ? AND resolve_alert = ?", yesterday.Format("2006-01-02"), now.Format("2006-01-02"), true).Count(&resolvedCount)
	
	// 生成统计信息
	stats := fmt.Sprintf("# 每日告警统计 (%s)\n\n", yesterday.Format("2006-01-02"))
	stats += fmt.Sprintf("- 总告警数量: %d\n", alertCount)
	stats += fmt.Sprintf("- 已解决告警: %d\n", resolvedCount)
	stats += fmt.Sprintf("- 未解决告警: %d\n", alertCount-resolvedCount)
	stats += fmt.Sprintf("- 解决率: %.2f%%\n", float64(resolvedCount)/float64(alertCount)*100)
	
	return stats
}