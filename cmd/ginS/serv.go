package main

import (
	. "github.com/hohice/gin-web/pkg/util/log"
	"github.com/hohice/gin-web/server"

	"github.com/spf13/cobra"
)

const servDesc = `
This command enable a API server.

$ ginS serv 

Before to start serv ,you need to config the conf file 

The file is named conf.yaml and it's path is define by  $WALM_CONF_PATH

and the default path is /etc/ginS/conf

`

type ServCmd struct {
	oauth bool
}

func newServCmd() *cobra.Command {
	inst := &ServCmd{}

	cmd := &cobra.Command{
		Use:   "serv",
		Short: "enable a Walm Web Server",
		Long:  servDesc,

		RunE: func(cmd *cobra.Command, args []string) error {
			return inst.run()
		},
	}

	return cmd
}

func (sc *ServCmd) run() error {
	apiErrCh := make(chan error)

	serv := server.NewServer(apiErrCh)

	if err := serv.StartServer(); err != nil {
		Log.Errorf("start API server failed:%s exiting\n", err)
		return err
	} else {
		select {
		case err := <-apiErrCh:
			Log.Errorf("API server exited:%s \n", err)
			return err
		}
	}

}
