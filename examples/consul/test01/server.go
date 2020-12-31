package main

import (
	"context"
	"io"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world:"+r.RemoteAddr)
	})
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return &Server{srv: server}
}

func (s *Server) Start(ctx context.Context) error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
