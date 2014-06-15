package main

import (
	"os/exec"
	"strings"
)

func ExecCmdWithSplitWindow(splitDirection string, cmdName string) {
	cmd := exec.Command("tmux", "split-window", splitDirection, cmdName)
	RunCommand(cmd)
}

func RenameSession(sessionName string) {
	cmd := exec.Command("tmux", "rename-session", sessionName)
	RunCommand(cmd)
}

func KillPane(paneName string) {
	cmd := exec.Command("tmux", "kill-pane", "-t", paneName)
	RunCommand(cmd)
}

func RetrievePaneNumWithCmd(cmdName string) string {
	tmuxFormat := "#S:#I.#P;#{pane_start_command}"
	cmd := exec.Command("tmux", "list-panes", "-F", tmuxFormat)
	debugCmd(cmd)

	buf, _ := cmd.Output()
	listPanes := strings.TrimRight(string(buf), "\000")

	var paneNum = ""
	var paneInfo = strings.Split(listPanes, "\n")
	for _, info := range paneInfo {

		if !strings.Contains(info, ";") {
			continue
		}

		paneStartCmd := strings.Split(info, ";")[1]
		debug(paneStartCmd)

		if strings.HasPrefix(paneStartCmd, cmdName) {
			paneNum = strings.Split(info, ";")[0]
			return paneNum
		}
	}

	return paneNum
}
