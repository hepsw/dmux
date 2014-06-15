package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
	"strings"
)

var Commands = []cli.Command{
	commandInit,
	commandStop,
	commandStart,
	commandDelete,
	commandSave,
}

var commandInit = cli.Command{
	Name:  "init",
	Usage: "Start new container",
	Description: `
Start new container with new tmux pane
`,
	Flags:  commandInitFlags,
	Action: doInit,
}

var commandInitFlags = []cli.Flag{
	cli.StringFlag{"i, image", "ubuntu", "Select docker image"},
}

var commandStop = cli.Command{
	Name:  "stop",
	Usage: "Stop container",
	Description: `
Stop(Pause) container and kill tmux pane which container is working.
`,
	Action: doStop,
}

var commandStart = cli.Command{
	Name:  "start",
	Usage: "Start paused container",
	Description: `
Restart container which is paused with new tmux pane.
`,
	Action: doStart,
}

var commandDelete = cli.Command{
	Name:  "delete",
	Usage: "Delete container",
	Description: `
Kill container and kill tmux pane pane which container is working.
`,
	Action: doDelete,
}

var commandSave = cli.Command{
	Name:  "save",
	Usage: "Save container as Image",
	Description: `
Save currently working container as Docker image
`,
	Action: doSave,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func doInit(c *cli.Context) {
	// Change Session name
	cmd := exec.Command("tmux", "rename-session", "dmux")
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	cmd.Run()

	// Start new container with tmux new window
	StartNewContainerCmd := "exec docker run -t -i --name dmux-playground ubuntu /bin/bash"
	splitDirection := "-v" // Should be configrable
	cmd = exec.Command("tmux", "split-window", splitDirection, StartNewContainerCmd)
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	cmd.Run()
}

func doStop(c *cli.Context) {
	// Pause container
	cmd := exec.Command("docker", "pause", "dmux-playground")
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	cmd.Run()

	// Retrieve Docker running pane
	cmd = exec.Command("tmux", "list-panes", "-F", "#I.#P:#{pane_start_command}")
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	buf, _ := cmd.Output()
	listPanes := strings.TrimRight(string(buf), "\000")

	var paneNum = ""
	var values = strings.Split(listPanes, "\n")
	for _, info := range values {

		if !strings.Contains(info, ":") {
			continue
		}
		paneStartCmd := strings.Split(info, ":")[1]
		debug(paneStartCmd)
		if strings.HasPrefix(paneStartCmd, "exec docker run") {
			paneNum = strings.Split(info, ":")[0]
			debug(paneNum)
		}
	}

	debug("docker running pane:", "dmux:"+paneNum)

	// Kill docker running pane
	cmd = exec.Command("tmux", "kill-pane", "-t", "dmux:"+paneNum)
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	cmd.Run()
}

func doStart(c *cli.Context) {
	// Unpause container
	cmd := exec.Command("docker", "unpause", "dmux-playground")
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	cmd.Run()

	// Split window
	attachContainerCmd := "exec docker attach dmux-playground"
	splitDirection := "-v" // Should be configrable
	cmd = exec.Command("tmux", "split-window", splitDirection, attachContainerCmd)
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	cmd.Run()
}

func doDelete(c *cli.Context) {
	// ToDO: # docker inspect  - format '{{.State.Paused }}' dmux-playground
	cmd := exec.Command("docker", "rm", "-f", "dmux-playground")
	cmd.Run()

	cmd = exec.Command("tmux", "kill-pane", "-t", "dmux:1.3")
	cmd.Run()
}

func doSave(c *cli.Context) {
	//ToDo
}
