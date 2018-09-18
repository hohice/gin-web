package main

import (
	"github.com/hohice/gin-web/db"
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const migrateDesc = `
This command enable a auto migrate.

$ ginS migrate

Before to start migrate ,you need to config the conf file 

The file is named conf.yaml and it's path is define by  $WALM_CONF_PATH

and the default path is /etc/ginS/conf

`

type MigrateCmd struct{}

func newMigrateCmd(fs *pflag.FlagSet) *cobra.Command {
	inst := &MigrateCmd{}

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "enable a db auto migrate",
		Long:  migrateDesc,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := setting.Init(); err != nil {
				return err
			}
			return inst.run()
		},
	}
	return cmd
}

func (mc *MigrateCmd) run() error {
	return db.AutoMigrate()
}
