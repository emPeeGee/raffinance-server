package connection

import (
	"context"
	"net/http"

	"github.com/emPeeGee/raffinance/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(serverConf config.Server, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           serverConf.Addr,
		MaxHeaderBytes: serverConf.MaxHeaderBytes,
		ReadTimeout:    serverConf.ReadTimeout,
		WriteTimeout:   serverConf.WriteTimeout,
		Handler:        handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
