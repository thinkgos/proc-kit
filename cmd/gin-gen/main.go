package main

import (
	"os"

	_ "github.com/thinkgos/carp/driver/mysql"
	"github.com/thinkgos/proc-kit/cmd/gin-gen/command"
)

func main() {
	err := command.NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
