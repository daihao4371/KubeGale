package response

import "KubeGale/model/system/request"

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
