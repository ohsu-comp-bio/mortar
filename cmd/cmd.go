package cmd

import (
  "github.com/ohsu-comp-bio/funnel/logger"
	"github.com/spf13/cobra"
)

var log = logger.NewLogger("test", logger.DefaultConfig())

// Cmd represents the root "mortar" command
var Cmd = &cobra.Command{
	Use:           "mortar",
	SilenceErrors: true,
	SilenceUsage:  true,
}
