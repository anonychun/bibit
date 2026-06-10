package auth

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/util"
	pb "github.com/anonychun/bibit/pkg/gen/proto/api/v1/app/auth"
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

func (h *GrpcHandler) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	usecaseReq := SignUpRequest{
		IpAddress:    util.GrpcPeerAddress(ctx),
		UserAgent:    util.GrpcMetadataValue(ctx, "user-agent"),
		Name:         req.GetName(),
		EmailAddress: req.GetEmailAddress(),
		Password:     req.GetPassword(),
	}

	res, err := h.usecase.SignUp(ctx, usecaseReq)
	if err != nil {
		return nil, err
	}

	return &pb.SignUpResponse{Token: res.Token}, nil
}

func (h *GrpcHandler) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	usecaseReq := SignInRequest{
		IpAddress:    util.GrpcPeerAddress(ctx),
		UserAgent:    util.GrpcMetadataValue(ctx, "user-agent"),
		EmailAddress: req.GetEmailAddress(),
		Password:     req.GetPassword(),
	}

	res, err := h.usecase.SignIn(ctx, usecaseReq)
	if err != nil {
		return nil, err
	}

	return &pb.SignInResponse{Token: res.Token}, nil
}

func (h *GrpcHandler) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	token := req.GetToken()
	if token == "" {
		token = util.GrpcMetadataValue(ctx, "authorization")
	}

	usecaseReq := SignOutRequest{
		Token: token,
	}

	err := h.usecase.SignOut(ctx, usecaseReq)
	if err != nil {
		return nil, err
	}

	return &pb.SignOutResponse{}, nil
}

func (h *GrpcHandler) Me(ctx context.Context, _ *pb.MeRequest) (*pb.MeResponse, error) {
	res, err := h.usecase.Me(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.MeResponse{
		User: &pb.MeResponse_User{
			Id:           res.User.Id.String(),
			Name:         res.User.Name,
			EmailAddress: res.User.EmailAddress,
		},
	}, nil
}
