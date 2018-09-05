package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hohice/gin-web/pkg/setting"
	. "github.com/hohice/gin-web/pkg/util/log"
	"github.com/hohice/gin-web/server/middleware"
)

type Server struct {
	ServiceName             string
	Account                 map[string]string
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
		ServiceName: conf.Service,
		ApiErrCh:    errch,

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

func (s *Server) StartServer() error {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		if err := middleware.Init(); err != nil {
			s.ApiErrCh <- err
		} else {
			defer middleware.Close()
		}

		router := InitRouter(s)
		s.server.Handler = router

		if s.TlsEnable {
			if err := s.server.ListenAndServeTLS(s.TlsCertFile, s.TlsKeyFile); err != nil {
				s.ApiErrCh <- err
			}
		} else {
			if err := s.server.ListenAndServe(); err != nil {
				s.ApiErrCh <- err
			}
		}

	}()

	<-sigs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		Log.Fatalln("Server Shutdown:", err)
	}
	s.ApiErrCh <- errors.New("Recv Signal Interrupt, Shutdown Server ...")

	return nil

}
