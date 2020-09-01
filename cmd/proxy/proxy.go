package main

import (
	"os"
	"strconv"
	"time"
)

import (
	"github.com/urfave/cli"
)

// Version proxy version
var Version = "0.1.0"

func main() {
	app := newProxyApp(&cmdStart)

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Run(os.Args)
}

func newProxyApp(startCmd *cli.Command) *cli.App {
	app := cli.NewApp()
	app.Name = "dubbogo proxy"
	app.Version = Version
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo"
	app.Usage = "Dubbogo proxy is a lightweight gateway."
	app.Flags = cmdStart.Flags

	//commands
	app.Commands = []cli.Command{
		cmdStart,
		cmdStop,
		cmdReload,
	}

	//action
	app.Action = func(c *cli.Context) error {
		if c.NumFlags() == 0 {
			return cli.ShowAppHelp(c)
		}

		return startCmd.Action.(func(c *cli.Context) error)(c)
	}

	return app
}
