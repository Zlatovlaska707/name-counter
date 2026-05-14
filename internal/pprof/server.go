package pprof

import (
	"context"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type Server struct {
	addr string
	srv  *http.Server
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Start() error {
	if s.addr == "" {
		return nil
	}

	if s.addr[0] != ':' {
		log.Printf("Адрес pprof должен начинаться с ':', например ':6060'")
	}

	s.srv = &http.Server{
		Addr:              s.addr,
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		log.Printf("pprof сервер запущен на http://localhost/%s/debug/pprof/", s.addr)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("pprof сервер: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.srv.Shutdown(ctx)
	}
	return nil
}
