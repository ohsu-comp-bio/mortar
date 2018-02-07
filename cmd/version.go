package cmd

import (
	"fmt"

	"github.com/ohsu-comp-bio/mortar/version"
	"github.com/spf13/cobra"
)

func init() {
	vCmd := cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.String())
		},
	}
	cmd.AddCommand(&vCmd)
}
