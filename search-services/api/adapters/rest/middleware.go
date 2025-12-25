package rest

import (
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
	"yadro.com/course/api/adapters/auth"
)

type Middleware struct {
	log  *slog.Logger
	auth auth.Authorizer
}

func NewMiddleware(log *slog.Logger, auth auth.Authorizer) *Middleware {
	return &Middleware{
		log:  log,
		auth: auth,
	}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Warn("missing authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			if len(parts) != 2 || parts[0] != "Token" {
				m.log.Warn("invalid authorization header format", "header", authHeader)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		tokenString := parts[1]
		claims, err := m.auth.Validate(tokenString)
		if err != nil {
			m.log.Warn("invalid token", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		m.log.Debug("request authorized", "user", claims.User)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ConcurrencyLimitMiddleware(limit int, next http.Handler) http.Handler {
	sem := make(chan struct{}, limit)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
			next.ServeHTTP(w, r)
		default:
			m.log.Warn("concurrency limit reached", "limit", limit)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
	})
}

func (m *Middleware) RateLimitMiddleware(rps int, next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), 1)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := limiter.Wait(r.Context()); err != nil {
			m.log.Error("rate limiter wait failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
