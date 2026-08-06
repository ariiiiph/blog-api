package middlewares

import (
	"blogapi/internal/utils"
	"context"
	"net/http"
	"strings"
)

const (
	CtxUserID        string = "userId"
	CtxUserEmail     string = "email"
	CtxUserRole      string = "role"
	CtxAuthorization string = "authorization"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get(CtxAuthorization))
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}
		accessToken := strings.TrimSpace(authHeader[7:])

		userId, email, role, err := utils.VerifyJWT(accessToken)
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, false, "Unathorized", nil)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxUserID, userId)
		ctx = context.WithValue(ctx, CtxUserEmail, email)
		ctx = context.WithValue(ctx, CtxUserRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
