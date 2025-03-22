package request

import (
	"KubeGale/model/common/request"
	"KubeGale/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
