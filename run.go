package main

import (
	"os/exec"
)

func RunCommand(cmd *exec.Cmd) {
	debugCmd(cmd)
	cmd.Run()
}
