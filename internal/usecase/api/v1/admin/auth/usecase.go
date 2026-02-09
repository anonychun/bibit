package auth

import (
	"context"
	"database/sql"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/consts"
	"github.com/anonychun/bibit/internal/current"
	"github.com/anonychun/bibit/internal/entity"
	repositoryAdmin "github.com/anonychun/bibit/internal/repository/admin"
	repositoryAdminSession "github.com/anonychun/bibit/internal/repository/admin_session"
	"github.com/anonychun/bibit/internal/validation"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewUsecase)
}

type IUsecase interface {
	SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error)
	SignOut(ctx context.Context, req SignOutRequest) error
	Me(ctx context.Context) (*MeResponse, error)
}

type Usecase struct {
	validator              validation.IValidator
	adminRepository        repositoryAdmin.IRepository
	adminSessionRepository repositoryAdminSession.IRepository
}

var _ IUsecase = (*Usecase)(nil)

func NewUsecase(i do.Injector) (*Usecase, error) {
	return &Usecase{
		validator:              do.MustInvoke[*validation.Validator](i),
		adminRepository:        do.MustInvoke[*repositoryAdmin.Repository](i),
		adminSessionRepository: do.MustInvoke[*repositoryAdminSession.Repository](i),
	}, nil
}

func (u *Usecase) SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error) {
	validationErr := u.validator.Struct(&req)
	if validationErr.IsFail() {
		return nil, validationErr
	}

	admin, err := u.adminRepository.FindByEmailAddress(ctx, req.EmailAddress)
	if err == sql.ErrNoRows {
		return nil, consts.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	err = admin.ComparePassword(req.Password)
	if err != nil {
		return nil, consts.ErrInvalidCredentials
	}

	adminSession := &entity.AdminSession{
		AdminId:   admin.Id,
		IpAddress: req.IpAddress,
		UserAgent: req.UserAgent,
	}
	adminSession.GenerateToken()

	err = u.adminSessionRepository.Create(ctx, adminSession)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{Token: adminSession.Token}, nil
}

func (u *Usecase) SignOut(ctx context.Context, req SignOutRequest) error {
	err := u.adminSessionRepository.DeleteByToken(ctx, req.Token)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Me(ctx context.Context) (*MeResponse, error) {
	admin := current.Admin(ctx)
	if admin == nil {
		return nil, consts.ErrUnauthorized
	}

	res := &MeResponse{}
	res.Admin.Id = admin.Id.String()
	res.Admin.Name = admin.Name

	return res, nil
}
