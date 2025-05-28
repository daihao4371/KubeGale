package cmdb

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	cmdbReq "KubeGale/model/cmdb/request"
	utils "KubeGale/utils/cmdb"

	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

type CmdbHostsService struct{}

var CmdbHostsServiceApp = new(CmdbHostsService)

// SSHTestCmdbHosts 测试本地是否可以免密连接远程主机
func (cmdbHostsService *CmdbHostsService) SSHTestCmdbHosts(req *cmdb.CmdbHosts) (err error) {
	port := strconv.Itoa(*req.Port)
	privateKeyPath, err := utils.GetDefaultPrivateKeyPath()
	if err != nil {
		return err
	}
	// 检查私钥文件是否存在
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		if err := utils.GenerateSSHKey(privateKeyPath); err != nil {
			return err
		}
	}
	canLogin, err := utils.CanSSHWithoutPassword(req.ServerHost, port, req.Username, privateKeyPath)
	if err != nil {
		return err
	}
	if canLogin {
		// 创建 SSH 客户端
		client, err := utils.CreateSSHClient(req.ServerHost, port, req.Username, privateKeyPath)
		if err != nil {
			return fmt.Errorf("创建 SSH 客户端失败: %v", err)
		}
		defer client.Close()
		// 获取主机信息
		if err := utils.PopulateHostInfo(client, req); err != nil {
			return fmt.Errorf("获取主机信息失败: %v", err)
		}
		req.Status = "已验证"
		// 存储到数据库
		if err := global.KUBEGALE_DB.Create(req).Error; err != nil {
			return fmt.Errorf("存储到数据库失败: %v", err)
		}
		return nil
	}
	if strings.Contains(err.Error(), "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain") {
		return fmt.Errorf("auth failed")
	}
	return err

}

// CreateCmdbHosts 创建CMDB主机
func (cmdbHostsService *CmdbHostsService) CreateCmdbHosts(req *cmdb.CmdbHosts) (err error) {
	port := strconv.Itoa(*req.Port)

	// 验证 SSH 连接
	client, err := utils.ValidateSSHConnection(req.ServerHost, port, req.Username, req.Password)
	if err != nil {
		return fmt.Errorf("SSH 验证失败: %v", err)
	}
	defer client.Close()

	// 获取主机信息
	if err := utils.PopulateHostInfo(client, req); err != nil {
		return fmt.Errorf("获取主机信息失败: %v", err)
	}
	// 连接成功后，设置免密登录
	if err := utils.EnablePasswordlessSSH(client, req.Username); err != nil {
		return fmt.Errorf("设置免密登录失败: %v", err)
	}
	req.Status = "已验证"

	// 连接成功，存储到数据库
	if err := global.KUBEGALE_DB.Create(req).Error; err != nil {
		return fmt.Errorf("存储到数据库失败: %v", err)
	}
	return nil
}

// DeleteCmdbHosts 删除主机
func (cmdbHostsService *CmdbHostsService) DeleteCmdbHosts(ID string, userID uint) (err error) {
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 先检查记录是否存在
		var host cmdb.CmdbHosts
		if err := tx.Where("id = ?", ID).First(&host).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("主机不存在")
			}
			return err
		}

		// 执行物理删除
		if err := tx.Unscoped().Delete(&cmdb.CmdbHosts{}, "id = ?", ID).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteCmdbHostsByIds 批量删除cmdbHosts表记录
func (cmdbHostsService *CmdbHostsService) DeleteCmdbHostsByIds(IDs []string, deleted_by uint) (err error) {
	return global.KUBEGALE_DB.Transaction(func(tx *gorm.DB) error {
		// 检查所有记录是否存在
		var count int64
		if err := tx.Model(&cmdb.CmdbHosts{}).Where("id IN ?", IDs).Count(&count).Error; err != nil {
			return err
		}
		if int(count) != len(IDs) {
			return fmt.Errorf("部分主机不存在")
		}

		// 执行批量物理删除
		if err := tx.Unscoped().Delete(&cmdb.CmdbHosts{}, "id IN ?", IDs).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateCmdbHosts 更新主机信息
func (cmdbHostsService *CmdbHostsService) UpdateCmdbHosts(cmdbHosts cmdb.CmdbHosts) (err error) {
	err = global.KUBEGALE_DB.Model(&cmdb.CmdbHosts{}).Where("id = ?", cmdbHosts.ID).Updates(&cmdbHosts).Error
	return err
}

// GetCmdbHosts 根据ID获取cmdbHosts表记录
func (cmdbHostsService *CmdbHostsService) GetCmdbHosts(ID string) (cmdbHosts cmdb.CmdbHosts, err error) {
	err = global.KUBEGALE_DB.Where("id = ?", ID).First(&cmdbHosts).Error
	return
}

// ImportHosts 导入主机信息
func (cmdbHostsService *CmdbHostsService) ImportHosts(filePath string, projectId int) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("无法读取行: %v", err)
	}
	for _, row := range rows[1:] { // 忽略标题行
		name := row[0]
		serverHost := row[1]
		port, _ := strconv.Atoi(row[2])
		username := row[3]
		password := row[4]
		note := row[5]

		host := &cmdb.CmdbHosts{
			Name:       name,
			ServerHost: serverHost,
			Port:       &port,
			Username:   username,
			Password:   password,
			Note:       note,
			Project:    projectId,
		}

		if err := cmdbHostsService.CreateCmdbHosts(host); err != nil {
			global.KUBEGALE_LOG.Error(fmt.Sprintf("创建主机 %s 失败: %v\n", name, err))
		} else {
			global.KUBEGALE_LOG.Info(fmt.Sprintf("主机 %s 创建成功\n", name))
		}
	}

	return nil
}

// GetCmdbHostsInfoList 分页获取cmdbHosts表记录
func (cmdbHostsService *CmdbHostsService) GetCmdbHostsInfoList(info cmdbReq.CmdbHostsSearch) (list []cmdb.CmdbHosts, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 创建db
	db := global.KUBEGALE_DB.Model(&cmdb.CmdbHosts{})
	var cmdbHostss []cmdb.CmdbHosts

	// 添加项目ID查询条件
	if info.Project > 0 {
		db = db.Where("project = ?", info.Project)
	}

	// 只查询未删除的记录
	db = db.Unscoped().Where("deleted_at IS NULL")

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 执行查询
	err = db.Find(&cmdbHostss).Error
	return cmdbHostss, total, err
}

// CreateSSHSession 创建SSH会话
func (cmdbHostsService *CmdbHostsService) CreateSSHSession(host *cmdb.CmdbHosts) (*ssh.Session, error) {
	port := strconv.Itoa(*host.Port)
	privateKeyPath, err := utils.GetDefaultPrivateKeyPath()
	if err != nil {
		return nil, fmt.Errorf("获取默认私钥路径失败: %v", err)
	}

	// 检查私钥文件是否存在
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("私钥文件不存在: %s", privateKeyPath)
	}

	// 创建SSH客户端
	client, err := utils.CreateSSHClient(host.ServerHost, port, host.Username, privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("创建SSH客户端失败: %v", err)
	}

	// 创建新的会话
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}

	return session, nil
}
