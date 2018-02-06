package main

import (
  "os"
  "github.com/ohsu-comp-bio/funnel/logger"
  "github.com/ohsu-comp-bio/mortar/cmd"
)

func main() {
	if err := cmd.Cmd.Execute(); err != nil {
		logger.PrintSimpleError(err)
		os.Exit(1)
	}
}
