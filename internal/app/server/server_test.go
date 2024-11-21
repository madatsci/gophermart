package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/store/database/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
	s := testServer(m)
	defer s.Close()

	path := s.URL + "/api/user/register"
	validRequestBody := `{"login":"john_doe","password":"my_secret_password"}`

	t.Run("positive case", func(t *testing.T) {
		m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(models.User{}, nil)
		m.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(models.Account{}, nil)

		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(validRequestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp := sendRequest(t, req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected response code")
	})

	t.Run("bad request (empty body)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, path, nil)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp := sendRequest(t, req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Unexpected response code")
	})

	t.Run("bad request (empty fields)", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(`{"login":"","password":""}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp := sendRequest(t, req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Unexpected response code")
	})

	t.Run("user already exists", func(t *testing.T) {
		m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(models.User{}, &createUserError{})

		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(validRequestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp := sendRequest(t, req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode, "Unexpected response code")
	})
}

func TestLoginUserHandler(t *testing.T) {

}

func testServer(m *mocks.MockStore) *httptest.Server {
	config := &config.Config{}
	logger := zap.NewNop().Sugar()
	s := New(config, m, logger)

	return httptest.NewServer(s.mux)
}

func sendRequest(t *testing.T, req *http.Request) *http.Response {
	cli := &http.Client{}
	resp, err := cli.Do(req)
	require.NoError(t, err)

	return resp
}
