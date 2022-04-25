package server

import (
	"context"
	"fmt"
	"github.com/Ypxd/diplom/auth/utils"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewServer(handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", utils.GetConfig().Server.Host, utils.GetConfig().Server.Port),
			ReadTimeout:  utils.GetConfig().Server.RTimeout * time.Second,
			WriteTimeout: utils.GetConfig().Server.WTimeout * time.Second,
			Handler:      handler,
		},
	}
}
