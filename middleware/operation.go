package middleware

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"KubeGale/service"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userId int
		var userName string
		var userRealName string
		
		// 从上下文中获取当前用户ID和用户信息
		userIdInterface, exists := c.Get("user_id")
		if exists {
			// 如果存在用户ID，则设置到操作记录中
			if userIdInt, ok := userIdInterface.(int); ok {
				userId = userIdInt
				
				// 尝试获取用户名和真实姓名
				if userNameInterface, ok := c.Get("user_name"); ok {
					if name, ok := userNameInterface.(string); ok {
						userName = name
					}
				}
				
				if userRealNameInterface, ok := c.Get("real_name"); ok {
					if realName, ok := userRealNameInterface.(string); ok {
						userRealName = realName
					}
				}
				
				// 如果上下文中没有用户名和真实姓名，尝试从数据库获取
				if userName == "" || userRealName == "" {
					var user system.User
					if err := global.KUBEGALE_DB.Where("id = ?", userId).First(&user).Error; err == nil {
						userName = user.Username
						userRealName = user.RealName
					}
				}
			} else if userIdUint, ok := userIdInterface.(uint); ok {
				userId = int(userIdUint)
				// 同样尝试获取用户信息...
			} else if userIdFloat, ok := userIdInterface.(float64); ok {
				userId = int(userIdFloat)
				// 同样尝试获取用户信息...
			} else {
				// 尝试其他可能的类型转换
				global.KUBEGALE_LOG.Warn("无法将用户ID转换为整数类型")
			}
		}
		
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.KUBEGALE_LOG.Error("read body from request error:", zap.Error(err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		}

		record := system.SysOperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   "",
			UserID: userId, // 设置用户ID
		}

		// 上传文件时候 中间件日志进行裁断操作
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-system") {
			record.Body = "[文件]"
		} else {
			if len(body) > bufferSize {
				record.Body = "[超出记录长度]"
			} else {
				record.Body = string(body)
			}
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Since(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		record.Resp = writer.body.String()

		if strings.Contains(c.Writer.Header().Get("Pragma"), "public") ||
			strings.Contains(c.Writer.Header().Get("Expires"), "0") ||
			strings.Contains(c.Writer.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/force-download") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/octet-stream") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/download") ||
			strings.Contains(c.Writer.Header().Get("Content-Disposition"), "attachment") ||
			strings.Contains(c.Writer.Header().Get("Content-Transfer-Encoding"), "binary") {
			if len(record.Resp) > bufferSize {
				// 截断
				record.Body = "超出记录长度"
			}
		}

		if err := operationRecordService.CreateSysOperationRecord(record); err != nil {
			global.KUBEGALE_LOG.Error("create operation record error:", zap.Error(err))
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
