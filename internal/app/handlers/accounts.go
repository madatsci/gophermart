package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
