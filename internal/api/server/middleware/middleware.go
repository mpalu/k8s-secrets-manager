package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

// Middleware represents a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Chain applies a sequence of middleware handlers
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func Logging(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("request started")
			next.ServeHTTP(w, r)
		})
	}
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
