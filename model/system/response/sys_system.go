package response

import "KubeGale/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
