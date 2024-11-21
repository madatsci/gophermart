package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/server/middleware"
	"github.com/madatsci/gophermart/internal/app/store/database/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBalanceHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStore(ctrl)
	h := newTestHandlers(m)

	path := "/api/user/balance"
	userID := uuid.NewString()

	t.Run("positive case", func(t *testing.T) {
		acc := models.Account{
			CurrentPointsTotal: 500,
			WithdrawnTotal:     1000,
		}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)

		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)

		r := httptest.NewRecorder()

		h.GetBalance(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected response code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "unexpected content type")

		respStr, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, `{"current":500,"withdrawn":1000}`+"\n", string(respStr), "unexpected response body")
	})

	t.Run("unauthorized user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)

		r := httptest.NewRecorder()

		h.GetBalance(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "unexpected response code")
	})

	t.Run("account not found", func(t *testing.T) {
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(models.Account{}, errors.New("account not found"))

		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)

		r := httptest.NewRecorder()

		h.GetBalance(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "unexpected response code")
	})
}
