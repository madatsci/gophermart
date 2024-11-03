package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/pkg/hash"
	"github.com/uptrace/bun/driver/pgdriver"
)

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request models.UserReristerRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		h.handleError("RegisterUser", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if request.Login == "" || request.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pwdHash, err := hash.HashPassword(request.Password)
	if err != nil {
		h.handleError("RegisterUser", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:        uuid.NewString(),
		Login:     request.Login,
		Password:  pwdHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.s.CreateUser(r.Context(), user)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			h.handleError("RegisterUser", errors.New("user already exists"))
			w.WriteHeader(http.StatusConflict)
			return
		}

		h.handleError("RegisterUser", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	h.log.With("ID", user.ID, "login", user.Login).Info("new user registered")

	w.WriteHeader(http.StatusCreated)
}
