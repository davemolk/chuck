package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davemolk/chuck/internal/tests/fixture"
	"github.com/davemolk/chuck/internal/tests/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	var gotCtx context.Context
	var gotEmail, gotPW string

	userService := &mock.UserService{
		CreateUserFn: func(ctx context.Context, email, password string) (int64, error) {
			return 0, errors.New("oops")
		},
	}
	h := NewUserHandlers(fixture.TestLogger(t), userService)

	t.Run("email required", func(t *testing.T) {
		body := strings.NewReader(`{"password":"blah"}`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/users", body)

		h.CreateUser(w, r)

		require.False(t, userService.CreateUserCalled)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("pw required", func(t *testing.T) {
		body := strings.NewReader(`{"email":"blah@google"}`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/users", body)

		h.CreateUser(w, r)

		require.False(t, userService.CreateUserCalled)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("handle service error", func(t *testing.T) {
		body := strings.NewReader(`{"email":"blah@google", "password":"roundhouse"}`)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/users", body)

		h.CreateUser(w, r)

		require.True(t, userService.CreateUserCalled)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		userService.ResetCalls()
	})

	t.Run("success", func(t *testing.T) {
		email := "walker@ranger"
		pw := "roundhouse"
		userService = &mock.UserService{
			CreateUserFn: func(ctx context.Context, email, password string) (int64, error) {
				gotCtx = ctx
				gotEmail = email
				gotPW = pw
				return 1, nil
			},
		}
		h.userService = userService
		body := strings.NewReader(`{"email":"walker@ranger", "password":"roundhouse"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/users", body)

		h.CreateUser(w, r)

		require.True(t, userService.CreateUserCalled)
		require.NotNil(t, gotCtx)

		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(t, email, gotEmail)
		require.Equal(t, pw, gotPW)

		var got map[string]any
		err := json.NewDecoder(w.Body).Decode(&got)
		require.NoError(t, err)

		require.Equal(t, float64(1), got["user_id"])
	})
}
