package cloudCmdb

import "KubeGale/utils"

var (
	CloudVerify = utils.Rules{"Name": {utils.NotEmpty()}, "AccessKeyId": {utils.NotEmpty()}, "AccessKeySecret": {utils.NotEmpty()}, "Platform": {utils.NotEmpty()}}
)
