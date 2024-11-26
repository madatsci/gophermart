package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

type createOrderError struct{}

func (e *createOrderError) Error() string {
	return "create order error"
}

func (e *createOrderError) IntegrityViolation() bool {
	return true
}

func TestCreateOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStore(ctrl)
	h := newTestHandlers(m)

	path := "/api/user/orders"
	userID := uuid.NewString()

	t.Run("positive case", func(t *testing.T) {
		order := "1234567890003"
		acc := models.Account{
			UserID:             userID,
			CurrentPointsTotal: 500,
			WithdrawnTotal:     1000,
		}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)
		m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)

		req, err := http.NewRequest(http.MethodGet, path, strings.NewReader(order))
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "text/plain")

		r := httptest.NewRecorder()

		h.CreateOrder(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusAccepted, resp.StatusCode, "unexpected response code")
	})

	t.Run("unauthorized user", func(t *testing.T) {
		order := "1234567890003"
		req, err := http.NewRequest(http.MethodGet, path, strings.NewReader(order))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "text/plain")

		r := httptest.NewRecorder()

		h.CreateOrder(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "unexpected response code")
	})

	t.Run("bad request", func(t *testing.T) {
		order := ""
		req, err := http.NewRequest(http.MethodGet, path, strings.NewReader(order))
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "text/plain")

		r := httptest.NewRecorder()

		h.CreateOrder(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "unexpected response code")
	})

	t.Run("unprocessable entity", func(t *testing.T) {
		order := "123123"
		req, err := http.NewRequest(http.MethodGet, path, strings.NewReader(order))
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "text/plain")

		r := httptest.NewRecorder()

		h.CreateOrder(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "unexpected response code")
	})

	t.Run("order already created", func(t *testing.T) {
		order := "1234567890003"
		acc := models.Account{
			UserID:             userID,
			CurrentPointsTotal: 500,
			WithdrawnTotal:     1000,
		}
		existingOrder := models.Order{
			Number: order,
			Account: models.Account{
				UserID: userID,
			},
		}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)
		m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(&createOrderError{})
		m.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(existingOrder, nil)

		req, err := http.NewRequest(http.MethodGet, path, strings.NewReader(order))
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "text/plain")

		r := httptest.NewRecorder()

		h.CreateOrder(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected response code")
	})

	t.Run("order already created by other user", func(t *testing.T) {
		order := "1234567890003"
		acc := models.Account{
			UserID:             userID,
			CurrentPointsTotal: 500,
			WithdrawnTotal:     1000,
		}
		existingOrder := models.Order{
			Number: order,
			Account: models.Account{
				UserID: uuid.NewString(),
			},
		}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)
		m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(&createOrderError{})
		m.EXPECT().GetOrderByNumber(gomock.Any(), order).Return(existingOrder, nil)

		req, err := http.NewRequest(http.MethodGet, path, strings.NewReader(order))
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "text/plain")

		r := httptest.NewRecorder()

		h.CreateOrder(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode, "unexpected response code")
	})
}

func TestGetOrdersHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStore(ctrl)
	h := newTestHandlers(m)

	path := "/api/user/orders"
	userID := uuid.NewString()

	t.Run("positive case", func(t *testing.T) {
		acc := models.Account{
			ID:     uuid.NewString(),
			UserID: userID,
		}
		createdAt := time.Now()
		orders := []models.Order{
			{
				Number:    "1111",
				Status:    models.OrderStatusNew,
				Accrual:   0,
				AccountID: acc.ID,
				CreatedAt: createdAt,
			},
			{
				Number:    "2222",
				Status:    models.OrderStatusProcessing,
				Accrual:   0,
				AccountID: acc.ID,
				CreatedAt: createdAt,
			},
			{
				Number:    "3333",
				Status:    models.OrderStatusInvalid,
				Accrual:   0,
				AccountID: acc.ID,
				CreatedAt: createdAt,
			},
			{
				Number:    "4444",
				Status:    models.OrderStatusProcessed,
				Accrual:   100.5,
				AccountID: acc.ID,
				CreatedAt: createdAt,
			},
		}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)
		m.EXPECT().ListOrdersByAccountID(gomock.Any(), acc.ID, listOrdersLimit).Return(orders, nil)

		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)

		r := httptest.NewRecorder()

		h.GetOrders(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "unexpected response code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "unexpected content type")

		respStr, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		ft := createdAt.Format(time.RFC3339Nano)
		expectedBody := fmt.Sprintf(`[{"number":"1111","status":"NEW","accrual":0,"uploaded_at":"%s"},{"number":"2222","status":"PROCESSING","accrual":0,"uploaded_at":"%s"},{"number":"3333","status":"INVALID","accrual":0,"uploaded_at":"%s"},{"number":"4444","status":"PROCESSED","accrual":100.5,"uploaded_at":"%s"}]`+"\n", ft, ft, ft, ft)
		assert.Equal(t, expectedBody, string(respStr), "unexpected response body")
	})

	t.Run("unauthorized user", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)

		r := httptest.NewRecorder()

		h.GetOrders(r, req)
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

		h.GetOrders(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "unexpected response code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "unexpected content type")
	})

	t.Run("no orders found", func(t *testing.T) {
		acc := models.Account{
			ID:     uuid.NewString(),
			UserID: userID,
		}
		orders := []models.Order{}
		m.EXPECT().GetAccountByUserID(gomock.Any(), userID).Return(acc, nil)
		m.EXPECT().ListOrdersByAccountID(gomock.Any(), acc.ID, listOrdersLimit).Return(orders, nil)

		req, err := http.NewRequest(http.MethodGet, path, http.NoBody)
		require.NoError(t, err)
		ctx := context.WithValue(req.Context(), middleware.AuthenticatedUserKey, userID)
		req = req.WithContext(ctx)

		r := httptest.NewRecorder()

		h.GetOrders(r, req)
		resp := r.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "unexpected response code")
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"), "unexpected content type")

		respStr, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "", string(respStr), "unexpected response body")
	})
}
