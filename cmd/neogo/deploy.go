package main

import (
	"path/filepath"

	cli "gopkg.in/urfave/cli.v2"
)

var deployCommand = &cli.Command{
	Name:      "deploy",
	Usage:     "deploy neo smart contract",
	Action:    deploy,
	ArgsUsage: "contract_root_path",
}

type projectConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

func deploy(c *cli.Context) error {
	if c.Args().Len() != 1 {
		cli.ShowCommandHelpAndExit(c, "deploy", 1)
	}

	rootPath, err := filepath.Abs(c.Args().First())

	if err != nil {
		return err
	}

	logger.InfoF("contract root path: %s", rootPath)

	avms, err := filepath.Glob(filepath.Join(rootPath, "*.avm"))

	if err != nil {
		return err
	}

	for _, avm := range avms {
		logger.Debug(avm)
	}

	return nil
}
