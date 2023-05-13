package main

import (
	"github.com/linux-immutability-tools/FsGuard/cmd"
	"os"
	"strings"
)

var (
	Version = "0.1.0"
)

func main() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	
	if strings.TrimSpace(executable) == "/usr/bin/init" { // TODO: Configuring init location
		cmd.ValidateCommand(nil, []string{"/FsGuard/hashList"}) // TODO: allow configuring path to the hash list
		return
	}
	cmd.Execute()
}
