package main

import (
	"github.com/urfave/cli/v2"

	"rootPrj/apps/sample-using-alias"
)

func sampleAliasCommand() *cli.Command {
	return &cli.Command{
		Name:        "sample-alias",
		Usage:       "Sample using alias",
		Description: `Sample using alias, without long project name`,
		Action: func(_ *cli.Context) error {
			sample_using_alias.Run()
			return nil
		},
	}
}
