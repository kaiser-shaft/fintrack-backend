package httpserver

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Port string `envconfig:"HTTP_PORT" default:"8080"`
}

type Server struct {
	server *http.Server
}

func New(handler http.Handler, c Config) *Server {
	return &Server{
		server: &http.Server{
			Addr:         net.JoinHostPort("", c.Port),
			Handler:      handler,
			ReadTimeout:  20 * time.Second,
			WriteTimeout: 20 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("httpserver.ListenAndServe: %v", err)
	}
	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
