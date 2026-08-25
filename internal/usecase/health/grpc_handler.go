package health

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	pb "github.com/anonychun/bibit/pkg/pb/health"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewGrpcHandler)
}

type IGrpcHandler interface {
	pb.ServiceServer
}

type GrpcHandler struct {
	pb.UnimplementedServiceServer
	usecase IUsecase
}

var _ IGrpcHandler = (*GrpcHandler)(nil)

func NewGrpcHandler(i do.Injector) (*GrpcHandler, error) {
	return &GrpcHandler{
		usecase: do.MustInvoke[*Usecase](i),
	}, nil
}

func (h *GrpcHandler) Up(ctx context.Context, req *pb.UpRequest) (*pb.UpResponse, error) {
	res, err := h.usecase.Up(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.UpResponse{Status: res.Status}, nil
}
