package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server interface {
	WithAddr(addr string) Server
	WithErrLogger(l *log.Logger) Server
	WithRouter(router *mux.Router) Server
	Start() error
	Close(ctx context.Context) error
}

type server struct {
	srv *http.Server
}

func GetServer() Server {
	return &server{
		srv: &http.Server{},
	}
}

func (s *server) WithAddr(addr string) Server {
	s.srv.Addr = addr
	return s
}

func (s *server) WithErrLogger(l *log.Logger) Server {
	s.srv.ErrorLog = l
	return s
}

func (s *server) WithRouter(router *mux.Router) Server {
	s.srv.Handler = router
	return s
}

func (s server) Start() error {
	if len(s.srv.Addr) == 0 {
		return errors.New("Server missing address")
	}

	if s.srv.Handler == nil {
		return errors.New("Server missing handler")
	}

	return s.srv.ListenAndServe()
}

func (s *server) Close(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
