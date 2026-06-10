package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anonychun/bibit/internal/consts"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHttpHandler_SignUp(t *testing.T) {
	t.Run("binds the request, creates a session cookie, and returns ok", func(t *testing.T) {
		e := echo.New()
		body := `{"name":"Ada Lovelace","emailAddress":"ada@example.com","password":"correct horse battery staple"}`
		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("User-Agent", "Go test")
		req.Header.Set("X-Real-IP", "192.0.2.1")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}
		expectedReq := SignUpRequest{
			IpAddress:    "192.0.2.1",
			UserAgent:    "Go test",
			Name:         "Ada Lovelace",
			EmailAddress: "ada@example.com",
			Password:     "correct horse battery staple",
		}

		usecase.EXPECT().SignUp(mock.Anything, expectedReq).Return(&SignUpResponse{Token: "session-token"}, nil).Once()

		err := httpHandler.SignUp(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"ok":true,"meta":null,"data":{"message":"ok"},"errors":null}`, rec.Body.String())

		cookies := rec.Result().Cookies()
		require.Len(t, cookies, 1)
		assert.Equal(t, consts.CookieUserSession, cookies[0].Name)
		assert.Equal(t, "session-token", cookies[0].Value)
		assert.Equal(t, "/", cookies[0].Path)
		assert.True(t, cookies[0].HttpOnly)
	})

	t.Run("returns bind errors", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(`{"name"`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		httpHandler := &HttpHandler{usecase: NewMockIUsecase(t)}

		err := httpHandler.SignUp(ctx)

		require.Error(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns usecase errors", func(t *testing.T) {
		e := echo.New()
		body := `{"name":"Ada Lovelace","emailAddress":"ada@example.com","password":"correct horse battery staple"}`
		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}
		expectedErr := errors.New("sign up")

		usecase.EXPECT().SignUp(mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

		err := httpHandler.SignUp(ctx)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, rec.Result().Cookies())
	})
}

func TestHttpHandler_SignIn(t *testing.T) {
	t.Run("binds the request, creates a session cookie, and returns ok", func(t *testing.T) {
		e := echo.New()
		body := `{"emailAddress":"ada@example.com","password":"correct horse battery staple"}`
		req := httptest.NewRequest(http.MethodPost, "/sign-in", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("User-Agent", "Go test")
		req.Header.Set("X-Real-IP", "192.0.2.1")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}
		expectedReq := SignInRequest{
			IpAddress:    "192.0.2.1",
			UserAgent:    "Go test",
			EmailAddress: "ada@example.com",
			Password:     "correct horse battery staple",
		}

		usecase.EXPECT().SignIn(mock.Anything, expectedReq).Return(&SignInResponse{Token: "session-token"}, nil).Once()

		err := httpHandler.SignIn(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"ok":true,"meta":null,"data":{"message":"ok"},"errors":null}`, rec.Body.String())

		cookies := rec.Result().Cookies()
		require.Len(t, cookies, 1)
		assert.Equal(t, consts.CookieUserSession, cookies[0].Name)
		assert.Equal(t, "session-token", cookies[0].Value)
		assert.Equal(t, "/", cookies[0].Path)
		assert.True(t, cookies[0].HttpOnly)
	})

	t.Run("returns usecase errors", func(t *testing.T) {
		e := echo.New()
		body := `{"emailAddress":"ada@example.com","password":"correct horse battery staple"}`
		req := httptest.NewRequest(http.MethodPost, "/sign-in", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}
		expectedErr := errors.New("sign in")

		usecase.EXPECT().SignIn(mock.Anything, mock.Anything).Return(nil, expectedErr).Once()

		err := httpHandler.SignIn(ctx)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, rec.Result().Cookies())
	})
}

func TestHttpHandler_SignOut(t *testing.T) {
	t.Run("deletes the session and clears the cookie", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sign-out", nil)
		req.AddCookie(&http.Cookie{Name: consts.CookieUserSession, Value: "session-token"})
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}

		usecase.EXPECT().SignOut(mock.Anything, SignOutRequest{Token: "session-token"}).Return(nil).Once()

		err := httpHandler.SignOut(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		setCookie := rec.Header().Get(echo.HeaderSetCookie)
		assert.Contains(t, setCookie, consts.CookieUserSession+"=")
		assert.Contains(t, setCookie, "Path=/")
		assert.Contains(t, setCookie, "HttpOnly")
		assert.Contains(t, setCookie, "Max-Age=0")
	})

	t.Run("returns missing cookie errors", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sign-out", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		httpHandler := &HttpHandler{usecase: NewMockIUsecase(t)}

		err := httpHandler.SignOut(ctx)

		require.Error(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns usecase errors", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/sign-out", nil)
		req.AddCookie(&http.Cookie{Name: consts.CookieUserSession, Value: "session-token"})
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}
		expectedErr := errors.New("sign out")

		usecase.EXPECT().SignOut(mock.Anything, SignOutRequest{Token: "session-token"}).Return(expectedErr).Once()

		err := httpHandler.SignOut(ctx)

		require.ErrorIs(t, err, expectedErr)
		assert.Empty(t, rec.Result().Cookies())
	})
}

func TestHttpHandler_Me(t *testing.T) {
	t.Run("returns the current user response", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}
		res := &MeResponse{}
		res.User.Id = uuid.MustParse("019e925f-3f42-76a0-8518-cb8e51c0b8e2")
		res.User.Name = "Ada Lovelace"
		res.User.EmailAddress = "ada@example.com"

		usecase.EXPECT().Me(mock.Anything).Return(res, nil).Once()

		err := httpHandler.Me(ctx)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"ok":true,"meta":null,"data":{"user":{"id":"019e925f-3f42-76a0-8518-cb8e51c0b8e2","name":"Ada Lovelace","emailAddress":"ada@example.com"}},"errors":null}`, rec.Body.String())
	})

	t.Run("returns usecase errors", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		usecase := NewMockIUsecase(t)
		httpHandler := &HttpHandler{usecase: usecase}
		expectedErr := errors.New("me")

		usecase.EXPECT().Me(mock.Anything).Return(nil, expectedErr).Once()

		err := httpHandler.Me(ctx)

		require.ErrorIs(t, err, expectedErr)
	})
}
