package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/pkg/luhn"
	"github.com/pkg/errors"
	"github.com/uptrace/bun/driver/pgdriver"
)

const listOrdersLimit = 100

// CreateOrder registers a new order number.
func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := ensureUserID(r)
	if err != nil {
		h.handleError("CreateOrder", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	orderNumber := string(body)
	if orderNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !luhn.VerifyLuhn(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	acc, err := h.s.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			h.handleError("CreateOrder", fmt.Errorf("account not found for user %s", userID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.handleError("CreateOrder", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	order := models.Order{
		ID:        uuid.NewString(),
		AccountID: acc.ID,
		Number:    orderNumber,
		Status:    models.OrderStatusNew,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.s.CreateOrder(r.Context(), order)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			order, err = h.s.GetOrderByNumber(r.Context(), orderNumber)
			if err != nil {
				h.handleError("CreateOrder", err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			if order.Account.UserID == userID {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.handleError("CreateOrder", errors.New("order already created by other user"))
			w.WriteHeader(http.StatusConflict)

			return
		}

		h.handleError("CreateOrder", errors.Wrap(err, "could not create new order"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// GetOrders returns all orders created by the authorized user.
func (h *Handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, err := ensureUserID(r)
	if err != nil {
		h.handleError("GetOrders", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	acc, err := h.s.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		h.handleError("GetOrders", fmt.Errorf("account not found for user %s", userID))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orders, err := h.s.ListOrdersByAccountID(r.Context(), acc.ID, listOrdersLimit)
	if err != nil {
		h.handleError("GetOrders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(orders); err != nil {
		h.handleError("GetOrders", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
