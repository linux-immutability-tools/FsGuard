package main

import (
	"fmt"
	"github.com/linux-immutability-tools/FsGuard/cmd"
	"github.com/linux-immutability-tools/FsGuard/config"
	"os"
	"strings"
	"syscall"
)

var (
	Version = "0.1.0"
)

func main() {

	pid := os.Getpid()
	fmt.Println("PID of this process:", pid)

	// This cannot be used until we find a way to execute a command while replacing the current process similiar to the exec command on linux
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(executable) == config.InitLocation {
		cmd.ValidateCommand(nil, []string{config.FileListPath})
		if config.RunPostInit {
			fmt.Println("here")
			execErr := syscall.Exec(config.PostInitExec, []string{}, os.Environ())
			if execErr != nil {
				panic(execErr)
			}
			fmt.Printf("here2")
			return
		} else {
			return
		}
	}

	cmd.Execute()
}
