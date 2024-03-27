package cmd

import (
	"github.com/spf13/cobra"

	"rootPrj/apps/sample-using-alias"
)

var sampleAliasCmd = &cobra.Command{
	Use:   "sample-alias",
	Short: "Sample using alias",
	Long:  `Sample using alias, without long project name`,
	Run: func(_ *cobra.Command, _ []string) {
		sample_using_alias.Run()
	},
}

func init() {
	rootCmd.AddCommand(sampleAliasCmd)
}
