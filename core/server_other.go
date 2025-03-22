package core

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"time"
)

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Minute
	s.WriteTimeout = 10 * time.Minute
	s.MaxHeaderBytes = 1 << 20
	return s
}
