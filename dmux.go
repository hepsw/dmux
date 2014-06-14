package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var mainFlags = []cli.Flag{
	cli.BoolFlag{"debug", "Run as DEBUG mode"},
}

func main() {
	app := cli.NewApp()
	app.Name = "dmux"
	app.Version = Version
	app.Usage = "Create Docker environment with tmux"
	app.Author = "tcnksm"
	app.Email = "nsd22843@gmail.com"
	app.Flags = mainFlags
	app.Commands = Commands

	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			os.Setenv("DEBUG", "1")
		}

		// ToDo: check command docker and tmux
		return nil
	}

	app.Run(os.Args)
}
