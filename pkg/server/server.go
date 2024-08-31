package server

import (
	"context"
	"net/http"
)

type HTTPServer struct {
	http *http.Server
}

func New(handler http.Handler, port string) *HTTPServer {
	s := &HTTPServer{
		http: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}

	return s
}

func (s *HTTPServer) Start() error {
	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
