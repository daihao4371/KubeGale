package cmdb

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	cmdbReq "KubeGale/model/cmdb/request"
	"KubeGale/model/common/response"
	"KubeGale/service"
	"KubeGale/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源的连接
		},
		Subprotocols: []string{"terminal"}, // 添加子协议支持
	}
	jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
)

type CmdbHostsApi struct{}

// CreateCmdbHosts 创建cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) CreateCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.CreatedBy = utils.GetUserID(c)
	err = cmdbHostsService.CreateCmdbHosts(&cmdbHosts)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// AuthenticationCmdbHosts 验证cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) AuthenticationCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.CreatedBy = utils.GetUserID(c)
	err = cmdbHostsService.SSHTestCmdbHosts(&cmdbHosts)
	if err != nil {
		if err.Error() == "auth failed" {
			response.Result(177, nil, "auth failed", c)
			return
		}
		global.KUBEGALE_LOG.Error("验证失败!", zap.Error(err))
		response.FailWithMessage("验证失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("验证成功", c)
}

// DeleteCmdbHosts 删除cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) DeleteCmdbHosts(c *gin.Context) {
	var req cmdbReq.DeleteCmdbHostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	err := cmdbHostsService.DeleteCmdbHosts(fmt.Sprintf("%d", req.ID), userID)
	if err != nil {
		global.KUBEGALE_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCmdbHostsByIds 批量删除cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) DeleteCmdbHostsByIds(c *gin.Context) {
	var req cmdbReq.DeleteCmdbHostsIdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 将 uint 数组转换为 string 数组
	ids := make([]string, len(req.IDs))
	for i, id := range req.IDs {
		ids[i] = fmt.Sprintf("%d", id)
	}

	userID := utils.GetUserID(c)
	err := cmdbHostsService.DeleteCmdbHostsByIds(ids, userID)
	if err != nil {
		global.KUBEGALE_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCmdbHosts 更新cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) UpdateCmdbHosts(c *gin.Context) {
	var cmdbHosts cmdb.CmdbHosts
	err := c.ShouldBindJSON(&cmdbHosts)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	cmdbHosts.UpdatedBy = utils.GetUserID(c)
	err = cmdbHostsService.UpdateCmdbHosts(cmdbHosts)
	if err != nil {
		global.KUBEGALE_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCmdbHosts 用id查询cmdbHosts表
func (cmdbHostsApi *CmdbHostsApi) FindCmdbHosts(c *gin.Context) {
	ID := c.Query("id")
	recmdbHosts, err := cmdbHostsService.GetCmdbHosts(ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(recmdbHosts, c)
}

// GetCmdbHostsList 获取cmdbHosts表列表
func (cmdbHostsApi *CmdbHostsApi) GetCmdbHostsList(c *gin.Context) {
	var pageInfo cmdbReq.CmdbHostsSearch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		global.KUBEGALE_LOG.Error("参数绑定失败!", zap.Error(err))
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	list, total, err := cmdbHostsService.GetCmdbHostsInfoList(pageInfo)
	if err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ImportHosts 根据模板批量创建主机
func (cmdbHostsApi *CmdbHostsApi) ImportHosts(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file") // "file" 是前端上传文件的字段名
	if err != nil {
		response.FailWithMessage("获取文件失败: "+err.Error(), c)
		return
	}
	projectIdStr := c.PostForm("projectId")
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		response.FailWithMessage("无效的 projectId: "+err.Error(), c)
		return
	}

	// 保存上传的文件到临时目录
	dst := "/tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.FailWithMessage("保存文件失败: "+err.Error(), c)
		return
	}

	// 调用服务层处理上传逻辑
	if err := cmdbHostsService.ImportHosts(dst, projectId); err != nil {
		global.KUBEGALE_LOG.Error("导入失败!", zap.Error(err))
		response.FailWithMessage("导入失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("导入成功", c)
}

// WebTerminal 处理Web终端连接
func (cmdbHostsApi *CmdbHostsApi) WebTerminal(c *gin.Context) {
	// 升级HTTP连接为WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.KUBEGALE_LOG.Error("WebSocket升级失败", zap.Error(err))
		return
	}

	// 创建一个通道用于同步关闭
	done := make(chan struct{})
	// 创建一个通道用于处理错误
	errChan := make(chan error, 1)

	// 确保资源清理
	defer func() {
		close(done)
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.Close()
	}()

	// 读取第一条消息，做认证
	_, msg, err := ws.ReadMessage()
	if err != nil {
		global.KUBEGALE_LOG.Error("读取认证消息失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "auth_result", "success": false, "message": "未收到认证消息"})
		return
	}

	var authMsg struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(msg, &authMsg); err != nil {
		global.KUBEGALE_LOG.Error("解析认证消息失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "auth_result", "success": false, "message": "认证消息格式错误"})
		return
	}

	if authMsg.Type != "auth" {
		global.KUBEGALE_LOG.Error("认证消息类型错误", zap.String("type", authMsg.Type))
		ws.WriteJSON(gin.H{"type": "auth_result", "success": false, "message": "认证消息类型错误"})
		return
	}

	if authMsg.Token == "" {
		global.KUBEGALE_LOG.Error("token为空")
		ws.WriteJSON(gin.H{"type": "auth_result", "success": false, "message": "token不能为空"})
		return
	}

	// 校验 token
	j := utils.NewJWT()
	_, err = j.ParseToken(authMsg.Token)
	if err != nil {
		global.KUBEGALE_LOG.Error("token解析失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "auth_result", "success": false, "message": "token无效"})
		return
	}

	// 检查token是否在黑名单中
	if global.KUBEGALE_CONFIG.System.UseMultipoint && jwtService.IsBlacklist(authMsg.Token) {
		global.KUBEGALE_LOG.Error("token在黑名单中")
		ws.WriteJSON(gin.H{"type": "auth_result", "success": false, "message": "token已失效"})
		return
	}

	// 认证通过
	ws.WriteJSON(gin.H{"type": "auth_result", "success": true})

	// 从URL参数获取主机ID
	hostID := c.Query("id")
	if hostID == "" {
		global.KUBEGALE_LOG.Error("主机ID为空")
		ws.WriteJSON(gin.H{"type": "error", "message": "主机ID不能为空"})
		return
	}

	// 将字符串ID转换为uint
	id, err := strconv.ParseUint(hostID, 10, 32)
	if err != nil {
		global.KUBEGALE_LOG.Error("主机ID格式错误", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "error", "message": "无效的主机ID"})
		return
	}

	// 获取主机信息
	var host cmdb.CmdbHosts
	if err := global.KUBEGALE_DB.First(&host, id).Error; err != nil {
		global.KUBEGALE_LOG.Error("获取主机信息失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "error", "message": "主机不存在"})
		return
	}

	// 创建SSH会话
	session, err := cmdbHostsService.CreateSSHSession(&host)
	if err != nil {
		global.KUBEGALE_LOG.Error("创建SSH会话失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "error", "message": fmt.Sprintf("SSH连接失败: %v", err)})
		return
	}
	defer session.Close()

	// 设置终端模式
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // 启用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
		ssh.ECHOCTL:       1,     // 控制字符回显
		ssh.ECHOKE:        1,     // 视觉擦除
		ssh.ECHONL:        1,     // 换行回显
		ssh.ICANON:        1,     // 规范模式
		ssh.ISIG:          1,     // 启用信号
		ssh.IEXTEN:        1,     // 启用扩展
		ssh.IXON:          1,     // 启用输出流控制
		ssh.IXOFF:         1,     // 启用输入流控制
	}

	// 请求伪终端
	if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
		global.KUBEGALE_LOG.Error("请求伪终端失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "error", "message": "请求伪终端失败"})
		return
	}

	// 获取stdin和stdout管道
	stdin, err := session.StdinPipe()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取stdin管道失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "error", "message": "获取stdin管道失败"})
		return
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		global.KUBEGALE_LOG.Error("获取stdout管道失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "error", "message": "获取stdout管道失败"})
		return
	}

	// 启动shell
	if err := session.Shell(); err != nil {
		global.KUBEGALE_LOG.Error("启动shell失败", zap.Error(err))
		ws.WriteJSON(gin.H{"type": "error", "message": "启动shell失败"})
		return
	}

	// 发送欢迎消息
	welcomeMsg := fmt.Sprintf("\r\n欢迎连接到 %s (%s)\r\n", host.Name, host.ServerHost)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(welcomeMsg)); err != nil {
		global.KUBEGALE_LOG.Error("发送欢迎消息失败", zap.Error(err))
		return
	}

	// 处理WebSocket和SSH会话之间的数据转发
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("WebSocket读取协程panic: %v", r)
			}
		}()

		for {
			select {
			case <-done:
				return
			default:
				_, message, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
						global.KUBEGALE_LOG.Error("读取WebSocket消息失败", zap.Error(err))
						errChan <- err
					}
					return
				}
				if _, err := stdin.Write(message); err != nil {
					global.KUBEGALE_LOG.Error("写入SSH stdin失败", zap.Error(err))
					errChan <- err
					return
				}
			}
		}
	}()

	// 从SSH会话读取数据并发送到WebSocket
	buffer := make([]byte, 1024)
	for {
		select {
		case <-done:
			return
		case err := <-errChan:
			global.KUBEGALE_LOG.Error("发生错误", zap.Error(err))
			ws.WriteJSON(gin.H{"type": "error", "message": fmt.Sprintf("连接错误: %v", err)})
			return
		default:
			n, err := stdout.Read(buffer)
			if err != nil {
				if err != io.EOF {
					global.KUBEGALE_LOG.Error("读取SSH stdout失败", zap.Error(err))
					ws.WriteJSON(gin.H{"type": "error", "message": fmt.Sprintf("读取SSH输出失败: %v", err)})
				}
				return
			}
			if err := ws.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					global.KUBEGALE_LOG.Error("写入WebSocket消息失败", zap.Error(err))
				}
				return
			}
		}
	}
}
