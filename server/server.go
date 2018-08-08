package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hohice/gin-web/pkg/setting"
	. "github.com/hohice/gin-web/pkg/util/log"
	"github.com/hohice/gin-web/server/middleware"
)

type Server struct {
	ApiErrCh                chan error
	Port                    int
	TlsEnable               bool
	TlsCertFile, TlsKeyFile string
	ReadTimeout             time.Duration
	WriteTimeout            time.Duration
	Debug                   bool
	ZipkinUrl               string
	server                  *http.Server
}

func NewServer(errch chan error) *Server {
	if setting.Config.Http.HTTPPort == 0 {
		Log.Fatalln("start API server failed, please spec Http port")
	}
	conf := setting.Config
	return &Server{
		ApiErrCh: errch,

		TlsEnable:   conf.Secret.Tls,
		TlsCertFile: conf.Secret.TlsCert,
		TlsKeyFile:  conf.Secret.TlsKey,

		Debug: conf.Debug,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", conf.Http.HTTPPort),
			ReadTimeout:    conf.Http.ReadTimeout,
			WriteTimeout:   conf.Http.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (server *Server) StartServer() error {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := middleware.Init(); err != nil {
			server.ApiErrCh <- err
		} else {
			defer middleware.Close()
		}

		router := InitRouter(server.Debug)
		server.server.Handler = router

		if server.TlsEnable {
			if err := server.server.ListenAndServeTLS(server.TlsCertFile, server.TlsKeyFile); err != nil {
				server.ApiErrCh <- err
			}
		} else {
			if err := server.server.ListenAndServe(); err != nil {
				server.ApiErrCh <- err
			}
		}

	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.server.Shutdown(ctx); err != nil {
		Log.Fatalln("Server Shutdown:", err)
	}
	server.ApiErrCh <- errors.New("Recv Signal Interrupt, Shutdown Server ...")

	return nil

}
