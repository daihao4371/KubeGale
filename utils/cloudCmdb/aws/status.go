package aws

var (
	EC2Status = map[string]string{
		"pending":                "创建中",
		"running":                "运行中",
		"shutting-down":          "关闭中",
		"terminated":             "已终止",
		"stopping":               "停止中",
		"stopped":                "已停止",
		"rebooting":              "重启中",
		"pending-instance-stop":  "实例停止中",
		"pending-instance-start": "实例启动中",
	}

	RDSStatus = map[string]string{
		"creating":                     "创建中",
		"available":                    "可用",
		"deleting":                     "删除中",
		"failed":                       "失败",
		"modifying":                    "修改中",
		"rebooting":                    "重启中",
		"renaming":                     "重命名中",
		"resetting-master-credentials": "重置主凭证中",
		"storage-full":                 "存储已满",
		"upgrading":                    "升级中",
		"backing-up":                   "备份中",
		"restoring":                    "恢复中",
	}

	LoadBalancerStatus = map[string]string{
		"active":          "运行中",
		"inactive":        "已停止",
		"failed":          "失败",
		"provisioning":    "配置中",
		"active_impaired": "运行中但受损",
	}
)
