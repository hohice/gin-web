package main

import (
	"os"
	"runtime/debug"

	"github.com/hohice/gin-web/pkg/setting"
	"github.com/spf13/cobra"
)

var globalUsage = `The Gin web API server starter

To begin working with ginS, run the 'serv' command:

	$ ginS [serv| migrate| version]

`

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ginS",
		Short:        "The Gin Web API Server Starter.",
		Long:         globalUsage,
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()
	setting.AddFlags(flags)

	cmd.AddCommand(
		newServCmd(flags),
		newVersionCmd(),
		newMigrateCmd(flags),
	)

	flags.Parse(args)
	return cmd
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			debug.PrintStack()
		}
	}()

	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
