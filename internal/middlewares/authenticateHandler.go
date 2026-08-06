package middlewares

import (
	"blogapi/internal/utils"
	"net/http"
)

func AuthenticateHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Authenticate(
				func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})(w, r)
		})
}

func AdminOnlyHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			Authenticate(
				func(w http.ResponseWriter, r *http.Request) {
					role, ok := r.Context().Value(CtxUserRole).(string)
					if !ok || role != "admin" {
						utils.JSON(w, http.StatusForbidden, false, "Admin access required", nil)
						return
					}
					next.ServeHTTP(w, r)
				})(w, r)
		})
}
