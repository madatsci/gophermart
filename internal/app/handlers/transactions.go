package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madatsci/gophermart/internal/app/models"
)

const listWithdrawalsLimit = 100

func (h *Handlers) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	userID, err := ensureUserID(r)
	if err != nil {
		h.handleError("GetWithdrawals", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	acc, err := h.s.GetAccountByUserID(r.Context(), userID)
	if err != nil {
		h.handleError("GetWithdrawals", fmt.Errorf("account not found for user %s", userID))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	txs, err := h.s.GetWithdrawals(r.Context(), acc.ID, models.TxDirectionWithdrawal, listWithdrawalsLimit)
	if err != nil {
		h.handleError("GetWithdrawals", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(txs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(txs); err != nil {
		h.handleError("GetWithdrawals", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
