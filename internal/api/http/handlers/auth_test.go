package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service/auth"
	"github.com/davemolk/chuck/internal/tests/fixture"
	"github.com/davemolk/chuck/internal/tests/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	var gotCtx context.Context
	var gotEmail, gotPW string

	authService := &mock.AuthService{
		LoginFn: func(ctx context.Context, email, password string) (*domain.Token, error) {
			return nil, auth.ErrInvalidCredentials
		},
	}
	h := NewAuthHandlers(fixture.TestLogger(t), authService)

	t.Run("email required", func(t *testing.T) {
		body := strings.NewReader(`{"password":"blah"}`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/auth/login", body)

		h.Login(w, r)

		require.False(t, authService.LoginFnCalled)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("pw required", func(t *testing.T) {
		body := strings.NewReader(`{"email":"blah@google"}`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/auth/login", body)

		h.Login(w, r)

		require.False(t, authService.LoginFnCalled)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("handle service error", func(t *testing.T) {
		body := strings.NewReader(`{"email":"blah@google", "password":"roundhouse"}`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/auth/login", body)

		h.Login(w, r)

		require.True(t, authService.LoginFnCalled)
		require.Equal(t, http.StatusUnauthorized, w.Code)
		authService.ResetCalls()
	})

	t.Run("success", func(t *testing.T) {
		email := "walker@ranger"
		pw := "roundhouse"
		authService := &mock.AuthService{
			LoginFn: func(ctx context.Context, email, password string) (*domain.Token, error) {
				gotCtx = ctx
				gotEmail = email
				gotPW = pw
				return &domain.Token{
					Plaintext: "blah",
				}, nil
			},
		}

		h.authService = authService
		body := strings.NewReader(`{"email":"walker@ranger", "password":"roundhouse"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/auth/login", body)

		h.Login(w, r)

		require.True(t, authService.LoginFnCalled)
		require.NotNil(t, gotCtx)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, email, gotEmail)
		require.Equal(t, pw, gotPW)

		var got domain.Token
		err := json.NewDecoder(w.Body).Decode(&got)
		require.NoError(t, err)

		require.Equal(t, "blah", got.Plaintext)
	})
}
