// +build db

package main

import (
	"github.com/spf13/cobra"
	"gopkg.in/check.v1"
)

//func Test(t *testing.T) { check.TestingT(t) }

type migrateSuite struct{}

var _ = check.Suite(&migrateSuite{})

func (ms *migrateSuite) Test(c *check.C) {
	args := []string{"ginS", "migrate"}

	cmd := &cobra.Command{}

	flags := cmd.PersistentFlags()
	flags.Parse(args)

	cmd = newMigrateCmd(flags)
	err := cmd.RunE(cmd, args)
	c.Assert(err, check.IsNil)
}
