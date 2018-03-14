package main

import (
	"os"

	"github.com/dynamicgo/slf4go"
	cli "gopkg.in/urfave/cli.v2"
)

var logger = slf4go.Get("neogo")

func main() {
	app := cli.App{}

	app.Commands = []*cli.Command{
		deployCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.ErrorF("%s", err)
		return
	}
}
