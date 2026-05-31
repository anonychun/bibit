package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/anonychun/bibit/internal/api"
	"github.com/anonychun/bibit/internal/consts"
	"github.com/anonychun/bibit/internal/current"
	"github.com/anonychun/bibit/internal/entity"
	repositoryUser "github.com/anonychun/bibit/internal/repository/user"
	repositoryUserSession "github.com/anonychun/bibit/internal/repository/user_session"
	"github.com/anonychun/bibit/internal/validation"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUsecase_SignUp(t *testing.T) {
	t.Run("returns validation errors from the validator", func(t *testing.T) {
		ctx := context.Background()
		req := SignUpRequest{
			Name:         "Ada Lovelace",
			EmailAddress: "ada@example.com",
			Password:     "short",
		}
		validationErr := api.ValidationError{"password": []string{"Password must be at least 8 characters"}}
		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		usecase := &Usecase{
			validator:      validator,
			userRepository: userRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(validationErr).Once()
		userRepository.EXPECT().ExistsByEmailAddress(ctx, req.EmailAddress).Return(false, nil).Once()

		res, err := usecase.SignUp(ctx, req)

		require.Error(t, err)
		assert.Nil(t, res)

		actualValidationErr, ok := err.(api.ValidationError)
		require.True(t, ok)
		assert.Equal(t, validationErr, actualValidationErr)
	})

	t.Run("returns a validation error when the email address is already registered", func(t *testing.T) {
		ctx := context.Background()
		req := SignUpRequest{
			Name:         "Ada Lovelace",
			EmailAddress: "ada@example.com",
			Password:     "correct horse battery staple",
		}
		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		usecase := &Usecase{
			validator:      validator,
			userRepository: userRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().ExistsByEmailAddress(ctx, req.EmailAddress).Return(true, nil).Once()

		res, err := usecase.SignUp(ctx, req)

		require.Error(t, err)
		assert.Nil(t, res)

		validationErr, ok := err.(api.ValidationError)
		require.True(t, ok)
		assert.Equal(t, []string{consts.ErrEmailAddressAlreadyRegistered.Error()}, validationErr["emailAddress"])
	})

	t.Run("returns an error when the uniqueness check fails", func(t *testing.T) {
		ctx := context.Background()
		req := SignUpRequest{
			Name:         "Ada Lovelace",
			EmailAddress: "ada@example.com",
			Password:     "correct horse battery staple",
		}
		expectedErr := errors.New("check email address")
		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		usecase := &Usecase{
			validator:      validator,
			userRepository: userRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().ExistsByEmailAddress(ctx, req.EmailAddress).Return(false, expectedErr).Once()

		res, err := usecase.SignUp(ctx, req)

		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, res)
	})

	t.Run("returns an error when password hashing fails", func(t *testing.T) {
		ctx := context.Background()
		req := SignUpRequest{
			Name:         "Ada Lovelace",
			EmailAddress: "ada@example.com",
			Password:     strings.Repeat("x", 73),
		}
		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		usecase := &Usecase{
			validator:      validator,
			userRepository: userRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().ExistsByEmailAddress(ctx, req.EmailAddress).Return(false, nil).Once()

		res, err := usecase.SignUp(ctx, req)

		require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
		assert.Nil(t, res)
	})
}

func TestUsecase_SignIn(t *testing.T) {
	t.Run("creates a session when the credentials are valid", func(t *testing.T) {
		ctx := context.Background()
		req := SignInRequest{
			IpAddress:    "127.0.0.1",
			UserAgent:    "Go test",
			EmailAddress: "ada@example.com",
			Password:     "correct horse battery staple",
		}
		userID := uuid.New()
		user := &entity.User{
			Base:         entity.Base{Id: userID},
			Name:         "Ada Lovelace",
			EmailAddress: req.EmailAddress,
		}
		require.NoError(t, user.HashPassword(req.Password))

		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		userSessionRepository := repositoryUserSession.NewMockIRepository(t)
		usecase := &Usecase{
			validator:             validator,
			userRepository:        userRepository,
			userSessionRepository: userSessionRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().FindByEmailAddress(ctx, req.EmailAddress).Return(user, nil).Once()

		var createdSession *entity.UserSession
		userSessionRepository.EXPECT().Create(ctx, mock.AnythingOfType("*entity.UserSession")).Run(func(ctx context.Context, actual *entity.UserSession) {
			createdSession = actual
		}).Return(nil).Once()

		res, err := usecase.SignIn(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, createdSession)
		assert.Equal(t, userID, createdSession.UserId)
		assert.Equal(t, req.IpAddress, createdSession.IpAddress)
		assert.Equal(t, req.UserAgent, createdSession.UserAgent)
		assert.Equal(t, createdSession.Token, res.Token)

		_, err = ulid.ParseStrict(res.Token)
		require.NoError(t, err)
	})

	t.Run("returns validation errors before checking credentials", func(t *testing.T) {
		ctx := context.Background()
		req := SignInRequest{EmailAddress: "not-an-email", Password: "short"}
		validationErr := api.ValidationError{"emailAddress": []string{"Email address is invalid"}}
		validator := validation.NewMockIValidator(t)
		usecase := &Usecase{validator: validator}

		validator.EXPECT().Struct(mock.Anything).Return(validationErr).Once()

		res, err := usecase.SignIn(ctx, req)

		require.Error(t, err)
		assert.Nil(t, res)

		actualValidationErr, ok := err.(api.ValidationError)
		require.True(t, ok)
		assert.Equal(t, validationErr, actualValidationErr)
	})

	t.Run("returns invalid credentials when the user does not exist", func(t *testing.T) {
		ctx := context.Background()
		req := SignInRequest{EmailAddress: "ada@example.com", Password: "correct horse battery staple"}
		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		usecase := &Usecase{
			validator:      validator,
			userRepository: userRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().FindByEmailAddress(ctx, req.EmailAddress).Return(nil, sql.ErrNoRows).Once()

		res, err := usecase.SignIn(ctx, req)

		require.ErrorIs(t, err, consts.ErrInvalidCredentials)
		assert.Nil(t, res)
	})

	t.Run("returns repository errors", func(t *testing.T) {
		ctx := context.Background()
		req := SignInRequest{EmailAddress: "ada@example.com", Password: "correct horse battery staple"}
		expectedErr := errors.New("find user")
		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		usecase := &Usecase{
			validator:      validator,
			userRepository: userRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().FindByEmailAddress(ctx, req.EmailAddress).Return(nil, expectedErr).Once()

		res, err := usecase.SignIn(ctx, req)

		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, res)
	})

	t.Run("returns invalid credentials when the password is wrong", func(t *testing.T) {
		ctx := context.Background()
		req := SignInRequest{EmailAddress: "ada@example.com", Password: "wrong password"}
		user := &entity.User{EmailAddress: req.EmailAddress}
		require.NoError(t, user.HashPassword("correct horse battery staple"))

		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		usecase := &Usecase{
			validator:      validator,
			userRepository: userRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().FindByEmailAddress(ctx, req.EmailAddress).Return(user, nil).Once()

		res, err := usecase.SignIn(ctx, req)

		require.ErrorIs(t, err, consts.ErrInvalidCredentials)
		assert.Nil(t, res)
	})

	t.Run("returns an error when session creation fails", func(t *testing.T) {
		ctx := context.Background()
		req := SignInRequest{EmailAddress: "ada@example.com", Password: "correct horse battery staple"}
		expectedErr := errors.New("create user session")
		user := &entity.User{Base: entity.Base{Id: uuid.New()}, EmailAddress: req.EmailAddress}
		require.NoError(t, user.HashPassword(req.Password))

		validator := validation.NewMockIValidator(t)
		userRepository := repositoryUser.NewMockIRepository(t)
		userSessionRepository := repositoryUserSession.NewMockIRepository(t)
		usecase := &Usecase{
			validator:             validator,
			userRepository:        userRepository,
			userSessionRepository: userSessionRepository,
		}

		validator.EXPECT().Struct(mock.Anything).Return(api.ValidationError{}).Once()
		userRepository.EXPECT().FindByEmailAddress(ctx, req.EmailAddress).Return(user, nil).Once()
		userSessionRepository.EXPECT().Create(ctx, mock.AnythingOfType("*entity.UserSession")).Return(expectedErr).Once()

		res, err := usecase.SignIn(ctx, req)

		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, res)
	})
}

func TestUsecase_SignOut(t *testing.T) {
	t.Run("deletes the session token", func(t *testing.T) {
		ctx := context.Background()
		req := SignOutRequest{Token: "session-token"}
		userSessionRepository := repositoryUserSession.NewMockIRepository(t)
		usecase := &Usecase{userSessionRepository: userSessionRepository}

		userSessionRepository.EXPECT().DeleteByToken(ctx, req.Token).Return(nil).Once()

		err := usecase.SignOut(ctx, req)

		require.NoError(t, err)
	})

	t.Run("returns an error when deleting the session fails", func(t *testing.T) {
		ctx := context.Background()
		req := SignOutRequest{Token: "session-token"}
		expectedErr := errors.New("delete session")
		userSessionRepository := repositoryUserSession.NewMockIRepository(t)
		usecase := &Usecase{userSessionRepository: userSessionRepository}

		userSessionRepository.EXPECT().DeleteByToken(ctx, req.Token).Return(expectedErr).Once()

		err := usecase.SignOut(ctx, req)

		require.ErrorIs(t, err, expectedErr)
	})
}

func TestUsecase_Me(t *testing.T) {
	t.Run("returns the current user", func(t *testing.T) {
		userID := uuid.New()
		ctx := current.SetUser(context.Background(), &entity.User{
			Base:         entity.Base{Id: userID},
			Name:         "Ada Lovelace",
			EmailAddress: "ada@example.com",
		})
		usecase := &Usecase{}

		res, err := usecase.Me(ctx)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, userID.String(), res.User.Id)
		assert.Equal(t, "Ada Lovelace", res.User.Name)
		assert.Equal(t, "ada@example.com", res.User.EmailAddress)
	})

	t.Run("returns unauthorized when there is no current user", func(t *testing.T) {
		usecase := &Usecase{}

		res, err := usecase.Me(context.Background())

		require.ErrorIs(t, err, consts.ErrUnauthorized)
		assert.Nil(t, res)
	})
}
