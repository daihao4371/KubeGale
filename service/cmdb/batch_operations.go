package cmdb

import (
	"KubeGale/global"
	"KubeGale/model/cmdb/request"
	"KubeGale/utils/cmdb"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type BatchOperationsService struct{}

// 批量执行命令并返回 ExecuteResponse
func (b *BatchOperationsService) CreateBatchOperations(req request.ExecuteRequest) (*request.ExecuteResponse, error) {
	keyPath, err := cmdb.GetDefaultPrivateKeyPath() // 替换为实际的私钥路径
	if err != nil {
		return nil, err
	}

	// 检查主机和用户列表长度是否一致
	if len(req.Hosts) != len(req.Users) {
		return nil, fmt.Errorf("hosts and users list must have the same length")
	}

	// 如果只提供了一个端口，则所有主机使用相同端口
	if len(req.Ports) == 1 {
		port := req.Ports[0]
		req.Ports = make([]int, len(req.Hosts))
		for i := range req.Hosts {
			req.Ports[i] = port
		}
	} else if len(req.Ports) != len(req.Hosts) {
		return nil, fmt.Errorf("ports list must have length 1 or match hosts length")
	}

	// 构建完整的命令
	fullCommand, err := buildCommand(req)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make(chan request.HostExecResult, len(req.Hosts)) // 缓冲通道来存储结果

	// 记录所有主机
	var allHosts []string
	for i := range req.Hosts {
		allHosts = append(allHosts, fmt.Sprintf("%s:%d", req.Hosts[i], req.Ports[i]))
	}

	// 并行处理每对主机和用户
	for i := range req.Hosts {
		host := fmt.Sprintf("%s:%d", req.Hosts[i], req.Ports[i])
		user := req.Users[i]

		wg.Add(1)
		go func(host, user string) {
			defer wg.Done()
			output, err := executeRemoteCommandWithKey(host, user, keyPath, fullCommand)
			if err != nil {
				results <- request.HostExecResult{
					Host:   host,
					Error:  err.Error(),
					Output: output, // 即使有错误也包含输出
				}
			} else {
				results <- request.HostExecResult{
					Host:   host,
					Output: output,
				}
			}
		}(host, user)
	}

	// 等待所有 Goroutines 执行完毕
	wg.Wait()
	close(results)

	// 分别收集成功和失败的主机以及执行结果
	var successHosts []string
	var failureHosts []string
	var execResults []request.HostExecResult
	for res := range results {
		execResults = append(execResults, res)
		if res.Error != "" {
			failureHosts = append(failureHosts, res.Host)
		} else {
			successHosts = append(successHosts, res.Host)
		}
	}

	// 如果有任何失败的主机，将状态设置为失败
	responseStatus := "success"
	if len(failureHosts) > 0 {
		responseStatus = "failed"
	}

	// 构建响应体
	response := &request.ExecuteResponse{
		Status:        responseStatus,
		AllHosts:      allHosts,     // 包含所有主机
		SuccessHosts:  successHosts, // 成功的主机
		FailureHosts:  failureHosts, // 失败的主机
		ExecutionLogs: execResults,  // 执行结果日志
	}

	// 将执行日志序列化为 JSON 字符串
	executionLogsJSON, err := json.Marshal(execResults)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal execution logs: %v", err)
	}

	// 保存执行记录到数据库
	executionLog := request.CommandExecutionLog{
		Command:       fullCommand,
		UserId:        req.UserId,
		AllHosts:      strings.Join(allHosts, ","),     // 将所有主机列表转换为逗号分隔的字符串
		SuccessHosts:  strings.Join(successHosts, ","), // 成功主机列表
		FailureHosts:  strings.Join(failureHosts, ","), // 失败主机列表
		ExecutionLogs: string(executionLogsJSON),       // 执行日志的 JSON
		Status:        responseStatus,
	}

	// 使用 gorm 保存到数据库
	if err := global.KUBEGALE_DB.Create(&executionLog).Error; err != nil {
		return nil, fmt.Errorf("failed to save execution log: %v", err)
	}

	return response, nil
}

func (b *BatchOperationsService) GetUserRecentExecutionRecords(userId uint) ([]request.CommandExecutionLog, error) {
	var logs []request.CommandExecutionLog

	// 查询数据库，获取当前用户的最近10条记录
	err := global.KUBEGALE_DB.Model(&request.CommandExecutionLog{}).
		Where("user_id = ?", userId). // 根据用户ID查询
		Order("created_at desc").     // 按创建时间倒序排列
		Find(&logs).                  // 查询记录
		Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// 根据请求构建完整命令
func buildCommand(req request.ExecuteRequest) (string, error) {
	switch req.Language {
	case "shell":
		return strings.Join(req.Commands, " && "), nil
	case "python":
		pythonCommands := make([]string, len(req.Commands))
		for i, cmd := range req.Commands {
			pythonCommands[i] = fmt.Sprintf("python3 -c '%s'", cmd)
		}
		return strings.Join(pythonCommands, " && "), nil
	default:
		return "", fmt.Errorf("unsupported language: %s", req.Language)
	}
}

// SSH连接与执行命令，使用私钥认证
func executeRemoteCommandWithKey(host, user, keyPath, command string) (string, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("unable to read private key: %v", err)
	}

	// 解析私钥
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("unable to parse private key: %v", err)
	}

	// SSH 客户端配置
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 建立连接
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// 设置命令执行超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建管道来捕获输出
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	// 使用 bash -c 执行命令
	fullCommand := fmt.Sprintf("bash -c '%s'", command)
	if err := session.Start(fullCommand); err != nil {
		return "", fmt.Errorf("failed to start command: %v", err)
	}

	// 读取输出
	var outputBuf, errorBuf bytes.Buffer
	outputDone := make(chan struct{})
	errorDone := make(chan struct{})

	go func() {
		io.Copy(&outputBuf, stdout)
		close(outputDone)
	}()

	go func() {
		io.Copy(&errorBuf, stderr)
		close(errorDone)
	}()

	// 等待命令完成或超时
	done := make(chan error)
	go func() {
		done <- session.Wait()
	}()

	select {
	case err := <-done:
		<-outputDone
		<-errorDone

		// 获取标准输出和错误输出
		output := outputBuf.String()
		errorOutput := errorBuf.String()

		// 构建完整的输出结果
		var result strings.Builder

		// 如果有标准输出，添加到结果中
		if output != "" {
			result.WriteString(output)
			if !strings.HasSuffix(output, "\n") {
				result.WriteString("\n")
			}
		}

		// 如果有错误输出，添加到结果中
		if errorOutput != "" {
			if output != "" {
				result.WriteString("\n")
			}
			result.WriteString(errorOutput)
			if !strings.HasSuffix(errorOutput, "\n") {
				result.WriteString("\n")
			}
		}

		// 如果命令执行出错或stderr有输出，返回错误
		if err != nil || errorOutput != "" {
			if result.Len() == 0 {
				return "", fmt.Errorf("command failed: %v", err)
			}
			return result.String(), fmt.Errorf("command failed: %v", err)
		}

		// 如果没有任何输出，返回空字符串
		if result.Len() == 0 {
			return "", nil
		}

		return result.String(), nil

	case <-ctx.Done():
		session.Close()
		return "", fmt.Errorf("command execution timed out")
	}
}
