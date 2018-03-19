package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/inwecrypto/neogo/tx"
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

	configFile := filepath.Join(rootPath, "projec.json")

	data, err := ioutil.ReadFile(configFile)

	if err != nil {
		return err
	}

	var config *projectConfig

	err = json.Unmarshal(data, &config)

	if err != nil {
		return err
	}

	_, err = ioutil.ReadFile(filepath.Join(rootPath, config.Name+".json"))

	if err != nil {
		return err
	}

	tx.NewInvocationTx(nil, 0)

	return nil
}
