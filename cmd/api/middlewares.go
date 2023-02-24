package main

import (
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

		if accessToken == "" {
			app.UserUnauthorizedResponse(w, r)
		}

		accessTokenMap, err := data.DecodeAccessToken(accessToken)

		if err != nil {
			app.UserUnauthorizedResponse(w, r)
		}

		userId := accessTokenMap["user_id"].(float64)
		_, err = app.models.Users.GetById(int64(userId))

		if err != nil {
			app.UserUnauthorizedResponse(w, r)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) roleMiddleware(roles []data.RoleType, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.UserUnauthorizedResponse(w, r)
		}

		accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if accessToken == "" {
			app.UserUnauthorizedResponse(w, r)
		}

		accessTokenMap, err := data.DecodeAccessToken(accessToken)

		if err != nil {
			app.UserUnauthorizedResponse(w, r)
		}

		userId := accessTokenMap["user_id"].(float64)
		user, err := app.models.Users.GetById(int64(userId))
		if err != nil {
			app.UserUnauthorizedResponse(w, r)
		}

		for _, s := range roles {
			if user.Role == s {
				if s == "OWNER" {
					id, err := app.readIDParam(r)
					if err != nil {
						app.badRequestResponse(w, r, err)
					}
					if id == user.ID {
						next.ServeHTTP(w, r)
					}
				}
				next.ServeHTTP(w, r)
			}
		}

		app.NotEnoughPermissionResponse(w, r)
	})
}
