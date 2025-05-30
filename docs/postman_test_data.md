# KubeGale API 测试数据

## 通知模块 (Notification)

### 1. 创建飞书通知
**POST** `/notification/createFeiShu`

```json
{
    "name": "测试飞书通知",
    "type": "feishu",
    "enabled": true,
    "webhook_url": "https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
    "description": "这是一个测试飞书通知配置",
    "tags": ["告警", "监控"],
    "notify_events": ["alert", "warning"],
    "receivers": ["张三", "李四"],
    "send_daily_stats": true
}
```

### 2. 更新飞书通知
**PUT** `/notification/updateFeiShu`

```json
{
    "id": 1,
    "name": "更新后的飞书通知",
    "notification_policy": "alert,warning",
    "robot_url": "https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
    "send_daily_stats": true,
    "card_content": {
        "alert_level": "warning",
        "alert_name": "测试告警",
        "notification_policy": "immediate",
        "alert_content": "这是一个测试告警内容",
        "alert_time": "2024-03-20T10:00:00Z",
        "notified_users": ["张三", "李四"],
        "last_similar_alert": "2024-03-19T10:00:00Z",
        "alert_handler": "系统管理员",
        "claim_alert": false,
        "resolve_alert": false,
        "mute_alert": false,
        "unresolved_alert": true
    }
}
```

### 3. 创建钉钉通知
**POST** `/notification/createDingTalk`

```json
{
    "name": "测试钉钉通知",
    "notify_events": ["alert", "warning"],
    "send_daily_stats": true,
    "webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=xxx",
    "secret": "your_secret"
}
```

### 4. 更新钉钉通知
**PUT** `/notification/updateDingTalk`

```json
{
    "id": 1,
    "name": "更新后的钉钉通知",
    "notification_policy": "alert,warning",
    "send_daily_stats": true,
    "webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=xxx",
    "secret": "your_secret",
    "card_content": {
        "alert_level": "warning",
        "alert_name": "测试告警",
        "notification_policy": "immediate",
        "alert_content": "这是一个测试告警内容",
        "alert_time": "2024-03-20T10:00:00Z",
        "notified_users": ["张三", "李四"],
        "last_similar_alert": "2024-03-19T10:00:00Z",
        "alert_handler": "系统管理员",
        "claim_alert": false,
        "resolve_alert": false,
        "mute_alert": false,
        "unresolved_alert": true
    }
}
```

### 5. 创建卡片内容
**POST** `/notification/createCardContent`

```json
{
    "notification_id": 1,
    "alert_level": "critical",
    "alert_name": "测试告警",
    "notification_policy": "immediate",
    "alert_content": "这是一个测试告警内容",
    "alert_time": "2024-03-20T10:00:00Z",
    "notified_users": [
        {
            "name": "张三",
            "email": "zhangsan@example.com"
        },
        {
            "name": "李四",
            "email": "lisi@example.com"
        }
    ],
    "last_similar_alert": "2024-03-19T10:00:00Z",
    "alert_handler": "系统管理员",
    "claim_alert": false,
    "resolve_alert": false,
    "mute_alert": false,
    "unresolved_alert": true
}
```

### 6. 更新卡片内容
**PUT** `/notification/updateCardContent`

```json
{
    "id": 1,
    "notification_id": 1,
    "alert_level": "warning",
    "alert_name": "更新后的测试告警",
    "notification_policy": "delayed",
    "alert_content": "这是更新后的测试告警内容",
    "alert_time": "2024-03-20T11:00:00Z",
    "notified_users": [
        {
            "name": "张三",
            "email": "zhangsan@example.com"
        }
    ],
    "last_similar_alert": "2024-03-19T11:00:00Z",
    "alert_handler": "运维团队",
    "claim_alert": true,
    "resolve_alert": false,
    "mute_alert": false,
    "unresolved_alert": false
}
```

### 7. 测试通知发送
**POST** `/notification/testNotification`

```json
{
    "id": 1,
    "type": "feishu",
    "test_message": "这是一条测试消息"
}
```

### 8. 获取通知列表
**POST** `/notification/getNotificationList`

```json
{
    "page": 1,
    "page_size": 10,
    "name": "",
    "type": "",
    "status": null
}
```

## 查询参数说明

### 1. 获取通知配置
**GET** `/notification/getNotificationById?id=1&type=feishu`

### 2. 获取卡片内容
**GET** `/notification/getCardContent?notification_id=1`

### 3. 删除通知配置
**DELETE** `/notification/deleteNotification?id=1&type=feishu`

## 注意事项

1. 所有时间字段使用 ISO 8601 格式
2. 布尔值使用 true/false
3. 数字类型使用整数
4. 字符串类型需要根据实际情况填写
5. 测试时请替换示例中的 URL 和密钥为实际值 