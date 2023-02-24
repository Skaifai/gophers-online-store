package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Skaifai/gophers-online-store/internal/data"
)

func (app *application) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.UserUnauthorizedResponse(w, r)
		}

		accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

		accessToken, err := data.DecodeAccessToken(accessToken)
		if err != nil {
			app.UserUnauthorizedResponse(w, r)
		}

		fmt.Println(accessToken["user_id"].(int64))

		next.ServeHTTP(w, r)
	})
}
