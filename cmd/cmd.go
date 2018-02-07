package cmd

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/spf13/cobra"
)

var log = logger.NewLogger("test", logger.DefaultConfig())

var cmd = &cobra.Command{
	Use:           "mortar",
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute executes the "mortar" CLI command against os.Args
// and returns any error raised during execution.
func Execute() error {
	return cmd.Execute()
}
