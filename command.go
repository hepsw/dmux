package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
	"strings"
)

const dmuxContainerName = "dmux-playground"

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

func debugCmd(cmd *exec.Cmd) {
	debug(cmd.Args[0], strings.Join(cmd.Args[1:], " "))
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doInit(c *cli.Context) {

	imageName := "ubuntu"
	if c.String("image") != "" {
		imageName = c.String("image")
	}
	debug("imageName:", imageName)

	splitDirection := "-v" // Should be configrable

	if HasContainer(dmuxContainerName) {
		fmt.Fprintf(os.Stderr, "Container %s is already initialized\n", dmuxContainerName)
		os.Exit(1)
	}

	// Start new container with tmux new pane
	StartNewContainerCmd := fmt.Sprintf("exec docker run -t -i --name %s %s /bin/bash", dmuxContainerName, imageName)
	ExecCmdWithSplitWindow(splitDirection, StartNewContainerCmd)

}

func doStop(c *cli.Context) {

	if !HasContainer(dmuxContainerName) {
		fmt.Fprintf(os.Stderr, "Container is not initialized, initialize it with `dmux init` command\n")
		os.Exit(1)
	}

	if IsPaused(dmuxContainerName) {
		fmt.Fprintf(os.Stderr, "Container is already stoped, start it with `dmux start` command\n")
		os.Exit(1)
	}

	PauseContainer(dmuxContainerName)

	// Kill docker running pane
	paneName := RetrievePaneNumWithCmd("exec docker")
	if paneName == "" {
		os.Exit(0)
	}

	debug("docker running pane:", paneName)
	KillPane(paneName)
	os.Exit(0)
}

func doStart(c *cli.Context) {

	if !HasContainer(dmuxContainerName) {
		fmt.Fprintf(os.Stderr, "Container is not initialized, initialize it with `dmux init` command\n")
		os.Exit(1)
	}

	if !IsPaused(dmuxContainerName) {
		fmt.Fprintf(os.Stderr, "Container is already working\n")
		os.Exit(1)
	}

	UnpauseContainer(dmuxContainerName)

	// Attach container with tmux new pane
	attachContainerCmd := "exec docker attach dmux-playground"
	splitDirection := "-v" // Should be configrable
	ExecCmdWithSplitWindow(splitDirection, attachContainerCmd)

	os.Exit(0)
}

func doDelete(c *cli.Context) {

	if !HasContainer(dmuxContainerName) {
		fmt.Fprintf(os.Stderr, "Container is not initialized, initialize it with `dmux init` command\n")
		os.Exit(1)
	}

	if IsPaused(dmuxContainerName) {
		UnpauseContainer(dmuxContainerName)
	}

	DeleteContainer(dmuxContainerName)
}

func doSave(c *cli.Context) {

	if len(c.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Set docker image name for save\n")
		os.Exit(1)
	}

	imageName := c.Args().First()

	if !HasContainer(dmuxContainerName) {
		fmt.Fprintf(os.Stderr, "Container is not initialized, initialize it with `dmux init` command\n")
		os.Exit(1)
	}

	if IsPaused(dmuxContainerName) {
		// Unpause container
		UnpauseContainer(dmuxContainerName)
	}

	SaveContainer(dmuxContainerName, imageName)
}
