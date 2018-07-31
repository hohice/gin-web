package main

import (
	"github.com/hohice/gin-web/pkg/version"

	"github.com/spf13/cobra"
)

const versionDesc = `
This command print version of pm.`

type vCmd struct {
}

func newVersionCmd() *cobra.Command {
	vc := &vCmd{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "print version",
		Long:  versionDesc,

		RunE: func(cmd *cobra.Command, args []string) error {
			defer vc.run()
			return nil
		},
	}
	return cmd
}

func (vc *vCmd) run() {
	version.PrintVersionInfo()
}
