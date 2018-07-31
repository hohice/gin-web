package router

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
	"github.com/hohice/gin-web/router/middleware"
)

type Server struct {
	ApiErrCh                chan error
	Port                    int
	TlsEnable               bool
	TlsCertFile, TlsKeyFile string
	OauthEnable             bool
	ReadTimeout             time.Duration
	WriteTimeout            time.Duration
	Debug                   bool
	ZipkinUrl               string
	server                  *http.Server
}

func NewServer(errch chan error) *Server {
	conf := setting.Config
	return &Server{
		ApiErrCh: errch,

		OauthEnable: conf.Auth.Enable,
		TlsEnable:   conf.Secret.Tls,
		TlsCertFile: conf.Secret.TlsCert,
		TlsKeyFile:  conf.Secret.TlsKey,

		Debug:     conf.Debug,
		ZipkinUrl: conf.Trace.ZipkinUrl,
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

		if !server.Debug {
			//EndTrac will be called when close the server
			//so the init need be placed here
			if err, closeble := middleware.InitTracer("ginS", server.ZipkinUrl, server.Port); err != nil {
				server.ApiErrCh <- err
				return
			} else {
				defer closeble()
			}
		}

		router := InitRouter(server.OauthEnable, server.Debug)

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
