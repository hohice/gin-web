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
	"github.com/hohice/gin-web/server/middleware"
)

type Server struct {
	ServiceName             string
	Account                 map[string]string
	Port                    int
	TlsEnable               bool
	TlsCertFile, TlsKeyFile string
	ReadTimeout             time.Duration
	WriteTimeout            time.Duration
	Debug                   bool
	ZipkinUrl               string
	server                  *http.Server
}

func NewServer() *Server {
	conf := setting.Config

	return &Server{
		ServiceName: conf.Service,

		TlsEnable:   conf.Secret.Tls,
		TlsCertFile: conf.Secret.TlsCert,
		TlsKeyFile:  conf.Secret.TlsKey,

		Account: conf.Secret.Account,

		Debug: conf.Debug,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", conf.Http.HTTPPort),
			ReadTimeout:    conf.Http.ReadTimeout,
			WriteTimeout:   conf.Http.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *Server) StartServer() (<-chan struct{}, <-chan error) {
	errchan := make(chan error, 1)
	donechan := make(chan struct{}, 1)

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	startup := func() {
		if err := middleware.Init(&setting.Config); err != nil {
			errchan <- errors.New(fmt.Sprintf("Init Middleware: %s", err))
		} else {
			defer middleware.Close()
		}

		router := s.InitRouter()
		s.server.Handler = router

		if s.TlsEnable {
			if err := s.server.ListenAndServeTLS(s.TlsCertFile, s.TlsKeyFile); err != nil {
				errchan <- err
			}
		} else {
			if err := s.server.ListenAndServe(); err != nil {
				errchan <- err
			}
		}

	}

	go startup()

	shutdown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			errchan <- errors.New(fmt.Sprintf("Server Shutdown: %s", err))
		}

		donechan <- struct{}{}
	}

	go func() {
		select {
		case <-sigs:
			shutdown()
		}
	}()

	return donechan, errchan

}
