package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/store"
	"github.com/madatsci/gophermart/pkg/luhn"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// GetBalance returns user account balance.
func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, err := ensureUserID(r)
	if err != nil {
		h.handleError("GetBalance", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	acc, err := h.s.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			h.handleError("GetBalance", fmt.Errorf("account not found for user %s", userID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.handleError("GetBalance", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(acc); err != nil {
		h.handleError("GetBalance", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// WithdrawPoints withdraws points from user account
func (h *Handlers) WithdrawPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, err := ensureUserID(r)
	if err != nil {
		h.handleError("WithdrawPoints", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request models.BalanceWithdrawRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		h.handleError("WithdrawPoints", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if request.Order == "" || request.Sum.Cmp(decimal.NewFromInt(0)) != 1 {
		h.handleError("WithdrawPoints", errors.New("invalid parameters"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !luhn.VerifyLuhn(request.Order) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	acc, err := h.s.WithdrawBalance(r.Context(), userID, request.Order, request.Sum)
	if err != nil {
		h.handleError("WithdrawPoints", err)

		var balanceErr *store.NotEnoughBalanceError
		if errors.As(err, &balanceErr) {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(acc); err != nil {
		h.handleError("WithdrawPoints", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
