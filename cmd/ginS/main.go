package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/hohice/gin-web/pkg/setting"
	. "github.com/hohice/gin-web/pkg/util/log"
)

var globalUsage = `The Gin web API server starter

To begin working with ginS, run the 'serv' command:

	$ ginS serv

`
var conf setting.Config

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ginS",
		Short:        "The Gin web API server starter.",
		Long:         globalUsage,
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()

	cmd.AddCommand(
		newServCmd(),
		newVersionCmd(),
	)

	flags.Parse(args)
	//setting.Init()

	return cmd
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %v", p)
			debug.PrintStack()
		}
	}()

	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		Log.Errorln(err)
		os.Exit(1)
	}
}
