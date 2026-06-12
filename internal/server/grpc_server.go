package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/config"
	"github.com/anonychun/bibit/internal/observability"
	usecaseApiV1AppAuth "github.com/anonychun/bibit/internal/usecase/api/v1/app/auth"
	pbApiV1AppAuth "github.com/anonychun/bibit/pkg/pb/api/v1/app/auth"
	"github.com/samber/do/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	do.Provide(bootstrap.Injector, NewGrpcServer)
}

type IGrpcServer interface {
	Start(ctx context.Context) error
}

type GrpcServer struct {
	server        *grpc.Server
	listener      net.Listener
	observability observability.IObservability
}

var _ IGrpcServer = (*GrpcServer)(nil)

func NewGrpcServer(i do.Injector) (*GrpcServer, error) {
	cfg := do.MustInvoke[*config.Config](i)
	o11y := do.MustInvoke[*observability.Observability](i)

	srv := grpc.NewServer()
	registerGrpcHandlers(i, srv)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))
	if err != nil {
		return nil, err
	}

	return &GrpcServer{
		server:        srv,
		listener:      lis,
		observability: o11y,
	}, nil
}

func (s *GrpcServer) Start(ctx context.Context) error {
	s.observability.Logger().Info("starting grpc server", slog.String("addr", s.listener.Addr().String()))
	err := s.server.Serve(s.listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *GrpcServer) Shutdown(ctx context.Context) error {
	s.observability.Logger().Info("shutting down grpc server")
	s.server.GracefulStop()
	return s.listener.Close()
}

func registerGrpcHandlers(i do.Injector, srv *grpc.Server) {
	pbApiV1AppAuth.RegisterServiceServer(srv, do.MustInvoke[*usecaseApiV1AppAuth.GrpcHandler](i))
}
