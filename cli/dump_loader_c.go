package main

import (
	"fmt"

	"github.com/sliverarmory/beignet"
	"github.com/spf13/cobra"
)

var dumpLoaderCCmd = &cobra.Command{
	Use:          "dump-loader-c",
	Short:        "Print the embedded loader C source",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Print(beignet.LoaderCSource())
		return err
	},
}
