package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/madatsci/gophermart/internal/app/models"
	"github.com/madatsci/gophermart/internal/app/store"
	"github.com/madatsci/gophermart/pkg/hash"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
)

// RegisterUser handles user registration.
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request models.UserReristerRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		h.handleError("RegisterUser", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.Login == "" || request.Password == "" {
		h.handleError("RegisterUser", errors.New("invalid login and/or password"))
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
		var sErr store.StoreError
		if errors.As(err, &sErr) && sErr.IntegrityViolation() {
			h.handleError("RegisterUser", errors.New("user already exists"))
			w.WriteHeader(http.StatusConflict)
			return
		}

		h.handleError("RegisterUser", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	account := models.Account{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.s.CreateAccount(r.Context(), account)
	if err != nil {
		h.handleError("RegisterUser", errors.Wrap(err, "could not create new account"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	h.log.With("ID", user.ID, "login", user.Login).Info("new user registered")
	h.log.With("ID", account.ID, "userID", user.ID).Info("new account created")

	if err = h.authenticateUser(w, user); err != nil {
		h.handleError("RegisterUser", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// LoginUser handles user authentication.
func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request models.UserLoginRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&request); err != nil {
		h.handleError("LoginUser", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.Login == "" || request.Password == "" {
		h.handleError("LoginUser", errInvalidCredentials)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.s.GetUserByLogin(r.Context(), request.Login)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			h.handleError("LoginUser", errInvalidCredentials)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.handleError("LoginUser", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !hash.VerifyPassword(request.Password, user.Password) {
		h.handleError("LoginUser", errInvalidCredentials)
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	if err = h.authenticateUser(w, user); err != nil {
		h.handleError("LoginUser", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) authenticateUser(w http.ResponseWriter, user models.User) error {
	token, err := h.jwt.GetString(user.ID)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: h.c.AuthCookieName, Value: token})

	h.log.With("ID", user.ID, "login", user.Login).Info("user authenticated")

	return nil
}
