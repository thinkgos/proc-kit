package synq

import (
	"context"

	"github.com/go-kratos/kratos/v3/transport"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

var _ transport.Server = &Server{}

type Server struct {
	mux *asynq.ServeMux
	srv *asynq.Server
}

func NewServer(mux *asynq.ServeMux, srv *asynq.Server) *Server {
	return &Server{mux: mux, srv: srv}
}

func NewServerFromRedisClient(mux *asynq.ServeMux, rdb redis.UniversalClient, config asynq.Config) *Server {
	return NewServer(mux, asynq.NewServerFromRedisClient(rdb, config))
}

// Start implements [transport.Server].
func (s *Server) Start(context.Context) error {
	return s.srv.Start(s.mux)
}

// Stop implements [transport.Server].
func (s *Server) Stop(context.Context) error {
	s.srv.Shutdown()
	return nil
}
