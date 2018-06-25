// Package middlewares provides common middleware handlers.
package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/kubService"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/redisService"
)

func SetRedisClient(RedisCache *redisService.RedisCache) func(handlerFunc http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "RedisCache", RedisCache))

			next.ServeHTTP(res, req)
		})
	}
}

func SetKubernetesService(serviceKubernetes *kubService.ServiceKubernetes) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "serviceKubernetes", serviceKubernetes))

			next.ServeHTTP(res, req)
		})
	}
}

func SetDB(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "db", db))

			next.ServeHTTP(res, req)
		})
	}
}

func SetSessionStore(sessionStore sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			req = req.WithContext(context.WithValue(req.Context(), "sessionStore", sessionStore))

			next.ServeHTTP(res, req)
		})
	}
}

// MustLogin is a middleware that checks existence of current user.
func MustLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		sessionStore := req.Context().Value("sessionStore").(sessions.Store)
		session, _ := sessionStore.Get(req, "kuberweb-session")
		userRowInterface := session.Values["user"]

		if userRowInterface == nil {
			http.Redirect(res, req, "/login", 302)
			return
		}

		next.ServeHTTP(res, req)
	})
}
