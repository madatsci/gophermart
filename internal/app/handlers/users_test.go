package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/store/database/mocks"
	"github.com/madatsci/gophermart/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createUserError struct{}

func (e *createUserError) Error() string {
	return "create user error"
}

func (e *createUserError) IntegrityViolation() bool {
	return true
}

func TestRegisterUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStore(ctrl)
	h := newTestHandlers(m)

	path := "/api/user/register"
	validRequestBody := `{"login":"john_doe","password":"my_secret_password"}`

	t.Run("positive case", func(t *testing.T) {
		m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
		m.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(models.Account{}, nil)

		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(validRequestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.RegisterUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected response code")

		cookies := resp.Cookies()
		assert.Equal(t, 1, len(cookies), "expected auth cookie in response")
		found := false
		for _, c := range cookies {
			if c.Name == h.c.AuthCookieName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %s cookie in response", h.c.AuthCookieName)
		}
	})

	t.Run("bad request (empty body)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, path, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.RegisterUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "unexpected response code")
	})

	t.Run("bad request (empty fields)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(`{"login":"","password":""}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.RegisterUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "unexpected response code")
	})

	t.Run("user already exists", func(t *testing.T) {
		m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(models.User{}, &createUserError{})

		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(validRequestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.RegisterUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode, "unexpected response code")
	})
}

func TestLoginUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStore(ctrl)
	h := newTestHandlers(m)

	path := "/api/user/login"
	validRequestBody := `{"login":"john_doe","password":"my_secret_password"}`

	t.Run("positive case", func(t *testing.T) {
		pwdHash, err := hash.HashPassword("my_secret_password")
		if err != nil {
			t.Fatal(err)
		}
		user := models.User{
			Login:    "john_doe",
			Password: pwdHash,
		}
		m.EXPECT().GetUserByLogin(gomock.Any(), "john_doe").Return(user, nil)

		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(validRequestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.LoginUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected response code")

		cookies := resp.Cookies()
		assert.Equal(t, 1, len(cookies), "expected auth cookie in response")
		found := false
		for _, c := range cookies {
			if c.Name == h.c.AuthCookieName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %s cookie in response", h.c.AuthCookieName)
		}
	})

	t.Run("bad request (empty body)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, path, http.NoBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.LoginUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "unexpected response code")
	})

	t.Run("bad request (empty fields)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(`{"login":"","password":""}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.LoginUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "unexpected response code")
	})

	t.Run("user not found", func(t *testing.T) {
		err := errors.New("sql: no rows in result set")
		m.EXPECT().GetUserByLogin(gomock.Any(), "john_doe").Return(models.User{}, err)

		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(validRequestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.LoginUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "unexpected response code")
	})

	t.Run("invalid password", func(t *testing.T) {
		user := models.User{
			Login:    "john_doe",
			Password: "some_hash",
		}
		m.EXPECT().GetUserByLogin(gomock.Any(), "john_doe").Return(user, nil)

		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(validRequestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()

		h.LoginUser(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "unexpected response code")
	})
}
