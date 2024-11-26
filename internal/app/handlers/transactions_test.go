package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/server/middleware"
	"github.com/madatsci/gophermart/internal/app/store/database/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWithdrawalsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStore(ctrl)
	h := newTestHandlers(m)

	path := "/api/user/withdrawals"
	userID := uuid.NewString()

	t.Run("positive case", func(t *testing.T) {
		acc := models.Account{
			ID:     uuid.NewString(),
			UserID: userID,
		}
		createdAt := time.Now()
		txs := []models.Transaction{
			{
				Amount:      100,
				OrderNumber: "1111",
				Direction:   models.TxDirectionWithdrawal,
				AccountID:   acc.ID,
				CreatedAt:   createdAt,
			},
			{
				Amount:      0,
				OrderNumber: "2222",
				Direction:   models.TxDirectionWithdrawal,
				AccountID:   acc.ID,
				CreatedAt:   createdAt,
			},
		}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)
		m.EXPECT().GetWithdrawals(gomock.Any(), acc.ID, models.TxDirectionWithdrawal, listWithdrawalsLimit).Return(txs, nil)

		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)

		r := httptest.NewRecorder()

		h.GetWithdrawals(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected response code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "unexpected content type")

		respStr, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		ft := createdAt.Format(time.RFC3339Nano)
		expectedBody := fmt.Sprintf(`[{"sum":100,"order":"1111","processed_at":"%s"},{"sum":0,"order":"2222","processed_at":"%s"}]`+"\n", ft, ft)
		assert.Equal(t, expectedBody, string(respStr), "unexpected response body")
	})

	t.Run("unauthorized user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)

		r := httptest.NewRecorder()

		h.GetWithdrawals(r, req)
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

		h.GetWithdrawals(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "unexpected response code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "unexpected content type")
	})

	t.Run("no transactions found", func(t *testing.T) {
		acc := models.Account{
			ID:     uuid.NewString(),
			UserID: userID,
		}
		txs := []models.Transaction{}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)
		m.EXPECT().GetWithdrawals(gomock.Any(), acc.ID, models.TxDirectionWithdrawal, listWithdrawalsLimit).Return(txs, nil)

		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)

		r := httptest.NewRecorder()

		h.GetWithdrawals(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "unexpected response code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "unexpected content type")

		respStr, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "", string(respStr), "unexpected response body")
	})
}
