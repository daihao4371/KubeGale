package main

import (
	"KubeGale/core"
	"os"
)

func main() {
	cmd := core.NewServerCommand()
	if err := cmd.Execute(); err != nil {
		panic(err)
		os.Exit(1)
	}
}
