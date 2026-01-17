package auth

import (
	"context"

	"github.com/anonychun/bibit/internal/bootstrap"
	"github.com/anonychun/bibit/internal/consts"
	"github.com/anonychun/bibit/internal/current"
	"github.com/anonychun/bibit/internal/entity"
	"github.com/anonychun/bibit/internal/repository"
	repositoryUser "github.com/anonychun/bibit/internal/repository/user"
	repositoryUserSession "github.com/anonychun/bibit/internal/repository/user_session"
	"github.com/anonychun/bibit/internal/validation"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewUsecase)
}

type Usecase interface {
	SignUp(ctx context.Context, req SignUpRequest) (*SignUpResponse, error)
	SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error)
	SignOut(ctx context.Context, req SignOutRequest) error
	Me(ctx context.Context) (*MeResponse, error)
}

type UsecaseImpl struct {
	validator             validation.Validator
	userRepository        repositoryUser.Repository
	userSessionRepository repositoryUserSession.Repository
}

func NewUsecase(i do.Injector) (Usecase, error) {
	return &UsecaseImpl{
		validator:             do.MustInvoke[validation.Validator](i),
		userRepository:        do.MustInvoke[repositoryUser.Repository](i),
		userSessionRepository: do.MustInvoke[repositoryUserSession.Repository](i),
	}, nil
}

func (u *UsecaseImpl) SignUp(ctx context.Context, req SignUpRequest) (*SignUpResponse, error) {
	validationErr := u.validator.Struct(&req)
	isEmailAddressExists, err := u.userRepository.ExistsByEmailAddress(ctx, req.EmailAddress)
	if err != nil {
		return nil, err
	}

	if isEmailAddressExists {
		validationErr.AddError("emailAddress", consts.ErrEmailAddressAlreadyRegistered)
	}

	if validationErr.IsFail() {
		return nil, validationErr
	}

	user := &entity.User{
		Name:         req.Name,
		EmailAddress: req.EmailAddress,
	}

	err = user.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	res := &SignUpResponse{}
	repository.Transaction(ctx, func(ctx context.Context) error {
		err = u.userRepository.Create(ctx, user)
		if err != nil {
			return err
		}

		userSession := &entity.UserSession{
			UserId:    user.Id,
			IpAddress: req.IpAddress,
			UserAgent: req.UserAgent,
		}
		userSession.GenerateToken()

		err = u.userSessionRepository.Create(ctx, userSession)
		if err != nil {
			return err
		}

		res.Token = userSession.Token
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *UsecaseImpl) SignIn(ctx context.Context, req SignInRequest) (*SignInResponse, error) {
	validationErr := u.validator.Struct(&req)
	if validationErr.IsFail() {
		return nil, validationErr
	}

	user, err := u.userRepository.FindByEmailAddress(ctx, req.EmailAddress)
	if err == consts.ErrRecordNotFound {
		return nil, consts.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	err = user.ComparePassword(req.Password)
	if err != nil {
		return nil, consts.ErrInvalidCredentials
	}

	userSession := &entity.UserSession{
		UserId:    user.Id,
		IpAddress: req.IpAddress,
		UserAgent: req.UserAgent,
	}
	userSession.GenerateToken()

	err = u.userSessionRepository.Create(ctx, userSession)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{Token: userSession.Token}, nil
}

func (u *UsecaseImpl) SignOut(ctx context.Context, req SignOutRequest) error {
	err := u.userSessionRepository.DeleteByToken(ctx, req.Token)
	if err != nil {
		return err
	}

	return nil
}

func (u *UsecaseImpl) Me(ctx context.Context) (*MeResponse, error) {
	user := current.User(ctx)
	if user == nil {
		return nil, consts.ErrUnauthorized
	}

	res := &MeResponse{}
	res.User.Id = user.Id.String()
	res.User.Name = user.Name
	res.User.EmailAddress = user.EmailAddress

	return res, nil
}
