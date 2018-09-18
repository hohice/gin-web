package main

import (
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/hohice/gin-web/pkg/util/logger"
	"github.com/hohice/gin-web/server"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const servDesc = `
This command enable a API server.

$ ginS serv 

Before to start serv ,you need to config the config file 

The file is named config.yaml and it's path is define by  $GINS_CONF_PATH

and the default path is /etc/ginS/config.yaml

`

type ServCmd struct{}

func newServCmd(fs *pflag.FlagSet) *cobra.Command {
	inst := &ServCmd{}

	cmd := &cobra.Command{
		Use:   "serv",
		Short: "Enable Web Server",
		Long:  servDesc,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := setting.Init(); err != nil {
				return err
			}
			return inst.run()
		},
	}

	return cmd
}

func (sc *ServCmd) run() error {
	log := logger.DefaultLogger
	done, errchan := server.NewServer().StartServer()
	select {
	case <-done:
		log.Infow("Recv Signal Interrupt, Shutdown Server ...")
		return nil
	case err := <-errchan:
		log.Errorw("Server error exited", "error:", err)
		return err
	}
}
