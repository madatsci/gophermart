package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/madatsci/gophermart/internal/app/config"
	"github.com/madatsci/gophermart/pkg/jwt"
	"go.uber.org/zap"
)

// AuthenticatedUserKey should be used to read userID from context.
const AuthenticatedUserKey ctxKey = 0

type (
	Auth struct {
		cookieName string
		jwt        *jwt.JWT
		log        *zap.SugaredLogger
		userID     string
	}

	Options struct {
		Config *config.Config
		JWT    *jwt.JWT
		Log    *zap.SugaredLogger
	}

	ctxKey int
)

// NewAuth creates new auth middleware.
func NewAuth(opts Options) *Auth {
	return &Auth{
		cookieName: opts.Config.AuthCookieName,
		jwt:        opts.JWT,
		log:        opts.Log,
	}
}

// PrivateAPIAuth ensures that user is authenticated.
func (a *Auth) PrivateAPIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(a.cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				a.handleUnauthorized(w, errors.New("no authorisation cookie"))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		userID, err := a.jwt.GetUserID(cookie.Value)
		if err != nil {
			a.handleUnauthorized(w, err)
			return
		}
		if userID == "" {
			a.handleUnauthorized(w, errors.New("token does not contain user ID"))
			return
		}

		a.userID = userID
		a.continueWithUser(w, r, next)
	})
}

func (a *Auth) handleUnauthorized(w http.ResponseWriter, err error) {
	a.log.Debugf("unauthorized attempt to access private API: %s", err)
	w.WriteHeader(http.StatusUnauthorized)
}

func (a *Auth) continueWithUser(w http.ResponseWriter, r *http.Request, next http.Handler) {
	a.log.With("userID", a.userID).Debug("request authorized with user")
	ctx := context.WithValue(r.Context(), AuthenticatedUserKey, a.userID)
	next.ServeHTTP(w, r.WithContext(ctx))
}
